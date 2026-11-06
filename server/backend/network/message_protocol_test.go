package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

type fragmentedReader struct{ io.Reader }

func (r fragmentedReader) Read(p []byte) (int, error) {
	if len(p) > 7 {
		p = p[:7]
	}
	return r.Reader.Read(p)
}

func TestReadBinaryMessageStream(t *testing.T) {
	// A command response larger than the old 4096-byte socket buffer, followed
	// immediately by another frame. Fragmentation can occur inside any field.
	payload := bytes.Repeat([]byte("adapter details\r\n"), 1000)
	first := models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_COMMAND_RESULT, MessageID: 99, DataSize: uint32(len(payload))}, Data: payload}
	second := models.InternalMessage{Header: models.MessageHeader{MessageType: models.MESSAGE_HEARTBEAT, MessageID: 100}}
	data := append(serializeBinaryMessage(first), serializeBinaryMessage(second)...)
	for _, fragmented := range []bool{false, true} {
		var reader io.Reader = bytes.NewReader(data)
		if fragmented {
			reader = fragmentedReader{reader}
		}
		got, err := readBinaryMessage(reader)
		if err != nil || got.Header != first.Header || !bytes.Equal(got.Data, payload) {
			t.Fatalf("fragmented=%t: bad first frame: %v", fragmented, err)
		}
		got, err = readBinaryMessage(reader)
		if err != nil || got.Header != second.Header || len(got.Data) != 0 {
			t.Fatalf("fragmented=%t: bad second frame: %v", fragmented, err)
		}
		if _, err := readBinaryMessage(reader); !errors.Is(err, io.EOF) {
			t.Fatalf("expected EOF: %v", err)
		}
	}
}

func TestReadBinaryMessageRejectsTruncatedAndOversizedFrames(t *testing.T) {
	frame := serializeBinaryMessage(models.InternalMessage{Header: models.MessageHeader{DataSize: 3}, Data: []byte("abc")})
	for _, size := range []int{5, len(frame) - 1} {
		if _, err := readBinaryMessage(bytes.NewReader(frame[:size])); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("size %d: expected unexpected EOF, got %v", size, err)
		}
	}
	var header [12]byte
	binary.LittleEndian.PutUint32(header[4:8], maxMessageDataSize+1)
	if _, err := readBinaryMessage(bytes.NewReader(header[:])); err == nil || errors.Is(err, io.EOF) {
		t.Fatalf("expected length rejection before reading payload: %v", err)
	}
}

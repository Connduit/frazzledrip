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
	MESSAGE_LIST_DIRECTORY     = 13
	MESSAGE_DIRECTORY_LISTING  = 14
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

// Attributes associated with Clients
type Client struct {
	BasicInfoStatus string `json:"BasicInfoStatus,omitempty"`
	Type            string
	Listener        string
	ID              string
	Conn            net.Conn
	LastSeen        time.Time
	SystemInfo      SystemInfo
	IPAddress       string
	IsActive        bool
	Messages        []ClientMessage // Store messages for this client
	CommandResults  []CommandResult `json:"CommandResults,omitempty"`
}

type CommandResult struct {
	Command    string    `json:"command"`
	Output     string    `json:"output"`
	ReceivedAt time.Time `json:"received_at"`
}

type SystemInfo struct {
	Hostname                string   `json:"hostname"`
	OS                      string   `json:"os"`
	Arch                    string   `json:"arch"`
	User                    string   `json:"user"`
	Domain                  string   `json:"domain,omitempty"`
	IPAddresses             []string `json:"ip_addresses,omitempty"`
	MACAddresses            []string `json:"mac_addresses,omitempty"`
	Manufacturer            string   `json:"manufacturer,omitempty"`
	Model                   string   `json:"model,omitempty"`
	BootTime                string   `json:"boot_time,omitempty"`
	MajorVersion            *uint32  `json:"major_version,omitempty"`
	MinorVersion            *uint32  `json:"minor_version,omitempty"`
	BuildNumber             *uint32  `json:"build_number,omitempty"`
	ProductType             *uint32  `json:"product_type,omitempty"`
	TotalPhysicalMemory     *uint64  `json:"total_physical_memory,omitempty"`
	AvailablePhysicalMemory *uint64  `json:"available_physical_memory,omitempty"`
	NumberOfProcessors      *uint32  `json:"number_of_processors,omitempty"`
	ProcessorArchitecture   *uint32  `json:"processor_architecture,omitempty"`
}

type CommandRequest struct {
	Client  string `json:"client"`
	Command string `json:"command"`
}

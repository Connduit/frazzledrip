package network

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

// The current Windows client returns key/value text as COMMAND_RESULT rather
// than SYS_INFO. Require its version signature before interpreting that text
// as metadata, so ordinary command output doesn't overwrite the session.
func mergeSystemInfo(current models.SystemInfo, data []byte, requireVersion bool) models.SystemInfo {
	fields := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if ok {
			fields[strings.TrimSpace(key)] = strings.TrimSpace(strings.TrimRight(strings.TrimSpace(value), "\\\x00"))
		}
	}
	number := func(key string, bits int) (uint64, bool) {
		value, err := strconv.ParseUint(fields[key], 10, bits)
		return value, err == nil
	}
	if requireVersion {
		for _, key := range []string{"dwMajorVersion", "dwMinorVersion", "dwBuildNumber"} {
			if _, ok := number(key, 32); !ok {
				return current
			}
		}
	}
	for _, field := range []struct {
		key    string
		target **uint32
	}{
		{"dwMajorVersion", &current.MajorVersion},
		{"dwMinorVersion", &current.MinorVersion},
		{"dwBuildNumber", &current.BuildNumber},
		{"numberOfProcessors", &current.NumberOfProcessors},
		{"processorArchitecture", &current.ProcessorArchitecture},
	} {
		if value, ok := number(field.key, 32); ok {
			v := uint32(value)
			*field.target = &v
		}
	}
	if value, ok := number("wProductType", 32); ok && value >= 1 && value <= 3 {
		v := uint32(value)
		current.ProductType = &v
	}
	for _, field := range []struct {
		key    string
		target **uint64
	}{
		{"totalPhysicalMemory", &current.TotalPhysicalMemory},
		{"availablePhysicalMemory", &current.AvailablePhysicalMemory},
	} {
		if value, ok := number(field.key, 64); ok {
			v := value
			*field.target = &v
		}
	}
	if current.MajorVersion != nil && current.MinorVersion != nil {
		current.OS = fmt.Sprintf("Windows %d.%d", *current.MajorVersion, *current.MinorVersion)
		if current.BuildNumber != nil {
			current.OS += fmt.Sprintf(" (build %d)", *current.BuildNumber)
		}
	}
	if current.ProcessorArchitecture != nil {
		switch *current.ProcessorArchitecture {
		case 0:
			current.Arch = "x86"
		case 5:
			current.Arch = "ARM"
		case 6:
			current.Arch = "IA64"
		case 9:
			current.Arch = "x64"
		case 12:
			current.Arch = "ARM64"
		default:
			current.Arch = fmt.Sprintf("Unknown (%d)", *current.ProcessorArchitecture)
		}
	}
	return current
}

package network

import (
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

type pendingCommand struct {
	TaskID    string
	Automatic bool
	ID        uint32
	Command   string
}

// Protected by clientsMutex. Legacy clients return 105 for every shell result.
var pendingCommands = make(map[string]pendingCommand)
var identityPattern = regexp.MustCompile(`^(?:[\p{L}\p{N}_.@$ -]+\\[\p{L}\p{N}_.@$ -]+|[\p{L}\p{N}_.@$-]+)$`)
var versionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)\b`)

func commandKey(command string) string {
	return strings.ToLower(strings.Join(strings.Fields(command), " "))
}

// Only parse known commands with known argument forms, never arbitrary shell
// output or pipelines. Unknown/localized output remains available verbatim.
func mergeCommandInfo(info models.SystemInfo, command, output string) models.SystemInfo {
	key := commandKey(command)
	text := strings.TrimSpace(strings.TrimRight(output, "\x00"))
	if key == "whoami" {
		if identityPattern.MatchString(text) {
			info.User = text
		}
		return info
	}
	if key != "ipconfig" && key != "ipconfig /all" && key != "systeminfo" {
		return info
	}
	var ips, macs []string
	add := func(values *[]string, value string) {
		for _, existing := range *values {
			if existing == value {
				return
			}
		}
		*values = append(*values, value)
	}
	for _, line := range strings.Split(text, "\n") {
		label, value, ok := strings.Cut(strings.TrimSpace(line), ":")
		if !ok {
			continue
		}
		label = strings.ToLower(strings.TrimSpace(strings.TrimRight(strings.TrimSpace(label), ". ")))
		value = strings.TrimSpace(strings.TrimRight(strings.TrimSpace(value), "\x00"))
		if value == "" {
			continue
		}
		if key == "ipconfig" || key == "ipconfig /all" {
			switch {
			case label == "host name":
				info.Hostname = value
			case label == "primary dns suffix":
				info.Domain = value
			case strings.Contains(label, "ipv4 address") || strings.Contains(label, "ipv6 address"):
				address := strings.Fields(value)[0]
				address = strings.Split(address, "(")[0]
				address = strings.Split(address, "%")[0]
				if ip := net.ParseIP(address); ip != nil {
					add(&ips, ip.String())
				}
			case label == "physical address":
				if mac, err := net.ParseMAC(value); err == nil {
					add(&macs, mac.String())
				}
			}
			continue
		}
		switch label {
		case "host name":
			info.Hostname = value
		case "os name":
			info.OS = value
		case "os version":
			if parts := versionPattern.FindStringSubmatch(value); parts != nil {
				targets := []**uint32{&info.MajorVersion, &info.MinorVersion, &info.BuildNumber}
				for i, target := range targets {
					if n, err := strconv.ParseUint(parts[i+1], 10, 32); err == nil {
						v := uint32(n)
						*target = &v
					}
				}
			}
		case "system type":
			switch strings.ToLower(value) {
			case "x64-based pc":
				info.Arch = "x64"
			case "x86-based pc":
				info.Arch = "x86"
			case "arm64-based pc":
				info.Arch = "ARM64"
			}
		case "domain":
			info.Domain = value
		case "system manufacturer":
			info.Manufacturer = value
		case "system model":
			info.Model = value
		case "system boot time":
			info.BootTime = value
		case "total physical memory":
			if n, ok := parseMemory(value); ok {
				info.TotalPhysicalMemory = &n
			}
		case "available physical memory":
			if n, ok := parseMemory(value); ok {
				info.AvailablePhysicalMemory = &n
			}
		case "processor(s)":
			if n, err := strconv.ParseUint(strings.Fields(value)[0], 10, 32); err == nil {
				v := uint32(n)
				info.NumberOfProcessors = &v
			}
		}
	}
	if len(ips) > 0 {
		info.IPAddresses = ips
	}
	if len(macs) > 0 {
		info.MACAddresses = macs
	}
	return info
}

func parseMemory(value string) (uint64, bool) {
	fields := strings.Fields(value)
	if len(fields) != 2 {
		return 0, false
	}
	n, err := strconv.ParseUint(strings.ReplaceAll(fields[0], ",", ""), 10, 64)
	multipliers := map[string]uint64{"KB": 1 << 10, "MB": 1 << 20, "GB": 1 << 30, "BYTES": 1}
	multiplier, ok := multipliers[strings.ToUpper(fields[1])]
	if err != nil || !ok || n > ^uint64(0)/multiplier {
		return 0, false
	}
	return n * multiplier, true
}

func storeCommandResult(client *models.Client, command, output string, at time.Time) {
	// Latest result per exact command, capped to 20 command variants per session.
	results := make([]models.CommandResult, 0, len(client.CommandResults)+1)
	results = append(results, models.CommandResult{Command: command, Output: output, ReceivedAt: at})
	for _, previous := range client.CommandResults {
		if commandKey(previous.Command) != commandKey(command) && len(results) < 20 {
			results = append(results, previous)
		}
	}
	client.CommandResults = results
}

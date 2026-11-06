package network

import (
	"fmt"
	"sort"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

// LinkSavedSession is an explicit operator decision, never an IP-based match.
// The live connection and its ID remain intact; the consumed offline record is removed.
func LinkSavedSession(savedID, liveID string) error {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	var saved, live *models.Client
	for _, client := range clients {
		if client.ID == savedID {
			saved = client
		}
		if client.ID == liveID {
			live = client
		}
	}
	if saved == nil || live == nil {
		return fmt.Errorf("session no longer exists; refresh and select again")
	}
	if saved == live || saved.IsActive || !live.IsActive || live.Conn == nil {
		return fmt.Errorf("select an offline saved session and a different active connection")
	}
	restoreSessionData(saved, live)
	delete(clients, saved.IPAddress)
	delete(pendingCommands, saved.IPAddress)
	return nil
}

// Call with clientsMutex held. Preserve the incoming connection and its activity.
func restoreSessionData(saved, live *models.Client) {
	if live.BasicInfoStatus == "" {
		live.BasicInfoStatus = saved.BasicInfoStatus
	}
	live.SystemInfo = restoreSystemInfo(saved.SystemInfo, live.SystemInfo)
	messages := append(append([]models.ClientMessage(nil), saved.Messages...), live.Messages...)
	sort.SliceStable(messages, func(i, j int) bool { return messages[i].Timestamp.Before(messages[j].Timestamp) })
	if len(messages) > 100 {
		messages = messages[len(messages)-100:]
	}
	live.Messages = messages
	results := append(append([]models.CommandResult(nil), live.CommandResults...), saved.CommandResults...)
	sort.SliceStable(results, func(i, j int) bool { return results[i].ReceivedAt.After(results[j].ReceivedAt) })
	seen := make(map[string]bool)
	live.CommandResults = nil
	for _, result := range results {
		key := commandKey(result.Command)
		if !seen[key] && len(live.CommandResults) < 20 {
			live.CommandResults = append(live.CommandResults, result)
			seen[key] = true
		}
	}
}

// Keep facts already collected on the live connection, filling only missing fields.
func restoreSystemInfo(saved, live models.SystemInfo) models.SystemInfo {
	if live.Hostname == "" {
		live.Hostname = saved.Hostname
	}
	if live.OS == "" {
		live.OS = saved.OS
	}
	if live.Arch == "" {
		live.Arch = saved.Arch
	}
	if live.User == "" {
		live.User = saved.User
	}
	if live.Domain == "" {
		live.Domain = saved.Domain
	}
	if len(live.IPAddresses) == 0 {
		live.IPAddresses = saved.IPAddresses
	}
	if len(live.MACAddresses) == 0 {
		live.MACAddresses = saved.MACAddresses
	}
	if live.Manufacturer == "" {
		live.Manufacturer = saved.Manufacturer
	}
	if live.Model == "" {
		live.Model = saved.Model
	}
	if live.BootTime == "" {
		live.BootTime = saved.BootTime
	}
	if live.MajorVersion == nil {
		live.MajorVersion = saved.MajorVersion
	}
	if live.MinorVersion == nil {
		live.MinorVersion = saved.MinorVersion
	}
	if live.BuildNumber == nil {
		live.BuildNumber = saved.BuildNumber
	}
	if live.ProductType == nil {
		live.ProductType = saved.ProductType
	}
	if live.TotalPhysicalMemory == nil {
		live.TotalPhysicalMemory = saved.TotalPhysicalMemory
	}
	if live.AvailablePhysicalMemory == nil {
		live.AvailablePhysicalMemory = saved.AvailablePhysicalMemory
	}
	if live.NumberOfProcessors == nil {
		live.NumberOfProcessors = saved.NumberOfProcessors
	}
	if live.ProcessorArchitecture == nil {
		live.ProcessorArchitecture = saved.ProcessorArchitecture
	}
	return live
}

package network

import (
	"fmt"
	"github.com/Connduit/frazzledrip/server/backend/models"
	"github.com/Connduit/frazzledrip/server/backend/state"
	"log"
)

// Assigned only at startup, before accepting connections.
var persistentState *state.Store

func ConfigurePersistentState(store *state.Store) error {
	records, err := store.LoadSessions()
	if err != nil {
		return err
	}
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if len(clients) != 0 {
		return fmt.Errorf("persistent state must be configured before accepting clients")
	}
	persistentState = store
	for _, client := range records {
		clients[client.IPAddress] = client
	}
	return nil
}
func PersistentState() *state.Store { return persistentState }
func persistSessionLocked(client *models.Client, event string) {
	if persistentState != nil {
		if err := persistentState.SaveSession(client, event); err != nil {
			log.Printf("Session persistence error: %v", err)
		}
	}
}
func newTrackedTaskLocked(client *models.Client, command, origin string, id uint32) (string, error) {
	if persistentState == nil {
		return "", nil
	}
	return persistentState.NewTask(client.ID, command, origin, id)
}
func taskTransition(id, status, output string) {
	if id != "" && persistentState != nil {
		if err := persistentState.TransitionTask(id, status, output); err != nil {
			log.Printf("Task persistence error: %v", err)
		}
	}
}
func interruptTaskLocked(addr string) {
	if task, ok := pendingCommands[addr]; ok {
		taskTransition(task.TaskID, "interrupted", "Connection closed before a result was recorded")
	}
}

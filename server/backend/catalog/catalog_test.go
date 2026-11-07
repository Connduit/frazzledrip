package catalog

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Connduit/frazzledrip/server/backend/state"
)

func TestSimulationsAreLabeledAndNeverDispatched(t *testing.T) {
	store, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	// No connection, client object, command executor, or script runner is passed.
	for _, tool := range List() {
		task, err := Simulate(store, "test-agent", tool.ID, "test note")
		if err != nil {
			t.Fatal(err)
		}
		if !task.Simulation || task.Status != "completed" || task.Origin != "simulation" || task.RequestID != 0 || !strings.Contains(task.Output, "No client was contacted") || !strings.HasPrefix(task.Command, "[SIMULATION]") {
			t.Fatalf("unlabeled simulation: %+v", task)
		}
	}
	events, err := store.Audit(100)
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events {
		if event.Event == "task.sent" {
			t.Fatal("simulation was dispatched")
		}
	}
	if _, err := Simulate(store, "test-agent", "unknown", ""); err == nil {
		t.Fatal("unknown tool accepted")
	}
	if _, err := Simulate(store, "test-agent", List()[0].ID, strings.Repeat("x", 501)); err == nil {
		t.Fatal("oversized note accepted")
	}
	changed := List()
	changed[0].ID = "modified"
	if List()[0].ID == "modified" {
		t.Fatal("caller mutated catalog")
	}
}

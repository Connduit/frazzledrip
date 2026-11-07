package state

import (
	"net"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Connduit/frazzledrip/server/backend/models"
)

func TestStateSurvivesRestartWithoutReplay(t *testing.T) {
	file := filepath.Join(t.TempDir(), "server.db")
	store, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	conn, peer := net.Pipe()
	defer conn.Close()
	defer peer.Close()
	client := &models.Client{ID: "test-agent", IPAddress: "127.0.0.1:1000", Conn: conn, IsActive: true, BasicInfoStatus: "collecting", SystemInfo: models.SystemInfo{Hostname: "TEST-COMPUTER"}, Messages: []models.ClientMessage{{Content: "test output"}}}
	if err := store.SaveSession(client, "session.connected"); err != nil {
		t.Fatal(err)
	}
	pending, err := store.NewTask(client.ID, "test-command", "test", 12)
	if err != nil {
		t.Fatal(err)
	}
	store.TransitionTask(pending, "sent", "")
	completed, _ := store.NewTask(client.ID, "test-complete", "test", 13)
	store.TransitionTask(completed, "completed", "test result")
	// A fast result must win over the writer's delayed sent notification.
	store.TransitionTask(completed, "sent", "")
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	restored, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer restored.Close()
	sessions, err := restored.LoadSessions()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].IsActive || sessions[0].Conn != nil || sessions[0].BasicInfoStatus != "interrupted" || sessions[0].SystemInfo.Hostname != "TEST-COMPUTER" || len(sessions[0].Messages) != 1 {
		t.Fatalf("bad restored session: %+v", sessions)
	}
	task, _ := restored.Task(pending)
	if task.Status != "interrupted" {
		t.Fatal("unfinished task not interrupted")
	}
	task, _ = restored.Task(completed)
	if task.Status != "completed" || task.Output != "test result" {
		t.Fatal("completed task regressed")
	}
	events, err := restored.Audit(100)
	if err != nil || len(events) < 5 {
		t.Fatalf("audit not preserved: %d %v", len(events), err)
	}
	for i := 1; i < len(events); i++ {
		if events[i].ID >= events[i-1].ID {
			t.Fatal("audit not newest-first")
		}
	}
}
func TestTaskOutputLimitAndSessionMerge(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	store.SaveSession(&models.Client{ID: "old", IPAddress: "test:1"}, "")
	id, _ := store.NewTask("old", "test", "fixture", 99)
	store.TransitionTask(id, "failed", strings.Repeat("x", (1<<20)+1))
	if err := store.MergeSession("old", &models.Client{ID: "new", IPAddress: "test:2"}); err != nil {
		t.Fatal(err)
	}
	tasks, _ := store.Tasks("new")
	if len(tasks) != 1 || len(tasks[0].Output) != 1<<20 || !tasks[0].OutputTruncated || tasks[0].Status != "failed" {
		t.Fatal("merge/output bounds failed")
	}
	sessions, _ := store.LoadSessions()
	if len(sessions) != 1 || sessions[0].ID != "new" {
		t.Fatal("saved session not consumed")
	}
}

func TestUnknownSchemaIsNotOverwritten(t *testing.T) {
	file := filepath.Join(t.TempDir(), "future.db")
	store, err := Open(file)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := store.db.Begin(true)
	if err != nil {
		t.Fatal(err)
	}
	if err := tx.Bucket(schemaBucket).Put([]byte("version"), []byte("999")); err != nil {
		tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	store.Close()
	for i := 0; i < 2; i++ {
		if reopened, err := Open(file); err == nil {
			reopened.Close()
			t.Fatal("unsupported schema accepted or overwritten")
		}
	}
}

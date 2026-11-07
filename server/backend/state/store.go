package state

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Connduit/frazzledrip/server/backend/models"
	bolt "go.etcd.io/bbolt"
)

var sessionBucket = []byte("sessions")
var taskBucket = []byte("tasks")
var auditBucket = []byte("audit")
var schemaBucket = []byte("schema")

type Store struct{ db *bolt.DB }
type Task struct {
	ID              string    `json:"id"`
	AgentID         string    `json:"agent_id"`
	RequestID       uint32    `json:"request_id"`
	Command         string    `json:"command"`
	Origin          string    `json:"origin"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Output          string    `json:"output,omitempty"`
	OutputTruncated bool      `json:"output_truncated,omitempty"`
}
type AuditEvent struct {
	ID      uint64    `json:"id"`
	At      time.Time `json:"at"`
	Event   string    `json:"event"`
	AgentID string    `json:"agent_id,omitempty"`
	TaskID  string    `json:"task_id,omitempty"`
	Source  string    `json:"source,omitempty"`
	Detail  string    `json:"detail,omitempty"`
}

func Open(filename string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return nil, err
	}
	db, err := bolt.Open(filename, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	store := &Store{db: db}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, name := range [][]byte{sessionBucket, taskBucket, auditBucket, schemaBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		version := tx.Bucket(schemaBucket).Get([]byte("version"))
		if version != nil && string(version) != "1" {
			return fmt.Errorf("unsupported state database schema %q", version)
		}
		return tx.Bucket(schemaBucket).Put([]byte("version"), []byte("1"))
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	if err = store.recoverTasks(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}
func (store *Store) Close() error { return store.db.Close() }
func put(bucket *bolt.Bucket, key []byte, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return bucket.Put(key, data)
}
func audit(tx *bolt.Tx, event AuditEvent) error {
	bucket := tx.Bucket(auditBucket)
	id, err := bucket.NextSequence()
	if err != nil {
		return err
	}
	event.ID = id
	if event.At.IsZero() {
		event.At = time.Now().UTC()
	}
	var key [8]byte
	binary.BigEndian.PutUint64(key[:], id)
	return put(bucket, key[:], event)
}
func (store *Store) SaveSession(client *models.Client, event string) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		if err := put(tx.Bucket(sessionBucket), []byte(client.ID), client); err != nil {
			return err
		}
		if event != "" {
			return audit(tx, AuditEvent{Event: event, AgentID: client.ID, Detail: client.IPAddress})
		}
		return nil
	})
}
func (store *Store) LoadSessions() ([]*models.Client, error) {
	sessions := []*models.Client{}
	err := store.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(sessionBucket).ForEach(func(_, data []byte) error {
			var client models.Client
			if err := json.Unmarshal(data, &client); err != nil {
				return err
			}
			client.Conn = nil
			client.IsActive = false
			if client.BasicInfoStatus == "collecting" {
				client.BasicInfoStatus = "interrupted"
			}
			sessions = append(sessions, &client)
			return nil
		})
	})
	return sessions, err
}
func (store *Store) MergeSession(savedID string, live *models.Client) error {
	return store.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(sessionBucket).Delete([]byte(savedID)); err != nil {
			return err
		}
		if err := put(tx.Bucket(sessionBucket), []byte(live.ID), live); err != nil {
			return err
		}
		bucket := tx.Bucket(taskBucket)
		updates := []Task{}
		if err := bucket.ForEach(func(_, data []byte) error {
			var task Task
			if err := json.Unmarshal(data, &task); err != nil {
				return err
			}
			if task.AgentID == savedID {
				task.AgentID = live.ID
				updates = append(updates, task)
			}
			return nil
		}); err != nil {
			return err
		}
		for _, task := range updates {
			if err := put(bucket, []byte(task.ID), task); err != nil {
				return err
			}
		}
		return audit(tx, AuditEvent{Event: "session.linked", AgentID: live.ID, Detail: "Saved session: " + savedID})
	})
}
func (store *Store) NewTask(agentID, command, origin string, requestID uint32) (string, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", err
	}
	task := Task{ID: "task-" + hex.EncodeToString(id[:]), AgentID: agentID, Command: command, Origin: origin, RequestID: requestID, Status: "queued", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	err := store.db.Update(func(tx *bolt.Tx) error {
		if err := put(tx.Bucket(taskBucket), []byte(task.ID), task); err != nil {
			return err
		}
		return audit(tx, AuditEvent{Event: "task.queued", AgentID: agentID, TaskID: task.ID, Source: origin})
	})
	return task.ID, err
}
func (store *Store) TransitionTask(id, status, output string) error {
	if status != "sent" && status != "completed" && status != "failed" && status != "interrupted" {
		return fmt.Errorf("invalid task status")
	}
	return store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("unknown task")
		}
		var task Task
		if err := json.Unmarshal(data, &task); err != nil {
			return err
		}
		// A response can arrive before the socket writer returns; never regress
		// a terminal state back to sent or overwrite it with a late failure.
		if task.Status != "queued" && task.Status != "sent" {
			return nil
		}
		if status == "sent" && task.Status != "queued" {
			return nil
		}
		task.Status = status
		task.UpdatedAt = time.Now().UTC()
		if status != "sent" {
			if len(output) > 1<<20 {
				output = output[:1<<20]
				task.OutputTruncated = true
			}
			task.Output = output
		}
		if err := put(bucket, []byte(id), task); err != nil {
			return err
		}
		return audit(tx, AuditEvent{Event: "task." + status, AgentID: task.AgentID, TaskID: id, Source: task.Origin})
	})
}
func (store *Store) recoverTasks() error {
	pending := []string{}
	err := store.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(taskBucket).ForEach(func(_, data []byte) error {
			var task Task
			if err := json.Unmarshal(data, &task); err != nil {
				return err
			}
			if task.Status == "queued" || task.Status == "sent" {
				pending = append(pending, task.ID)
			}
			return nil
		})
	})
	if err != nil {
		return err
	}
	for _, id := range pending {
		if err := store.TransitionTask(id, "interrupted", "Server restarted before a result was recorded"); err != nil {
			return err
		}
	}
	return nil
}
func (store *Store) Tasks(agentID string) ([]Task, error) {
	result := []Task{}
	err := store.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(taskBucket).ForEach(func(_, data []byte) error {
			var task Task
			if err := json.Unmarshal(data, &task); err != nil {
				return err
			}
			if agentID == "" || task.AgentID == agentID {
				result = append(result, task)
			}
			return nil
		})
	})
	return result, err
}
func (store *Store) Record(event AuditEvent) error {
	return store.db.Update(func(tx *bolt.Tx) error { return audit(tx, event) })
}
func (store *Store) Audit(limit int) ([]AuditEvent, error) {
	events := []AuditEvent{}
	err := store.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(auditBucket).Cursor()
		for _, data := cursor.Last(); data != nil && len(events) < limit; _, data = cursor.Prev() {
			var event AuditEvent
			if err := json.Unmarshal(data, &event); err != nil {
				return err
			}
			events = append(events, event)
		}
		return nil
	})
	return events, err
}

func (store *Store) Task(id string) (*Task, error) {
	var task Task
	err := store.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(taskBucket).Get([]byte(id))
		if data == nil {
			return fmt.Errorf("task not found")
		}
		return json.Unmarshal(data, &task)
	})
	return &task, err
}

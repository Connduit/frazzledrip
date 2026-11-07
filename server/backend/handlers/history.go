package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/state"
)

func taskHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodGet) {
		return
	}
	store := network.PersistentState()
	if store == nil {
		controlError(w, http.StatusServiceUnavailable, "persistent store unavailable")
		return
	}
	if id := r.URL.Query().Get("task"); id != "" {
		task, err := store.Task(id)
		if err != nil {
			controlError(w, http.StatusNotFound, "task not found")
			return
		}
		json.NewEncoder(w).Encode(task)
		return
	}
	tasks, err := store.Tasks(r.URL.Query().Get("agent"))
	if err != nil {
		controlError(w, http.StatusInternalServerError, "could not read tasks")
		return
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].CreatedAt.After(tasks[j].CreatedAt) })
	if len(tasks) > 500 {
		tasks = tasks[:500]
	}
	for i := range tasks {
		tasks[i].Output = ""
	}
	json.NewEncoder(w).Encode(tasks)
}
func auditHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if !requireControlMethod(w, r, http.MethodGet) {
		return
	}
	store := network.PersistentState()
	if store == nil {
		controlError(w, http.StatusServiceUnavailable, "persistent store unavailable")
		return
	}
	events, err := store.Audit(300)
	if err != nil {
		controlError(w, http.StatusInternalServerError, "could not read audit history")
		return
	}
	json.NewEncoder(w).Encode(events)
}

type auditResponse struct {
	http.ResponseWriter
	status int
}

func (w *auditResponse) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w *auditResponse) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

// Source records the actual HTTP peer, not an authenticated operator identity.
// Request bodies and query strings are deliberately excluded from the audit.
func AuditHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}
		response := &auditResponse{ResponseWriter: w}
		next.ServeHTTP(response, r)
		if response.status == 0 {
			response.status = http.StatusOK
		}
		if store := network.PersistentState(); store != nil {
			if err := store.Record(state.AuditEvent{Event: "http.request", Source: r.RemoteAddr, Detail: fmt.Sprintf("%s %s status=%d", r.Method, r.URL.Path, response.status)}); err != nil {
				log.Printf("Audit persistence error: %v", err)
			}
		}
	})
}

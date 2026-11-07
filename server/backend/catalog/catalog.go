// Package catalog contains non-executable development placeholders only.
package catalog

import (
	"fmt"
	"unicode/utf8"

	"github.com/Connduit/frazzledrip/server/backend/state"
)

type Tool struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

var tools = []Tool{
	{"host-survey", "Host survey", "Discovery", "attack", "Placeholder for presenting host discovery tasks. No host queries are issued."},
	{"network-survey", "Network survey", "Discovery", "attack", "Mock network discovery result workflow. No scanning or connections occur."},
	{"process-survey", "Process inventory", "Discovery", "attack", "Simulated process discovery task. No processes are inspected."},
	{"privilege-assessment", "Privilege assessment", "Privilege escalation", "attack", "Task/result placeholder only. No privilege changes or exploit code."},
	{"credential-assessment", "Credential access assessment", "Credential access", "attack", "Simulates the assessment workflow. No credentials, tokens, or secrets are accessed."},
	{"lateral-exercise", "Lateral movement exercise", "Lateral movement", "attack", "Mock task presentation. No remote logins, execution, or network access."},
	{"persistence-exercise", "Persistence exercise", "Persistence", "attack", "Simulation entry only. No services, tasks, startup entries, or files are created on a client."},
	{"defense-validation", "Detection validation", "Defense validation", "attack", "Placeholder for documenting expected detection events. Does not bypass or modify defenses."},
	{"powershell-template", "PowerShell inventory template", "Script templates", "script", "Non-executable script placeholder. Returns a TEST result without running PowerShell."},
	{"python-template", "Python inventory template", "Script templates", "script", "Non-executable script placeholder. Returns a TEST result without running Python."},
	{"bash-template", "Bash inventory template", "Script templates", "script", "Non-executable script placeholder. Returns a TEST result without running a shell."},
	{"report-template", "Engagement report template", "Reporting", "script", "Mock report task for testing history and audit presentation. No report file is written."},
}

func List() []Tool { return append([]Tool(nil), tools...) }
func Find(id string) (Tool, bool) {
	for _, tool := range tools {
		if tool.ID == id {
			return tool, true
		}
	}
	return Tool{}, false
}
func Validate(id, note string) error {
	if _, ok := Find(id); !ok {
		return fmt.Errorf("unknown placeholder tool")
	}
	if !utf8.ValidString(note) || utf8.RuneCountInString(note) > 500 {
		return fmt.Errorf("note must be valid text of at most 500 characters")
	}
	return nil
}

// This function has no transport, filesystem, interpreter, or command executor.
// It records a simulated task directly in the server database.
func Simulate(store *state.Store, agentID, id, note string) (*state.Task, error) {
	if err := Validate(id, note); err != nil {
		return nil, err
	}
	if store == nil {
		return nil, fmt.Errorf("persistent store unavailable")
	}
	tool, _ := Find(id)
	taskID, err := store.NewTask(agentID, "[SIMULATION] "+tool.Name, "simulation", 0)
	if err != nil {
		return nil, err
	}
	output := fmt.Sprintf("TEST OUTPUT — SIMULATION ONLY\nTool: %s\nCategory: %s\nSession: %s\nNo client was contacted. No commands or scripts were executed. No client state was changed.\n", tool.Name, tool.Category, agentID)
	if note != "" {
		output += "Operator note: " + note + "\n"
	}
	if err := store.TransitionTask(taskID, "completed", output); err != nil {
		return nil, err
	}
	return store.Task(taskID)
}

package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// ErrNoProject is returned by Load when the project root has no
// .agentism/config.json.
var ErrNoProject = errors.New("no agentism project found")

// Project is an agentism project loaded from disk.
type Project struct {
	Phase   string
	Version string
	Tags    []string
	Tasks   []Task
}

// Task is one agentism/tasks/Task_NNNN_<slug>/plan.md.
type Task struct {
	ID        string
	Title     string
	Slug      string
	Status    string
	Priority  int
	Tags      []string
	DependsOn []string
	Goal      string
	Contract  string
	Tickets   []Ticket
}

// Ticket is one T-NNNN-<slug>.md file inside a task folder.
type Ticket struct {
	ID           string
	Title        string
	TaskID       string
	Status       string
	Priority     int
	Tags         []string
	DependsOn    []string
	ContractHash string
	Contract     string
	Work         string
	Acceptance   string
}

var (
	taskDirRe    = regexp.MustCompile(`^Task_\d{4}_(.+)$`)
	ticketFileRe = regexp.MustCompile(`^T-\d{4}-.+\.md$`)

	validTicketStatus = map[string]bool{
		"TODO": true, "IN_PROGRESS": true, "IN_REVIEW": true, "DONE": true,
		"NOT_ACCEPTED": true, "BLOCKED": true, "STALE": true,
	}
)

// Load reads an agentism project rooted at root. It never writes a file
// and never runs the agentism binary.
func Load(root string) (*Project, error) {
	configPath := filepath.Join(root, ".agentism", "config.json")
	raw, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNoProject
	}
	if err != nil {
		return nil, err
	}

	var config struct {
		Phase   string   `json:"phase"`
		Version string   `json:"version"`
		Tags    []string `json:"tags"`
	}
	if err := json.Unmarshal(raw, &config); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", relPath(root, configPath), err)
	}

	tasks, err := loadTasks(root)
	if err != nil {
		return nil, err
	}

	return &Project{
		Phase:   config.Phase,
		Version: config.Version,
		Tags:    config.Tags,
		Tasks:   tasks,
	}, nil
}

func loadTasks(root string) ([]Task, error) {
	tasksDir := filepath.Join(root, "agentism", "tasks")
	entries, err := os.ReadDir(tasksDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for _, entry := range entries {
		m := taskDirRe.FindStringSubmatch(entry.Name())
		if !entry.IsDir() || m == nil {
			continue
		}

		taskDir := filepath.Join(tasksDir, entry.Name())
		task, err := loadTask(root, taskDir, m[1])
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Priority != tasks[j].Priority {
			return tasks[i].Priority < tasks[j].Priority
		}
		return tasks[i].ID < tasks[j].ID
	})
	return tasks, nil
}

func loadTask(root, taskDir, slug string) (Task, error) {
	planPath := filepath.Join(taskDir, "plan.md")
	doc, err := readDoc(root, planPath)
	if err != nil {
		return Task{}, err
	}

	goal, _ := doc.Section("Goal")
	contract, _ := doc.Section("Contracts")

	tickets, err := loadTickets(root, taskDir)
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:        metaString(doc.Meta, "id"),
		Title:     metaString(doc.Meta, "title"),
		Slug:      slug,
		Status:    metaString(doc.Meta, "status"),
		Priority:  metaInt(doc.Meta, "priority"),
		Tags:      metaStringSlice(doc.Meta, "tags"),
		DependsOn: metaStringSlice(doc.Meta, "depends_on"),
		Goal:      goal,
		Contract:  contract,
		Tickets:   tickets,
	}, nil
}

func loadTickets(root, taskDir string) ([]Ticket, error) {
	entries, err := os.ReadDir(taskDir)
	if err != nil {
		return nil, err
	}

	var tickets []Ticket
	for _, entry := range entries {
		if entry.IsDir() || !ticketFileRe.MatchString(entry.Name()) {
			continue
		}

		doc, err := readDoc(root, filepath.Join(taskDir, entry.Name()))
		if err != nil {
			return nil, err
		}

		status := metaString(doc.Meta, "status")
		if !validTicketStatus[status] {
			status = "TODO"
		}
		contract, _ := doc.Section("Contract")
		work, _ := doc.Section("Work")
		acceptance, _ := doc.Section("Acceptance")

		tickets = append(tickets, Ticket{
			ID:           metaString(doc.Meta, "id"),
			Title:        metaString(doc.Meta, "title"),
			TaskID:       metaString(doc.Meta, "task"),
			Status:       status,
			Priority:     metaInt(doc.Meta, "priority"),
			Tags:         metaStringSlice(doc.Meta, "tags"),
			DependsOn:    metaStringSlice(doc.Meta, "depends_on"),
			ContractHash: metaString(doc.Meta, "contract_hash"),
			Contract:     contract,
			Work:         work,
			Acceptance:   acceptance,
		})
	}

	sort.Slice(tickets, func(i, j int) bool {
		if tickets[i].Priority != tickets[j].Priority {
			return tickets[i].Priority < tickets[j].Priority
		}
		return tickets[i].ID < tickets[j].ID
	})
	return tickets, nil
}

// readDoc reads and parses a markdown file, naming path relative to root
// in any error it returns.
func readDoc(root, path string) (Doc, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Doc{}, err
	}
	doc, err := ParseDoc(string(raw))
	if err != nil {
		return Doc{}, fmt.Errorf("parsing %s: %w", relPath(root, path), err)
	}
	return doc, nil
}

func relPath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return rel
}

func metaString(meta map[string]any, key string) string {
	s, _ := meta[key].(string)
	return s
}

func metaInt(meta map[string]any, key string) int {
	switch v := meta[key].(type) {
	case int:
		return v
	default:
		return 0
	}
}

func metaStringSlice(meta map[string]any, key string) []string {
	list, ok := meta[key].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, item := range list {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

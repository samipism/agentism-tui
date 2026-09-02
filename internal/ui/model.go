// Package ui is the Bubble Tea dashboard for an agentism project.
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samipism/agentism-tui/internal/store"
)

// nonDoneStatuses lists the ticket statuses the status bar reports besides
// DONE, in report order.
var nonDoneStatuses = []string{"TODO", "IN_PROGRESS", "IN_REVIEW", "NOT_ACCEPTED", "BLOCKED", "STALE"}

// Model is the top half of the dashboard: the status bar and the
// task/ticket tree. T-0007 adds the detail pane and the log view.
type Model struct {
	project  *store.Project
	cursor   int
	expanded map[string]bool // task ID -> whether its tickets show
}

// New builds the dashboard model for project.
func New(project *store.Project) tea.Model {
	return Model{project: project, expanded: map[string]bool{}}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if last := len(m.rows()) - 1; m.cursor < last {
				m.cursor++
			}
		case "enter", " ":
			m.toggleCursorRow()
		}
	}
	return m, nil
}

// toggleCursorRow expands or collapses the task under the cursor. It does
// nothing when the cursor sits on a ticket row.
func (m *Model) toggleCursorRow() {
	rows := m.rows()
	if m.cursor >= len(rows) || !rows[m.cursor].isTask {
		return
	}
	id := rows[m.cursor].taskID
	m.expanded[id] = !m.expanded[id]
	if last := len(m.rows()) - 1; m.cursor > last {
		m.cursor = last
	}
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(m.statusBar())
	b.WriteString("\n")
	for i, r := range m.rows() {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		indent := ""
		if !r.isTask {
			indent = "    "
		}
		fmt.Fprintf(&b, "%s%s%s %s %s\n", cursor, indent, r.id, r.title, r.status)
	}
	return b.String()
}

// row is one visible line of the task/ticket tree.
type row struct {
	isTask bool
	id     string
	title  string
	status string
	taskID string // owning task's ID, for toggling expand state
}

// rows flattens the project into the tree's visible lines: every task, plus
// a task's tickets right after it when that task is expanded.
func (m Model) rows() []row {
	var rows []row
	for _, task := range m.project.Tasks {
		rows = append(rows, row{isTask: true, id: task.ID, taskID: task.ID, title: task.Title, status: task.Status})
		if !m.expanded[task.ID] {
			continue
		}
		for _, ticket := range task.Tickets {
			rows = append(rows, row{id: ticket.ID, taskID: task.ID, title: ticket.Title, status: ticket.Status})
		}
	}
	return rows
}

// statusBar renders the phase, the version, the DONE/total ticket count,
// and any non-zero non-DONE status count.
func (m Model) statusBar() string {
	counts := m.project.Counts()
	total := 0
	for _, n := range counts {
		total += n
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s v%s  DONE %d/%d", m.project.Phase, m.project.Version, counts["DONE"], total)
	for _, status := range nonDoneStatuses {
		if n := counts[status]; n > 0 {
			fmt.Fprintf(&b, "  %s %d", status, n)
		}
	}
	return b.String()
}

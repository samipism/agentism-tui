// Package ui is the Bubble Tea dashboard for an agentism project.
package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/samipism/agentism-tui/internal/store"
)

// nonDoneStatuses lists the ticket statuses the status bar reports besides
// DONE, in report order.
var nonDoneStatuses = []string{"TODO", "IN_PROGRESS", "IN_REVIEW", "NOT_ACCEPTED", "BLOCKED", "STALE"}

// Model is the dashboard: the status bar, the task/ticket tree, and
// (T-0007) the detail pane and the log view.
type Model struct {
	project  *store.Project
	root     string // project root, so 'r' can reload
	load     func(root string) (*store.Project, error)
	cursor   int
	expanded map[string]bool // task ID -> whether its tickets show
	selected string          // ticket ID shown in the detail pane, "" for none
	showLog  bool
	err      string // last refresh error, shown in the status bar
}

// New builds the dashboard model for project, rooted at root.
func New(project *store.Project, root string) tea.Model {
	return Model{project: project, root: root, expanded: map[string]bool{}, load: store.Load}
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
			m.selectOrToggleCursorRow()
		case "r":
			m.refresh()
		case "l":
			m.showLog = !m.showLog
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

// selectOrToggleCursorRow expands or collapses the task under the cursor,
// or, when the cursor sits on a ticket row, selects that ticket for the
// detail pane.
func (m *Model) selectOrToggleCursorRow() {
	rows := m.rows()
	if m.cursor >= len(rows) {
		return
	}
	r := rows[m.cursor]
	if !r.isTask {
		m.selected = r.id
		return
	}
	m.expanded[r.taskID] = !m.expanded[r.taskID]
	if last := len(m.rows()) - 1; m.cursor > last {
		m.cursor = last
	}
}

// refresh reloads the project from disk. A load error is kept in m.err and
// the previous project stays in place; it never replaces m.project with
// nil and never panics.
func (m *Model) refresh() {
	p, err := m.load(m.root)
	if err != nil {
		m.err = err.Error()
		return
	}
	m.project = p
	m.err = ""
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(m.statusBar())
	if m.err != "" {
		fmt.Fprintf(&b, "  ! %s", m.err)
	}
	b.WriteString("\n")

	right := m.detailPane()
	if m.showLog {
		right = m.logPane()
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.treeView(), "  ", right))
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

// treeView renders the task/ticket tree, cursor included.
func (m Model) treeView() string {
	var b strings.Builder
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

// findTicket looks up a ticket by ID across every task. It returns nil,
// nil when id is empty or no ticket matches.
func (m Model) findTicket(id string) (*store.Ticket, *store.Task) {
	if id == "" {
		return nil, nil
	}
	for i := range m.project.Tasks {
		task := &m.project.Tasks[i]
		for j := range task.Tickets {
			if task.Tickets[j].ID == id {
				return &task.Tickets[j], task
			}
		}
	}
	return nil, nil
}

// detailPane renders the selected ticket's title, status, task, and its
// Contract, Work, and Acceptance text as styled markdown. It is empty when
// no ticket is selected.
func (m Model) detailPane() string {
	t, task := m.findTicket(m.selected)
	if t == nil {
		return ""
	}

	var md strings.Builder
	fmt.Fprintf(&md, "# %s %s\n\n**Status:** %s **Task:** %s\n\n", t.ID, t.Title, t.Status, task.Title)
	for _, section := range []struct {
		heading string
		body    string
	}{
		{"Contract", t.Contract},
		{"Work", t.Work},
		{"Acceptance", t.Acceptance},
	} {
		if section.body == "" {
			continue
		}
		fmt.Fprintf(&md, "## %s\n\n%s\n\n", section.heading, section.body)
	}

	out, err := glamour.Render(md.String(), "auto")
	if err != nil {
		return md.String() // raw markdown beats a blank pane
	}
	return out
}

// logPane renders Project.Log, newest entry first.
func (m Model) logPane() string {
	var b strings.Builder
	b.WriteString("Activity Log\n\n")
	for i := len(m.project.Log) - 1; i >= 0; i-- {
		e := m.project.Log[i]
		fmt.Fprintf(&b, "%s  %s\n", e.At, e.Kind)
	}
	return b.String()
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

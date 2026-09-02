// Package ui is the Bubble Tea dashboard for an agentism project.
package ui

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/lipgloss"

	"github.com/samipism/agentism-tui/internal/store"
)

// nonDoneStatuses lists the ticket statuses the status bar reports besides
// DONE, in report order.
var nonDoneStatuses = []string{"TODO", "IN_PROGRESS", "IN_REVIEW", "NOT_ACCEPTED", "BLOCKED", "STALE"}

// statusColors maps a status value to the one fixed color it always
// renders in. A status not listed here (TODO, PLANNED, STALE, and so on)
// gets neutralColor.
var statusColors = map[string]lipgloss.Color{
	"DONE":         lipgloss.Color("2"), // green
	"IN_PROGRESS":  lipgloss.Color("3"), // yellow
	"IN_REVIEW":    lipgloss.Color("3"), // yellow
	"BLOCKED":      lipgloss.Color("1"), // red
	"NOT_ACCEPTED": lipgloss.Color("1"), // red
}

const neutralColor = lipgloss.Color("245")

var (
	boxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	dimStyle = lipgloss.NewStyle().Foreground(neutralColor)
)

// Layout constants. Regions are sized from the terminal window; these
// control the fixed parts of that split.
const (
	defaultWidth  = 100 // used until the first tea.WindowSizeMsg arrives
	defaultHeight = 30

	headerContentLines = 2 // status line + error line
	footerContentLines = 1 // keybinding hints
	sidebarWidthPct    = 30
	minRegionWidth     = 16
)

// Model is the dashboard's status bar, task/ticket tree, detail pane, and
// log view.
type Model struct {
	project *store.Project
	root    string                               // project root, for 'r' to reload
	load    func(string) (*store.Project, error) // swappable in tests

	cursor     int
	expanded   map[string]bool // task ID -> whether its tickets show
	selectedID string          // ticket ID shown in the detail pane, "" for none
	showLog    bool
	errMsg     string // set by a failed 'r' refresh, cleared on success
	mainScroll int    // top line of the main region's scrolled content

	width, height int
	darkBG        bool // guessed once at startup; see darkBackground
}

// New builds the dashboard model for project, loaded from root.
func New(project *store.Project, root string) tea.Model {
	return Model{
		project:  project,
		root:     root,
		load:     store.Load,
		expanded: map[string]bool{},
		width:    defaultWidth,
		height:   defaultHeight,
		darkBG:   darkBackground(),
	}
}

// darkBackground guesses whether the terminal has a dark background from
// the COLORFGBG environment variable (set by most terminal emulators),
// defaulting to dark when it's absent. This deliberately never queries the
// terminal directly (as termenv.HasDarkBackground or glamour.WithAutoStyle
// do): that query writes an OSC escape sequence and blocks reading the
// reply from the tty, and once Bubble Tea owns stdin in raw mode for its
// own key-reading loop, that reply gets swallowed and the read hangs the
// whole program well past its own timeout.
func darkBackground() bool {
	fgbg := os.Getenv("COLORFGBG")
	parts := strings.Split(fgbg, ";")
	bg, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return true // unknown: default to dark, the common terminal convention
	}
	return bg != 7 && bg != 15 // xterm convention: 7 and 15 are light backgrounds
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.activateCursorRow()
	case "r":
		m.refresh()
	case "l":
		m.showLog = !m.showLog
		m.mainScroll = 0
	case "pgdown":
		m.mainScroll += m.innerHeight()
	case "pgup":
		m.mainScroll -= m.innerHeight()
		if m.mainScroll < 0 {
			m.mainScroll = 0
		}
	case "q", "ctrl+c":
		// Invariant: no keybinding writes a file or spawns a process. Quit
		// only stops the program.
		return m, tea.Quit
	}
	return m, nil
}

// activateCursorRow expands or collapses the task under the cursor, or, for
// a ticket row, selects it for the detail pane.
func (m *Model) activateCursorRow() {
	rows := m.rows()
	if m.cursor >= len(rows) {
		return
	}
	r := rows[m.cursor]
	if !r.isTask {
		m.selectedID = r.id
		m.mainScroll = 0
		return
	}
	m.expanded[r.taskID] = !m.expanded[r.taskID]
	if last := len(m.rows()) - 1; m.cursor > last {
		m.cursor = last
	}
}

// refresh reloads the project from root. A load error leaves the previous
// project on screen and reports the error instead of losing it.
func (m *Model) refresh() {
	project, err := m.load(m.root)
	if err != nil {
		m.errMsg = err.Error()
		return
	}
	m.project = project
	m.errMsg = ""
}

// selectedTicket resolves selectedID against the current project, so a
// refresh picks up any status change on the selected ticket.
func (m Model) selectedTicket() *store.Ticket {
	if m.selectedID == "" {
		return nil
	}
	for _, task := range m.project.Tasks {
		for i := range task.Tickets {
			if task.Tickets[i].ID == m.selectedID {
				return &task.Tickets[i]
			}
		}
	}
	return nil
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

func (m Model) View() string {
	sidebarOuter, mainOuter := m.regionWidths()
	innerH := m.innerHeight()

	sidebar := boxStyle.Width(sidebarOuter - 2).Height(innerH).Render(clipLines(m.treeView(), innerH))
	main := boxStyle.Width(mainOuter - 2).Height(innerH).Render(scrollLines(m.mainView(mainOuter-2), m.mainScroll, innerH))
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)

	return lipgloss.JoinVertical(lipgloss.Left, m.headerView(), body, m.footerView())
}

// innerHeight is the content height available inside the sidebar/main
// boxes, border excluded. It also sizes a pgup/pgdown scroll step.
func (m Model) innerHeight() int {
	h := m.bodyHeight() - 2
	if h < 1 {
		h = 1
	}
	return h
}

// regionWidths splits the window into the sidebar (the tree) and the main
// region (the detail pane or the log view), each including its own border.
func (m Model) regionWidths() (sidebarOuter, mainOuter int) {
	sidebarOuter = m.width * sidebarWidthPct / 100
	if sidebarOuter < minRegionWidth {
		sidebarOuter = minRegionWidth
	}
	if max := m.width - minRegionWidth; sidebarOuter > max {
		sidebarOuter = max
	}
	if sidebarOuter < 4 {
		sidebarOuter = 4
	}
	mainOuter = m.width - sidebarOuter
	if mainOuter < 4 {
		mainOuter = 4
	}
	return sidebarOuter, mainOuter
}

func (m Model) bodyHeight() int {
	h := m.height - (headerContentLines + 2) - (footerContentLines + 2)
	if h < 3 {
		h = 3
	}
	return h
}

// headerView renders the phase, version, the overall completion bar, the
// non-zero status counts, and any refresh error.
func (m Model) headerView() string {
	counts := m.project.Counts()
	total := 0
	for _, n := range counts {
		total += n
	}

	var line1 strings.Builder
	fmt.Fprintf(&line1, "%s %s  %s  DONE %d/%d", m.project.Phase, m.project.Version, progressBar(counts["DONE"], total), counts["DONE"], total)
	for _, status := range nonDoneStatuses {
		if n := counts[status]; n > 0 {
			fmt.Fprintf(&line1, "  %s %d", styledStatus(status), n)
		}
	}

	line2 := ""
	if m.errMsg != "" {
		line2 = lipgloss.NewStyle().Foreground(statusColors["BLOCKED"]).Render("error: " + m.errMsg)
	}

	return boxStyle.Width(m.width - 2).Height(headerContentLines).Render(line1.String() + "\n" + line2)
}

func (m Model) footerView() string {
	hint := "up/down move   enter/space select   pgup/pgdown scroll   r refresh   l log   q quit"
	return boxStyle.Width(m.width - 2).Height(footerContentLines).Render(hint)
}

// treeView renders the task/ticket tree, cursor, colored statuses, and each
// task's own completion bar included.
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
		fmt.Fprintf(&b, "%s%s%s %s %s", cursor, indent, r.id, r.title, styledStatus(r.status))
		if r.isTask {
			if bar := m.taskProgressBar(r.taskID); bar != "" {
				fmt.Fprintf(&b, "  %s", bar)
			}
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// taskProgressBar builds the completion bar for one task, from that task's
// own ticket counts.
func (m Model) taskProgressBar(taskID string) string {
	for _, task := range m.project.Tasks {
		if task.ID != taskID {
			continue
		}
		done := 0
		for _, t := range task.Tickets {
			if t.Status == "DONE" {
				done++
			}
		}
		return progressBar(done, len(task.Tickets))
	}
	return ""
}

// mainView renders whichever of the log view or the detail pane is active.
func (m Model) mainView(width int) string {
	if m.showLog {
		return m.logView()
	}
	ticket := m.selectedTicket()
	if ticket == nil {
		return dimStyle.Render("Select a ticket to view its details.")
	}
	return m.detailView(ticket, width)
}

// detailView renders a ticket's title, status, task, and Contract/Work/
// Acceptance text as one markdown document through glamour.
func (m Model) detailView(ticket *store.Ticket, width int) string {
	md := fmt.Sprintf(
		"# %s\n\n**Status:** %s  **Task:** %s\n\n## Contract\n\n%s\n\n## Work\n\n%s\n\n## Acceptance\n\n%s\n",
		ticket.Title, ticket.Status, ticket.TaskID, ticket.Contract, ticket.Work, ticket.Acceptance,
	)
	return m.renderMarkdown(md, width)
}

// renderMarkdown runs md through glamour, matched to the terminal
// background Model.darkBG guessed at startup (see darkBackground). A
// renderer error falls back to the raw markdown rather than crashing the
// program.
func (m Model) renderMarkdown(md string, width int) string {
	if width < 20 {
		width = 20
	}
	style := styles.LightStyleConfig
	if m.darkBG {
		style = styles.DarkStyleConfig
	}
	// The default h2 style keeps a literal "## " prefix as a decoration;
	// clear it so a heading renders bold/colored only, not as raw markdown.
	style.H2.Prefix = ""
	r, err := glamour.NewTermRenderer(glamour.WithStyles(style), glamour.WithWordWrap(width))
	if err != nil {
		return md
	}
	out, err := r.Render(md)
	if err != nil {
		return md
	}
	return out
}

// logView renders Project.Log, newest entry first. View scrolls it with
// pgup/pgdown, same as the detail pane.
func (m Model) logView() string {
	if len(m.project.Log) == 0 {
		return dimStyle.Render("No log entries.")
	}
	var b strings.Builder
	for i := len(m.project.Log) - 1; i >= 0; i-- {
		entry := m.project.Log[i]
		fmt.Fprintf(&b, "%s  %s", entry.At, entry.Kind)

		keys := make([]string, 0, len(entry.Fields))
		for k := range entry.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&b, "  %s=%v", k, entry.Fields[k])
		}
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// styledStatus renders a status value in its one fixed color.
func styledStatus(status string) string {
	color, ok := statusColors[status]
	if !ok {
		color = neutralColor
	}
	return lipgloss.NewStyle().Foreground(color).Render(status)
}

// progressBar builds the block-and-percentage string used by both the
// header's overall bar and each task row's bar. "" for a task with no
// tickets.
func progressBar(done, total int) string {
	if total == 0 {
		return ""
	}
	const barWidth = 10
	filled := done * barWidth / total
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	return fmt.Sprintf("%s %d%%", bar, done*100/total)
}

// clipLines keeps only the first maxLines lines of s, so rendered content
// never overflows a bordered box's fixed height.
func clipLines(s string, maxLines int) string {
	return scrollLines(s, 0, maxLines)
}

// scrollLines returns the maxLines-line window of s starting at offset, so
// rendered content never overflows a bordered box's fixed height while
// still reaching every line via pgup/pgdown. offset clamps to [0, the
// largest offset that still fills the window], so scrolling can't run past
// either end of the content.
func scrollLines(s string, offset, maxLines int) string {
	if maxLines <= 0 {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) <= maxLines {
		return s
	}
	if max := len(lines) - maxLines; offset > max {
		offset = max
	}
	if offset < 0 {
		offset = 0
	}
	return strings.Join(lines[offset:offset+maxLines], "\n")
}

package ui

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samipism/agentism-tui/internal/store"
)

// ansiCode matches a terminal escape sequence, so a styled render can be
// checked for its plain text content.
var ansiCode = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string {
	return ansiCode.ReplaceAllString(s, "")
}

func fixtureProject() *store.Project {
	return &store.Project{
		Phase:   "TASKS",
		Version: "v0",
		Tasks: []store.Task{
			{ID: "Task_0001", Title: "First", Status: "COMPLETE", Tickets: []store.Ticket{
				{ID: "T-0001", Title: "One", TaskID: "Task_0001", Status: "DONE",
					Contract: "## Rules\n\nDo the thing.", Work: "Build it.", Acceptance: "It works."},
			}},
			{ID: "Task_0002", Title: "Second", Status: "PLANNED", Tickets: []store.Ticket{
				{ID: "T-0002", Title: "Two", TaskID: "Task_0002", Status: "TODO"},
			}},
		},
		Log: []store.LogEntry{
			{At: "2026-09-01T10:00:00Z", Kind: "created", Fields: map[string]any{"ticket": "T-0001"}},
			{At: "2026-09-02T10:00:00Z", Kind: "status_changed", Fields: map[string]any{"ticket": "T-0001", "to": "DONE"}},
		},
	}
}

func lineWith(view, substr string) string {
	for _, line := range strings.Split(view, "\n") {
		if strings.Contains(line, substr) {
			return line
		}
	}
	return ""
}

// column returns id's index within line, ignoring the box-drawing border
// characters lipgloss puts at the start of a row.
func column(line, id string) int {
	return strings.Index(line, id)
}

func TestViewExpandShowsTicketsIndentedCollapseHidesThem(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)
	m.expanded["Task_0001"] = true

	view := m.View()

	taskLine := lineWith(view, "Task_0001")
	ticketLine := lineWith(view, "T-0001")
	if taskLine == "" || ticketLine == "" {
		t.Fatalf("expected both task and ticket lines, got:\n%s", view)
	}
	if column(ticketLine, "T-0001") <= column(taskLine, "Task_0001") {
		t.Errorf("expanded ticket line not indented deeper than its task:\ntask:   %q\nticket: %q", taskLine, ticketLine)
	}
	if strings.Contains(view, "T-0002") {
		t.Errorf("collapsed task's ticket should not appear:\n%s", view)
	}
}

func TestCursorMovesDownAndStopsAtLastRow(t *testing.T) {
	m := New(fixtureProject(), "root").(Model) // 2 collapsed task rows: indices 0,1

	for i := 0; i < 5; i++ {
		mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = mm.(Model)
	}
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1 (clamped at the last visible row)", m.cursor)
	}
}

func TestCursorMovesUpAndStopsAtFirstRow(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)

	for i := 0; i < 5; i++ {
		mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
		m = mm.(Model)
	}
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 (no wraparound above the first row)", m.cursor)
	}
}

func TestEnterOnTaskTogglesExpandCollapse(t *testing.T) {
	m := New(fixtureProject(), "root").(Model) // cursor starts on the first task row

	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm.(Model)
	if !strings.Contains(m.View(), "T-0001") {
		t.Fatal("enter on a task row should expand it")
	}

	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm.(Model)
	if strings.Contains(m.View(), "T-0001") {
		t.Fatal("enter again should collapse it")
	}
}

func TestStatusBarShowsPhaseVersionAndCounts(t *testing.T) {
	p := fixtureProject()
	m := New(p, "root").(Model)

	view := m.View()

	if !strings.Contains(view, p.Phase) {
		t.Errorf("view missing phase %q", p.Phase)
	}
	if !strings.Contains(view, p.Version) {
		t.Errorf("view missing version %q", p.Version)
	}

	counts := p.Counts()
	total := 0
	for _, n := range counts {
		total += n
	}
	want := fmt.Sprintf("%d/%d", counts["DONE"], total)
	if !strings.Contains(view, want) {
		t.Errorf("view missing done/total count %q", want)
	}
	if !strings.Contains(view, fmt.Sprintf("%d", counts["TODO"])) {
		t.Errorf("view missing non-zero TODO count: %q", view)
	}
}

func TestPressingQReturnsQuitCommand(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected a command from pressing q")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.Quit, got %T", cmd())
	}
}

func TestPressingCtrlCReturnsQuitCommand(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected a command from pressing ctrl+c")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.Quit, got %T", cmd())
	}
}

// --- T-0007: detail pane, log view, and keybindings ---

func TestSelectingTicketShowsStyledContractInDetailPane(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)
	// expand Task_0001, move cursor onto its ticket T-0001, select it
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // expand task
	m = mm.(Model)
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown}) // move onto T-0001
	m = mm.(Model)
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // select ticket
	m = mm.(Model)

	view := stripANSI(m.View())

	if !strings.Contains(view, "One") {
		t.Errorf("detail pane missing ticket title, view:\n%s", view)
	}
	if !strings.Contains(view, "Do the thing.") {
		t.Errorf("detail pane missing contract text, view:\n%s", view)
	}
	if strings.Contains(view, "## Rules") {
		t.Errorf("contract heading should render styled, not as literal markdown:\n%s", view)
	}
}

func TestRefreshErrorKeepsPreviousProjectAndShowsErrorInStatusBar(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)
	m.load = func(string) (*store.Project, error) { return nil, errors.New("boom disk error") }

	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m = mm.(Model)

	view := m.View()
	if !strings.Contains(view, "boom disk error") {
		t.Errorf("expected error text in view, got:\n%s", view)
	}
	if !strings.Contains(view, "Task_0001") {
		t.Errorf("expected previous project's tasks to remain, got:\n%s", view)
	}
}

func TestRefreshWithChangedProjectUpdatesCounts(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)
	changed := fixtureProject()
	changed.Tasks[1].Tickets[0].Status = "DONE"
	m.load = func(string) (*store.Project, error) { return changed, nil }

	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m = mm.(Model)

	view := m.View()
	if !strings.Contains(view, "2/2") {
		t.Errorf("expected updated DONE/total count 2/2, got:\n%s", view)
	}
}

func TestPressingLTwiceReturnsToPreToggleState(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)
	before := m.View()

	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = mm.(Model)
	during := m.View()
	if during == before {
		t.Fatal("pressing 'l' should change the view (log appears)")
	}

	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = mm.(Model)
	after := m.View()
	if after != before {
		t.Errorf("pressing 'l' twice should return to the pre-toggle view\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func TestLogViewShowsLogEntries(t *testing.T) {
	m := New(fixtureProject(), "root").(Model)
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = mm.(Model)

	view := m.View()
	if !strings.Contains(view, "status_changed") {
		t.Errorf("log view missing log entry kind, got:\n%s", view)
	}
}

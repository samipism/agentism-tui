package ui

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samipism/agentism-tui/internal/store"
)

func fixtureProject() *store.Project {
	return &store.Project{
		Phase:   "TASKS",
		Version: "v0",
		Tasks: []store.Task{
			{ID: "Task_0001", Title: "First", Status: "COMPLETE", Tickets: []store.Ticket{
				{ID: "T-0001", Title: "One", Status: "DONE", Contract: "- Alpha\n- Beta"},
			}},
			{ID: "Task_0002", Title: "Second", Status: "PLANNED", Tickets: []store.Ticket{
				{ID: "T-0002", Title: "Two", Status: "TODO"},
			}},
		},
		Log: []store.LogEntry{
			{At: "2026-09-01T00:00:00Z", Kind: "ticket_status_changed"},
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

func TestViewExpandShowsTicketsIndentedCollapseHidesThem(t *testing.T) {
	m := New(fixtureProject(), "").(Model)
	m.expanded["Task_0001"] = true

	view := m.View()

	taskLine := lineWith(view, "Task_0001")
	ticketLine := lineWith(view, "T-0001")
	if taskLine == "" || ticketLine == "" {
		t.Fatalf("expected both task and ticket lines, got:\n%s", view)
	}
	if leadingSpace(ticketLine) <= leadingSpace(taskLine) {
		t.Errorf("expanded ticket line not indented deeper than its task:\ntask:   %q\nticket: %q", taskLine, ticketLine)
	}
	if strings.Contains(view, "T-0002") {
		t.Errorf("collapsed task's ticket should not appear:\n%s", view)
	}
}

func leadingSpace(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}

func TestCursorMovesDownAndStopsAtLastRow(t *testing.T) {
	m := New(fixtureProject(), "").(Model) // 2 collapsed task rows: indices 0,1

	for i := 0; i < 5; i++ {
		mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = mm.(Model)
	}
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1 (clamped at the last visible row)", m.cursor)
	}
}

func TestCursorMovesUpAndStopsAtFirstRow(t *testing.T) {
	m := New(fixtureProject(), "").(Model)

	for i := 0; i < 5; i++ {
		mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
		m = mm.(Model)
	}
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0 (no wraparound above the first row)", m.cursor)
	}
}

func TestEnterOnTaskTogglesExpandCollapse(t *testing.T) {
	m := New(fixtureProject(), "").(Model) // cursor starts on the first task row

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
	m := New(p, "").(Model)

	statusBar := strings.SplitN(m.View(), "\n", 2)[0]

	if !strings.Contains(statusBar, p.Phase) {
		t.Errorf("status bar missing phase %q: %q", p.Phase, statusBar)
	}
	if !strings.Contains(statusBar, p.Version) {
		t.Errorf("status bar missing version %q: %q", p.Version, statusBar)
	}

	counts := p.Counts()
	total := 0
	for _, n := range counts {
		total += n
	}
	want := fmt.Sprintf("%d/%d", counts["DONE"], total)
	if !strings.Contains(statusBar, want) {
		t.Errorf("status bar missing done/total count %q: %q", want, statusBar)
	}
	if !strings.Contains(statusBar, fmt.Sprintf("TODO %d", counts["TODO"])) {
		t.Errorf("status bar missing non-zero TODO count: %q", statusBar)
	}
}

// selectFirstTicket expands Task_0001, moves the cursor onto its ticket
// row, and selects it.
func selectFirstTicket(m Model) Model {
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // expand Task_0001
	m = mm.(Model)
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown}) // move onto T-0001
	m = mm.(Model)
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // select it
	return mm.(Model)
}

func TestSelectingTicketShowsStyledContractInDetailPane(t *testing.T) {
	m := selectFirstTicket(New(fixtureProject(), "").(Model))

	view := m.View()
	if !strings.Contains(view, "One") {
		t.Fatalf("detail pane missing selected ticket's title:\n%s", view)
	}
	if !strings.Contains(view, "•") {
		t.Errorf("contract's list should render through glamour as bullets:\n%s", view)
	}
	if strings.Contains(view, "- Alpha") {
		t.Errorf("contract list item shown as raw markdown source, not rendered:\n%s", view)
	}
}

func keyRune(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func TestRefreshErrorKeepsPreviousDataAndShowsErrorInStatusBar(t *testing.T) {
	m := New(fixtureProject(), "").(Model)
	wantErr := errors.New("boom: disk read failed")
	m.load = func(string) (*store.Project, error) { return nil, wantErr }

	mm, _ := m.Update(keyRune('r'))
	m = mm.(Model)

	view := m.View()
	if !strings.Contains(view, "Task_0001") || !strings.Contains(view, "Task_0002") {
		t.Fatalf("previous tasks should still show after a failed refresh:\n%s", view)
	}
	if !strings.Contains(view, wantErr.Error()) {
		t.Errorf("status bar should show the refresh error, got:\n%s", view)
	}
}

func TestRefreshSuccessReplacesProjectAndUpdatesCounts(t *testing.T) {
	m := New(fixtureProject(), "").(Model)
	newProject := fixtureProject()
	newProject.Tasks[1].Tickets[0].Status = "DONE" // one more DONE ticket
	m.load = func(string) (*store.Project, error) { return newProject, nil }

	mm, _ := m.Update(keyRune('r'))
	m = mm.(Model)

	counts := newProject.Counts()
	total := 0
	for _, n := range counts {
		total += n
	}
	want := fmt.Sprintf("DONE %d/%d", counts["DONE"], total)
	statusBar := strings.SplitN(m.View(), "\n", 2)[0]
	if !strings.Contains(statusBar, want) {
		t.Errorf("status bar should reflect the refreshed counts %q, got %q", want, statusBar)
	}
}

func TestToggleLogTwiceReturnsToPreToggleView(t *testing.T) {
	m := New(fixtureProject(), "").(Model)
	before := m.View()

	mm, _ := m.Update(keyRune('l'))
	m = mm.(Model)
	afterOne := m.View()
	if afterOne == before {
		t.Fatal("'l' should change the view to show the log")
	}
	if !strings.Contains(afterOne, "ticket_status_changed") {
		t.Errorf("log view should show the log's entries:\n%s", afterOne)
	}

	mm, _ = m.Update(keyRune('l'))
	m = mm.(Model)
	afterTwo := m.View()
	if afterTwo != before {
		t.Errorf("pressing 'l' twice should return to the pre-toggle view\nbefore: %q\nafter:  %q", before, afterTwo)
	}
}

func TestPressingQReturnsQuitCommand(t *testing.T) {
	m := New(fixtureProject(), "").(Model)

	_, cmd := m.Update(keyRune('q'))
	if cmd == nil {
		t.Fatal("expected a command from pressing q")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.Quit, got %T", cmd())
	}
}

func TestPressingCtrlCReturnsQuitCommand(t *testing.T) {
	m := New(fixtureProject(), "").(Model)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected a command from pressing ctrl+c")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("expected tea.Quit, got %T", cmd())
	}
}

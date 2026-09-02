package store

import (
	"errors"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("two tasks, one with two tickets and one with none", func(t *testing.T) {
		proj, err := Load("testdata/valid")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(proj.Tasks) != 2 {
			t.Fatalf("len(Tasks) = %d, want 2", len(proj.Tasks))
		}
		if proj.Phase != "EXECUTION" || proj.Version != "v0" {
			t.Errorf("Phase/Version = %q/%q, want EXECUTION/v0", proj.Phase, proj.Version)
		}

		byID := map[string]Task{}
		for _, task := range proj.Tasks {
			byID[task.ID] = task
		}
		if len(byID["Task_0001"].Tickets) != 0 {
			t.Errorf("Task_0001 tickets = %d, want 0", len(byID["Task_0001"].Tickets))
		}
		if len(byID["Task_0002"].Tickets) != 2 {
			t.Errorf("Task_0002 tickets = %d, want 2", len(byID["Task_0002"].Tickets))
		}
	})

	t.Run("task and ticket order matches priority then ID", func(t *testing.T) {
		proj, err := Load("testdata/valid")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if proj.Tasks[0].ID != "Task_0002" || proj.Tasks[1].ID != "Task_0001" {
			t.Errorf("task order = %v, want [Task_0002 Task_0001]", []string{proj.Tasks[0].ID, proj.Tasks[1].ID})
		}

		var second Task
		for _, task := range proj.Tasks {
			if task.ID == "Task_0002" {
				second = task
			}
		}
		if second.Tickets[0].ID != "T-0005" || second.Tickets[1].ID != "T-0003" {
			t.Errorf("ticket order = %v, want [T-0005 T-0003]", []string{second.Tickets[0].ID, second.Tickets[1].ID})
		}
	})

	t.Run("missing config.json returns ErrNoProject", func(t *testing.T) {
		_, err := Load("testdata/noconfig")
		if !errors.Is(err, ErrNoProject) {
			t.Errorf("err = %v, want ErrNoProject", err)
		}
	})

	t.Run("a task folder that does not match Task_NNNN_ is skipped", func(t *testing.T) {
		proj, err := Load("testdata/valid")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, task := range proj.Tasks {
			if task.ID == "NotATask" {
				t.Errorf("NotATask folder should be skipped, got task %+v", task)
			}
		}
	})

	t.Run("a ticket file that does not match T-NNNN- is skipped", func(t *testing.T) {
		proj, err := Load("testdata/valid")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		for _, task := range proj.Tasks {
			for _, ticket := range task.Tickets {
				if ticket.ID != "T-0003" && ticket.ID != "T-0005" {
					t.Errorf("unexpected ticket loaded: %+v", ticket)
				}
			}
		}
	})

	t.Run("an unknown ticket status falls back to TODO", func(t *testing.T) {
		proj, err := Load("testdata/valid")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var weird Ticket
		for _, task := range proj.Tasks {
			for _, ticket := range task.Tickets {
				if ticket.ID == "T-0005" {
					weird = ticket
				}
			}
		}
		if weird.Status != "TODO" {
			t.Errorf("Status = %q, want TODO", weird.Status)
		}
	})

	t.Run("broken plan.md frontmatter returns an error naming the file", func(t *testing.T) {
		_, err := Load("testdata/badplan")
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		want := "Task_0001_broken/plan.md"
		if !strings.Contains(err.Error(), want) {
			t.Errorf("err = %v, want it to contain %q", err, want)
		}
	})
}

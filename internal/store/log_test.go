package store

import "testing"

func TestLoadLog(t *testing.T) {
	t.Run("three valid lines load in file order with extra fields", func(t *testing.T) {
		proj, err := Load("testdata/log")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(proj.Log) != 3 {
			t.Fatalf("len(Log) = %d, want 3", len(proj.Log))
		}
		want := []string{"created", "status_changed", "status_changed"}
		for i, kind := range want {
			if proj.Log[i].Kind != kind {
				t.Errorf("Log[%d].Kind = %q, want %q", i, proj.Log[i].Kind, kind)
			}
		}
		if proj.Log[0].At != "2026-09-01T10:00:00Z" {
			t.Errorf("Log[0].At = %q, want 2026-09-01T10:00:00Z", proj.Log[0].At)
		}
		if proj.Log[1].Fields["ticket"] != "T-0001" || proj.Log[1].Fields["to"] != "IN_PROGRESS" {
			t.Errorf("Log[1].Fields = %+v, want ticket=T-0001 to=IN_PROGRESS", proj.Log[1].Fields)
		}
	})

	t.Run("a malformed line among valid ones contributes nothing", func(t *testing.T) {
		proj, err := Load("testdata/badlog")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(proj.Log) != 2 {
			t.Fatalf("len(Log) = %d, want 2", len(proj.Log))
		}
	})

	t.Run("absent log.jsonl loads an empty Log", func(t *testing.T) {
		proj, err := Load("testdata/valid")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(proj.Log) != 0 {
			t.Errorf("len(Log) = %d, want 0", len(proj.Log))
		}
	})
}

func TestCounts(t *testing.T) {
	proj := &Project{
		Tasks: []Task{
			{ID: "Task_0001", Tickets: []Ticket{
				{ID: "T-0001", Status: "TODO"},
				{ID: "T-0002", Status: "TODO"},
				{ID: "T-0003", Status: "IN_PROGRESS"},
			}},
			{ID: "Task_0002", Tickets: []Ticket{
				{ID: "T-0004", Status: "DONE"},
			}},
		},
	}

	counts := proj.Counts()
	want := map[string]int{
		"TODO": 2, "IN_PROGRESS": 1, "IN_REVIEW": 0, "DONE": 1,
		"NOT_ACCEPTED": 0, "BLOCKED": 0, "STALE": 0,
	}
	for status, wantN := range want {
		if counts[status] != wantN {
			t.Errorf("Counts()[%q] = %d, want %d", status, counts[status], wantN)
		}
	}
	if len(counts) != 7 {
		t.Errorf("len(Counts()) = %d, want 7", len(counts))
	}
}

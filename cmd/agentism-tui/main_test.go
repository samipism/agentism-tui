package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestErrMsg(t *testing.T) {
	got := errMsg(errors.New("boom"))
	want := "agentism-tui: boom\n"
	if got != want {
		t.Errorf("errMsg = %q, want %q", got, want)
	}
}

// TestMainSmoke builds the binary and runs it against the store package's
// "valid" fixture project. It sends 'q' on stdin and expects the program to
// quit with exit code 0 within a short timeout.
func TestMainSmoke(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "agentism-tui")
	build := exec.Command("go", "build", "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %s\n%s", err, out)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin)
	cmd.Dir = filepath.Join("..", "..", "internal", "store", "testdata", "valid")
	cmd.Stdin = bytes.NewBufferString("q")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("binary exited with error: %s\nstderr: %s", err, stderr.String())
	}
}

func TestHasProject(t *testing.T) {
	t.Run("config.json present returns true", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(dir, ".agentism"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".agentism", "config.json"), []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}
		if !hasProject(dir) {
			t.Error("hasProject = false, want true")
		}
	})

	t.Run("empty directory returns false", func(t *testing.T) {
		dir := t.TempDir()
		if hasProject(dir) {
			t.Error("hasProject = true, want false")
		}
	})

	t.Run(".agentism as a file, not a directory, returns false", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ".agentism"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if hasProject(dir) {
			t.Error("hasProject = true, want false")
		}
	})
}

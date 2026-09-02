package main

import (
	"os"
	"path/filepath"
	"testing"
)

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

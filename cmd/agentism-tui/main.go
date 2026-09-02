// Command agentism-tui is a terminal dashboard for an agentism project.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samipism/agentism-tui/internal/store"
	"github.com/samipism/agentism-tui/internal/ui"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, errMsg(err))
		os.Exit(1)
	}

	if !hasProject(dir) {
		fmt.Fprintln(os.Stderr, `No agentism project found in this folder. Run "agentism init" first.`)
		os.Exit(1)
	}

	project, err := store.Load(dir)
	if err != nil {
		fmt.Fprint(os.Stderr, errMsg(err))
		os.Exit(1)
	}

	// tea.WithInput pins the input to os.Stdin as given, instead of letting
	// Bubble Tea fall back to opening /dev/tty when stdin isn't a terminal
	// (piped input in tests, for example).
	if _, err := tea.NewProgram(ui.New(project, dir), tea.WithInput(os.Stdin)).Run(); err != nil {
		fmt.Fprint(os.Stderr, errMsg(err))
		os.Exit(1)
	}
}

// errMsg formats err as this program's standard stderr message, newline
// included.
func errMsg(err error) string {
	return fmt.Sprintf("agentism-tui: %s\n", err)
}

// hasProject reports whether dir has an agentism project directly under it,
// i.e. <dir>/.agentism/config.json exists as a regular file. Any os.Stat
// failure, including a permission error, counts as false.
func hasProject(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, ".agentism", "config.json"))
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

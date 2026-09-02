// Command agentism-tui is a terminal dashboard for an agentism project.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "agentism-tui: %s\n", err)
		os.Exit(1)
	}

	if !hasProject(dir) {
		fmt.Fprintln(os.Stderr, `No agentism project found in this folder. Run "agentism init" first.`)
		os.Exit(1)
	}
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

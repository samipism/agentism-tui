---
id: Task_0002
kind: task
title: CLI entry point and project discovery
status: PLANNED
priority: 20
tags:
  - cli
locked: true
depends_on: []
version: v0
updated_at: "2026-09-02T10:20:34Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Task created from the intent and the architecture
  - date: 2026-09-02
    kind: locked
    summary: Planning finished
---

# CLI entry point and project discovery

## Goal

The `agentism-tui` binary's entry point. It checks the current working
directory for an agentism project, prints a clear error and exits 1
when there is none, and otherwise loads the project and starts the
Bubble Tea program.

## Scope

- `cmd/agentism-tui/main.go`.
- Check that `.agentism/config.json` exists directly under the current
  working directory (no parent-directory search).
- Call `store.Load` (from Task_0001) on success.
- Turn a load error into a stderr message and exit code 1.
- Start the Bubble Tea program (`tea.NewProgram(...).Run()`) with the
  loaded `store.Project`, once Task_0003's model exists.

## Out Of Scope

- Any command-line flag or argument. The tool takes none.
- Any project path other than the current working directory.
- The dashboard's own screens and keybindings (Task_0003).

## Contracts

<!-- The source of truth. A ticket must not contradict this section.
     A change here marks finished tickets as STALE. -->

### Inputs

- The process's current working directory. No command-line arguments.

### Outputs

- Exit code 0 after the user quits the TUI normally.
- Exit code 1, with a message on stderr, when the current directory
  has no project or the project fails to load.

### Errors

- No `.agentism/config.json` in the current working directory: print
  `No agentism project found in this folder. Run "agentism init"
  first.` to stderr, exit 1. No TUI screen opens.
- `store.Load` returns any other error: print `agentism-tui: ` followed
  by the error's message to stderr, exit 1. No TUI screen opens.

### Invariants

- The program never writes a file under `.agentism/` or `agentism/`.
- The program never searches a parent directory for a project.

## Acceptance

Running the built binary in a folder with no `.agentism/config.json`
prints the exact error above to stderr and exits 1. Running it in this
project's own root, once Task_0003 exists, opens the dashboard.

## Design Notes

Keep `main.go` thin: find the project, load it, hand the result to the
`ui` package. All decision logic beyond "does config.json exist" lives
in `internal/store`.

## Resolved Questions

- Q: Should the tool search parent directories for a project, the way
  the `agentism` CLI does?
  A: No. FR-10 in INTENT.md fixes this: the current working directory
  only, no path argument, no upward search.

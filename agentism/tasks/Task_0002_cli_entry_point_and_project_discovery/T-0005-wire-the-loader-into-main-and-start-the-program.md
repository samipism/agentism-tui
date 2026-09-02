---
id: T-0005
kind: ticket
title: Wire the loader into main and start the program
task: Task_0002
status: BLOCKED
priority: 20
tags:
  - cli
depends_on:
  - T-0004
  - T-0002
contract_hash: f335024c7f1074f4
commits: []
updated_at: "2026-09-02T12:41:38Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0002
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
  - date: 2026-09-02
    kind: "status:BLOCKED"
    summary: "Ticket work item requires handing the loaded *store.Project to 'the Bubble Tea program' and running it, but that program doesn't exist yet: Task_0003 (Bubble Tea dashboard) is still PLANNED with 0/2 tickets done, no bubbletea dependency in go.mod, no tea.Program anywhere in the codebase. Cannot implement without either building Task_0003's dashboard inside this ticket (breaks its contract) or guessing an API that isn't designed yet."
---

# Wire the loader into main and start the program

## Work

`main.go`: when T-0004's check finds a project, call `store.Load`,
turn any error into the stderr message and exit 1, otherwise hand the
loaded `*store.Project` to the Bubble Tea program and run it.

## Contract

### Inputs

- The current working directory (already confirmed to hold a project
  by T-0004).

### Outputs

- Exit code 0 once the Bubble Tea program returns normally (the user
  quit).
- Exit code 1, with `agentism-tui: <message>` on stderr, when
  `store.Load` returns an error.

### Errors

- `store.Load` error: print `agentism-tui: ` plus the error's message
  to stderr, exit 1, never start the program.
- `tea.Program.Run()` error: print `agentism-tui: ` plus the error's
  message to stderr, exit 1.

## Tests First

- A `main` helper function that maps a `store.Load` error to the exact
  stderr string and exit code, tested directly (without spawning the
  real binary) for a sample error.
- An end-to-end smoke test that runs the built binary (via
  `os/exec`) against a fixture project directory and checks it exits 0
  within a short timeout when fed a 'q' keypress on stdin.

## Acceptance

Running the built binary in a folder with a valid agentism project
opens the dashboard from Task_0003. Pressing 'q' exits with status 0.
A `store.Load` error (for example a corrupt `plan.md`) prints to
stderr and exits 1 without opening the dashboard.

## Human Verification Steps

1. Run the built binary in this repository's root. Expect the
   dashboard to open (once Task_0003 ships; before that, expect it to
   reach the point of trying to start the program without crashing on
   the loader step).
2. Press 'q'. Expect the terminal to return to the shell prompt and
   `echo $?` to print `0`.

---
id: T-0004
kind: ticket
title: Find the project and print the no-project error
task: Task_0002
status: DONE
priority: 10
tags:
  - cli
depends_on: []
contract_hash: f335024c7f1074f4
commits:
  - 60bb7e811ca42036b8fefe5fedcfe5314e80626d
updated_at: "2026-09-02T12:40:20Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0002
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
  - date: 2026-09-02
    kind: "status:IN_REVIEW"
    summary: Moved from IN_PROGRESS
  - date: 2026-09-02
    kind: "status:DONE"
    summary: Human accepted the result
---

# Find the project and print the no-project error

## Work

A function that checks whether `.agentism/config.json` exists directly
under a given directory, for `main.go` to call against the current
working directory.

## Contract

### Inputs

- A directory path (in production, `os.Getwd()`'s result).

### Outputs

- `true` when `<dir>/.agentism/config.json` exists as a regular file,
  `false` otherwise.

### Errors

- None. A permission error or any other `os.Stat` failure counts as
  `false` (no project found), not a program error.

## Tests First

- A temp directory with `.agentism/config.json` present returns
  `true`.
- An empty temp directory returns `false`.
- A temp directory with a `.agentism` file (not a directory) returns
  `false`.

## Acceptance

`go test ./cmd/agentism-tui/... -run TestHasProject` passes. Running
the built binary in a folder without `.agentism/config.json` prints
`No agentism project found in this folder. Run "agentism init"
first.` to stderr and exits 1.

## Human Verification Steps

1. `cd` into an empty temp folder and run the built `agentism-tui`
   binary. Expect the exact stderr message above and exit code 1
   (check with `echo $?`).
2. `cd` back into this repository's root and confirm the same command
   does not print that message (it may fail later, on wiring not yet
   built by T-0005).

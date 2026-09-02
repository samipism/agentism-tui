---
id: T-0003
kind: ticket
title: Load the activity log and compute ticket-status counts
task: Task_0001
status: TODO
priority: 30
tags: [parser]
depends_on: [T-0002]
contract_hash: 58780ac705de007e
commits: []
updated_at: 2026-09-02T10:16:25Z
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0001
---

# Load the activity log and compute ticket-status counts

## Work

Add the activity log to `Project`, and a `Counts` method that tallies
every ticket by status for the dashboard's status bar.

## Contract

### Inputs

- `.agentism/log.jsonl` (may be absent): zero or more lines, each a
  JSON object with `at` and `kind` string fields plus any other
  fields.
- The `Project` that T-0002's `Load` already built.

### Outputs

- `Project.Log []LogEntry`, `LogEntry{ At, Kind string; Fields
  map[string]any }`, in file order (oldest first).
- `Project.Counts() map[string]int`, one entry for each of `TODO`,
  `IN_PROGRESS`, `IN_REVIEW`, `DONE`, `NOT_ACCEPTED`, `BLOCKED`,
  `STALE`, counting every ticket across every task.

### Errors

- A missing `log.jsonl`: `Project.Log` is an empty slice, not an
  error.
- A line that fails to parse as JSON: `Load` skips it and keeps
  reading the rest of the file.

## Tests First

- A `log.jsonl` with three valid lines loads three `LogEntry` values
  in file order, with `Fields` holding the extra keys.
- A `log.jsonl` with one malformed line among two valid ones loads
  the two valid entries; the malformed line contributes nothing.
- An absent `log.jsonl` loads an empty `Project.Log`.
- `Counts()` on a project with tickets in four different statuses
  returns the exact count for each of the seven statuses, zero for a
  status no ticket uses.

## Acceptance

`go test ./internal/store/... -run TestLoadLog` and `-run TestCounts`
pass.

## Human Verification Steps

1. Run `go test ./internal/store/... -run TestLoadLog -v`. Expect
   every subtest to print `PASS`.
2. Run `go test ./internal/store/... -run TestCounts -v`. Expect every
   subtest to print `PASS`.

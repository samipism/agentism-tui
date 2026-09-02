---
id: T-0002
kind: ticket
title: Load config, tasks, and tickets into a Project
task: Task_0001
status: IN_PROGRESS
priority: 20
tags:
  - parser
depends_on:
  - T-0001
contract_hash: 58780ac705de007e
commits: []
updated_at: "2026-09-02T12:15:34Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0001
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
---

# Load config, tasks, and tickets into a Project

## Work

`store.Load(root string) (*Project, error)`: read
`.agentism/config.json`, walk `agentism/tasks/`, and assemble every
task and its tickets, using T-0001's parser for each markdown file.

## Contract

### Inputs

- A project root directory path.
- `.agentism/config.json` (may be absent).
- `agentism/tasks/Task_NNNN_<slug>/plan.md` and its
  `T-NNNN-<slug>.md` ticket files.

### Outputs

- `Project{ Phase, Version string; Tags []string; Tasks []Task }`.
- `Task{ ID, Title, Slug, Status string; Priority int; Tags,
  DependsOn []string; Goal, Contract string; Tickets []Ticket }`.
- `Ticket{ ID, Title, TaskID, Status string; Priority int; Tags,
  DependsOn []string; ContractHash, Contract, Work, Acceptance
  string }`.
- `Project.Tasks` sorted by priority then ID. Each task's `Tickets`
  sorted by priority then ID.

### Errors

- No `.agentism/config.json`: `Load` returns `ErrNoProject`.
- A task folder name that does not match `Task_\d{4}_.+`: skipped, no
  error.
- A ticket file name that does not match `T-\d{4}-.+\.md`: skipped, no
  error.
- A `plan.md` or ticket file whose frontmatter fails to parse (using
  T-0001's `ParseDoc`): `Load` returns an error naming the file's path
  relative to the project root.
- A ticket `status` value outside `TODO`, `IN_PROGRESS`, `IN_REVIEW`,
  `DONE`, `NOT_ACCEPTED`, `BLOCKED`, `STALE`: falls back to `TODO`.

## Tests First

- A project with two tasks, one holding two tickets and one holding
  none, loads into a `Project` with both tasks and the right ticket
  counts.
- Task and ticket order in the result matches priority-then-ID, given
  fixtures whose files are not already in that order.
- A missing `.agentism/config.json` returns `ErrNoProject`.
- A task folder named `NotATask` is absent from `Project.Tasks`.
- A ticket file named `notaticket.md` inside a task folder does not
  appear in that task's `Tickets`.
- A ticket file with `status: WEIRD` loads with `Status == "TODO"`.
- A `plan.md` with broken frontmatter makes `Load` return an error
  that names the file's path.

## Acceptance

`go test ./internal/store/... -run TestLoad` passes against fixture
projects under `internal/store/testdata/`.

## Human Verification Steps

1. Run `go test ./internal/store/... -run TestLoad -v`. Expect every
   subtest to print `PASS`.
2. In this repository (`agentism-tui`), run a short throwaway
   `go run` snippet that calls `store.Load(".")` and prints the
   returned task count. Expect it to match `agentism task list`'s
   count run in the same folder.

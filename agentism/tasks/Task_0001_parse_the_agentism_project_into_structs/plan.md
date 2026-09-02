---
id: Task_0001
kind: task
title: Parse the agentism project into structs
status: PLANNED
priority: 10
tags:
  - parser
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

# Parse the agentism project into structs

## Goal

A package that reads an agentism project from disk and returns a plain
struct tree: the phase, the tags, every task, every task's tickets, and
the activity log. Nothing in this package writes a file or runs the
`agentism` CLI.

## Scope

- Read `.agentism/config.json` for the phase, the version, and the tag
  list.
- Read `.agentism/log.jsonl` for the activity log.
- Walk `agentism/tasks/Task_NNNN_<slug>/` folders. Read each `plan.md`
  for the task's frontmatter and its `Contracts` and `Goal` sections.
- Read each `T-NNNN-<slug>.md` ticket file in a task folder for the
  ticket's frontmatter and its `Contract`, `Work`, and `Acceptance`
  sections.
- Parse the frontmatter YAML subset with `gopkg.in/yaml.v3`.
- Extract a named markdown section the same way agentism's own
  `mdfile.ts` does: match a heading case-insensitively, ignore leading
  numbers and punctuation, and capture every line until the next
  heading at the same or a shallower level.
- Sort tasks by priority then ID. Sort each task's tickets by priority
  then ID. This matches the order `agentism status` shows.
- Compute a ticket-status count (`TODO`, `IN_PROGRESS`, `IN_REVIEW`,
  `DONE`, `NOT_ACCEPTED`, `BLOCKED`, `STALE`) across every ticket, for
  the status bar.

## Out Of Scope

- `agentism/versions/` (archived tasks from a finished version). Only
  `agentism/tasks/` loads in this pass.
- `.agentism/state.json`. This package never reads it; it rebuilds the
  same data from the markdown files every time.
- Any write, or any call to the `agentism` binary.
- Rendering markdown to a display string. This package returns raw
  section text; the ui task decides how to show it.

## Contracts

<!-- The source of truth. A ticket must not contradict this section.
     A change here marks finished tickets as STALE. -->

### Inputs

- A project root directory path (an absolute path to the folder that
  holds `.agentism/` and `agentism/`).
- `.agentism/config.json` - JSON object with `phase`, `version`, and
  `tags` (a string array). May be absent.
- `.agentism/log.jsonl` - zero or more lines, each a JSON object with
  an `at` (ISO-8601 string) and a `kind` (string) field, plus any other
  fields. May be absent.
- `agentism/tasks/Task_NNNN_<slug>/plan.md` - YAML frontmatter with
  `id`, `title`, `status`, `priority`, `tags`, `depends_on`; a `##
  Contracts` section; a `## Goal` section.
- `agentism/tasks/Task_NNNN_<slug>/T-NNNN-<slug>.md` - YAML frontmatter
  with `id`, `title`, `task`, `status`, `priority`, `tags`,
  `depends_on`, `contract_hash`; a `## Contract` section; a `## Work`
  section; a `## Acceptance` section.

### Outputs

- `Project{ Phase, Version string; Tags []string; Tasks []Task; Log
  []LogEntry }`.
- `Task{ ID, Title, Slug, Status string; Priority int; Tags,
  DependsOn []string; Goal, Contract string; Tickets []Ticket }`.
- `Ticket{ ID, Title, TaskID, Status string; Priority int; Tags,
  DependsOn []string; ContractHash, Contract, Work, Acceptance
  string }`.
- `LogEntry{ At, Kind string; Fields map[string]any }` where `Fields`
  holds every JSON key other than `at` and `kind`.
- `Counts() map[string]int`, one entry for each of the seven ticket
  statuses, computed from `Project.Tasks`.

### Errors

- The project root has no `.agentism/config.json`: return a
  `ErrNoProject` sentinel error. The caller (the cli task) turns this
  into the exit-1 message.
- The loader skips a folder under `agentism/tasks/` whose name does
  not match `Task_\d{4}_.+`. It raises no error for it.
- The loader skips a file under a task folder whose name does not
  match `T-\d{4}-.+\.md`. It raises no error for it.
- A `plan.md` or a ticket file with frontmatter that fails to parse
  returns an error naming the file's path relative to the project
  root.
- The loader skips a line in `log.jsonl` that fails to parse as JSON,
  so one bad line never blocks the rest of the log.
- A ticket `status` value outside the seven known statuses falls back
  to `TODO`, matching the CLI's own behaviour.

### Invariants

- The loader never writes a file and never runs the `agentism` binary.
- Reading the same unchanged project twice returns equal data.
- A task with zero tickets still appears in `Project.Tasks`.

## Acceptance

`go test ./internal/store/...` passes against fixture projects under
`internal/store/testdata/`, covering: a project with two tasks, one
with two tickets and one with none; an absent `log.jsonl`; a
`log.jsonl` with one malformed line among valid ones; a project
folder with no `.agentism/config.json`; and a ticket with an unknown
status value.

## Design Notes

The loader is one function, `store.Load(root string) (*Project,
error)`. No interface: there is one data source, so an interface
would only add a name to jump through.

## Resolved Questions

- Q: Should this package read `agentism/versions/` (archived tasks)?
  A: Not in v1. The intent scopes browsing to the active plan; nothing
  in this project has an archived version yet, and the scan logic can
  extend to it later without changing the struct shapes.

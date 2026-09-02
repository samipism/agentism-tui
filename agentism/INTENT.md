---
kind: intent
version: v0
status: FINAL
updated_at: "2026-09-02T10:06:42Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: First version from the intention questions
  - date: 2026-09-02
    kind: finalised
    summary: Human approved the document
---

# Intent

## Problem

A developer checks an agentism project's plan by running many separate CLI
commands: `status`, `task list`, `ticket list`, `task show`. The person must
piece the state together by hand. There is no single view of the plan tree
and its progress. This wastes time and hides the overall picture behind a
sequence of commands.

## Users

**Developer / operator.** Runs `agentism-tui` in a project folder to see the
current phase, task tree, and ticket statuses at a glance. This person
already knows agentism concepts (INTENT, ARCH, tasks, tickets, phases) from
using the CLI.

## Use Cases

- View the project's current phase and progress.
- Browse the task list and expand a task to see its tickets.
- View a ticket's full detail: title, status, task, and contract.
- Browse the activity log.
- Refresh the view to pick up changes made outside the TUI.

## User Path Flows

**Browse tasks and tickets**

1. The user runs `agentism-tui` in a project folder.
2. The tool reads the agentism project data and shows one dashboard: a
   status bar (phase, progress) on top, a task/ticket tree in the middle,
   and a detail pane.
3. The user moves the selection with arrow keys or j/k.
4. The user selects a task. The tree expands or collapses that task's
   tickets.
5. The user selects a ticket. The detail pane shows its title, status,
   task, and full contract.
6. The user presses 'r' to re-read the project data and redraw the view.
7. The user presses 'q' or Ctrl+C to quit.

**No project found**

1. The user runs `agentism-tui` in a folder with no agentism project.
2. The tool checks for the project and finds none.
3. The tool prints "No agentism project found in this folder. Run
   'agentism init' first." to stderr.
4. The tool exits with status code 1. No TUI screen opens.

## Functional Requirements

- FR-1: The tool shows the project's current phase and progress on launch.
- FR-2: The tool lists every task and shows each task's status.
- FR-3: The tool lists every ticket under a task and shows each ticket's
  status.
- FR-4: The tool shows a selected ticket's title, status, task, and full
  contract.
- FR-5: The tool shows the activity log.
- FR-6: The tool reads project data once at launch as a snapshot. The 'r'
  key re-reads the data and redraws the view. The tool does not watch
  files in the background.
- FR-7: The tool reads the agentism project folder's files directly (the
  markdown files with YAML frontmatter under the agentism folder, for
  example INTENT.md, ArchDecision.md, task plans, and ticket files). It
  does not depend on the `agentism` CLI, because that CLI only runs
  inside a plugin host and is not on a normal user's PATH.
- FR-8: The tool never writes to any file under the agentism project
  folder. It opens every file read-only.
- FR-9: If no agentism project exists in the current working directory,
  the tool prints an error to stderr and exits with status code 1,
  without opening the TUI.
- FR-10: The tool reads the project in the current working directory
  only. It takes no project path argument and offers no in-app project
  switch.

## Non Functional Requirements

- NFR-1: The tool starts and renders the dashboard in under 200 ms when the
  project has up to 100 tasks and 1000 tickets.
- NFR-2: The tool uses under 100 MB of memory for such a project.
- NFR-3: The project ships release binaries for darwin/amd64,
  darwin/arm64, linux/amd64, and linux/arm64.

## Out Of Scope

- Editing or writing any agentism data. The tool is view-only.
- Accepting or rejecting tickets, changing ticket status, or running
  `agentism init` from within the tool.
- Switching between multiple projects within one running session.
- Reading a project over the network or from a remote source. Local
  filesystem only.
- Windows binaries.

## Open Questions

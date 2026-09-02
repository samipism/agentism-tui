---
id: T-0007
kind: ticket
title: Detail pane, log view, and keybindings
task: Task_0003
status: IN_PROGRESS
priority: 20
tags:
  - ui
depends_on:
  - T-0006
  - T-0003
contract_hash: 8c379b03406644e5
commits:
  - ec31fb3
updated_at: "2026-09-02T14:02:51Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0003
  - date: 2026-09-02
    kind: contract-change
    summary: Detail pane now renders section markdown through glamour instead of printing it raw, per the architecture update adding that library. Updated Scope, Out Of Scope, Contracts, and Design Notes; T-0007's contract restamped to match.
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
  - date: 2026-09-02
    kind: "status:IN_REVIEW"
    summary: Moved from IN_PROGRESS
  - date: 2026-09-02
    kind: "status:DONE"
    summary: Human accepted the result
  - date: 2026-09-02
    kind: contract-change
    summary: "Dashboard read as plain text, not a dashboard. Replaced the status-bar-plus-two-columns layout with four bordered regions (header, sidebar, main, footer), added a fixed status-to-color lookup, and added a completion bar in the header and per task row. No new library: lipgloss already covered borders and color."
  - date: 2026-09-02
    kind: "status:TODO"
    summary: contract changed
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
---

# Detail pane, log view, and keybindings

## Work

Extend T-0006's model with a detail pane that fills in when the user
selects a ticket, rendering that ticket's markdown text through
`glamour` instead of printing it raw. Add a toggle-able log view built
from `Project.Log`, and the remaining keys: 'r' to refresh, 'l' to
toggle the log view, 'q' and Ctrl+C to quit.

## Contract

### Inputs

- `tea.KeyMsg` for enter/space on a ticket line (select it), 'r'
  (refresh), 'l' (toggle log view), 'q' and Ctrl+C (quit).
- The project root path, so 'r' can call `store.Load` again.

### Outputs

- Selecting a ticket line fills the detail pane with that ticket's
  title, status, task, and its `Contract`, `Work`, and `Acceptance`
  text, run through `glamour.Render` before display.
- 'r' replaces the model's `*store.Project` with a freshly loaded one
  and redraws; a load error shows as one line in the status bar
  without losing the previous data.
- 'l' shows or hides a scrollable view of `Project.Log`, newest entry
  first.
- 'q' or Ctrl+C ends the Bubble Tea program (`tea.Quit`).

### Errors

- A `store.Load` error on refresh never crashes the program or clears
  the screen; it shows as a status bar message and keeps the last good
  `*store.Project`.

## Tests First

- Selecting a ticket line then calling `View()` includes that ticket's
  title and its contract text, styled (for example a heading in the
  contract renders bold, not as a literal `##`).
- Simulating 'r' with a loader that returns an error leaves the tree
  showing the previous project's tasks and shows the error text in the
  status bar.
- Simulating 'r' with a loader that returns a changed project (for
  example one more DONE ticket) updates the counts shown.
- Pressing 'l' twice returns the view to the pre-toggle state.
- Pressing 'q' returns a `tea.Quit` command from `Update`.

## Acceptance

Selecting a ticket shows its contract in the detail pane. Pressing 'r'
after an on-disk change (a ticket's status file hand-edited between
runs of the loader) shows the new status. Pressing 'l' shows and hides
the activity log. Pressing 'q' quits.

## Human Verification Steps

1. Run the binary against a fixture project. Select a ticket. Confirm
   the detail pane shows its contract text formatted (a heading stands
   out, a list shows bullets), not as raw `##` and `-` markdown source.
2. Edit that ticket's `status:` field on disk to a different value,
   save, then press 'r' in the running program. Confirm the tree shows
   the new status without restarting the program.
3. Press 'l'. Confirm a log view appears. Press 'l' again. Confirm it
   disappears.
4. Press 'q'. Confirm the program exits and the shell prompt returns.

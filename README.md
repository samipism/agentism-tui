# agentism-tui

A read-only terminal viewer for an [agentism](https://github.com/samipism/agentism)
project's plan. It shows the project's phase, its tasks and tickets, and
each ticket's contract, in a scrollable dashboard. It never writes a file
and never spawns a process.

## Install

Download a release binary for your platform from the
[Releases page](https://github.com/samipism/agentism-tui/releases), or
build it with Go:

```sh
go install github.com/samipism/agentism-tui/cmd/agentism-tui@latest
```

## Usage

Run `agentism-tui` inside a project folder that has a `.agentism/`
directory:

```sh
cd your-agentism-project
agentism-tui
```

### Keys

| Key | Action |
|---|---|
| `↑` / `k` | Move the cursor up |
| `↓` / `j` | Move the cursor down |
| `Enter` / `Space` | Expand or collapse a task, or select a ticket |
| `PgUp` / `PgDn` | Scroll the detail pane or log |
| `r` | Refresh the project from disk |
| `l` | Toggle the log view |
| `q` / `Ctrl+C` | Quit |

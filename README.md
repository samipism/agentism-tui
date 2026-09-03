# agentism-tui

A read-only terminal viewer for an [agentism](https://github.com/samipism/agentism)
project's plan. It shows the project's phase, its tasks and tickets, and
each ticket's contract, in a scrollable dashboard. It never writes a file
and never spawns a process.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/samipism/agentism-tui/main/install.sh | sh
```

Downloads the latest release for your platform (darwin/linux, amd64/arm64),
verifies its checksum, and installs `agentism-tui` to `/usr/local/bin`
(override with `INSTALL_DIR`).

Or do it by hand: download a release archive from the
[Releases page](https://github.com/samipism/agentism-tui/releases), extract it,
and put the `agentism-tui` binary on your `PATH`:

```sh
tar -xzf agentism-tui_*.tar.gz
mv agentism-tui /usr/local/bin/
```

Or build it with Go:

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

# Boba Text 

A modern, Vim-inspired terminal code editor with a built-in AI agent. Built with Go, Bubble Tea, and Lipgloss.


## Features

- **Neon/Dark Theme**: A sleek interface with pink, purple, and neon blue accents.
- **Vim-like Editing**: Modal editing with Normal and Insert modes.
- **File Tree**: Collapsible sidebar for navigating your project.
- **AI Agent**: Dedicated tab for chatting with an AI assistant about your code.
- **Command Bar**: Save files and quit using `:w`, `:s`, `:q`.
- **Instant Startup**: Pre-compiled binary execution for zero-latency launch.

## Quick Start

**Windows:**

Run the included shortcut script to build (if needed) and launch instantly:

```powershell
.\boba
```

**Manual Build:**

```powershell
go build -o boba-text.exe main.go
./boba-text.exe
```

## Keybindings

| Context | Key | Action |
| :--- | :--- | :--- |
| **Global** | `Ctrl+b` | Toggle File Tree sidebar |
| **Global** | `Ctrl+e` | Focus / Open File Tree |
| **Global** | `Ctrl+a` | Focus AI Agent |
| **Global** | `Tab` | Cycle focus (Tree -> Editor -> Agent) |
| **Global** | `Ctrl+c` / `q` | Quit |
| **File Tree** | `j` / `k` | Navigate Down / Up |
| **File Tree** | `Enter` | Open File / Enter Dir / Go Up (`..`) |
| **Editor** | `i` | Enter **Insert Mode** |
| **Editor** | `Esc` | Enter **Normal Mode** |
| **Editor** | `:` | Open **Command Bar** (`:w` to save) |
| **Agent** | `Enter` | Send message to AI |

## Stack

- **Language**: [Go](https://go.dev/)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles)

## ⚙️ Configuration

Create a `config.toml` (or `~/.boba-config.toml`) to customize colors, keys, AI, and commands:

```toml
[colors]
primary = "#F25D94"
text = "#FAFAFA"

[keys]
toggle_tree = "ctrl+b"
focus_agent = "ctrl+a"

[ai]
name = "BobaBot"

[commands]
save = ["w", "s", "save"]
quit = ["q", "quit"]
```

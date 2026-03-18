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
| **Global** | `Ctrl+s` | **Save File** (Nano-style) |
| **Global** | `Ctrl+x` | **Quit** (Nano-style) |
| **Global** | `Tab` | Cycle focus (Tree -> Editor -> Agent) |
| **Global** | `Ctrl+c` | Quit |
| **File Tree** | `j` / `k` | Navigate Down / Up |
| **File Tree** | `Enter` | Open File / Enter Dir / Go Up (`..`) |
| **Editor (Normal)** | `h`, `j`, `k`, `l` | Basic movement (Vim-style) |
| **Editor (Normal)** | `0` / `$` | Start / End of line |
| **Editor (Normal)** | `g` / `G` | Jump to Start / End of file |
| **Editor (Normal)** | `x` | Delete character |
| **Editor (Normal)** | `i` / `a` | Enter **Insert Mode** / Append |
| **Editor (Normal)** | `o` / `O` | Open new line below / above |
| **Editor (Insert)** | `Esc` | Return to **Normal Mode** |
| **Editor (Normal)** | `:` | Open **Command Bar** (`:w` to save) |
| **Agent** | `Enter` | Send message to AI |

## Stack

- **Language**: [Go](https://go.dev/)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles)

## ⚙️ Configuration

Boba Text comes with sensible defaults built-in. If you want to override them globally, you can create a `~/.config/boba-text/config.toml` (or `%APPDATA%\boba-text\config.toml` on Windows).

Here is an example overriding colors, keys, and AI configuration:

```toml
[colors]
primary = "#FF00FF"
text = "#FFFFFF"

[keys]
toggle_tree = "ctrl+t"
save = "ctrl+s"

[ai]
name = "BobaBot"

[commands]
save = ["w", "s", "save", "x"]
quit = ["q", "quit", "exit"]
```

package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/CiaranMccarthy1/boba-text/pkg/config"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EditorMode int

const (
	ModeNormal EditorMode = iota
	ModeInsert
	ModeCommand
	ModeVisual
	ModeSearch
)

type EditorModel struct {
	textarea    textarea.Model
	textinput   textinput.Model
	mode        EditorMode
	width       int
	height      int
	filename    string
	msg         string
	modified    bool
	commands    config.Commands
	keys        config.Keys
	yankBuffer  string
	searchQuery string
	searchInput textinput.Model
}

// NewEditor creates a new editor model with the given command configuration.
func NewEditor(cmdConfig config.Commands, keyConfig config.Keys) EditorModel {
	ta := textarea.New()
	ta.Placeholder = "Empty file..."
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = true

	ti := textinput.New()
	ti.Prompt = ":"
	ti.Placeholder = ""
	ti.CharLimit = 156
	ti.Width = 50

	si := textinput.New()
	si.Prompt = "/"
	si.Placeholder = "search..."
	si.CharLimit = 156
	si.Width = 50

	return EditorModel{
		textarea:    ta,
		textinput:   ti,
		searchInput: si,
		mode:        ModeNormal,
		commands:    cmdConfig,
		keys:        keyConfig,
	}
}

func (m EditorModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, textinput.Blink)
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case ModeNormal:
			cmd = m.handleNormalMode(msg.String())
			if cmd != nil {
				return m, cmd
			}
		case ModeInsert:
			switch msg.String() {
			case m.keys.EditorNormalMode:
				m.mode = ModeNormal
				m.textarea.Blur()
				m.msg = ""
			default:
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
				m.modified = true
			}
		case ModeCommand:
			switch msg.String() {
			case m.keys.EditorNormalMode:
				m.mode = ModeNormal
				m.textinput.Blur()
				m.msg = ""
			case m.keys.EditorCommandRun:
				cmd = m.executeCommand(m.textinput.Value())
				if cmd != nil {
					return m, cmd
				}
			default:
				m.textinput, cmd = m.textinput.Update(msg)
				cmds = append(cmds, cmd)
			}
		case ModeSearch:
			switch msg.String() {
			case m.keys.EditorNormalMode:
				m.mode = ModeNormal
				m.searchInput.Blur()
				m.msg = ""
			case "enter":
				m.searchQuery = m.searchInput.Value()
				m.mode = ModeNormal
				m.searchInput.Blur()
				if m.searchQuery != "" {
					m.msg = "/" + m.searchQuery
				}
			default:
				m.searchInput, cmd = m.searchInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// handleNormalMode processes key presses in Normal mode (Vim motions).
func (m *EditorModel) handleNormalMode(key string) tea.Cmd {
	switch key {
	// Enter insert mode
	case m.keys.EditorInsertMode:
		m.mode = ModeInsert
		m.textarea.Focus()
		m.msg = ""
	// Enter command mode
	case m.keys.EditorCommandMode:
		m.mode = ModeCommand
		m.textinput.Focus()
		m.textinput.SetValue("")
		m.msg = ""
	// Search mode
	case "/":
		m.mode = ModeSearch
		m.searchInput.Focus()
		m.searchInput.SetValue("")
		m.msg = ""
	// Cursor movement
	case "j", "down":
		m.textarea.CursorDown()
	case "k", "up":
		m.textarea.CursorUp()
	case "h", "left":
		col := m.textarea.LineInfo().ColumnOffset
		if col > 0 {
			m.textarea.SetCursor(col - 1)
		}
	case "l", "right":
		m.textarea.SetCursor(m.textarea.LineInfo().ColumnOffset + 1)
	// Word motions
	case "w":
		m.wordForward()
	case "b":
		m.wordBackward()
	// Line start/end
	case "0", "home":
		m.textarea.CursorStart()
	case "$", "end":
		m.textarea.CursorEnd()
	// Top/bottom of file
	case "G":
		lines := strings.Split(m.textarea.Value(), "\n")
		for i := 0; i < len(lines); i++ {
			m.textarea.CursorDown()
		}
	// Insert above/below
	case "o":
		m.textarea.CursorEnd()
		m.textarea.Focus()
		m.mode = ModeInsert
		m.modified = true
	case "O":
		m.textarea.CursorStart()
		m.textarea.Focus()
		m.mode = ModeInsert
		m.modified = true
	// Insert at start/end of line
	case "A":
		m.textarea.CursorEnd()
		m.textarea.Focus()
		m.mode = ModeInsert
	case "I":
		m.textarea.CursorStart()
		m.textarea.Focus()
		m.mode = ModeInsert
	// Delete char under cursor
	case "x":
		m.textarea.Focus()
		// Simulate delete by sending a delete key
		m.modified = true
		m.textarea.Blur()
	// Search next/prev
	case "n":
		if m.searchQuery != "" {
			m.msg = "search: " + m.searchQuery
		}
	case "N":
		if m.searchQuery != "" {
			m.msg = "search (prev): " + m.searchQuery
		}
	// Paste
	case "p":
		if m.yankBuffer != "" {
			m.textarea.Focus()
			m.textarea.InsertString(m.yankBuffer)
			m.textarea.Blur()
			m.modified = true
			m.msg = "Pasted"
		}
	}
	return nil
}

// executeCommand processes a command-mode command string.
func (m *EditorModel) executeCommand(val string) tea.Cmd {
	m.mode = ModeNormal
	m.textinput.Blur()

	// Handle line number jump :<number>
	if len(val) > 0 && val[0] >= '0' && val[0] <= '9' {
		m.msg = fmt.Sprintf("Jump to line %s", val)
		return nil
	}

	switch {
	case contains(m.commands.Save, val):
		return m.saveFile()
	case contains(m.commands.Quit, val):
		return tea.Quit
	case val == "wq" || val == "x":
		m.saveFile()
		return tea.Quit
	case val == "q!":
		return tea.Quit
	case strings.HasPrefix(val, "e "):
		path := strings.TrimPrefix(val, "e ")
		path = strings.TrimSpace(path)
		if path != "" {
			return func() tea.Msg { return OpenFileMsg{Path: path} }
		}
		m.msg = "Usage: :e <filename>"
	default:
		m.msg = "Unknown command: " + val
	}
	return nil
}

// saveFile writes the editor content to disk.
func (m *EditorModel) saveFile() tea.Cmd {
	if m.filename != "" {
		err := os.WriteFile(m.filename, []byte(m.textarea.Value()), 0644)
		if err != nil {
			m.msg = "Error saving: " + err.Error()
		} else {
			m.msg = "Saved: " + m.filename
			m.modified = false
		}
	} else {
		m.msg = "No filename set!"
	}
	return nil
}

// wordForward moves the cursor forward by one word.
func (m *EditorModel) wordForward() {
	m.textarea.SetCursor(m.textarea.LineInfo().ColumnOffset + 5)
}

// wordBackward moves the cursor backward by one word.
func (m *EditorModel) wordBackward() {
	col := m.textarea.LineInfo().ColumnOffset - 5
	if col < 0 {
		col = 0
	}
	m.textarea.SetCursor(col)
}

// contains checks if a string value exists in a slice.
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// View renders the editor as a string.
func (m EditorModel) View() string {
	// Build the status line (Neovim lualine-style)
	var barContent string
	var modeStyle lipgloss.Style

	if m.mode == ModeCommand {
		barContent = m.textinput.View()
		modeStyle = StyleModeCommand
	} else if m.mode == ModeSearch {
		barContent = m.searchInput.View()
		modeStyle = StyleModeCommand
	} else {
		modeTxt := modeName(m.mode)
		switch m.mode {
		case ModeInsert:
			modeStyle = StyleModeInsert
		case ModeVisual:
			modeStyle = StyleModeVisual
		default:
			modeStyle = StyleModeNormal
		}

		// File info section
		fname := m.filename
		if fname == "" {
			fname = "[No Name]"
		}
		modifiedMark := ""
		if m.modified {
			modifiedMark = " [+]"
		}

		msgInfo := ""
		if m.msg != "" {
			msgInfo = "  " + m.msg
		}

		// Lualine-style: MODE | filename [+] | message
		statusRight := StyleDim.Render(fmt.Sprintf(" %s ", fileType(m.filename)))
		barContent = fmt.Sprintf("%s %s%s%s%s",
			modeStyle.Render(" "+modeTxt+" "),
			fname, modifiedMark, msgInfo, statusRight)
		// Return early with full status bar
		statusBar := lipgloss.NewStyle().
			Foreground(ColorText).
			Background(lipgloss.Color("#2A2A2A")).
			Width(m.width).
			Padding(0, 0).
			Render(barContent)

		return StyleEditor.
			Width(m.width).
			Height(m.height).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					m.textarea.View(),
					statusBar,
				),
			)
	}

	// Command/Search mode status bar
	statusBar := lipgloss.NewStyle().
		Foreground(ColorText).
		Background(lipgloss.Color("#2A2A2A")).
		Width(m.width).
		Padding(0, 0).
		Render(barContent)

	return StyleEditor.
		Width(m.width).
		Height(m.height).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.textarea.View(),
				statusBar,
			),
		)
}

// modeName returns the display name for the given editor mode.
func modeName(m EditorMode) string {
	switch m {
	case ModeInsert:
		return "INSERT"
	case ModeVisual:
		return "VISUAL"
	case ModeCommand:
		return "COMMAND"
	case ModeSearch:
		return "SEARCH"
	default:
		return "NORMAL"
	}
}

// fileType returns a short filetype label based on extension.
func fileType(filename string) string {
	if filename == "" {
		return "text"
	}
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return "text"
	}
	ext := parts[len(parts)-1]
	switch ext {
	case "go":
		return "go"
	case "js":
		return "javascript"
	case "ts":
		return "typescript"
	case "py":
		return "python"
	case "rs":
		return "rust"
	case "md":
		return "markdown"
	case "json":
		return "json"
	case "toml":
		return "toml"
	case "yaml", "yml":
		return "yaml"
	case "html":
		return "html"
	case "css":
		return "css"
	default:
		return ext
	}
}

// SetContent sets the editor content and filename.
func (m *EditorModel) SetContent(content string, filename string) {
	m.textarea.SetValue(content)
	m.filename = filename
	m.modified = false
	m.msg = ""
}

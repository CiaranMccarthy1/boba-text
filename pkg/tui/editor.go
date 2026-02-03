package tui

import (
	"fmt"
	"os"

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
)

type EditorModel struct {
	textarea  textarea.Model
	textinput textinput.Model
	mode      EditorMode
	width     int
	height    int
	filename  string
	msg       string
	commands  config.Commands
}

// NewEditor creates a new editor model with the given command configuration.
func NewEditor(cmdConfig config.Commands) EditorModel {
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

	return EditorModel{
		textarea:  ta,
		textinput: ti,
		mode:      ModeNormal,
		commands:  cmdConfig,
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
			switch msg.String() {
			case "i":
				m.mode = ModeInsert
				m.textarea.Focus()
				m.msg = ""
			case ":":
				m.mode = ModeCommand
				m.textinput.Focus()
				m.textinput.SetValue("")
				m.msg = ""
			}
		case ModeInsert:
			switch msg.String() {
			case "esc":
				m.mode = ModeNormal
				m.textarea.Blur()
				m.msg = ""
			default:
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
			}
		case ModeCommand:
			switch msg.String() {
			case "esc":
				m.mode = ModeNormal
				m.textinput.Blur()
				m.textarea.Focus()
				m.msg = ""
			case "enter":
				val := m.textinput.Value()
				m.mode = ModeNormal
				m.textinput.Blur()

				if contains(m.commands.Save, val) {
					if m.filename != "" {
						err := os.WriteFile(m.filename, []byte(m.textarea.Value()), 0644)
						if err != nil {
							m.msg = "Error saving: " + err.Error()
						} else {
							m.msg = "Saved: " + m.filename
						}
					} else {
						m.msg = "No filename set!"
					}
				} else if contains(m.commands.Quit, val) {
					return m, tea.Quit
				} else {
					m.msg = "Unknown command: " + val
				}
			default:
				m.textinput, cmd = m.textinput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
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
	var barContent string
	var statusStyle lipgloss.Style

	if m.mode == ModeCommand {
		barContent = m.textinput.View()
		statusStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(ColorDark).
			Padding(0, 1)
	} else {
		statusMode := modeName(m.mode)
		var statusColor lipgloss.Color
		if m.mode == ModeInsert {
			statusColor = ColorPrimary
		} else {
			statusColor = ColorSecondary
		}

		statusStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(statusColor).
			Padding(0, 1)

		msgInfo := ""
		if m.msg != "" {
			msgInfo = "  " + m.msg
		}

		barContent = fmt.Sprintf("%s  %s%s", statusMode, m.filename, msgInfo)
	}

	return StyleEditor.
		Width(m.width).
		Height(m.height).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.textarea.View(),
				statusStyle.Width(m.width).Render(barContent),
			),
		)
}

// modeName returns the display name for the given editor mode.
func modeName(m EditorMode) string {
	if m == ModeInsert {
		return "INSERT"
	}
	return "NORMAL"
}

// SetContent sets the editor content and filename.
func (m *EditorModel) SetContent(content string, filename string) {
	m.textarea.SetValue(content)
	m.filename = filename
}

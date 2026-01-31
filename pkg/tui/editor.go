package tui

import (
	"fmt"
	"os"

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
	msg       string // For status messages like "Saved!"
}

func NewEditor() EditorModel {
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
				m.textarea.Focus() // Return focus to textarea (but in normal mode logic)
				m.msg = ""
			case "enter":
				val := m.textinput.Value()
				m.mode = ModeNormal
				m.textinput.Blur()

				// Execute command
				switch val {
				case "w", "s", "save", "write":
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
				case "q", "quit":
					return m, tea.Quit
				default:
					m.msg = "Unknown command: " + val
				}
			default:
				m.textinput, cmd = m.textinput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	// Always update blinking cursors if needed?
	// But usually only active component needs update.
	// We did specific updates above.

	return m, tea.Batch(cmds...)
}

func (m EditorModel) View() string {
	// Status/Command Bar
	var barContent string
	var statusStyle lipgloss.Style

	if m.mode == ModeCommand {
		// Show textinput
		barContent = m.textinput.View()
		statusStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(ColorDark). // Or different background for command?
			Padding(0, 1)
	} else {
		// Show Status
		statusMode := modeName(m.mode)
		var statusColor lipgloss.Color
		if m.mode == ModeInsert {
			statusColor = ColorPrimary // Pink
		} else {
			statusColor = ColorSecondary // Purple
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
				statusStyle.Width(m.width).Render(barContent), // Ensure bar stretches?
			),
		)
}

func modeName(m EditorMode) string {
	if m == ModeInsert {
		return "INSERT"
	}
	return "NORMAL"
}

func (m *EditorModel) SetContent(content string, filename string) {
	m.textarea.SetValue(content)
	m.filename = filename
}

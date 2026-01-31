package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EditorMode int

const (
	ModeNormal EditorMode = iota
	ModeInsert
)

type EditorModel struct {
	textarea textarea.Model
	mode     EditorMode
	width    int
	height   int
	filename string
}

func NewEditor() EditorModel {
	ta := textarea.New()
	ta.Placeholder = "Empty file..."
	ta.Focus()
	ta.CharLimit = 0
	ta.ShowLineNumbers = true

	return EditorModel{
		textarea: ta,
		mode:     ModeNormal,
	}
}

func (m EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m EditorModel) Update(msg tea.Msg) (EditorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case ModeNormal:
			switch msg.String() {
			case "i":
				m.mode = ModeInsert
				// Re-focus to enable typing
				m.textarea.Focus()
				// Navigation in normal mode could be mapped here if we supported custom viewport nav
				// For now, bubbles/textarea handles navigation even when focused, but we want to intercept interactions
			}
		case ModeInsert:
			switch msg.String() {
			case "esc":
				m.mode = ModeNormal
				m.textarea.Blur()
			default:
				m.textarea, cmd = m.textarea.Update(msg)
			}
		}
	}

	return m, cmd
}

func (m EditorModel) View() string {
	// Status Bar
	statusMode := modeName(m.mode)
	var statusColor lipgloss.Color
	if m.mode == ModeInsert {
		statusColor = ColorPrimary // Pink
	} else {
		statusColor = ColorSecondary // Purple
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(ColorText).
		Background(statusColor).
		Padding(0, 1)

	statusText := fmt.Sprintf("%s  %s", statusMode, m.filename)

	// Textarea with borders handled by layout, or subtle border
	// We'll remove the heavy border here and let layout handle main structure,
	// or use a very subtle one.

	return StyleEditor.
		Width(m.width).
		Height(m.height).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.textarea.View(),
				statusStyle.Render(statusText),
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

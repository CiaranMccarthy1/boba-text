package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Focus int

const (
	FocusFileTree Focus = iota
	FocusEditor
	FocusAgent
)

type Model struct {
	fileTree FileTreeModel
	editor   EditorModel
	agent    AgentModel
	focus    Focus
	width    int
	height   int
}

func InitialModel(startPath string) Model {
	return Model{
		fileTree: NewFileTree(startPath),
		editor:   NewEditor(),
		agent:    NewAgent(),
		focus:    FocusFileTree,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fileTree.Init(), m.editor.Init(), m.agent.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q": // Global quit
			return m, tea.Quit
		case "tab":
			// Cycle focus: FileTree -> Editor -> Agent -> FileTree
			m.focus = (m.focus + 1) % 3
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.fileTree.width = 30
		m.fileTree.height = msg.Height - 2

		contentWidth := msg.Width - 34
		contentHeight := msg.Height - 2

		m.editor.width = contentWidth
		m.editor.height = contentHeight
		m.editor.textarea.SetWidth(contentWidth - 2)
		m.editor.textarea.SetHeight(contentHeight - 2)

		m.agent.SetSize(contentWidth, contentHeight)

	case OpenFileMsg:
		content, err := os.ReadFile(msg.Path)
		if err == nil {
			m.editor.SetContent(string(content), msg.Path)
			m.focus = FocusEditor
		}
	}

	// Update components based on focus or broadcast if needed
	// For simplicity, we can update all if they don't consume unexpected keys,
	// but focusing is safer to route keys.

	// Always update Filetree if it's focused? Or maybe separate Key handling?
	// To keep it simple, we delegate based on focus.
	switch m.focus {
	case FocusFileTree:
		m.fileTree, cmd = m.fileTree.Update(msg)
		cmds = append(cmds, cmd)
	case FocusEditor:
		m.editor, cmd = m.editor.Update(msg)
		cmds = append(cmds, cmd)
	case FocusAgent:
		m.agent, cmd = m.agent.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var content string
	if m.focus == FocusAgent {
		content = m.agent.View()
	} else {
		content = m.editor.View()
	}

	// Dynamic Border Color
	var borderColor lipgloss.Color
	switch m.focus {
	case FocusAgent:
		borderColor = ColorSuccess
	case FocusEditor:
		borderColor = ColorPrimary
	default:
		borderColor = ColorSubText
	}

	// Layout: FileTree | Content
	// Use a clean vertical bar or just spacing?
	// Reference image has clean separation.
	// We'll use a border on the content pane to frame it, but "Rounded" style from styles.go

	// Calculate content style:
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width - 34).
		Height(m.height - 2)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.fileTree.View(),
		// Divider could go here if needed
		contentStyle.Render(content),
	)
}

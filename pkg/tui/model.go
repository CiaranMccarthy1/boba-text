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
	showTree bool
	width    int
	height   int
}

func InitialModel(startPath string) Model {
	return Model{
		fileTree: NewFileTree(startPath),
		editor:   NewEditor(),
		agent:    NewAgent(),
		focus:    FocusFileTree,
		showTree: true,
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
			// If Editor is focused and in Insert mode, pass the key through (don't switch focus)
			if m.focus == FocusEditor && m.editor.mode == ModeInsert {
				break
			}
			// Cycle focus: FileTree -> Editor -> Agent -> FileTree
			m.focus = (m.focus + 1) % 3
			return m, nil

		case "ctrl+b":
			m.showTree = !m.showTree
			m.resizePanes()

		case "ctrl+e":
			if m.focus == FocusFileTree {
				// If already focused, toggle back to editor
				m.focus = FocusEditor
			} else {
				// Otherwise focus tree and ensure visible
				m.focus = FocusFileTree
				if !m.showTree {
					m.showTree = true
					m.resizePanes()
				}
			}

		case "ctrl+a":
			if m.focus == FocusAgent {
				m.focus = FocusEditor
			} else {
				m.focus = FocusAgent
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizePanes()

	case OpenFileMsg:
		content, err := os.ReadFile(msg.Path)
		if err == nil {
			m.editor.SetContent(string(content), msg.Path)
			m.focus = FocusEditor
		}
	}

	// Update components based on focus
	switch m.focus {
	case FocusFileTree:
		if m.showTree {
			m.fileTree, cmd = m.fileTree.Update(msg)
			cmds = append(cmds, cmd)
		}
	case FocusEditor:
		m.editor, cmd = m.editor.Update(msg)
		cmds = append(cmds, cmd)
	case FocusAgent:
		m.agent, cmd = m.agent.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) resizePanes() {
	treeWidth := 0
	if m.showTree {
		treeWidth = 30
		m.fileTree.width = treeWidth
		m.fileTree.height = m.height - 2
	}

	// Content width depends on tree visibility
	// If tree is shown: width - treeWidth - padding (approx 4 chars total borders/gap)
	// If tree hidden: width - padding

	padding := 4
	if !m.showTree {
		padding = 2
	}

	contentWidth := m.width - treeWidth - padding
	if contentWidth < 10 {
		contentWidth = 10 // Minimum safety
	}
	contentHeight := m.height - 2

	m.editor.width = contentWidth
	m.editor.height = contentHeight
	m.editor.textarea.SetWidth(contentWidth - 2)
	m.editor.textarea.SetHeight(contentHeight - 2)

	m.agent.SetSize(contentWidth, contentHeight)
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

	// Content Pane
	contentWidth := m.width - 34 // Default if tree shown
	if !m.showTree {
		contentWidth = m.width - 2
	}

	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(contentWidth).
		Height(m.height - 2)

	if m.showTree {
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.fileTree.View(),
			contentStyle.Render(content),
		)
	}

	// Just render content if tree hidden
	return contentStyle.Render(content)
}

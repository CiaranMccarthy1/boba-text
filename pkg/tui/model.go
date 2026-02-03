package tui

import (
	"os"

	"github.com/CiaranMccarthy1/boba-text/pkg/config"
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
	keys     config.Keys
}

// InitialModel creates the initial application model with the given configuration.
func InitialModel(startPath string, cfg config.Config) Model {
	InitStyles(cfg.Colors)
	return Model{
		fileTree: NewFileTree(startPath),
		editor:   NewEditor(cfg.Commands),
		agent:    NewAgent(cfg.AI),
		focus:    FocusFileTree,
		showTree: true,
		keys:     cfg.Keys,
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
		case "ctrl+c", m.keys.Quit:
			return m, tea.Quit
		case m.keys.CycleFocus:
			if m.focus == FocusEditor && m.editor.mode == ModeInsert {
				break
			}
			m.focus = (m.focus + 1) % 3
			return m, nil

		case m.keys.ToggleTree:
			m.showTree = !m.showTree
			m.resizePanes()
			return m, nil

		case m.keys.FocusTree:
			if m.focus == FocusFileTree {
				m.focus = FocusEditor
			} else {
				m.focus = FocusFileTree
				if !m.showTree {
					m.showTree = true
					m.resizePanes()
				}
			}
			return m, nil

		case m.keys.FocusAgent:
			if m.focus == FocusAgent {
				m.focus = FocusEditor
			} else {
				m.focus = FocusAgent
			}
			return m, nil
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

	padding := 4
	if !m.showTree {
		padding = 2
	}

	contentWidth := m.width - treeWidth - padding
	if contentWidth < 10 {
		contentWidth = 10
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

	var borderColor lipgloss.Color
	switch m.focus {
	case FocusAgent:
		borderColor = ColorSuccess
	case FocusEditor:
		borderColor = ColorPrimary
	default:
		borderColor = ColorSubText
	}

	contentWidth := m.width - 34
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

	return contentStyle.Render(content)
}

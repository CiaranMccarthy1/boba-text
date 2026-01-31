package tui

import (
	"fmt"
	"strings"

	"github.com/CiaranMccarthy1/boba-text/pkg/config"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AgentModel struct {
	viewport    viewport.Model
	textarea    textarea.Model
	messages    []string
	senderStyle lipgloss.Style
	aiStyle     lipgloss.Style
	width       int
	height      int
	config      config.AI
}

func NewAgent(aiConfig config.AI) AgentModel {
	ta := textarea.New()
	ta.Placeholder = "Ask the AI agent..."
	ta.Focus()
	ta.CharLimit = 0
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	vp := viewport.New(0, 0)
	vp.SetContent("Welcome to the AI Agent Tab!\nAsk me anything about your code.\n")

	return AgentModel{
		textarea:    ta,
		viewport:    vp,
		messages:    []string{},
		senderStyle: lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true),
		aiStyle:     lipgloss.NewStyle().Foreground(ColorSuccess),
		config:      aiConfig,
	}
}

func (m AgentModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m AgentModel) Update(msg tea.Msg) (AgentModel, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textarea.Value() == "" {
				break
			}
			userMsg := m.senderStyle.Render("You: ") + m.textarea.Value()
			m.messages = append(m.messages, userMsg)

			// Simulate AI response
			aiName := m.config.Name
			if aiName == "" {
				aiName = "Agent"
			}
			aiMsg := m.aiStyle.Render(aiName+": ") + "I received: " + m.textarea.Value()
			m.messages = append(m.messages, aiMsg)

			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m AgentModel) View() string {
	return StyleAgent.Render(
		fmt.Sprintf(
			"%s\n\n%s",
			m.viewport.View(),
			m.textarea.View(),
		),
	)
}

func (m *AgentModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.textarea.SetWidth(w)
	m.viewport.Width = w
	m.viewport.Height = h - m.textarea.Height() - 2 // Subtract textarea and padding
}

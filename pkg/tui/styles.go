package tui

import "github.com/charmbracelet/lipgloss"

// Color Palette
var (
	ColorText      = lipgloss.Color("#FAFAFA")
	ColorSubText   = lipgloss.Color("#7D7D7D")
	ColorPrimary   = lipgloss.Color("#F25D94") // Neon Pink
	ColorSecondary = lipgloss.Color("#A550DF") // Purple
	ColorAccent    = lipgloss.Color("#61AFEF") // Blue/Cyan
	ColorSuccess   = lipgloss.Color("#98C379") // Green
	ColorWarning   = lipgloss.Color("#E5C07B") // Yellow
	ColorError     = lipgloss.Color("#E06C75") // Red
	ColorDark      = lipgloss.Color("#1E1E1E") // Dark Background (optional usage)
)

// UI Styles
var (
	// Text Styles
	StyleNormal = lipgloss.NewStyle().Foreground(ColorText)
	StyleDim    = lipgloss.NewStyle().Foreground(ColorSubText)
	StyleBold   = lipgloss.NewStyle().Bold(true).Foreground(ColorText)

	// Selection / Focus Styles
	StyleSelected = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	StyleActiveBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary)

	StyleInactiveBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorSubText)

	// Panes
	StyleFileTree = lipgloss.NewStyle().
			Padding(1, 2) // Clean padding, no border by default? Or subtle right border.

	StyleEditor = lipgloss.NewStyle().
			Padding(0, 1)

	StyleAgent = lipgloss.NewStyle().
			Padding(0, 1)
)

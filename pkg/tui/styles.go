package tui

import (
	"github.com/CiaranMccarthy1/boba-text/pkg/config"
	"github.com/charmbracelet/lipgloss"
)

var (
	ColorText      lipgloss.Color
	ColorSubText   lipgloss.Color
	ColorPrimary   lipgloss.Color
	ColorSecondary lipgloss.Color
	ColorAccent    lipgloss.Color
	ColorSuccess   lipgloss.Color
	ColorWarning   lipgloss.Color
	ColorError     lipgloss.Color
	ColorDark      lipgloss.Color

	StyleNormal         lipgloss.Style
	StyleDim            lipgloss.Style
	StyleBold           lipgloss.Style
	StyleSelected       lipgloss.Style
	StyleActiveBorder   lipgloss.Style
	StyleInactiveBorder lipgloss.Style
	StyleFileTree       lipgloss.Style
	StyleEditor         lipgloss.Style
	StyleAgent          lipgloss.Style

	// New styles for Neovim features
	StyleTabActive   lipgloss.Style
	StyleTabInactive lipgloss.Style
	StyleTabBar      lipgloss.Style
	StyleStatusLine  lipgloss.Style
	StyleModeNormal  lipgloss.Style
	StyleModeInsert  lipgloss.Style
	StyleModeVisual  lipgloss.Style
	StyleModeCommand lipgloss.Style
	StyleSearch      lipgloss.Style
	StyleWelcome     lipgloss.Style
	StyleWelcomeDim  lipgloss.Style
	StyleLineNumber  lipgloss.Style
	StyleCursorLine  lipgloss.Style
	StyleFileIcon    lipgloss.Style
	StyleDirIcon     lipgloss.Style
	StyleModified    lipgloss.Style
)

// InitStyles initializes all global styles with the given color configuration.
func InitStyles(c config.Colors) {
	ColorText = lipgloss.Color(c.Text)
	ColorSubText = lipgloss.Color(c.SubText)
	ColorPrimary = lipgloss.Color(c.Primary)
	ColorSecondary = lipgloss.Color(c.Secondary)
	ColorAccent = lipgloss.Color(c.Accent)
	ColorSuccess = lipgloss.Color(c.Success)
	ColorWarning = lipgloss.Color(c.Warning)
	ColorError = lipgloss.Color(c.Error)
	ColorDark = lipgloss.Color(c.Dark)

	StyleNormal = lipgloss.NewStyle().Foreground(ColorText)
	StyleDim = lipgloss.NewStyle().Foreground(ColorSubText)
	StyleBold = lipgloss.NewStyle().Bold(true).Foreground(ColorText)

	StyleSelected = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	StyleActiveBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary)

	StyleInactiveBorder = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSubText)

	StyleFileTree = lipgloss.NewStyle().
		Padding(1, 2)

	StyleEditor = lipgloss.NewStyle().
		Padding(0, 1)

	StyleAgent = lipgloss.NewStyle().
		Padding(0, 1)

	// Tab bar styles
	StyleTabActive = lipgloss.NewStyle().
		Foreground(ColorText).
		Background(ColorPrimary).
		Bold(true).
		Padding(0, 2)

	StyleTabInactive = lipgloss.NewStyle().
		Foreground(ColorSubText).
		Background(lipgloss.Color("#2A2A2A")).
		Padding(0, 2)

	StyleTabBar = lipgloss.NewStyle().
		Background(lipgloss.Color("#1A1A1A")).
		Padding(0, 0)

	// Status line
	StyleStatusLine = lipgloss.NewStyle().
		Foreground(ColorText).
		Background(lipgloss.Color("#2A2A2A")).
		Padding(0, 1)

	// Mode indicators (Neovim-style lualine colors)
	StyleModeNormal = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1E1E1E")).
		Background(ColorAccent).
		Bold(true).
		Padding(0, 1)

	StyleModeInsert = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1E1E1E")).
		Background(ColorSuccess).
		Bold(true).
		Padding(0, 1)

	StyleModeVisual = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1E1E1E")).
		Background(ColorSecondary).
		Bold(true).
		Padding(0, 1)

	StyleModeCommand = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1E1E1E")).
		Background(ColorWarning).
		Bold(true).
		Padding(0, 1)

	// Search highlight
	StyleSearch = lipgloss.NewStyle().
		Background(ColorWarning).
		Foreground(lipgloss.Color("#1E1E1E")).
		Bold(true)

	// Welcome screen
	StyleWelcome = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	StyleWelcomeDim = lipgloss.NewStyle().
		Foreground(ColorSubText).
		Italic(true)

	// Line numbers
	StyleLineNumber = lipgloss.NewStyle().
		Foreground(ColorSubText).
		Width(4).
		Align(lipgloss.Right)

	// Cursor line highlight
	StyleCursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("#2A2A2A"))

	// File tree icons
	StyleFileIcon = lipgloss.NewStyle().
		Foreground(ColorAccent)

	StyleDirIcon = lipgloss.NewStyle().
		Foreground(ColorWarning)

	// Modified indicator
	StyleModified = lipgloss.NewStyle().
		Foreground(ColorWarning).
		Bold(true)
}

package tui

import (
	"github.com/CiaranMccarthy1/boba-text/pkg/config"
	"github.com/charmbracelet/lipgloss"
)

// Global Styles Struct (populated on Init)
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
)

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
}

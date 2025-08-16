package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#FF9500") // AWS Orange
	secondaryColor = lipgloss.Color("#232F3E") // AWS Dark Blue
	accentColor    = lipgloss.Color("#146EB4") // AWS Light Blue
	successColor   = lipgloss.Color("#00C853")
	warningColor   = lipgloss.Color("#FFA000")
	errorColor     = lipgloss.Color("#D32F2F")
	mutedColor     = lipgloss.Color("#666666")
	highlightColor = lipgloss.Color("#FFE082")

	// Base styles
	BaseStyle = lipgloss.NewStyle()

	// Title and header styles
	TitleStyle = BaseStyle.
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	HeaderStyle = BaseStyle.
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(secondaryColor).
			Padding(0, 1)

	// List styles
	ListHeader = BaseStyle.
			Bold(true).
			Foreground(accentColor).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(mutedColor)

	ListItem = BaseStyle.
			PaddingLeft(1)

	SelectedItem = BaseStyle.
			Foreground(lipgloss.Color("#000000")).
			Background(highlightColor).
			PaddingLeft(1). // Same padding as ListItem for alignment
			Bold(true)

	// Status bar styles
	StatusBar = BaseStyle.
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(secondaryColor)

	StatusKey = BaseStyle.
			Bold(true).
			Foreground(primaryColor).
			Background(secondaryColor).
			Padding(0, 1)

	StatusValue = BaseStyle.
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(secondaryColor).
			Padding(0, 1)

	// Help styles
	HelpStyle = BaseStyle.
			Foreground(mutedColor)

	HelpKey = BaseStyle.
		Bold(true).
		Foreground(accentColor)

	HelpDesc = BaseStyle.
			Foreground(mutedColor)

	// Detail view styles
	DetailTitle = BaseStyle.
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(mutedColor)

	DetailLabel = BaseStyle.
			Bold(true).
			Foreground(accentColor).
			Width(20)

	DetailValue = BaseStyle.
			Foreground(lipgloss.Color("#FFFFFF"))

	// Code/JSON styles
	CodeBlock = BaseStyle.
			Background(lipgloss.Color("#1E1E1E")).
			Foreground(lipgloss.Color("#D4D4D4")).
			Padding(1).
			MarginTop(1).
			MarginBottom(1)

	// Search styles
	SearchPrompt = BaseStyle.
			Foreground(primaryColor).
			Bold(true)

	SearchInput = BaseStyle.
			Foreground(lipgloss.Color("#FFFFFF"))

	// Tab styles
	ActiveTab = BaseStyle.
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(primaryColor).
			Padding(0, 2)

	InactiveTab = BaseStyle.
			Foreground(mutedColor).
			Background(lipgloss.Color("#333333")).
			Padding(0, 2)

	// Error styles
	ErrorStyle = BaseStyle.
			Foreground(errorColor).
			Bold(true)

	// Loading styles
	LoadingStyle = BaseStyle.
			Foreground(accentColor).
			Bold(true)

	// Container with border (like k9s) - base style without width
	MainContainer = BaseStyle.
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1)

	// Header styles for k9s-like display
	ASCIIArtStyle = BaseStyle.
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Right)

	HeaderKey = BaseStyle.
			Foreground(mutedColor)

	HeaderValue = BaseStyle.
			Foreground(accentColor).
			Bold(true)
)

// Helper functions
func GetMainContainer(width, height int) lipgloss.Style {
	return MainContainer.
		Width(width - 2). // Account for terminal margins
		Height(height)
}

func RenderStatusBar(profile, region string, itemCount int, width int) string {
	left := lipgloss.JoinHorizontal(
		lipgloss.Top,
		StatusKey.Render("Profile:"),
		StatusValue.Render(profile),
		StatusKey.Render("Region:"),
		StatusValue.Render(region),
	)

	right := StatusValue.Render(fmt.Sprintf("%d items", itemCount))

	if width <= 0 {
		width = 80 // Fallback width
	}

	spaces := width - lipgloss.Width(left) - lipgloss.Width(right)
	if spaces < 0 {
		spaces = 0
	}

	return StatusBar.Width(width).Render(
		left + strings.Repeat(" ", spaces) + right,
	)
}

func RenderHelp() string {
	help := []string{
		HelpKey.Render("j/k") + " " + HelpDesc.Render("up/down"),
		HelpKey.Render("Enter") + " " + HelpDesc.Render("view"),
		HelpKey.Render("/") + " " + HelpDesc.Render("search"),
		HelpKey.Render("r") + " " + HelpDesc.Render("refresh"),
		HelpKey.Render("q") + " " + HelpDesc.Render("quit"),
		HelpKey.Render("?") + " " + HelpDesc.Render("help"),
	}
	return HelpStyle.Render(strings.Join(help, " | "))
}

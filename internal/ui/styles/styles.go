package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/johnoct/a3s/internal/aws/identity"
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

	SearchMatch = BaseStyle.
			Background(lipgloss.Color("#FFE082")).
			Foreground(lipgloss.Color("#000000"))

	SearchCurrentMatch = BaseStyle.
				Background(lipgloss.Color("#FF5722")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)

	SearchInfo = BaseStyle.
			Foreground(accentColor).
			Bold(true)

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
			Bold(true)

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

// RenderHeader renders the application header with AWS identity information and ASCII art
func RenderHeader(profile, region string, identity *identity.Identity, terminalWidth int) string {
	var header strings.Builder

	// Simple and readable a3s logo
	asciiArt := []string{
		"   __ _  _____  ___ ",
		"  / _` ||___ / / __|",
		" | (_| | |_ \\ \\__ \\",
		"  \\__,_||___/ |___/",
	}

	// Format AWS identity information (left side, like k9s)
	// Add padding to align with the main content border
	leftPadding := " " // Space to align with bordered content
	var infoLines []string
	if identity != nil {
		infoLines = []string{
			fmt.Sprintf("%s%s %s", leftPadding, HeaderKey.Render("Account:"), HeaderValue.Render(identity.Account)),
			fmt.Sprintf("%s%s %s", leftPadding, HeaderKey.Render("User:"), HeaderValue.Render(identity.DisplayName)),
			fmt.Sprintf("%s%s %s", leftPadding, HeaderKey.Render("Region:"), HeaderValue.Render(region)),
		}
		// Add profile if different from user
		if profile != "" && profile != "default" {
			infoLines = append(infoLines, fmt.Sprintf("%s%s %s", leftPadding, HeaderKey.Render("Profile:"), HeaderValue.Render(profile)))
		}
	} else {
		infoLines = []string{
			fmt.Sprintf("%s%s %s", leftPadding, HeaderKey.Render("Profile:"), HeaderValue.Render(profile)),
			fmt.Sprintf("%s%s %s", leftPadding, HeaderKey.Render("Region:"), HeaderValue.Render(region)),
		}
	}

	// Calculate dimensions for proper k9s-style layout
	asciiWidth := 19  // Actual width of the ASCII art
	rightPadding := 4 // Padding from right edge (like k9s)
	minSpacing := 12  // Minimum spacing for better separation

	// Find the maximum width of left content for consistent spacing
	maxLeftWidth := 0
	for _, line := range infoLines {
		if w := lipgloss.Width(line); w > maxLeftWidth {
			maxLeftWidth = w
		}
	}

	// Calculate available space (account for terminal width and right padding)
	availableWidth := terminalWidth - rightPadding
	totalRequiredWidth := maxLeftWidth + minSpacing + asciiWidth

	// Calculate spacing - prioritize right-alignment like k9s
	var spacing int
	if totalRequiredWidth <= availableWidth {
		// We have enough space - calculate spacing to right-align the ASCII art
		spacing = availableWidth - maxLeftWidth - asciiWidth
		// Ensure minimum spacing is maintained
		if spacing < minSpacing {
			spacing = minSpacing
		}
	} else {
		// Terminal too narrow - use minimum spacing and let ASCII art overflow gracefully
		spacing = minSpacing
	}

	// Combine info (left) and ASCII art (right) - k9s-style layout
	maxLines := len(asciiArt)
	if len(infoLines) > maxLines {
		maxLines = len(infoLines)
	}

	for i := 0; i < maxLines; i++ {
		var line strings.Builder

		// Left side - AWS info
		if i < len(infoLines) {
			line.WriteString(infoLines[i])
			// Pad to consistent width for alignment
			currentWidth := lipgloss.Width(infoLines[i])
			if padding := maxLeftWidth - currentWidth; padding > 0 {
				line.WriteString(strings.Repeat(" ", padding))
			}
		} else {
			// Empty left side - pad to max width
			line.WriteString(strings.Repeat(" ", maxLeftWidth))
		}

		// Add calculated spacing to position ASCII art properly
		line.WriteString(strings.Repeat(" ", spacing))

		// Right side - ASCII art with consistent right alignment
		if i < len(asciiArt) {
			// Apply styling and ensure consistent positioning
			artLine := ASCIIArtStyle.Render(asciiArt[i])
			line.WriteString(artLine)
		}

		header.WriteString(line.String())
		header.WriteString("\n")
	}

	return strings.TrimRight(header.String(), "\n")
}

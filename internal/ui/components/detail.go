package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/johnoct/a3s/internal/aws/iam"
	"github.com/johnoct/a3s/internal/ui/styles"
)

type DetailModel struct {
	role       *iam.Role
	profile    string
	region     string
	width      int
	height     int
	scrollY    int
	activeTab  int
	tabs       []string
}

func NewDetailModel(role *iam.Role, profile, region string) *DetailModel {
	return &DetailModel{
		role:    role,
		profile: profile,
		region:  region,
		tabs:    []string{"Overview", "Trust Policy", "Policies", "Tags"},
	}
}

func (m *DetailModel) Init() tea.Cmd {
	return nil
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "l":
			m.activeTab = (m.activeTab + 1) % len(m.tabs)
			m.scrollY = 0
		case "shift+tab", "h":
			m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
			m.scrollY = 0
		case "j", "down":
			m.scrollY++
		case "k", "up":
			if m.scrollY > 0 {
				m.scrollY--
			}
		case "g":
			m.scrollY = 0
		}
	}

	return m, nil
}

func (m *DetailModel) View() string {
	var content strings.Builder
	var fullView strings.Builder

	// Title (outside the border)
	title := fmt.Sprintf("🔍 Role: %s", m.role.Name)
	fullView.WriteString(styles.TitleStyle.Render(title))
	fullView.WriteString("\n")

	// Tabs (outside the border, just above it)
	var tabs []string
	for i, tab := range m.tabs {
		if i == m.activeTab {
			tabs = append(tabs, styles.ActiveTab.Render(tab))
		} else {
			tabs = append(tabs, styles.InactiveTab.Render(tab))
		}
	}
	fullView.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	fullView.WriteString("\n")

	// Content based on active tab (inside the border)
	tabContent := ""
	switch m.activeTab {
	case 0: // Overview
		tabContent = m.renderOverview()
	case 1: // Trust Policy
		tabContent = m.renderTrustPolicy()
	case 2: // Policies
		tabContent = m.renderPolicies()
	case 3: // Tags
		tabContent = m.renderTags()
	}

	// Apply scrolling
	lines := strings.Split(tabContent, "\n")
	
	// Calculate visible height accounting for border
	borderHeight := 4
	titleHeight := 2
	tabHeight := 2
	statusHeight := 2
	helpHeight := 1
	
	visibleHeight := m.height - borderHeight - titleHeight - tabHeight - statusHeight - helpHeight - 2
	if visibleHeight < 5 {
		visibleHeight = 5
	}

	endIdx := m.scrollY + visibleHeight
	if endIdx > len(lines) {
		endIdx = len(lines)
	}

	for i := m.scrollY; i < endIdx; i++ {
		if i < len(lines) {
			content.WriteString(lines[i])
			content.WriteString("\n")
		}
	}

	// Calculate available width
	availableWidth := m.width - 6 // Account for border and padding
	if availableWidth < 80 {
		availableWidth = 80
	}
	
	// Fill empty space inside the border with full width lines
	for i := endIdx - m.scrollY; i < visibleHeight; i++ {
		content.WriteString(strings.Repeat(" ", availableWidth))
		content.WriteString("\n")
	}

	// Calculate container height
	containerHeight := visibleHeight
	
	// Apply the border container to the content with dynamic sizing
	borderedContent := styles.GetMainContainer(m.width, containerHeight).Render(strings.TrimRight(content.String(), "\n"))
	fullView.WriteString(borderedContent)
	fullView.WriteString("\n")

	// Status bar (outside the border) with full width
	fullView.WriteString(styles.RenderStatusBar(m.profile, m.region, 1, m.width))
	fullView.WriteString("\n")

	// Help (outside the border)
	help := []string{
		styles.HelpKey.Render("Tab/l") + " " + styles.HelpDesc.Render("next tab"),
		styles.HelpKey.Render("Shift+Tab/h") + " " + styles.HelpDesc.Render("prev tab"),
		styles.HelpKey.Render("j/k") + " " + styles.HelpDesc.Render("scroll"),
		styles.HelpKey.Render("Esc") + " " + styles.HelpDesc.Render("back"),
	}
	fullView.WriteString(styles.HelpStyle.Render(strings.Join(help, " | ")))

	return fullView.String()
}

func (m *DetailModel) renderOverview() string {
	var s strings.Builder

	s.WriteString(styles.DetailTitle.Render("Role Information"))
	s.WriteString("\n\n")

	fields := []struct {
		label string
		value string
	}{
		{"ARN", m.role.ARN},
		{"Role ID", m.role.RoleID},
		{"Path", m.role.Path},
		{"Created", m.role.CreateDate.Format("2006-01-02 15:04:05")},
		{"Description", m.role.Description},
		{"Max Session", fmt.Sprintf("%d seconds", m.role.MaxSessionDuration)},
	}

	if m.role.LastUsed != nil {
		fields = append(fields, struct {
			label string
			value string
		}{"Last Used", m.role.LastUsed.Format("2006-01-02 15:04:05")})
	}

	for _, field := range fields {
		s.WriteString(styles.DetailLabel.Render(field.label + ":"))
		s.WriteString(" ")
		s.WriteString(styles.DetailValue.Render(field.value))
		s.WriteString("\n")
	}

	return s.String()
}

func (m *DetailModel) renderTrustPolicy() string {
	var s strings.Builder

	s.WriteString(styles.DetailTitle.Render("Trust Relationships"))
	s.WriteString("\n")
	s.WriteString(styles.CodeBlock.Render(m.role.TrustPolicy))

	return s.String()
}

func (m *DetailModel) renderPolicies() string {
	var s strings.Builder

	s.WriteString(styles.DetailTitle.Render("Attached Policies"))
	s.WriteString("\n\n")

	if len(m.role.ManagedPolicies) > 0 {
		s.WriteString(styles.DetailLabel.Render("Managed Policies:"))
		s.WriteString("\n")
		for _, policy := range m.role.ManagedPolicies {
			s.WriteString("  • " + policy)
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	if len(m.role.InlinePolicies) > 0 {
		s.WriteString(styles.DetailLabel.Render("Inline Policies:"))
		s.WriteString("\n")
		for _, policy := range m.role.InlinePolicies {
			s.WriteString("  • " + policy)
			s.WriteString("\n")
		}
	}

	if len(m.role.ManagedPolicies) == 0 && len(m.role.InlinePolicies) == 0 {
		s.WriteString(styles.HelpDesc.Render("No policies attached"))
	}

	return s.String()
}

func (m *DetailModel) renderTags() string {
	var s strings.Builder

	s.WriteString(styles.DetailTitle.Render("Tags"))
	s.WriteString("\n\n")

	if len(m.role.Tags) > 0 {
		for _, tag := range m.role.Tags {
			s.WriteString(styles.DetailLabel.Render(tag.Key + ":"))
			s.WriteString(" ")
			s.WriteString(styles.DetailValue.Render(tag.Value))
			s.WriteString("\n")
		}
	} else {
		s.WriteString(styles.HelpDesc.Render("No tags"))
	}

	return s.String()
}
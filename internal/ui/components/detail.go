package components

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/johnoct/a3s/internal/aws/iam"
	"github.com/johnoct/a3s/internal/aws/identity"
	"github.com/johnoct/a3s/internal/ui/styles"
)

type DetailModel struct {
	role           *iam.Role
	profile        string
	region         string
	identity       *identity.Identity
	roleService    *iam.RoleService
	width          int
	height         int
	scrollY        int
	activeTab      int
	tabs           []string
	viewState      viewState
	selectedPolicy int
	policyDocument string
	loadingPolicy  bool
}

type viewState int

const (
	viewNormal viewState = iota
	viewPolicyDocument
)

func (m *DetailModel) IsViewingPolicyDocument() bool {
	return m.viewState == viewPolicyDocument
}

func NewDetailModel(role *iam.Role, profile, region string, roleService *iam.RoleService) *DetailModel {
	return &DetailModel{
		role:        role,
		profile:     profile,
		region:      region,
		roleService: roleService,
		tabs:        []string{"Overview", "Trust Policy", "Policies", "Tags"},
		viewState:   viewNormal,
	}
}

// Message types for async policy loading
type policyDocumentLoadedMsg struct {
	document string
	err      error
}

func (m *DetailModel) Init() tea.Cmd {
	return nil
}

func (m *DetailModel) SetIdentity(id *identity.Identity) {
	m.identity = id
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case policyDocumentLoadedMsg:
		m.loadingPolicy = false
		if msg.err != nil {
			m.policyDocument = fmt.Sprintf("Error loading policy: %v", msg.err)
		} else {
			m.policyDocument = msg.document
		}
		m.viewState = viewPolicyDocument
		m.scrollY = 0
		return m, nil

	case tea.KeyMsg:
		switch m.viewState {
		case viewPolicyDocument:
			return m.updatePolicyDocumentView(msg)
		case viewNormal:
			return m.updateNormalView(msg)
		}
	}

	return m, nil
}

func (m *DetailModel) updatePolicyDocumentView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewState = viewNormal
		m.scrollY = 0
		return m, nil
	case "j", "down":
		m.scrollY++
	case "k", "up":
		if m.scrollY > 0 {
			m.scrollY--
		}
	case "g":
		m.scrollY = 0
	case "G":
		// Scroll to bottom
		lines := strings.Split(m.policyDocument, "\n")
		visibleHeight := m.calculateVisibleHeight()
		m.scrollY = max(0, len(lines)-visibleHeight)
	}
	return m, nil
}

func (m *DetailModel) updateNormalView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "l":
		m.activeTab = (m.activeTab + 1) % len(m.tabs)
		m.scrollY = 0
		m.selectedPolicy = 0
	case "shift+tab", "h":
		m.activeTab = (m.activeTab - 1 + len(m.tabs)) % len(m.tabs)
		m.scrollY = 0
		m.selectedPolicy = 0
	case "j", "down":
		if m.activeTab == 2 { // Policies tab
			totalPolicies := len(m.role.ManagedPolicies) + len(m.role.InlinePolicies)
			if totalPolicies > 0 {
				m.selectedPolicy = min(m.selectedPolicy+1, totalPolicies-1)
			}
		} else {
			m.scrollY++
		}
	case "k", "up":
		if m.activeTab == 2 { // Policies tab
			m.selectedPolicy = max(0, m.selectedPolicy-1)
		} else if m.scrollY > 0 {
			m.scrollY--
		}
	case "g":
		m.scrollY = 0
		m.selectedPolicy = 0
	case "G":
		if m.activeTab == 2 { // Policies tab
			totalPolicies := len(m.role.ManagedPolicies) + len(m.role.InlinePolicies)
			if totalPolicies > 0 {
				m.selectedPolicy = totalPolicies - 1
			}
		}
	case "enter":
		if m.activeTab == 2 { // Policies tab
			return m, m.loadSelectedPolicy()
		}
	}
	return m, nil
}

func (m *DetailModel) loadSelectedPolicy() tea.Cmd {
	if m.loadingPolicy {
		return nil
	}

	totalManagedPolicies := len(m.role.ManagedPolicies)

	if m.selectedPolicy < totalManagedPolicies {
		// It's a managed policy
		policy := m.role.ManagedPolicies[m.selectedPolicy]
		m.loadingPolicy = true
		return func() tea.Msg {
			doc, err := m.roleService.GetManagedPolicyDocument(context.Background(), policy.ARN)
			return policyDocumentLoadedMsg{document: doc, err: err}
		}
	} else {
		// It's an inline policy
		inlineIndex := m.selectedPolicy - totalManagedPolicies
		if inlineIndex < len(m.role.InlinePolicies) {
			policyName := m.role.InlinePolicies[inlineIndex]
			m.loadingPolicy = true
			return func() tea.Msg {
				doc, err := m.roleService.GetInlinePolicy(context.Background(), m.role.Name, policyName)
				return policyDocumentLoadedMsg{document: doc, err: err}
			}
		}
	}

	return nil
}

func (m *DetailModel) calculateVisibleHeight() int {
	borderHeight := 4
	titleHeight := 2
	tabHeight := 2
	statusHeight := 2
	helpHeight := 1
	return max(5, m.height-borderHeight-titleHeight-tabHeight-statusHeight-helpHeight-2)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *DetailModel) View() string {
	switch m.viewState {
	case viewPolicyDocument:
		return m.renderPolicyDocumentView()
	default:
		return m.renderNormalView()
	}
}

func (m *DetailModel) renderPolicyDocumentView() string {
	var content strings.Builder
	var fullView strings.Builder

	// Title for policy document view
	title := "📄 Policy Document"
	fullView.WriteString(styles.TitleStyle.Render(title))
	fullView.WriteString("\n\n")

	// Policy document content with scrolling
	lines := strings.Split(m.policyDocument, "\n")
	visibleHeight := m.calculateVisibleHeight()

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
	availableWidth := m.width - 6
	if availableWidth < 80 {
		availableWidth = 80
	}

	// Fill empty space
	for i := endIdx - m.scrollY; i < visibleHeight; i++ {
		content.WriteString(strings.Repeat(" ", availableWidth))
		content.WriteString("\n")
	}

	// Apply code block styling and border
	styledContent := styles.CodeBlock.Render(strings.TrimRight(content.String(), "\n"))
	borderedContent := styles.GetMainContainer(m.width, visibleHeight+4).Render(styledContent)
	fullView.WriteString(borderedContent)
	fullView.WriteString("\n")

	// Status bar
	fullView.WriteString(styles.RenderStatusBar(m.profile, m.region, 1, m.width))
	fullView.WriteString("\n")

	// Help for policy document view
	help := []string{
		styles.HelpKey.Render("j/k") + " " + styles.HelpDesc.Render("scroll"),
		styles.HelpKey.Render("g/G") + " " + styles.HelpDesc.Render("top/bottom"),
		styles.HelpKey.Render("Esc") + " " + styles.HelpDesc.Render("back to policies"),
	}
	fullView.WriteString(styles.HelpStyle.Render(strings.Join(help, " | ")))

	return fullView.String()
}

func (m *DetailModel) renderNormalView() string {
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
	visibleHeight := m.calculateVisibleHeight()

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

	// Apply the border container to the content with dynamic sizing
	borderedContent := styles.GetMainContainer(m.width, visibleHeight).Render(strings.TrimRight(content.String(), "\n"))
	fullView.WriteString(borderedContent)
	fullView.WriteString("\n")

	// Status bar (outside the border) with full width
	fullView.WriteString(styles.RenderStatusBar(m.profile, m.region, 1, m.width))
	fullView.WriteString("\n")

	// Help (outside the border)
	help := m.getHelpText()
	fullView.WriteString(styles.HelpStyle.Render(strings.Join(help, " | ")))

	return fullView.String()
}

func (m *DetailModel) getHelpText() []string {
	if m.activeTab == 2 { // Policies tab
		if len(m.role.ManagedPolicies) > 0 || len(m.role.InlinePolicies) > 0 {
			return []string{
				styles.HelpKey.Render("Tab/l") + " " + styles.HelpDesc.Render("next tab"),
				styles.HelpKey.Render("j/k") + " " + styles.HelpDesc.Render("navigate"),
				styles.HelpKey.Render("Enter") + " " + styles.HelpDesc.Render("view policy"),
				styles.HelpKey.Render("Esc") + " " + styles.HelpDesc.Render("back"),
			}
		}
	}

	return []string{
		styles.HelpKey.Render("Tab/l") + " " + styles.HelpDesc.Render("next tab"),
		styles.HelpKey.Render("Shift+Tab/h") + " " + styles.HelpDesc.Render("prev tab"),
		styles.HelpKey.Render("j/k") + " " + styles.HelpDesc.Render("scroll"),
		styles.HelpKey.Render("Esc") + " " + styles.HelpDesc.Render("back"),
	}
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

	if m.loadingPolicy {
		s.WriteString(styles.LoadingStyle.Render("Loading policy document..."))
		s.WriteString("\n")
		return s.String()
	}

	currentIndex := 0

	if len(m.role.ManagedPolicies) > 0 {
		s.WriteString(styles.DetailLabel.Render("Managed Policies:"))
		s.WriteString("\n")
		for _, policy := range m.role.ManagedPolicies {
			prefix := "  • "
			policyText := policy.Name

			if currentIndex == m.selectedPolicy {
				// Highlight selected policy
				s.WriteString(styles.SelectedItem.Render(prefix + policyText))
			} else {
				s.WriteString(styles.ListItem.Render(prefix + policyText))
			}
			s.WriteString("\n")
			currentIndex++
		}
		s.WriteString("\n")
	}

	if len(m.role.InlinePolicies) > 0 {
		s.WriteString(styles.DetailLabel.Render("Inline Policies:"))
		s.WriteString("\n")
		for _, policy := range m.role.InlinePolicies {
			prefix := "  • "
			policyText := policy

			if currentIndex == m.selectedPolicy {
				// Highlight selected policy
				s.WriteString(styles.SelectedItem.Render(prefix + policyText))
			} else {
				s.WriteString(styles.ListItem.Render(prefix + policyText))
			}
			s.WriteString("\n")
			currentIndex++
		}
	}

	if len(m.role.ManagedPolicies) == 0 && len(m.role.InlinePolicies) == 0 {
		s.WriteString(styles.HelpDesc.Render("No policies attached"))
	} else {
		s.WriteString("\n")
		s.WriteString(styles.HelpDesc.Render("Press Enter to view the selected policy document"))
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

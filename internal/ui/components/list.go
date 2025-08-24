package components

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/johnoct/a3s/internal/aws/iam"
	"github.com/johnoct/a3s/internal/aws/identity"
	"github.com/johnoct/a3s/internal/ui/styles"
)

type ListModel struct {
	roles         []iam.Role
	filteredRoles []iam.Role
	cursor        int
	searchMode    bool
	searchInput   textinput.Model
	width         int
	height        int
	profile       string
	region        string
	identity      *identity.Identity
	selectedRole  *iam.Role
	showDetail    bool
	detailView    *DetailModel
	roleService   *iam.RoleService
	loadingDetail bool
}

func NewListModel(roles []iam.Role, profile, region string) ListModel {
	return NewListModelWithSize(roles, profile, region, 80, 24)
}

func NewListModelWithSize(roles []iam.Role, profile, region string, width, height int) ListModel {
	ti := textinput.New()
	ti.Placeholder = "Search roles..."
	ti.CharLimit = 100

	m := ListModel{
		roles:         roles,
		filteredRoles: roles,
		searchInput:   ti,
		profile:       profile,
		region:        region,
		width:         width,
		height:        height,
	}

	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) SetIdentity(id *identity.Identity) {
	m.identity = id
}

func (m *ListModel) SetRoleService(rs *iam.RoleService) {
	m.roleService = rs
}

type roleDetailsLoadedMsg struct {
	role *iam.Role
}

func (m *ListModel) loadRoleDetails(roleName string) tea.Cmd {
	return func() tea.Msg {
		if m.roleService == nil {
			return nil
		}
		ctx := context.Background()
		role, err := m.roleService.GetRoleDetails(ctx, roleName)
		if err != nil {
			// For now, return nil on error
			return nil
		}
		return roleDetailsLoadedMsg{role: role}
	}
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle detail view updates
	if m.showDetail && m.detailView != nil {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			// Pass window size to detail view
			var detailModel tea.Model
			detailModel, cmd = m.detailView.Update(msg)
			m.detailView = detailModel.(*DetailModel)
			return m, cmd
		case tea.KeyMsg:
			// Only handle esc/q to close detail view if we're not viewing a policy document
			if msg.String() == "esc" && !m.detailView.IsViewingPolicyDocument() {
				m.showDetail = false
				m.detailView = nil
				m.loadingDetail = false
				return m, nil
			}
			if msg.String() == "q" {
				m.showDetail = false
				m.detailView = nil
				m.loadingDetail = false
				return m, nil
			}
		}

		var detailModel tea.Model
		detailModel, cmd = m.detailView.Update(msg)
		m.detailView = detailModel.(*DetailModel)
		return m, cmd
	}

	switch msg := msg.(type) {
	case roleDetailsLoadedMsg:
		m.loadingDetail = false
		if msg.role != nil {
			m.selectedRole = msg.role
			m.detailView = NewDetailModel(m.selectedRole, m.profile, m.region, m.roleService)
			// Set the window size and identity for detail view
			m.detailView.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			if m.identity != nil {
				m.detailView.SetIdentity(m.identity)
			}
			m.showDetail = true
			return m, m.detailView.Init()
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.searchMode {
			switch msg.String() {
			case "esc":
				m.searchMode = false
				m.searchInput.SetValue("")
				m.filteredRoles = m.roles
				m.cursor = 0
				return m, nil
			case "enter":
				m.searchMode = false
				m.filterRoles()
				return m, nil
			default:
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.filterRoles()
				return m, cmd
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(m.filteredRoles)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "g":
			m.cursor = 0
		case "G":
			if len(m.filteredRoles) > 0 {
				m.cursor = len(m.filteredRoles) - 1
			}
		case "/":
			m.searchMode = true
			m.searchInput.Focus()
			return m, textinput.Blink
		case "enter":
			if len(m.filteredRoles) > 0 && m.cursor < len(m.filteredRoles) && !m.loadingDetail {
				m.loadingDetail = true
				roleName := m.filteredRoles[m.cursor].Name
				return m, m.loadRoleDetails(roleName)
			}
		case "r":
			// TODO: Implement refresh
			return m, nil
		}
	}

	return m, cmd
}

func (m *ListModel) filterRoles() {
	searchTerm := strings.ToLower(m.searchInput.Value())
	if searchTerm == "" {
		m.filteredRoles = m.roles
		return
	}

	filtered := []iam.Role{}
	for _, role := range m.roles {
		if strings.Contains(strings.ToLower(role.Name), searchTerm) ||
			strings.Contains(strings.ToLower(role.Description), searchTerm) {
			filtered = append(filtered, role)
		}
	}
	m.filteredRoles = filtered
	if m.cursor >= len(m.filteredRoles) {
		m.cursor = 0
	}
}

func (m ListModel) View() string {
	if m.showDetail && m.detailView != nil {
		return m.detailView.View()
	}

	if m.loadingDetail {
		return "\n  Loading role details... ⚡\n"
	}

	var content strings.Builder
	var fullView strings.Builder

	// Create header with ASCII art and AWS info (no top margin needed)
	fullView.WriteString(styles.RenderHeader(m.profile, m.region, m.identity, m.width))
	fullView.WriteString("\n\n") // Extra line for spacing, will be occupied by title/tabs in detail view

	// Search bar - always reserve space to prevent layout shifts
	searchPrompt := styles.SearchPrompt.Render(" Search: ")
	if m.searchMode {
		fullView.WriteString(searchPrompt)
		fullView.WriteString(m.searchInput.View())
	} else {
		// Render invisible search bar to maintain layout consistency
		// Reserve space for both the prompt and a reasonable input width
		promptWidth := lipgloss.Width(searchPrompt)
		minInputWidth := 20 // Reserve space for input field
		invisibleSpace := strings.Repeat(" ", promptWidth+minInputWidth)
		fullView.WriteString(invisibleSpace)
	}
	fullView.WriteString("\n")

	// Calculate column widths based on terminal width
	// Account for smaller container (m.width-2) and border padding
	availableWidth := m.width - 8 // Account for border, padding, and left margin
	if availableWidth < 80 {
		availableWidth = 80
	}

	// Distribute width with fixed spacing between columns
	// Using fixed column widths with proper spacing
	roleWidth := 40
	createdWidth := 12
	lastUsedWidth := 12
	// Calculate remaining space for description
	descWidth := availableWidth - roleWidth - createdWidth - lastUsedWidth - 3 // 3 spaces between columns
	if descWidth < 20 {
		descWidth = 20
	}

	// Column headers (inside the border)
	headers := fmt.Sprintf("%-*s %-*s %-*s %s",
		roleWidth, "Role Name",
		createdWidth, "Created",
		lastUsedWidth, "Last Used",
		"Description",
	)
	content.WriteString(styles.ListHeader.Width(availableWidth).Render(headers))
	content.WriteString("\n")

	// Calculate visible height accounting for border and header
	borderHeight := 2 // Reduced from 4
	headerHeight := 9 // ASCII art (6 lines) + top margin (1) + spacing (2)
	searchHeight := 1 // Always reserve space for search bar to prevent layout shifts
	statusHeight := 1 // Reduced from 2
	helpHeight := 1

	visibleHeight := m.height - borderHeight - headerHeight - searchHeight - statusHeight - helpHeight
	if visibleHeight < 5 {
		visibleHeight = 5
	}

	startIdx := 0
	if m.cursor >= visibleHeight {
		startIdx = m.cursor - visibleHeight + 1
	}
	endIdx := startIdx + visibleHeight
	if endIdx > len(m.filteredRoles) {
		endIdx = len(m.filteredRoles)
	}

	// Role list (inside the border)
	for i := startIdx; i < endIdx; i++ {
		role := m.filteredRoles[i]

		created := role.CreateDate.Format("2006-01-02")
		lastUsed := "Never"
		if role.LastUsed != nil {
			lastUsed = role.LastUsed.Format("2006-01-02")
		}

		// Truncate fields to exact column widths
		roleName := truncate(role.Name, roleWidth-1) // -1 for spacing
		createdStr := truncate(created, createdWidth-1)
		lastUsedStr := truncate(lastUsed, lastUsedWidth-1)
		description := truncate(role.Description, descWidth)

		// Build the line with exact spacing
		line := fmt.Sprintf("%-*s %-*s %-*s %s",
			roleWidth, roleName,
			createdWidth, createdStr,
			lastUsedWidth, lastUsedStr,
			description,
		)

		// Ensure the entire line doesn't exceed available width
		line = truncate(line, availableWidth)

		if i == m.cursor {
			// Apply selection without padding to maintain alignment
			content.WriteString(styles.SelectedItem.Render(line))
		} else {
			// Regular items with padding
			content.WriteString(styles.ListItem.Render(line))
		}
		content.WriteString("\n")
	}

	// Fill empty space inside the border with full width lines
	for i := endIdx - startIdx; i < visibleHeight; i++ {
		content.WriteString(strings.Repeat(" ", availableWidth))
		content.WriteString("\n")
	}

	// Calculate container height
	containerHeight := visibleHeight + 2 // Content + header line

	// Apply the border container to the content with dynamic sizing
	borderedContent := styles.GetMainContainer(m.width, containerHeight).Render(strings.TrimRight(content.String(), "\n"))
	fullView.WriteString(borderedContent)
	fullView.WriteString("\n")

	// Help line (outside the border)
	fullView.WriteString(styles.RenderHelp())

	return fullView.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

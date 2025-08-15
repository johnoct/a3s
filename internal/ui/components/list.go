package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johnoct/a3s/internal/aws/iam"
	"github.com/johnoct/a3s/internal/ui/styles"
)

type ListModel struct {
	roles          []iam.Role
	filteredRoles  []iam.Role
	cursor         int
	searchMode     bool
	searchInput    textinput.Model
	width          int
	height         int
	profile        string
	region         string
	selectedRole   *iam.Role
	showDetail     bool
	detailView     *DetailModel
	err            error
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
			if msg.String() == "esc" || msg.String() == "q" {
				m.showDetail = false
				m.detailView = nil
				return m, nil
			}
		}
		
		var detailModel tea.Model
		detailModel, cmd = m.detailView.Update(msg)
		m.detailView = detailModel.(*DetailModel)
		return m, cmd
	}

	switch msg := msg.(type) {
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
			if len(m.filteredRoles) > 0 && m.cursor < len(m.filteredRoles) {
				m.selectedRole = &m.filteredRoles[m.cursor]
				m.detailView = NewDetailModel(m.selectedRole, m.profile, m.region)
				// Set the window size for detail view
				m.detailView.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				m.showDetail = true
				return m, m.detailView.Init()
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

	var content strings.Builder
	var fullView strings.Builder

	// Title (outside the border)
	fullView.WriteString(styles.TitleStyle.Render("🚀 a3s - AWS IAM Roles"))
	fullView.WriteString("\n")

	// Search bar (if in search mode) - outside the border
	if m.searchMode {
		fullView.WriteString(styles.SearchPrompt.Render("Search: "))
		fullView.WriteString(m.searchInput.View())
		fullView.WriteString("\n")
	}

	// Calculate column widths based on terminal width
	availableWidth := m.width - 6 // Account for border and padding
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

	// Calculate visible height accounting for border
	borderHeight := 4 // Border takes up space
	titleHeight := 2
	searchHeight := 0
	if m.searchMode {
		searchHeight = 2
	}
	statusHeight := 2
	helpHeight := 1
	
	visibleHeight := m.height - borderHeight - titleHeight - searchHeight - statusHeight - helpHeight - 3
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

	// Status bar (outside the border) with full width
	fullView.WriteString(styles.RenderStatusBar(m.profile, m.region, len(m.filteredRoles), m.width))
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
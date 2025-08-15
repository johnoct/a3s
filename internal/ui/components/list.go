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
	ti := textinput.New()
	ti.Placeholder = "Search roles..."
	ti.CharLimit = 100

	m := ListModel{
		roles:         roles,
		filteredRoles: roles,
		searchInput:   ti,
		profile:       profile,
		region:        region,
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

	var s strings.Builder

	// Title
	s.WriteString(styles.TitleStyle.Render("🚀 a3s - AWS IAM Roles"))
	s.WriteString("\n\n")

	// Search bar (if in search mode)
	if m.searchMode {
		s.WriteString(styles.SearchPrompt.Render("Search: "))
		s.WriteString(m.searchInput.View())
		s.WriteString("\n\n")
	}

	// Column headers
	headers := fmt.Sprintf("%-40s %-20s %-20s %s",
		"Role Name",
		"Created",
		"Last Used",
		"Description",
	)
	s.WriteString(styles.ListHeader.Render(headers))
	s.WriteString("\n")

	// Role list
	visibleHeight := m.height - 10 // Account for headers, status bar, etc.
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

	for i := startIdx; i < endIdx; i++ {
		role := m.filteredRoles[i]
		
		created := role.CreateDate.Format("2006-01-02")
		lastUsed := "Never"
		if role.LastUsed != nil {
			lastUsed = role.LastUsed.Format("2006-01-02")
		}

		description := role.Description
		if len(description) > 40 {
			description = description[:37] + "..."
		}

		line := fmt.Sprintf("%-40s %-20s %-20s %s",
			truncate(role.Name, 40),
			created,
			lastUsed,
			description,
		)

		if i == m.cursor {
			s.WriteString(styles.SelectedItem.Render(line))
		} else {
			s.WriteString(styles.ListItem.Render(line))
		}
		s.WriteString("\n")
	}

	// Fill empty space
	for i := endIdx - startIdx; i < visibleHeight; i++ {
		s.WriteString("\n")
	}

	// Status bar
	s.WriteString("\n")
	s.WriteString(styles.RenderStatusBar(m.profile, m.region, len(m.filteredRoles)))
	s.WriteString("\n")

	// Help line
	s.WriteString(styles.RenderHelp())

	return s.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
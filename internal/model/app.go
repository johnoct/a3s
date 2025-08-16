package model

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/johnoct/a3s/internal/aws/client"
	"github.com/johnoct/a3s/internal/aws/iam"
	"github.com/johnoct/a3s/internal/aws/identity"
	"github.com/johnoct/a3s/internal/ui/components"
)

type State int

const (
	StateLoading State = iota
	StateList
	StateError
)

type App struct {
	state       State
	awsClient   *client.AWSClient
	roleService *iam.RoleService
	listModel   components.ListModel
	identity    *identity.Identity
	err         error
	width       int
	height      int
}

func NewApp(profile, region string) (*App, error) {
	return NewAppWithSize(profile, region, 80, 24)
}

func NewAppWithSize(profile, region string, width, height int) (*App, error) {
	ctx := context.Background()

	awsClient, err := client.New(ctx, profile, region)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS client: %w", err)
	}

	app := &App{
		state:       StateLoading,
		awsClient:   awsClient,
		roleService: iam.NewRoleService(awsClient),
		width:       width,
		height:      height,
	}

	return app, nil
}

type rolesLoadedMsg struct {
	roles []iam.Role
}

type identityLoadedMsg struct {
	identity *identity.Identity
}

type errorMsg struct {
	err error
}

func (a *App) Init() tea.Cmd {
	// Return batch of initialization commands
	return tea.Batch(
		a.loadRoles(),
		a.loadIdentity(),
		tea.WindowSize(), // Request initial window size
	)
}

func (a *App) loadRoles() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		roles, err := a.roleService.ListRoles(ctx)
		if err != nil {
			return errorMsg{err: err}
		}
		return rolesLoadedMsg{roles: roles}
	}
}

func (a *App) loadIdentity() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		id, err := identity.GetCallerIdentity(ctx, a.awsClient)
		if err != nil {
			// Don't fail the app if we can't get identity
			// Just return empty identity
			return identityLoadedMsg{identity: nil}
		}
		return identityLoadedMsg{identity: id}
	}
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		if a.state == StateList {
			return a.listModel.Update(msg)
		}
		return a, nil

	case rolesLoadedMsg:
		a.listModel = components.NewListModelWithSize(msg.roles, a.awsClient.Profile, a.awsClient.Region, a.width, a.height)
		if a.identity != nil {
			a.listModel.SetIdentity(a.identity)
		}
		a.state = StateList
		return a, a.listModel.Init()

	case identityLoadedMsg:
		a.identity = msg.identity
		if a.state == StateList {
			a.listModel.SetIdentity(a.identity)
		}
		return a, nil

	case errorMsg:
		a.err = msg.err
		a.state = StateError
		return a, nil

	case tea.KeyMsg:
		if a.state == StateError && msg.String() == "q" {
			return a, tea.Quit
		}
	}

	// Forward messages to list model when in list state
	if a.state == StateList {
		updatedModel, cmd := a.listModel.Update(msg)
		a.listModel = updatedModel.(components.ListModel)
		return a, cmd
	}

	return a, nil
}

func (a *App) View() string {
	switch a.state {
	case StateLoading:
		return "\n  Loading IAM roles... ⚡\n"
	case StateError:
		return fmt.Sprintf("\n  Error: %v\n\n  Press 'q' to quit.\n", a.err)
	case StateList:
		return a.listModel.View()
	default:
		return ""
	}
}

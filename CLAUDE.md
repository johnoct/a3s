# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**a3s** is a terminal user interface (TUI) application for AWS resources, inspired by k9s (Kubernetes TUI). It provides a fast, keyboard-driven interface for viewing and managing AWS resources directly from the terminal.

**Current Status**: MVP complete with IAM Role viewer featuring:
- K9s-style header with AWS identity display and ASCII art logo
- List view with search and filtering
- Detailed role view with tabs (Overview, Trust Policy, Policies, Tags)
- Async loading of role details including attached policies
- Vim-like keyboard navigation

## Technology Stack

- **Language**: Go
- **TUI Framework**: Bubble Tea (github.com/charmbracelet/bubbletea)
- **Styling**: Lipgloss (github.com/charmbracelet/lipgloss)
- **AWS SDK**: AWS SDK for Go v2
- **Architecture**: Model-View-Update pattern (Bubble Tea)
- **Terminal Support**: golang.org/x/term for terminal size detection

## Development Setup

```bash
# Initialize Go module
go mod init github.com/johnoct/a3s

# Install dependencies
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/iam
```

## Common Commands

```bash
# Run the application
go run cmd/a3s/main.go

# Build the application
go build -o a3s cmd/a3s/main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code (after installing golangci-lint)
golangci-lint run
```

## Architecture Notes

### Project Structure
```
a3s/
├── .claude/
│   └── docs/         # Context session documentation
├── cmd/
│   └── a3s/          # Main application entry point
├── internal/
│   ├── aws/          # AWS SDK integration layer
│   │   ├── iam/      # IAM-specific operations (roles.go)
│   │   ├── identity/ # AWS identity management (STS)
│   │   └── client/   # AWS client management
│   ├── ui/           # Bubble Tea UI components
│   │   ├── components/ # UI components (list.go, detail.go)
│   │   └── styles/   # Lipgloss styles and themes
│   └── model/        # Application state management (app.go)
```

### Key Design Patterns

1. **Bubble Tea Model-View-Update**:
   - Each UI component implements `tea.Model` interface
   - State changes through `Update()` method
   - Rendering through `View()` method
   - Async operations via tea.Cmd messages

2. **AWS Service Abstraction**:
   - AWS SDK calls wrapped in service layer
   - Pagination handled transparently
   - Async loading of detailed role information
   - Separate lightweight list vs detailed role fetching

3. **Keyboard Navigation**:
   - Vim-like bindings (j/k for up/down, g/G for top/bottom)
   - Search with `/`
   - Tab navigation in detail view
   - ESC to go back, q to quit

## Important Resources

- **Bubble Tea Examples**: https://github.com/charmbracelet/bubbletea/tree/master/examples
- **Lipgloss Documentation**: https://github.com/charmbracelet/lipgloss
- **AWS SDK Go v2**: https://aws.github.io/aws-sdk-go-v2/docs/
- **k9s (inspiration)**: https://github.com/derailed/k9s

## Testing Approach

- Unit tests for AWS service layer
- Component tests for UI elements
- Integration tests with AWS using local credentials
- Use table-driven tests for comprehensive coverage

## Features Implemented

### IAM Role Viewer
- **List View**: Displays all IAM roles with columns for name, creation date, last used, and description
- **Search**: Real-time filtering of roles with `/` command
- **Detail View**: Multi-tab interface showing:
  - Overview: Role ARN, ID, path, creation date, description, max session duration
  - Trust Policy: Formatted JSON of trust relationships
  - Policies: List of attached managed and inline policies (with async loading)
  - Tags: Key-value pairs associated with the role
- **Header**: K9s-style header showing AWS account, user, region, and profile with ASCII art logo

### UI Components
- Dynamic terminal size detection and responsive layout
- Color-coded interface with consistent styling
- Loading indicators for async operations
- Status bar showing current context and role count
- Help text for keyboard shortcuts

## Next Steps / TODO
- Add support for other AWS resources (EC2, S3, Lambda, etc.)
- Implement refresh functionality (`r` key)
- Add inline policy viewing/editing
- Export functionality for role configurations
- Add CloudTrail event viewing for roles
- Implement command mode (`:` prefix commands)
- Add configuration file support for custom settings

## Rules
- Always document progress when making significant changes in `.claude/docs/context_session_x.md` and read it before making more significant changes
- Always use @agent-git-manager to push changes when asked to save and push
- Follow existing code patterns and conventions
- Maintain comprehensive error handling
- Keep UI responsive with async operations
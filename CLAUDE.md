# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**a3s** is a terminal user interface (TUI) application for AWS resources, inspired by k9s (Kubernetes TUI). It provides a fast, keyboard-driven interface for viewing and managing AWS resources directly from the terminal.

**MVP Focus**: IAM Role viewer with search, navigation, and detailed views.

## Technology Stack

- **Language**: Go
- **TUI Framework**: Bubble Tea (github.com/charmbracelet/bubbletea)
- **Styling**: Lipgloss (github.com/charmbracelet/lipgloss)
- **AWS SDK**: AWS SDK for Go v2
- **Architecture**: Model-View-Update pattern (Bubble Tea)

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
├── cmd/
│   └── a3s/          # Main application entry point
├── internal/
│   ├── aws/          # AWS SDK integration layer
│   │   ├── iam/      # IAM-specific operations
│   │   └── client/   # AWS client management
│   ├── ui/           # Bubble Tea UI components
│   │   ├── list/     # List view component
│   │   ├── detail/   # Detail view component
│   │   └── styles/   # Lipgloss styles
│   ├── model/        # Application state management
│   └── config/       # Configuration management
```

### Key Design Patterns

1. **Bubble Tea Model-View-Update**:
   - Each UI component implements `tea.Model` interface
   - State changes through `Update()` method
   - Rendering through `View()` method

2. **AWS Service Abstraction**:
   - AWS SDK calls wrapped in service layer
   - Handle pagination transparently
   - Cache responses when appropriate

3. **Keyboard Navigation**:
   - Vim-like bindings (j/k for up/down)
   - Search with `/`
   - Command mode with `:`

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

## Rules
- this is really good so far, lets always document our progress when we've made significant changes so far in ./claude/docs/context_session_x.md and read it before we decide to make more significant changes

- always use @agent-git-manager to push it up when i ask you to save and push things up
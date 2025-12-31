# CLAUDE.md - Context Guide for Claude Code

> **Context Engineering Note**: This file serves as Claude Code's primary project memory. It's structured to optimize Claude's comprehension and follows Claude Code memory management best practices. Each section provides progressive context building with clear hierarchical organization.

## Project Identity & Status

**a3s** is a terminal user interface (TUI) application for AWS resources, inspired by k9s (Kubernetes TUI). It provides a fast, keyboard-driven interface for viewing and managing AWS resources directly from the terminal.

**Current Development Phase**: MVP Complete - Production-Ready IAM Role Viewer
**Status**: ✅ Fully functional with comprehensive IAM role management capabilities
**Latest Release**: Automated releases via GoReleaser for darwin/linux (amd64/arm64)

### Core Features Delivered
- K9s-inspired header with AWS identity display and ASCII art logo
- Responsive list view with real-time search and filtering (`/` command)
- Multi-tab detail view (Overview, Trust Policy, Policies, Tags)
- Async role detail loading with loading indicators
- Vim-like keyboard navigation (j/k, g/G, Tab, ESC, q)
- Complete IAM policy display (managed + inline policies)
- **Interactive policy document viewer** with full JSON display and navigation
- Policy selection and viewing with async loading (Enter to view, ESC to return)
- **Consistent UI layout system** with standardized border alignment and terminal compatibility
- **Automated release workflow** with GoReleaser and GitHub Actions

## Technology Stack & Dependencies

**Primary Stack**:
- **Language**: Go 1.24.2
- **TUI Framework**: Bubble Tea v1.3.6 (github.com/charmbracelet/bubbletea)
- **Styling Engine**: Lipgloss v1.1.0 (github.com/charmbracelet/lipgloss)
- **AWS Integration**: AWS SDK for Go v2 (github.com/aws/aws-sdk-go-v2)
- **Release Automation**: GoReleaser with GitHub Actions

**Architecture Pattern**: Model-View-Update (Bubble Tea's Elm Architecture)

**Key Dependencies**:
```go
// Core UI framework
github.com/charmbracelet/bubbletea v1.3.6
github.com/charmbracelet/lipgloss v1.1.0
github.com/charmbracelet/bubbles v0.21.0

// AWS integration
github.com/aws/aws-sdk-go-v2 v1.38.0
github.com/aws/aws-sdk-go-v2/config v1.31.0
github.com/aws/aws-sdk-go-v2/service/iam v1.46.0
github.com/aws/aws-sdk-go-v2/service/sts v1.37.0

// Terminal utilities
golang.org/x/term v0.34.0
```

## Codebase Architecture & File Organization

### Project Structure Map
```
a3s/                          # Root project directory
├── .claude/
│   ├── docs/                 # Context session documentation (7 sessions)
│   └── settings.local.json   # Local Claude Code settings
├── .github/
│   └── workflows/
│       └── release.yml       # GitHub Actions release workflow
├── cmd/a3s/                  # Application entry points
│   ├── main.go              # Current main application
│   └── main_old.go          # Legacy main (reference)
├── internal/                 # Private application code
│   ├── aws/                 # AWS SDK abstraction layer
│   │   ├── client/          # AWS client configuration
│   │   ├── iam/             # IAM service operations
│   │   └── identity/        # STS identity management
│   ├── model/               # Application state management
│   │   └── app.go          # Main application model
│   └── ui/                  # Bubble Tea UI layer
│       ├── components/      # Reusable UI components
│       │   ├── list.go     # Role list component
│       │   └── detail.go   # Role detail component
│       └── styles/          # Lipgloss styling
├── bin/                     # Built binaries (gitignored)
├── .goreleaser.yml          # GoReleaser configuration
├── demo.gif                 # VHS-generated demo animation
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── Makefile                 # Build automation
└── README.md               # Public documentation
```

### Critical File Locations (for Claude Code navigation)
- **Main Entry**: `/cmd/a3s/main.go`
- **Core App Logic**: `/internal/model/app.go`
- **UI Components**: `/internal/ui/components/{list.go,detail.go}`
- **AWS Services**: `/internal/aws/{iam/roles.go,identity/identity.go,client/client.go}`
- **Styling**: `/internal/ui/styles/styles.go`
- **Release Config**: `/.goreleaser.yml`
- **CI/CD Workflow**: `/.github/workflows/release.yml`
- **Context Documentation**: `/.claude/docs/context_session_*.md`

## Development Workflow & Commands

### Essential Commands (Priority Order)
```bash
# Development cycle
go run cmd/a3s/main.go           # Run application directly
make build                       # Build binary to bin/a3s
make run                         # Build and run
make dev                         # Format, test, and build
go test ./...                    # Run all tests
go fmt ./...                     # Format code (always before commits)

# Quality assurance
golangci-lint run                # Lint code (install first)
go mod tidy                      # Clean dependencies
go vet ./...                     # Static analysis
make lint                        # Run linter via Makefile

# Testing
make test                        # Run all tests
make test-coverage               # Run tests with coverage report
```

### Build System
- **Makefile present**: Use `make` commands for consistent builds
- **Binary output**: `bin/a3s` (executable)
- **Go version**: 1.24.2 (ensure compatibility)

### Release Process
```bash
# Releases are automated via GitHub Actions
# Triggered by pushing version tags:
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

**GoReleaser Configuration**:
- Builds for: darwin/linux on amd64/arm64
- CGO disabled for cross-platform compatibility
- Optimized with `-trimpath` and `-s -w` ldflags
- Archives as tar.gz (zip for Windows)
- Auto-generates changelog excluding docs/test/chore commits

## Key Design Patterns & Implementation Details

### 1. Bubble Tea Architecture Implementation
```go
// Each UI component implements tea.Model interface
type Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() string
}
```

**Message Flow**: User Input → Update() → State Change → View() → UI Render
**Async Operations**: Handled via `tea.Cmd` messages (critical for AWS API calls)

### 2. AWS Service Layer Pattern
```go
// Service abstraction in internal/aws/
type RoleService struct {
    client *iam.Client
}

// Key methods
func (rs *RoleService) ListRoles() ([]Role, error)
func (rs *RoleService) GetRoleDetails(roleName string) (*RoleDetails, error)
func (rs *RoleService) GetManagedPolicyDocument(policyArn, version string) (string, error)
```

**Important**: Always use `GetRoleDetails()` for complete role information (includes policies, tags)

### 3. Async Loading Pattern
- **Trigger**: User selects role (Enter key) or policy document (Enter in Policies tab)
- **Flow**: Loading indicator → AWS API call → Update UI with detailed data
- **Implementation**: `roleDetailsLoadedMsg` and `policyDocumentLoadedMsg` message types in Bubble Tea pattern
- **Policy Document Loading**: Async fetching of IAM policy JSON documents with loading indicators

### 4. UI State Management
- **Navigation State**: Tracked in main app model
- **Component State**: Each UI component manages its own state
- **Loading States**: Explicit loading indicators for async operations
- **No Status Bar**: Clean interface without status bar for minimal design

### 5. UI Layout System
- **Border Alignment**: Standardized 2-space left padding for consistent border positioning
- **Container Spacing**: Tab spacing adjusted to 3 spaces (2 for border + 1 for padding)
- **Width Calculations**: Edge-to-edge content filling with proper container constraints
- **Terminal Compatibility**: Mouse capture disabled for improved text selection in modern terminals
- **Header Positioning**: Consistent header alignment across all views and components
- **Compact ASCII Art**: Optimized logo for minimal vertical space usage

## Recent Development History & Context

### Latest Changes: Release Infrastructure & UI Polish
**Commits**: `dbfd5ba` through `f724bcf`

**Key Updates**:
1. **GoReleaser Integration**: Minimal release workflow with automated builds for darwin/linux
2. **VHS Demo**: Added demo.gif showing full application functionality
3. **README Updates**: Comprehensive documentation with demo, keyboard shortcuts, and architecture
4. **UI Polish**: Compact ASCII art, removed status bar, improved header spacing

### Session 7: UI Layout and Border Alignment Fixes
**Problem Solved**: Border alignment inconsistencies between different views with misaligned headers, tabs, and borders
**Root Cause**: Inconsistent spacing calculations and mouse capture interfering with terminal text selection
**Solution**: Standardized UI layout system with consistent border positioning and improved terminal compatibility

**Key Changes**:
- **Border Alignment Consistency**: Standardized 2-space left padding across all components
- **Tab Positioning**: Adjusted tab spacing from 2 to 3 spaces
- **Mouse Capture Removal**: Removed `tea.WithMouseCellMotion()` for better terminal compatibility
- **Width Calculations**: Fixed JSON content width calculations for edge-to-edge container filling

### Session 6: IAM Policy Document Viewer Implementation
**Key Changes**:
- Added `PolicyInfo` struct with ARN support for managed policies
- Implemented `GetManagedPolicyDocument()` service method
- Created interactive policy navigation (j/k selection, Enter to view)
- Added policy document viewer with scroll support (j/k, g/G navigation)

### Critical Implementation Notes
- **AWS API Strategy**: Multi-tier loading (list → details → policy documents on demand)
- **Performance**: Async loading prevents UI blocking for all AWS API calls
- **Policy Document Viewing**: Interactive navigation with full JSON display and async loading
- **Navigation Flow**: Hierarchical navigation (Roles → Role Detail → Policy Document → back to Policy List)
- **Error Handling**: Comprehensive error handling for AWS API failures
- **UI Layout Consistency**: Standardized 2-space padding and 3-space tab positioning
- **Terminal Compatibility**: Mouse capture disabled to maintain text selection capabilities

## Testing Strategy & Quality Assurance

### Test Structure
```bash
internal/
├── aws/        # Unit tests for AWS service layer
├── ui/         # Component tests for UI elements
└── model/      # Integration tests for app state
```

### Testing Approach
- **Unit Tests**: AWS service layer with mocked AWS clients
- **Component Tests**: UI components with mock data
- **Integration Tests**: Full workflow with real AWS credentials
- **Table-Driven Tests**: Comprehensive coverage for data transformations

### Quality Standards
- **Error Handling**: Always handle AWS API errors gracefully
- **UI Responsiveness**: No blocking operations in UI thread
- **Memory Management**: Proper cleanup of AWS resources
- **Terminal Compatibility**: Support various terminal sizes and capabilities

## Context Session Documentation System

### Documentation Pattern
- **Location**: `/.claude/docs/context_session_N.md`
- **Purpose**: Track development progress and decisions
- **Current Count**: 7 sessions (context_session_1.md through context_session_7.md)

**Before making significant changes**, always:
1. Read the latest context session for recent developments
2. Document new changes in the next session file
3. Update CLAUDE.md if architectural changes occur

## Claude Code Interaction Guidelines

### Preferred Development Approach
1. **Always read latest context session** before making changes
2. **Use existing patterns** - follow established conventions
3. **Maintain async patterns** - keep UI responsive
4. **Test incrementally** - run application after each change
5. **Document progress** - update context sessions

### File Navigation Priorities
When asked to modify functionality:
1. Start with `/internal/model/app.go` for app-level changes
2. UI changes: `/internal/ui/components/{list.go,detail.go}`
3. AWS integration: `/internal/aws/{iam,identity,client}/`
4. Entry point: `/cmd/a3s/main.go`
5. Release config: `/.goreleaser.yml`

### Code Quality Expectations
- **Idiomatic Go**: Follow Go best practices and conventions
- **Error Handling**: Comprehensive error handling for all AWS operations
- **UI Patterns**: Consistent with Bubble Tea patterns
- **Performance**: Maintain async loading and responsive UI
- **Documentation**: Comment complex logic and AWS integrations

### Common Operations Map
- **Add new AWS resource type**: Extend `/internal/aws/` with new service
- **UI component changes**: Modify `/internal/ui/components/`
- **New navigation features**: Update main app model
- **Styling changes**: Modify `/internal/ui/styles/styles.go`
- **Key bindings**: Update component `Update()` methods
- **Policy document viewing**: Navigate with j/k in Policies tab, Enter to view, ESC to return
- **Document scrolling**: j/k for line-by-line, g/G for top/bottom navigation
- **UI layout fixes**: Maintain 2-space padding for borders, 3-space tab positioning
- **Terminal compatibility**: Consider mouse capture impact on text selection
- **Release updates**: Modify `.goreleaser.yml` for build configuration

## Future Development Roadmap

### High Priority Features
- **Resource Expansion**: EC2, S3, Lambda, RDS support
- **Refresh Functionality**: Real-time data updates (`r` key) - currently TODO
- **Configuration System**: User preferences and settings
- **Export Capabilities**: Save role configurations

### UI/UX Enhancements
- **Command Mode**: `:` prefix commands (vim-style)
- **Help System**: Built-in help and shortcuts (`?` key)
- **Theme Support**: Multiple color schemes
- **Search Enhancement**: Advanced filtering options

### Technical Improvements
- **Caching Layer**: Reduce AWS API calls
- **Performance Optimization**: Large dataset handling
- **Error Recovery**: Robust error handling and retry logic
- **Testing Coverage**: Comprehensive test suite
- **Windows Support**: Add Windows builds in GoReleaser

## Development Rules & Constraints

### Mandatory Practices
- **Context Documentation**: Always update `.claude/docs/context_session_N.md` for significant changes
- **Self Context Documentation**: Always use claude-context-engineer agent to update CLAUDE.md before you commit and push up
- **Leverage Agents**: Review agents you have access to and use it as needed. For example, when working on the tui, consult and work with charm-tui-developer agent.
- **Review your code**: Before committing, consult with code-reviewer agent to review your code.
- **Code Formatting**: Run `go fmt ./...` before any commits
- **Committing Changes**: When asked to save and push up or making any commits, always use the git-manager agent.
- **Error Handling**: Never ignore AWS API errors
- **UI Responsiveness**: Maintain async patterns for all blocking operations
- **Pattern Consistency**: Follow established Bubble Tea and Go patterns

### Quality Gates
- **Build Success**: All code must compile cleanly
- **Test Passing**: Existing tests must continue to pass
- **Lint Clean**: Code must pass golangci-lint checks
- **Functional Verification**: Test critical paths manually

### Architectural Constraints
- **No Blocking UI**: All AWS API calls must be async
- **Service Separation**: Maintain clear separation between AWS/UI/Model layers
- **Component Isolation**: UI components should be self-contained
- **Error Boundaries**: Each component handles its own error states
- **No Mouse Capture**: Keep terminal text selection working

---

**Meta-Context for Claude Code**: This CLAUDE.md is designed to provide hierarchical context that aligns with Claude Code's memory retrieval patterns. Each section builds upon the previous one, providing both high-level understanding and specific implementation details. The structure supports both quick reference and deep understanding based on the complexity of the requested task.

**Last Updated**: Reflects codebase state as of commit `dbfd5ba` (feat: add minimal release workflow with GoReleaser)

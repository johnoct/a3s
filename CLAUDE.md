# CLAUDE.md - Context Guide for Claude Code

> **Context Engineering Note**: This file serves as Claude Code's primary project memory. It's structured to optimize Claude's comprehension and follows Claude Code memory management best practices. Each section provides progressive context building with clear hierarchical organization.

## Project Identity & Status

**a3s** is a terminal user interface (TUI) application for AWS resources, inspired by k9s (Kubernetes TUI). It provides a fast, keyboard-driven interface for viewing and managing AWS resources directly from the terminal.

**Current Development Phase**: MVP Complete - Production-Ready IAM Role Viewer
**Status**: ✅ Fully functional with comprehensive IAM role management capabilities

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

## Technology Stack & Dependencies

**Primary Stack**:
- **Language**: Go 1.24.2
- **TUI Framework**: Bubble Tea v1.3.6 (github.com/charmbracelet/bubbletea)
- **Styling Engine**: Lipgloss v1.1.0 (github.com/charmbracelet/lipgloss)
- **AWS Integration**: AWS SDK for Go v2 (github.com/aws/aws-sdk-go-v2)

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
```

## Codebase Architecture & File Organization

### Project Structure Map
```
a3s/                          # Root project directory
├── .claude/
│   ├── docs/                 # Context session documentation (7 sessions)
│   └── settings.local.json   # Local Claude Code settings
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
├── bin/                     # Built binaries
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
- **Context Documentation**: `/.claude/docs/context_session_*.md`

## Development Workflow & Commands

### Essential Commands (Priority Order)
```bash
# Development cycle
go run cmd/a3s/main.go           # Run application directly
make build                       # Build binary to bin/a3s
make run                         # Build and run
go test ./...                    # Run all tests
go fmt ./...                     # Format code (always before commits)

# Quality assurance
golangci-lint run                # Lint code (install first)
go mod tidy                      # Clean dependencies
go vet ./...                     # Static analysis
```

### Build System
- **Makefile present**: Use `make` commands for consistent builds
- **Binary output**: `bin/a3s` (executable)
- **Go version**: 1.24.2 (ensure compatibility)

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
```

**Important**: Always use `GetRoleDetails()` for complete role information (includes policies, tags)

### 3. Async Loading Pattern (Recently Implemented)
- **Trigger**: User selects role (Enter key) or policy document (Enter in Policies tab)
- **Flow**: Loading indicator → AWS API call → Update UI with detailed data
- **Implementation**: `roleDetailsLoadedMsg` and `policyDocumentLoadedMsg` message types in Bubble Tea pattern
- **Policy Document Loading**: Async fetching of IAM policy JSON documents with loading indicators

### 4. UI State Management
- **Navigation State**: Tracked in main app model
- **Component State**: Each UI component manages its own state
- **Loading States**: Explicit loading indicators for async operations

### 5. UI Layout System (Session 7 Implementation)
- **Border Alignment**: Standardized 2-space left padding for consistent border positioning
- **Container Spacing**: Tab spacing adjusted to 3 spaces (2 for border + 1 for padding)
- **Width Calculations**: Edge-to-edge content filling with proper container constraints
- **Terminal Compatibility**: Mouse capture disabled for improved text selection in modern terminals
- **Header Positioning**: Consistent header alignment across all views and components

## Recent Development History & Context

### Session 7 (Latest): UI Layout and Border Alignment Fixes
**Problem Solved**: Border alignment inconsistencies between different views with misaligned headers, tabs, and borders
**Root Cause**: Inconsistent spacing calculations and mouse capture interfering with terminal text selection
**Solution**: Standardized UI layout system with consistent border positioning and improved terminal compatibility

**Key Changes**:
- **Border Alignment Consistency**: Standardized 2-space left padding across all components for proper border positioning
- **Tab Positioning**: Adjusted tab spacing from 2 to 3 spaces (2 for border position + 1 for MainContainer padding)
- **Header Positioning**: Removed inconsistent top margins and standardized header positioning across views
- **Mouse Capture Removal**: Removed `tea.WithMouseCellMotion()` from main.go to enable text selection in Kitty terminal
- **Width Calculations**: Fixed JSON content width calculations for edge-to-edge container filling
- **Container Padding**: Standardized spacing calculations across list and detail components

**Files Modified**:
- `/cmd/a3s/main.go` - Removed mouse capture for better terminal compatibility
- `/internal/ui/components/detail.go` - Multiple alignment and spacing fixes
- `/internal/ui/components/list.go` - Consistent border padding implementation

### Session 6: IAM Policy Document Viewer Implementation
**Problem Solved**: Users needed to view actual policy JSON documents without leaving the TUI
**Root Cause**: Policy list showed names but not document content
**Solution**: Implemented comprehensive policy document viewing with interactive navigation

**Key Changes**:
- Added `PolicyInfo` struct with ARN support for managed policies
- Implemented `GetManagedPolicyDocument()` service method
- Created interactive policy navigation (j/k selection, Enter to view)
- Added policy document viewer with scroll support (j/k, g/G navigation)
- Fixed navigation flow (ESC from policy document returns to policies tab)
- Async policy document loading with loading indicators

### Session 5: IAM Policy Display Fix
**Problem Solved**: Roles showing "No policies attached" when policies existed
**Root Cause**: Using `ListRoles()` basic data instead of `GetRoleDetails()`
**Solution**: Implemented async role detail loading with proper AWS API integration

**Key Changes**:
- Enhanced `ListModel` with `roleService` dependency injection
- Added `loadRoleDetails()` async command
- Proper integration of detailed role data including policies

### Critical Implementation Notes
- **AWS API Strategy**: Multi-tier loading (list → details → policy documents on demand)
- **Performance**: Async loading prevents UI blocking for all AWS API calls
- **Policy Document Viewing**: Interactive navigation with full JSON display and async loading
- **Navigation Flow**: Hierarchical navigation (Roles → Role Detail → Policy Document → back to Policy List)
- **Error Handling**: Comprehensive error handling for AWS API failures including policy document fetching
- **UI Layout Consistency**: Standardized 2-space padding and 3-space tab positioning for consistent borders
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

## Future Development Roadmap

### High Priority Features
- **Resource Expansion**: EC2, S3, Lambda, RDS support
- **Refresh Functionality**: Real-time data updates (`r` key)
- **Configuration System**: User preferences and settings
- **Export Capabilities**: Save role configurations

### UI/UX Enhancements
- **Command Mode**: `:` prefix commands (vim-style)
- **Help System**: Built-in help and shortcuts
- **Theme Support**: Multiple color schemes
- **Search Enhancement**: Advanced filtering options

### Technical Improvements
- **Caching Layer**: Reduce AWS API calls
- **Performance Optimization**: Large dataset handling
- **Error Recovery**: Robust error handling and retry logic
- **Testing Coverage**: Comprehensive test suite

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

---

**Meta-Context for Claude Code**: This CLAUDE.md is designed to provide hierarchical context that aligns with Claude Code's memory retrieval patterns. Each section builds upon the previous one, providing both high-level understanding and specific implementation details. The structure supports both quick reference and deep understanding based on the complexity of the requested task.
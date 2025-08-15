# Product Requirements Document: a3s
## AWS Terminal User Interface (TUI) Application

**Version**: 1.0  
**Date**: August 2025  
**Status**: Draft

---

## 1. Executive Summary

**a3s** is a terminal-based user interface for managing AWS resources, inspired by the popular Kubernetes TUI tool k9s. It provides a fast, keyboard-driven interface for viewing and managing AWS resources directly from the terminal, eliminating the need for context switching between the AWS Console and command line.

### Key Value Propositions
- **Speed**: Navigate AWS resources 10x faster than the web console
- **Efficiency**: Keyboard-driven interface with vim-like navigation
- **Focus**: Stay in the terminal without context switching
- **Overview**: Get a clear, dense view of resources at a glance

---

## 2. Problem Statement

### Current Pain Points
1. **AWS Console is slow**: Multiple page loads, excessive clicking, and poor information density
2. **Context switching**: Constantly switching between terminal and browser breaks flow
3. **AWS CLI limitations**: Great for automation but poor for exploration and overview
4. **Multi-account complexity**: Switching between accounts/regions is cumbersome
5. **IAM visibility**: Understanding role relationships and permissions requires multiple console pages

### Target Users
- **DevOps Engineers**: Managing infrastructure and troubleshooting issues
- **Security Engineers**: Auditing IAM permissions and compliance
- **Platform Engineers**: Understanding service relationships
- **Developers**: Debugging permission issues and understanding resource configurations

---

## 3. Solution Overview

a3s provides a k9s-like experience for AWS resources, starting with IAM roles as the MVP. Users can quickly navigate, search, and inspect AWS resources using familiar keyboard shortcuts in a clean, information-dense terminal interface.

### Core Principles
- **Keyboard-first**: All actions accessible via keyboard shortcuts
- **Information density**: Show maximum useful information in minimal space
- **Speed**: Sub-second response times for navigation
- **Familiarity**: Use k9s/vim navigation patterns

---

## 4. MVP Scope: IAM Role Viewer

### 4.1 Functional Requirements

#### Main List View
- Display all IAM roles in current account/region
- Columns: Role Name, Creation Date, Last Used, Trust Policy Summary, Policy Count
- Real-time search/filter as you type
- Sort by any column
- Color coding for role types (service roles, user roles, etc.)
- Pagination handling for large role lists

#### Detail View
- **Trust Relationships**: Full trust policy JSON with syntax highlighting
- **Managed Policies**: List of attached AWS managed and customer managed policies
- **Inline Policies**: View inline policy documents
- **Tags**: Display all tags
- **Metadata**: ARN, creation date, last modified, max session duration
- **Usage History**: Last assumed time and by whom (if available)

#### Navigation
- `j/k` or arrow keys: Navigate up/down
- `Enter`: View role details
- `Esc`: Go back to list view
- `/`: Start search
- `n/N`: Next/previous search result
- `g/G`: Go to top/bottom
- `:`: Command mode
- `q`: Quit application
- `r`: Refresh current view
- `?`: Show help/shortcuts

#### AWS Integration
- Use AWS SDK for Go v2
- Support AWS CLI profiles (~/.aws/credentials)
- Support environment variables (AWS_PROFILE, AWS_REGION)
- Profile switcher (`:profile <name>`)
- Region switcher (`:region <region>`)
- Handle AWS API rate limits gracefully

### 4.2 Non-Functional Requirements

#### Performance
- Initial load < 2 seconds
- Navigation response < 100ms
- Search results update in real-time
- Handle 1000+ roles efficiently

#### User Experience
- Clean, readable interface using Lipgloss styling
- Responsive layout adapting to terminal size
- Clear error messages with recovery suggestions
- Loading indicators for long operations

#### Technical Stack
- **Language**: Go
- **TUI Framework**: Bubble Tea (Model-View-Update architecture)
- **Styling**: Lipgloss
- **AWS SDK**: AWS SDK for Go v2
- **Configuration**: Viper for config management

---

## 5. Success Metrics

### Quantitative Metrics
- Time to view role details: < 2 seconds (vs ~30s in console)
- Number of keystrokes to common tasks: 50% reduction
- Application startup time: < 1 second
- Support for 100+ IAM roles without performance degradation

### Qualitative Metrics
- User can perform 80% of IAM viewing tasks without opening console
- Intuitive enough that k9s users need no documentation
- Clear enough that non-k9s users can use with minimal learning

---

## 6. Technical Architecture

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
├── go.mod
├── go.sum
└── README.md
```

### Key Components
1. **AWS Service Layer**: Abstracts AWS SDK calls, handles pagination
2. **Model Layer**: Manages application state using Bubble Tea patterns
3. **View Components**: Reusable UI components for list/detail views
4. **Style System**: Consistent Lipgloss styling across components
5. **Config Manager**: Handles profiles, regions, and user preferences

---

## 7. Development Roadmap

### Phase 1: MVP (Week 1-2)
- [x] Project setup with Bubble Tea/Lipgloss
- [ ] Basic IAM role list view
- [ ] Role detail view
- [ ] Keyboard navigation
- [ ] Search functionality

### Phase 2: Enhancement (Week 3-4)
- [ ] Profile/region switching
- [ ] Improved styling and layout
- [ ] Performance optimization
- [ ] Error handling
- [ ] Basic documentation

### Phase 3: Future Expansion
- [ ] EC2 instances viewer
- [ ] S3 bucket browser
- [ ] Lambda function viewer
- [ ] CloudFormation stack navigator
- [ ] Cross-resource navigation

---

## 8. Risks and Mitigations

| Risk | Impact | Mitigation |
|------|---------|------------|
| AWS API rate limits | High | Implement caching and request batching |
| Large resource lists | Medium | Implement virtual scrolling and pagination |
| Complex IAM policies | Medium | Provide both formatted and raw JSON views |
| Terminal compatibility | Low | Test on multiple terminal emulators |

---

## 9. Open Questions

1. Should the MVP support write operations (role modifications)?
2. How should cross-account role assumptions be handled?
3. Should we cache API responses for faster navigation?
4. What level of policy analysis should be included (effective permissions)?

---

## 10. Acceptance Criteria

The MVP is complete when:
- [ ] Users can list all IAM roles in their AWS account
- [ ] Users can view detailed information for any role
- [ ] Search works across role names and descriptions
- [ ] All primary keyboard shortcuts are functional
- [ ] The app works with standard AWS credential chains
- [ ] Basic documentation exists
- [ ] The app handles errors gracefully
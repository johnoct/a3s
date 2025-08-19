# Session 6: IAM Policy Document Viewer Implementation

## Overview
Implemented a comprehensive policy document viewing feature that allows users to view the actual JSON content of both AWS managed and inline policies attached to IAM roles.

## Problem Statement
Users could see the list of policies attached to a role but couldn't view the actual policy documents and their permission details without leaving the TUI to check the AWS Console.

## Features Implemented

### 1. Enhanced AWS Service Layer
**File**: `internal/aws/iam/roles.go`

#### Added PolicyInfo Structure
- Created `PolicyInfo` struct to store both policy name and ARN for managed policies
- Updated `Role` struct to use `[]PolicyInfo` instead of `[]string` for managed policies
- This change enables fetching policy documents using the ARN

#### New Service Method
- Implemented `GetManagedPolicyDocument(ctx, policyARN)` method
- Fetches the latest version of managed policy documents from AWS
- Returns formatted JSON for display

### 2. Interactive Policy Navigation
**File**: `internal/ui/components/detail.go`

#### Navigation Features
- **j/k keys**: Navigate up/down through the list of policies in the Policies tab
- **Enter key**: View the selected policy's full JSON document
- **Visual selection**: Currently selected policy is highlighted
- **Policy counter**: Shows which policy is selected (e.g., "1 of 3")

#### View States
- Added `viewState` enum with two states:
  - `viewNormal`: Standard tab view with policy list
  - `viewPolicyDocument`: Full-screen policy JSON viewer
- Proper state transitions maintain user context

### 3. Policy Document Viewer
**File**: `internal/ui/components/detail.go`

#### Document Display Features
- Full-screen JSON document view with syntax highlighting
- Scrollable content for long policy documents
- Navigation controls:
  - **j/k**: Scroll up/down line by line
  - **g**: Jump to top of document
  - **G**: Jump to bottom of document
  - **ESC**: Return to policies list (not to roles list)

#### Async Loading
- Non-blocking policy document fetching
- Loading indicator: "Loading policy document..."
- Error handling with user-friendly messages
- Supports both managed and inline policies seamlessly

### 4. Navigation Flow Fix
**Files**: `internal/ui/components/detail.go`, `internal/ui/components/list.go`

#### Issue Resolved
- Initial implementation: ESC from policy document incorrectly returned to roles list
- Fixed: ESC from policy document now returns to policies tab
- Added `IsViewingPolicyDocument()` helper method to DetailModel
- List component checks this state before handling ESC key

#### Navigation Hierarchy
1. Roles List → (Enter) → Role Detail View
2. Role Detail View → (Tab) → Policies Tab
3. Policies Tab → (j/k + Enter) → Policy Document View
4. Policy Document View → (ESC) → Policies Tab
5. Policies Tab → (ESC) → Roles List

## Technical Implementation Details

### Message Flow for Policy Loading
1. User presses Enter on selected policy
2. `loadSelectedPolicy()` command initiated
3. Async fetch from AWS (managed or inline policy)
4. `policyDocumentLoadedMsg` received with document
5. View state transitions to `viewPolicyDocument`
6. Document rendered with scroll position reset

### State Management
- `selectedPolicy`: Tracks which policy is highlighted (0-indexed)
- `viewState`: Controls which view is active (normal vs document)
- `policyDocument`: Stores the loaded JSON document
- `loadingPolicy`: Prevents duplicate loading requests
- `scrollY`: Maintains scroll position in document view

### AWS API Integration
- Managed policies: Require ARN to fetch document
- Inline policies: Use role name + policy name
- Both types handled transparently to the user
- Proper URL decoding for policy documents

## UI/UX Improvements

### Visual Feedback
- Selected policy highlighted with consistent app styling
- Loading states clearly indicated
- Error messages displayed inline
- Dynamic help text based on context

### Help Text Updates
- Policies tab with policies: "Press Enter to view the selected policy document"
- Policies tab without policies: "No policies attached"
- Policy document view: Updated help shows scroll controls
- Context-aware help improves discoverability

## Files Modified
- `internal/aws/iam/roles.go` - Added PolicyInfo struct and GetManagedPolicyDocument method
- `internal/ui/components/detail.go` - Complete policy viewer implementation
- `internal/ui/components/list.go` - Navigation state management fix

## Testing Performed
- Tested with roles containing both managed and inline policies
- Verified navigation flow (ESC returns to correct view)
- Tested scrolling in long policy documents
- Confirmed async loading doesn't block UI
- Validated error handling for failed policy fetches

## Next Steps
- Consider adding policy search/filter within document view
- Add policy version history viewing
- Implement policy comparison between roles
- Add export functionality for policy documents
- Consider syntax highlighting for specific policy elements

## Result
Users can now seamlessly browse and view full policy documents directly within the TUI, eliminating the need to switch to the AWS Console for policy inspection. The feature maintains the k9s-inspired keyboard-driven interface while providing comprehensive policy visibility.
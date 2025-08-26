# a3s - AWS Terminal User Interface

🚀 **a3s** is a terminal-based user interface for AWS resources, inspired by [k9s](https://k9scli.io/) for Kubernetes. Navigate and manage AWS resources with speed and efficiency, all from your terminal.

![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue)
![AWS SDK](https://img.shields.io/badge/AWS%20SDK-v2-orange)
![License](https://img.shields.io/badge/license-MIT-green)

## Demo

![a3s Demo](demo.gif)

**What you'll see in this demo:**
- 🎨 Beautiful k9s-inspired interface with AWS identity display and ASCII art logo
- 📋 Navigate through IAM roles using vim-like j/k keys
- 🔍 Real-time search filtering with `/` command (searching for "lambda" and "service")
- 📄 Detailed role views with tabbed interface (Overview, Trust Policy, Policies, Tags)
- 📜 Interactive policy document viewer with full JSON display
- ⌨️ Seamless navigation: Tab between sections, Enter to view details, ESC to navigate back
- 🏃 Responsive async loading for all AWS API calls

## Features

### MVP - IAM Role Viewer
- 📋 **List all IAM roles** with sortable columns
- 🔍 **Real-time search** filtering with `/` command
- 📄 **Detailed role view** with tabbed interface:
  - **Overview**: Role metadata and last usage information
  - **Trust Policy**: Trust relationships and assume role policies
  - **Policies**: Attached managed and inline policies with interactive JSON viewer
  - **Tags**: Role tags and metadata
- 📜 **Interactive policy document viewer** - select any policy to view full JSON with navigation
- ⌨️ **Vim-like keyboard navigation** (j/k, g/G, Tab, Enter, ESC)
- 🔄 **AWS profile and region switching**
- 🎨 **Beautiful k9s-inspired TUI** with AWS identity display and consistent styling
- ⚡ **Async loading** with loading indicators for responsive performance

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/johnoct/a3s.git
cd a3s

# Build the application
go build -o a3s cmd/a3s/main.go

# Optionally, move to your PATH
sudo mv a3s /usr/local/bin/
```

### Requirements
- Go 1.24 or higher
- AWS credentials configured (~/.aws/credentials)
- Terminal with 256 color support

## Usage

### Basic Usage

```bash
# Use default AWS profile and region
a3s

# Use specific profile
a3s -profile production

# Use specific region
a3s -region us-west-2

# Combine profile and region
a3s -profile dev -region eu-west-1
```

### Keyboard Shortcuts

#### List View
| Key | Action |
|-----|--------|
| `j`/`k` or `↑`/`↓` | Navigate up/down |
| `Enter` | View role details |
| `/` | Search roles |
| `g`/`G` | Go to top/bottom |
| `r` | Refresh list |
| `q` | Quit application |
| `?` | Show help |

#### Detail View  
| Key | Action |
|-----|--------|
| `Tab`/`l` | Next tab |
| `Shift+Tab`/`h` | Previous tab |
| `j`/`k` | Scroll content or navigate policies |
| `g`/`G` | Go to top/bottom |
| `Enter` | View selected policy document (in Policies tab) |
| `Esc` | Back to list or previous view |

#### Policy Document View
| Key | Action |
|-----|--------|
| `j`/`k` | Scroll up/down |
| `g`/`G` | Go to top/bottom of document |
| `Esc` | Back to policies tab |

## Configuration

### AWS Credentials

a3s uses standard AWS credential resolution:

1. Command line flags (`-profile`, `-region`)
2. Environment variables:
   - `AWS_PROFILE`
   - `AWS_REGION` or `AWS_DEFAULT_REGION`
3. Shared credentials file (`~/.aws/credentials`)
4. IAM role (when running on EC2/ECS/Lambda)

### Required IAM Permissions

For the MVP (IAM role viewer with interactive policy document viewing), you need:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iam:ListRoles",
        "iam:GetRole",
        "iam:ListRoleTags",
        "iam:ListAttachedRolePolicies",
        "iam:ListRolePolicies",
        "iam:GetRolePolicy",
        "iam:GetPolicy",
        "iam:GetPolicyVersion",
        "sts:GetCallerIdentity"
      ],
      "Resource": "*"
    }
  ]
}
```

## Architecture

a3s is built with:
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework using Model-View-Update pattern
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/)** - AWS API interactions

### Project Structure

```
a3s/
├── cmd/a3s/           # Application entry point
├── internal/
│   ├── aws/           # AWS service layers
│   │   ├── client/    # AWS client management
│   │   └── iam/       # IAM service operations
│   ├── ui/            # UI components
│   │   ├── components/# List and detail views
│   │   └── styles/    # Lipgloss styling
│   └── model/         # Application state
└── docs/              # Documentation
```

## Roadmap

### Phase 1: MVP ✅
- [x] IAM role listing and viewing
- [x] Search and navigation
- [x] Profile/region switching

### Phase 2: Enhanced IAM
- [ ] IAM users view
- [ ] IAM policies view
- [ ] Permission boundary analysis
- [ ] Cross-account role assumptions

### Phase 3: Core AWS Services  
- [ ] EC2 instances
- [ ] S3 buckets
- [ ] Lambda functions
- [ ] RDS databases

### Phase 4: Advanced Features
- [ ] CloudFormation stacks
- [ ] ECS/EKS resources
- [ ] CloudWatch logs integration
- [ ] Resource tagging operations

## Development

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
# Install golangci-lint if not already installed
brew install golangci-lint

# Run linter
golangci-lint run
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [k9s](https://github.com/derailed/k9s) - the amazing Kubernetes TUI
- Built with [Charm](https://charm.sh/) libraries
- AWS SDK for Go team

## Support

- 🐛 Report bugs via [GitHub Issues](https://github.com/johnoct/a3s/issues)
- 💡 Request features via [GitHub Issues](https://github.com/johnoct/a3s/issues)
- 📧 Contact: [your-email@example.com]

---

Made with ❤️ for the AWS community
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/johnoct/a3s/internal/model"
	"golang.org/x/term"
)

func main() {
	var (
		profile = flag.String("profile", "", "AWS profile to use")
		region  = flag.String("region", "", "AWS region to use")
		help    = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	// Use environment variables if flags not provided
	if *profile == "" {
		*profile = os.Getenv("AWS_PROFILE")
	}
	if *region == "" {
		*region = os.Getenv("AWS_REGION")
		if *region == "" {
			*region = os.Getenv("AWS_DEFAULT_REGION")
		}
	}

	// Get initial terminal size
	width, height := 80, 24 // defaults
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width, height = w, h
	}

	app, err := model.NewAppWithSize(*profile, *region, width, height)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(app,
		tea.WithAltScreen(),
		// Removed tea.WithMouseCellMotion() to allow terminal text selection
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	fmt.Println(`a3s - AWS Terminal User Interface

Usage:
  a3s [flags]

Flags:
  -profile string   AWS profile to use (default: from environment)
  -region string    AWS region to use (default: from environment)
  -help            Show this help message

Environment Variables:
  AWS_PROFILE       Default AWS profile
  AWS_REGION        Default AWS region
  AWS_DEFAULT_REGION Alternative for AWS region

Keyboard Shortcuts:
  j/k or ↑/↓       Navigate up/down
  Enter            View role details
  /                Search roles
  Tab/Shift+Tab    Switch between tabs in detail view
  Esc              Go back
  q                Quit
  r                Refresh
  ?                Show help

Examples:
  a3s                           # Use default profile and region
  a3s -profile prod             # Use 'prod' profile
  a3s -region us-west-2         # Use specific region
  a3s -profile dev -region eu-west-1`)
}

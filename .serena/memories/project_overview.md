# Redmine CLI Project Overview

## Purpose
This is a Go CLI application for interacting with Redmine project management system. It provides commands to:
- Authenticate with Redmine instances
- List and view issues
- Manage profiles for different Redmine servers

## Tech Stack
- **Language**: Go 1.24.3
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Configuration**: YAML (gopkg.in/yaml.v3)

## Project Structure
```
redmine-cli/
├── main.go           # Entry point, calls cmd.Execute()
├── cmd/              # CLI commands
│   ├── root.go       # Root command setup
│   ├── auth.go       # Authentication commands
│   ├── issues.go     # Issue-related commands
│   └── profile.go    # Profile management commands
├── config/           # Configuration management
│   └── config.go     # Config struct and methods
├── client/           # Redmine API client
│   └── client.go     # HTTP client for Redmine API
├── go.mod            # Go module definition
└── go.sum            # Go module checksums
```

## Entry Point
The application is started via `main.go` which calls `cmd.Execute()` from the Cobra framework.
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

```bash
# Build the project
go build -o redmine

# Install dependencies and clean up
go mod tidy

# Test the built binary
./redmine --help

# Run all tests (currently no tests exist)
go test -v ./...
```

## Project Architecture

This is a Redmine CLI tool written in Go using the Cobra framework. The application follows a standard CLI architecture with clear separation of concerns:

### Core Components

- **main.go**: Entry point that delegates to cmd.Execute()
- **cmd/**: Contains all CLI command definitions using Cobra framework
  - **root.go**: Root command setup with global profile flag
  - **issues.go**: Issue management commands (list, show)
  - **profile.go**: Profile management commands (add, list, use, remove, show)
  - **auth.go**: Authentication commands (legacy, prefer profile commands)
- **config/**: Configuration management with YAML persistence
  - **config.go**: Profile-based configuration with ~/.redminecli/config storage
- **client/**: Redmine API client with HTTP communication
  - **client.go**: HTTP client with complete Redmine API models (Issue, Project, User, etc.)

### Configuration System

The application uses a profile-based configuration system:
- Configuration stored in `~/.redminecli/config` as YAML
- Each profile contains: name, Redmine URL, and API key
- Supports default profile and per-command profile override via `--profile` flag
- Profile management through `profile` commands (add, remove, use, list, show)

### API Client Architecture

The HTTP client (`client/client.go`) provides:
- Structured Redmine API models (Issue, Journal, Project, User, etc.)
- Authentication via X-Redmine-API-Key header
- JSON response parsing with proper error handling
- Support for pagination and filtering parameters
- Include parameters for fetching related data (e.g., journals for comments)

### Command Structure

Commands follow a hierarchical structure:
- `redmine issues list` - List issues with filtering options
- `redmine issues show <id>` - Show issue details with optional comments
- `redmine profile add/list/use/remove/show` - Profile management
- Global `--profile` flag for per-command profile selection

### Dependencies

- **github.com/spf13/cobra**: CLI framework for command structure
- **gopkg.in/yaml.v3**: YAML configuration file handling
- Standard library for HTTP client and JSON processing
- Uses mise.toml for Go version management (latest)

### Development Notes

- No test suite currently exists
- Binary excluded from git via .gitignore (redmine, redmine-*)
- Configuration files (.yaml/.yml) excluded from git for security
- API keys are masked in profile display output
- 30-second HTTP timeout for API requests
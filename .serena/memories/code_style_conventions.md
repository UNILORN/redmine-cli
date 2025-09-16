# Code Style and Conventions

## Go Conventions
Based on the existing codebase, this project follows standard Go conventions:

### Package Structure
- `main` package in root for entry point
- Functional packages: `cmd`, `config`, `client`
- Local imports use module name: `redmine-cli/cmd`

### Naming Conventions
- **Variables**: camelCase (e.g., `rootCmd`, `profileFlag`)
- **Functions**: PascalCase for exported (e.g., `Execute`, `NewClient`)
- **Types/Structs**: PascalCase (e.g., `Client`, `Config`, `Profile`)
- **Methods**: PascalCase with receiver (e.g., `(*Client).GetIssues`)

### File Organization
- One main concept per file
- Commands grouped in `cmd/` package
- API client logic in `client/` package
- Configuration handling in `config/` package

### Error Handling
- Standard Go error handling patterns
- Exit with status 1 on errors
- Print errors to stderr

### Dependencies
- Minimal external dependencies
- Use of established libraries (Cobra for CLI)
- Standard library preferred where possible

## Project-Specific Patterns
- Cobra command structure with init() functions
- YAML for configuration files
- HTTP client for API interactions
- Profile-based configuration system
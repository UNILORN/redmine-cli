# Task Completion Checklist

When completing development tasks on this Go project, run the following commands:

## Code Quality Checks
1. **Format code**: `go fmt ./...`
2. **Vet code**: `go vet ./...` 
3. **Build check**: `go build .`

## Testing (if tests exist)
4. **Run tests**: `go test ./...`

## Dependencies
5. **Clean dependencies**: `go mod tidy`

## Final Verification
6. **Test binary**: `./redmine --help` (after building)

## Optional but Recommended
- Check for unused imports
- Verify all exported functions are documented
- Ensure error handling follows Go conventions
- Run `go mod verify` to verify dependencies

## Before Committing
- All above checks pass
- Code follows project conventions
- Commit messages are descriptive
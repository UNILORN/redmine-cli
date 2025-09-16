# Development Commands

## Build and Run
```bash
# Build the application
go build -o redmine .

# Run directly
go run main.go

# Run with arguments
go run main.go --help
./redmine --help
```

## Go Commands
```bash
# Install dependencies
go mod tidy

# Update dependencies
go get -u

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run tests (if any exist)
go test ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o redmine-linux .
GOOS=windows GOARCH=amd64 go build -o redmine.exe .
```

## Git Commands (Darwin/macOS)
```bash
git status
git add .
git commit -m "message"
git push
```

## System Commands (Darwin/macOS)
```bash
ls -la           # List files
find . -name "*.go"  # Find Go files
grep -r "pattern" .  # Search in files
```
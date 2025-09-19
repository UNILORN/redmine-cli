.PHONY: build build-all clean test

# Default build for current platform
build:
	go build -o redmine .

# Build for all platforms
build-all: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/redmine-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o dist/redmine-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/redmine-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/redmine-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/redmine-windows-amd64.exe .

# Clean build artifacts
clean:
	rm -rf dist/
	rm -f redmine

# Run tests
test:
	go test ./...

# Install dependencies
install:
	go mod download
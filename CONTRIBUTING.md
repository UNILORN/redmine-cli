# Contributing to Redmine CLI

> **⚠️ このプロジェクトは開発を終了しました / This project is deprecated**
>
> より優れたRedmine CLIツールが開発されたため、このプロジェクトの開発は終了しました。
> 新しいツールへのコントリビュートをご検討ください：
> 
> **[kqns91/redmine-go](https://github.com/kqns91/redmine-go)**
>
> ---
> 
> Development of this project has stopped as a better Redmine CLI tool has been developed.
> Please consider contributing to the new tool instead:
>
> **[kqns91/redmine-go](https://github.com/kqns91/redmine-go)**
>
> ---

## Historical Information

The information below is kept for historical reference only.

## Commit Message Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning and changelog generation.

### Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **feat**: A new feature (triggers a minor version bump)
- **fix**: A bug fix (triggers a patch version bump)
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **chore**: Changes to the build process or auxiliary tools

### Breaking Changes

To trigger a major version bump, add `BREAKING CHANGE:` in the commit body or add `!` after the type:

```
feat!: remove support for legacy config format

BREAKING CHANGE: The old YAML configuration format is no longer supported.
Please migrate to the new JSON format.
```

### Examples

```bash
# Minor version bump (new feature)
feat: add tracker selection to issues add command

# Patch version bump (bug fix)
fix: resolve user email lookup in issues add

# No version bump
docs: update README with new command examples

# Major version bump (breaking change)
feat!: change CLI argument structure for better usability
```

## Release Process

1. All changes should be made via Pull Requests to the `main` branch
2. When PR is merged to `main`, GitHub Actions will:
   - Run tests
   - Build cross-platform binaries
   - Analyze commit messages
   - Generate changelog
   - Create a GitHub release with version tag
   - Upload binary assets to the release

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test
```

### Testing

Make sure to add tests for new features and run the test suite before submitting PRs:

```bash
go test ./...
```
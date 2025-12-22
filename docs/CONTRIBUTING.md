# Contributing to TaskHub

Thank you for your interest in contributing to TaskHub! This guide will help you get started with contributing to the project.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [How to Contribute](#how-to-contribute)
4. [Development Workflow](#development-workflow)
5. [Pull Request Process](#pull-request-process)
6. [Coding Standards](#coding-standards)
7. [Testing Guidelines](#testing-guidelines)
8. [Documentation](#documentation)
9. [Release Process](#release-process)
10. [Community](#community)

## Code of Conduct

### Our Pledge

We are committed to making participation in this project a harassment-free experience for everyone, regardless of level of experience, gender, gender identity and expression, sexual orientation, disability, personal appearance, body size, race, ethnicity, age, religion, or nationality.

### Our Standards

**Positive behavior includes:**
- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

**Unacceptable behavior includes:**
- The use of sexualized language or imagery
- Personal attacks or political attacks
- Trolling or insulting/derogatory comments
- Public or private harassment
- Publishing others' private information without explicit permission
- Any other conduct which could reasonably be considered inappropriate

### Enforcement

Project maintainers have the right and responsibility to remove, edit, or reject comments, commits, code, wiki edits, issues, and other contributions that are not aligned with this Code of Conduct. Project maintainers who do not follow the Code of Conduct may be removed from the project team.

## Getting Started

### Prerequisites

Before you start contributing, make sure you have:

- Go 1.25+ installed
- Git configured with your name and email
- A GitHub account
- Basic understanding of Go, PostgreSQL, and REST APIs

### Setup Your Development Environment

1. **Fork the repository**
   ```bash
   # Fork the repository on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/task-hub.git
   cd task-hub
   ```

2. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/original-org/task-hub.git
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Set up development environment**
   ```bash
   cp .env.example .env
   # Edit .env with your development configuration
   ```

5. **Start development services**
   ```bash
   docker-compose up -d
   ```

6. **Run the application**
   ```bash
   go run ./cmd/main.go
   ```

## How to Contribute

### Reporting Bugs

Before creating a bug report, please check:

1. **Existing issues** - Search to see if the bug has already been reported
2. **FAQ and documentation** - Check if it's covered in our docs

When creating a bug report, include:

- **Title**: Clear and descriptive
- **Description**: Detailed explanation of the bug
- **Steps to reproduce**: Exact steps to reproduce the issue
- **Expected behavior**: What you expected to happen
- **Actual behavior**: What actually happened
- **Environment**: OS, Go version, browser version, etc.
- **Screenshots**: If applicable, include screenshots
- **Additional context**: Any other relevant information

### Suggesting Features

Feature suggestions are welcome! When suggesting a feature:

1. **Check existing issues** - Make sure it hasn't been suggested already
2. **Use the feature request template** - Provide all requested information
3. **Explain the use case** - Why would this feature be useful?
4. **Consider implementation** - How might this be implemented?

### Contributing Code

We welcome code contributions! Here are ways you can contribute:

- **Bug fixes** - Help us fix bugs
- **New features** - Implement new functionality
- **Performance improvements** - Make the application faster
- **Documentation** - Improve our documentation
- **Tests** - Add or improve tests
- **Refactoring** - Improve code quality

## Development Workflow

### 1. Create a Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create a new branch
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Changes

- Follow the coding standards (see below)
- Write tests for your changes
- Update documentation as needed
- Commit your changes with clear messages

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linting
golangci-lint run

# Run integration tests
go test -tags=integration ./tests/...
```

### 4. Commit Your Changes

```bash
# Stage your changes
git add .

# Commit with a clear message
git commit -m "feat: add user authentication with JWT"

# Push to your fork
git push origin feature/your-feature-name
```

## Pull Request Process

### Before Submitting

1. **Test thoroughly** - Make sure all tests pass
2. **Update documentation** - Update relevant documentation
3. **Check formatting** - Run `go fmt ./...` and `goimports -w .`
4. **Run linters** - Fix any linting issues
5. **Rebase if needed** - Keep your branch up to date with main

### Creating a Pull Request

1. **Go to your fork on GitHub**
2. **Click "New Pull Request"**
3. **Select the correct branch** - Choose your feature branch
4. **Fill out the PR template** - Provide all requested information
5. **Link relevant issues** - Reference any related issues
6. **Request reviewers** - Tag relevant team members

### PR Template

```markdown
## Description
Brief description of the changes made.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Added new tests for new functionality

## Checklist
- [ ] My code follows the project's coding standards
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published in downstream modules
```

### Review Process

1. **Automated checks** - CI/CD will run tests and checks
2. **Code review** - Maintainers will review your code
3. **Feedback** - You may receive feedback for changes
4. **Approval** - Once approved, your PR will be merged

## Coding Standards

### Go Code Style

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and these additional guidelines:

#### Naming

```go
// Good - Clear, descriptive names
var userService UserService
var taskRepository TaskRepository

// Bad - Abbreviations, unclear names
var us UserService
var tr TaskRepository
```

#### Error Handling

```go
// Good - Handle errors immediately
user, err := userRepo.GetByID(ctx, userID)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// Bad - Ignore errors
user, _ := userRepo.GetByID(ctx, userID)
```

#### Documentation

```go
// UserService handles user-related business logic.
type UserService struct {
    userRepo UserRepository
    logger   *slog.Logger
}

// CreateUser creates a new user with the given request data.
// It validates the request, hashes the password, and stores the user.
// Returns the created user without sensitive information.
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Implementation
}
```

### Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(auth): add JWT token refresh functionality

fix(tasks): resolve deadline filtering bug

docs(api): update authentication documentation
```

## Testing Guidelines

### Test Structure

```go
func TestTaskService_CreateTask(t *testing.T) {
    // Arrange
    mockRepo := &MockTaskRepository{}
    service := NewTaskService(mockRepo, nil, nil)
    
    // Act
    task, err := service.CreateTask(context.Background(), req, userID)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedTitle, task.Title)
}
```

### Test Coverage

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test component interactions
- **End-to-end tests**: Test complete user workflows

**Target coverage:**
- Unit tests: 80%+
- Integration tests: 60%+

### Test Categories

```go
// Unit tests
func TestUser_Validate(t *testing.T) { ... }

// Integration tests
func TestTaskHandler_CreateTask_Integration(t *testing.T) { ... }

// End-to-end tests
func TestUserRegistration_E2E(t *testing.T) { ... }
```

## Documentation

### Code Documentation

- Document all public functions and types
- Include examples for complex functions
- Document configuration options

### API Documentation

- Update API.md for new endpoints
- Include request/response examples
- Document error codes

### README Updates

- Update README.md for new features
- Keep installation instructions current
- Update examples as needed

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. **Update version** in go.mod and other files
2. **Update CHANGELOG.md** with release notes
3. **Create release tag** on GitHub
4. **Build and test** release artifacts
5. **Deploy** to production
6. **Announce** release

### Changelog Format

```markdown
## [1.2.0] - 2024-01-15

### Added
- User authentication with JWT
- Task filtering by status and priority
- Email notifications for task reminders

### Changed
- Improved database connection pooling
- Updated dependencies

### Fixed
- Fixed deadline filtering bug
- Resolved memory leak in task service

### Security
- Updated JWT library to latest version
- Added rate limiting to API endpoints
```

## Community

### Getting Help

- **GitHub Issues**: For bug reports and feature requests
- **Discussions**: For questions and general discussion
- **Discord**: Join our community Discord server
- **Email**: Contact maintainers at maintainers@taskhub.dev

### Communication Channels

- **GitHub**: Primary development platform
- **Discord**: Real-time discussion and help
- **Blog**: Updates and announcements
- **Twitter**: News and quick updates

### Recognition

Contributors are recognized in several ways:

- **AUTHORS file**: List of all contributors
- **Release notes**: Acknowledgment in release notes
- **Contributor badge**: Special badge on GitHub
- **Community spotlight**: Featured in blog posts

## Ways to Contribute

### Code Contributions

- **Bug fixes**: Help us squash bugs
- **Features**: Implement new functionality
- **Performance**: Optimize existing code
- **Tests**: Improve test coverage
- **Documentation**: Enhance project documentation

### Non-Code Contributions

- **Bug triage**: Help identify and categorize issues
- **User support**: Help other users in discussions
- **Translation**: Translate documentation
- **Design**: Improve UI/UX design
- **Writing**: Blog posts, tutorials, case studies

### Financial Support

- **Sponsorship**: Support the project financially
- **Bounties**: Fund specific features or fixes
- **Donations**: One-time or recurring donations

## Recognition and Rewards

### Contributor Benefits

- **Recognition**: Listed in AUTHORS and release notes
- **Influence**: Help shape project direction
- **Learning**: Gain experience with modern Go development
- **Networking**: Connect with other developers
- **Portfolio**: Build your open source portfolio

### Top Contributors

We recognize our most active contributors with:

- **Core maintainer status**: For consistent, high-quality contributions
- **Special recognition**: Featured in project communications
- **Swag**: Project stickers, t-shirts, and other merchandise
- **Speaking opportunities**: Invitations to present about the project

## Getting Help

If you need help with contributing:

1. **Check the documentation** - Start with existing docs
2. **Search issues** - See if your question has been answered
3. **Ask in discussions** - Get help from the community
4. **Contact maintainers** - Reach out directly if needed

### Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [TaskHub Architecture](ARCHITECTURE.md)
- [TaskHub API Documentation](API.md)

---

Thank you for contributing to TaskHub! Your contributions help make this project better for everyone. We look forward to working with you!
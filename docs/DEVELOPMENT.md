# Development Guide

## Table of Contents

1. [Overview](#overview)
2. [Development Environment Setup](#development-environment-setup)
3. [Project Structure](#project-structure)
4. [Coding Standards](#coding-standards)
5. [Testing](#testing)
6. [Debugging](#debugging)
7. [Performance Profiling](#performance-profiling)
8. [Database Development](#database-development)
9. [API Development](#api-development)
10. [Frontend Development](#frontend-development)
11. [Common Development Tasks](#common-development-tasks)
12. [Troubleshooting](#troubleshooting)

## Overview

This guide provides comprehensive information for developers working on TaskHub. It covers setup procedures, coding standards, testing practices, and development workflows.

## Development Environment Setup

### Prerequisites

- **Go**: 1.25 or higher
- **PostgreSQL**: 16 or higher
- **NATS Server**: 2.9 or higher
- **Docker**: 20.10 or higher (optional)
- **Git**: 2.30 or higher
- **Task**: Task runner for development automation
- **Fyne**: Desktop UI framework (for desktop development)

### Quick Setup

```bash
# Clone the repository
git clone https://github.com/your-org/task-hub.git
cd task-hub

# Install Go dependencies
go mod download

# Install Task runner
go install github.com/go-task/task/v3/cmd/task@latest

# Set up development environment
task build

# Start development services
task run

# For desktop development
task run-desktop

# For web-only development
task run-web
```

### Manual Setup

#### 1. Install Go

```bash
# macOS (using Homebrew)
brew install go

# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Verify installation
go version
```

#### 2. Configure Go Environment

```bash
# Add to ~/.bashrc or ~/.zshrc
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export GO111MODULE=on

# Reload shell
source ~/.bashrc  # or source ~/.zshrc
```

#### 3. Install Development Tools

```bash
# Install useful Go tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golang/mock/gomock@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

#### 4. Database Setup

```bash
# Start PostgreSQL with Docker
docker run --name taskhub-postgres \
  -e POSTGRES_USER=taskhub \
  -e POSTGRES_PASSWORD=dev_password \
  -e POSTGRES_DB=taskhub \
  -p 5432:5432 \
  -d postgres:16-alpine

# Start NATS with Docker
docker run --name taskhub-nats \
  -p 4222:4222 \
  -p 8222:8222 \
  -d nats:2.9-alpine
```

#### 5. Environment Configuration

```bash
# Create development environment file
cat > .env << EOF
APP_ENV=development
SERVER_PORT=8080
SERVER_HOST=localhost

DB_HOST=localhost
DB_PORT=5432
DB_USER=taskhub
DB_PASSWORD=dev_password
DB_NAME=taskhub
DB_SSLMODE=disable

JWT_SECRET=dev_jwt_secret_key_at_least_32_characters
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h

NATS_URL=nats://localhost:4222
NATS_SUBJECT_PREFIX=taskhub

LOG_LEVEL=debug
LOG_FORMAT=text

CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
EOF
```

## Project Structure

```
task-hub/
├── cmd/                    # Application entry points
│   ├── main.go            # Web application
│   └── desktop/           # Desktop application
│       └── main.go        # Desktop app entry point
├── internal/              # Private application code
│   ├── app/               # Application services
│   │   ├── auth_service.go
│   │   ├── task_service.go
│   │   ├── user_service.go
│   │   └── notification_service.go
│   ├── desktop/           # Desktop UI components
│   │   ├── app.go         # Desktop app logic
│   │   └── theme.go       # Custom themes
│   ├── domains/           # Domain models and logic
│   │   ├── user/
│   │   │   ├── user.go
│   │   │   ├── user_test.go
│   │   │   └── repo/
│   │   │       └── user.go
│   │   ├── task/
│   │   │   ├── task.go
│   │   │   ├── task_test.go
│   │   │   └── repo/
│   │   │       └── task.go
│   │   └── notification/
│   │       ├── notification.go
│   │       └── repo/
│   │           └── notification.go
│   ├── gateway/          # External service gateways
│   │   └── gateway.go
│   └── handler/          # HTTP handlers
│       ├── auth_handler.go
│       ├── task_handler.go
│       ├── user_handler.go
│       └── web_handler.go
├── pkg/                   # Public library code
│   ├── base/
│   │   ├── entity/
│   │   │   ├── entity.go
│   │   │   └── entity_test.go
│   │   └── repo/
│   │       └── repo.go
│   ├── db/
│   │   └── db.go
│   ├── logger/
│   │   └── logger.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── nats/
│   │   ├── nats.go
│   │   └── nats_test.go
│   └── utils/
│       ├── utils.go
│       └── utils_test.go
├── web/                   # Web assets
│   ├── static/
│   │   └── css/
│   │       └── style.css
│   └── templates/
│       ├── base.html
│       ├── dashboard.html
│       ├── login.html
│       └── register.html
├── config/                # Configuration
│   ├── config.go
│   └── config_test.go
├── docs/                  # Documentation
├── scripts/               # Development scripts
├── migrations/            # Database migrations
├── tests/                 # Integration tests
├── .env.example           # Environment template
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Coding Standards

### 1. Go Formatting

```bash
# Format all Go files
go fmt ./...

# Use goimports for import organization
goimports -w .
```

### 2. Linting

```bash
# Run golangci-lint
golangci-lint run

# Run specific linters
golangci-lint run --enable=govet,goconst,gocritic
```

### 3. Naming Conventions

#### Packages
- Use short, lowercase names
- Avoid abbreviations unless widely known
- Examples: `user`, `task`, `auth`, `gateway`

#### Variables and Functions
- Use camelCase for local variables
- Use PascalCase for exported names
- Use descriptive names

```go
// Good
var userService UserService
var taskRepository TaskRepository

func CreateUser(req *CreateUserRequest) (*User, error)

// Bad
var us UserService
var tr TaskRepository
func CreateUser(req *CreateUserReq) (*User, error)
```

#### Constants
- Use UPPER_SNAKE_CASE
- Group related constants

```go
const (
    StatusTodo        TaskStatus = "todo"
    StatusInProgress  TaskStatus = "in_progress"
    StatusDone        TaskStatus = "done"
)
```

### 4. Error Handling

```go
// Good - Handle errors immediately
user, err := userRepo.GetByID(ctx, userID)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// Good - Wrap errors with context
func (s *TaskService) CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    // ...
}

// Bad - Ignore errors
user, _ := userRepo.GetByID(ctx, userID)
```

### 5. Documentation

```go
// Package user provides user management functionality.
package user

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

### 6. Interface Design

```go
// Define interfaces in the client package
type TaskRepository interface {
    Create(ctx context.Context, task *Task) error
    GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
    GetByUserID(ctx context.Context, userID uuid.UUID, filters TaskFilters) ([]*Task, error)
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// Implement interfaces in separate package
type taskRepository struct {
    gateway *gateway.Gateway
}

func (r *taskRepository) Create(ctx context.Context, task *Task) error {
    // Implementation
}
```

## Testing

### 1. Unit Testing

```go
// Example unit test
func TestTaskService_CreateTask(t *testing.T) {
    // Arrange
    mockRepo := &MockTaskRepository{}
    mockUserRepo := &MockUserRepository{}
    mockGateway := &MockGateway{}
    
    service := NewTaskService(mockRepo, mockUserRepo, mockGateway)
    
    userID := uuid.New()
    req := &CreateTaskRequest{
        Title:       "Test Task",
        Description: "Test Description",
        Priority:    PriorityHigh,
    }
    
    expectedTask := &Task{
        ID:          uuid.New(),
        Title:       req.Title,
        Description: req.Description,
        Status:      StatusTodo,
        Priority:    req.Priority,
        UserID:      userID,
    }
    
    // Setup mocks
    mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(&User{ID: userID}, nil)
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    mockGateway.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
    
    // Act
    task, err := service.CreateTask(context.Background(), req, userID)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedTask.Title, task.Title)
    assert.Equal(t, expectedTask.UserID, task.UserID)
}
```

### 2. Integration Testing

```go
// Example integration test
func TestTaskHandler_CreateTask_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Setup test server
    app := setupTestApp(t, db)
    server := httptest.NewServer(app)
    defer server.Close()
    
    // Create test user
    user := createTestUser(t, db)
    token := generateTestToken(t, user.ID)
    
    // Test request
    reqBody := map[string]interface{}{
        "title":       "Integration Test Task",
        "description": "Test Description",
        "priority":    "high",
    }
    
    resp, err := http.Post(
        server.URL+"/api/tasks",
        "application/json",
        bytes.NewBuffer(marshalJSON(t, reqBody)),
    )
    require.NoError(t, err)
    defer resp.Body.Close()
    
    // Assertions
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    
    var response struct {
        Success bool `json:"success"`
        Data    struct {
            ID string `json:"id"`
        } `json:"data"`
    }
    
    err = json.NewDecoder(resp.Body).Decode(&response)
    require.NoError(t, err)
    assert.True(t, response.Success)
    assert.NotEmpty(t, response.Data.ID)
}
```

### 3. Test Utilities

```go
// tests/testutils.go
package tests

import (
    "context"
    "testing"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "taskhub/internal/domains/user"
    "taskhub/internal/domains/task"
)

// CreateTestUser creates a user for testing
func CreateTestUser(t *testing.T, db *sql.DB) *user.User {
    t.Helper()
    
    userID := uuid.New()
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
    
    query := `
        INSERT INTO users (id, name, email, password, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `
    
    var u user.User
    err := db.QueryRow(query, userID, "Test User", "test@example.com", 
        string(hashedPassword), time.Now(), time.Now()).Scan(
        &u.ID, &u.CreatedAt, &u.UpdatedAt)
    
    if err != nil {
        t.Fatalf("Failed to create test user: %v", err)
    }
    
    u.Name = "Test User"
    u.Email = "test@example.com"
    
    return &u
}

// GenerateTestToken creates a JWT token for testing
func GenerateTestToken(t *testing.T, userID uuid.UUID) string {
    t.Helper()
    
    claims := jwt.MapClaims{
        "user_id": userID.String(),
        "email":   "test@example.com",
        "exp":     time.Now().Add(time.Hour).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte("test_jwt_secret"))
    if err != nil {
        t.Fatalf("Failed to generate test token: %v", err)
    }
    
    return tokenString
}

// SetupTestDB creates a test database
func SetupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    
    db, err := sql.Open("postgres", "postgres://taskhub:test@localhost/taskhub_test?sslmode=disable")
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // Run migrations
    err = runMigrations(db)
    if err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }
    
    return db
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(t *testing.T, db *sql.DB) {
    t.Helper()
    
    // Clean up all tables
    tables := []string{"tasks", "users"}
    for _, table := range tables {
        _, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
        if err != nil {
            t.Logf("Failed to truncate table %s: %v", table, err)
        }
    }
    
    db.Close()
}
```

### 4. Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific package tests
go test ./internal/domains/task/...

# Run integration tests
go test -tags=integration ./tests/...

# Run benchmarks
go test -bench=. ./...

# Run tests with race detection
go test -race ./...
```

## Debugging

### 1. Using Delve Debugger

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug main application
dlv debug ./cmd/main.go

# Debug tests
dlv test ./internal/domains/task/

# Debug with specific arguments
dlv debug -- -port 8080 -env development
```

### 2. Remote Debugging

```go
// Add debug endpoint in development
if config.AppEnv == "development" {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

### 3. Logging for Debug

```go
// Add detailed logging in development
if config.AppEnv == "development" {
    s.logger.Debug("Creating task",
        "user_id", userID,
        "title", req.Title,
        "priority", req.Priority,
    )
}
```

### 4. Debug Environment Variables

```bash
# Print all environment variables
env | grep TASKHUB

# Debug configuration
go run ./cmd/main.go -config-debug
```

## Performance Profiling

### 1. CPU Profiling

```go
// Add CPU profiling
func main() {
    if os.Getenv("PROFILE_CPU") == "true" {
        f, err := os.Create("cpu.prof")
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()
        
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    
    // Application code
}
```

### 2. Memory Profiling

```go
// Add memory profiling
if os.Getenv("PROFILE_MEMORY") == "true" {
    f, err := os.Create("mem.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    runtime.GC()
    pprof.WriteHeapProfile(f)
}
```

### 3. HTTP Profiling

```go
// Enable HTTP profiling in development
if config.AppEnv == "development" {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

### 4. Using Profiling Tools

```bash
# Analyze CPU profile
go tool pprof cpu.prof

# Analyze memory profile
go tool pprof mem.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
```

## Database Development

### 1. Database Migrations

```go
// migrations/001_create_users_table.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID
);

-- migrations/002_create_tasks_table.sql
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'todo',
    priority VARCHAR(50) NOT NULL DEFAULT 'medium',
    deadline TIMESTAMP WITH TIME ZONE,
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID
);
```

### 2. Running Migrations

```bash
# Using Goose
goose -dir ./migrations postgres "user=taskhub password=dev_password dbname=taskhub sslmode=disable" up

# Create new migration
goose -dir ./migrations create add_task_tags sql

# Rollback migration
goose -dir ./migrations postgres "user=taskhub password=dev_password dbname=taskhub sslmode=disable" down
```

### 3. Database Seeding

```go
// scripts/seed.go
package main

import (
    "database/sql"
    "log"
    
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    db, err := sql.Open("postgres", "postgres://taskhub:dev_password@localhost/taskhub?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Seed users
    seedUsers(db)
    
    // Seed tasks
    seedTasks(db)
}

func seedUsers(db *sql.DB) {
    users := []struct {
        Name     string
        Email    string
        Password string
    }{
        {"John Doe", "john@example.com", "password123"},
        {"Jane Smith", "jane@example.com", "password123"},
    }
    
    for _, user := range users {
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        
        query := `
            INSERT INTO users (id, name, email, password, created_at, updated_at)
            VALUES ($1, $2, $3, $4, NOW(), NOW())
            ON CONFLICT (email) DO NOTHING
        `
        
        _, err := db.Exec(query, uuid.New(), user.Name, user.Email, string(hashedPassword))
        if err != nil {
            log.Printf("Failed to seed user %s: %v", user.Email, err)
        }
    }
}
```

## API Development

### 1. Handler Structure

```go
// Example handler
type TaskHandler struct {
    taskService *app.TaskService
    logger      *slog.Logger
}

func NewTaskHandler(taskService *app.TaskService, logger *slog.Logger) *TaskHandler {
    return &TaskHandler{
        taskService: taskService,
        logger:      logger,
    }
}

// CreateTask handles task creation
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Parse request
    var req CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.handleError(w, err, "invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate request
    if err := req.Validate(); err != nil {
        h.handleError(w, err, "validation failed", http.StatusUnprocessableEntity)
        return
    }
    
    // Get user from context
    user, ok := getUserFromContext(ctx)
    if !ok {
        h.handleError(w, nil, "unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Create task
    task, err := h.taskService.CreateTask(ctx, &req, user.ID)
    if err != nil {
        h.handleError(w, err, "failed to create task", http.StatusInternalServerError)
        return
    }
    
    // Send response
    h.sendResponse(w, task, http.StatusCreated)
}
```

### 2. Request/Response Models

```go
// Request models
type CreateTaskRequest struct {
    Title       string     `json:"title" validate:"required,max=200"`
    Description string     `json:"description" validate:"max=1000"`
    Priority    TaskPriority `json:"priority" validate:"omitempty,oneof=low medium high"`
    Deadline    *time.Time `json:"deadline" validate:"omitempty"`
}

func (r *CreateTaskRequest) Validate() error {
    if r.Title == "" {
        return errors.New("title is required")
    }
    if len(r.Title) > 200 {
        return errors.New("title must be less than 200 characters")
    }
    if r.Description != nil && len(*r.Description) > 1000 {
        return errors.New("description must be less than 1000 characters")
    }
    return nil
}

// Response models
type TaskResponse struct {
    ID          uuid.UUID  `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      TaskStatus `json:"status"`
    Priority    TaskPriority `json:"priority"`
    Deadline    *time.Time `json:"deadline"`
    UserID      uuid.UUID  `json:"user_id"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

func NewTaskResponse(task *task.Task) *TaskResponse {
    return &TaskResponse{
        ID:          task.ID,
        Title:       task.Title,
        Description: task.Description,
        Status:      task.Status,
        Priority:    task.Priority,
        Deadline:    task.Deadline,
        UserID:      task.UserID,
        CreatedAt:   task.CreatedAt,
        UpdatedAt:   task.UpdatedAt,
    }
}
```

### 3. Middleware Development

```go
// Authentication middleware
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := extractTokenFromRequest(r)
            if token == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            claims, err := validateJWTToken(token, secret)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            ctx := context.WithValue(r.Context(), "user", claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Request logging middleware
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Create response writer wrapper to capture status code
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            
            next.ServeHTTP(wrapped, r)
            
            duration := time.Since(start)
            logger.Info("HTTP Request",
                "method", r.Method,
                "path", r.URL.Path,
                "status", wrapped.statusCode,
                "duration", duration,
                "remote_addr", r.RemoteAddr,
                "user_agent", r.UserAgent(),
            )
        })
    }
}
```

## Desktop Development

### 1. Fyne Desktop Framework

The desktop application uses the Fyne cross-platform framework for native UI development.

```go
// Desktop app structure
type DesktopApp struct {
    fyneApp     fyne.App
    authService *app.AuthService
    mainWindow  fyne.Window
    currentUser *app.LoginResponse
}

// Theme customization
type CustomTheme struct{}

func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
    switch name {
    case theme.ColorNamePrimary:
        return color.RGBA{R: 99, G: 102, B: 241, A: 255} // Modern indigo
    case theme.ColorNameBackground:
        return color.RGBA{R: 248, G: 250, B: 252, A: 255} // Light slate
    default:
        return theme.DefaultTheme().Color(name, variant)
    }
}
```

### 2. Running Desktop Application

```bash
# Development mode (services in Docker, app native)
task run-desktop

# Pure Docker deployment
docker compose -f docker-compose.desktop.yml up -d
```

### 3. Desktop UI Components

```go
// Create styled form with cards
form := container.NewVBox(
    titleLabel,
    subtitleLabel,
    container.NewPadded(
        container.NewVBox(
            emailEntry,
            passwordEntry,
        ),
    ),
    container.NewPadded(loginBtn),
)

card := widget.NewCard("", "", form)
return container.NewCenter(card)
```

## Frontend Development

### 1. HTMX Development

```html
<!-- templates/dashboard.html -->
{{template "base" .}}

<div class="container" hx-get="/api/tasks" hx-trigger="load">
    <div id="task-list">
        <!-- Tasks will be loaded here -->
    </div>
</div>

<!-- Task creation form -->
<form hx-post="/api/tasks" hx-target="#task-list" hx-swap="innerHTML">
    <div class="form-group">
        <label for="title">Title</label>
        <input type="text" id="title" name="title" required>
    </div>
    <div class="form-group">
        <label for="description">Description</label>
        <textarea id="description" name="description"></textarea>
    </div>
    <div class="form-group">
        <label for="priority">Priority</label>
        <select id="priority" name="priority">
            <option value="low">Low</option>
            <option value="medium" selected>Medium</option>
            <option value="high">High</option>
        </select>
    </div>
    <button type="submit">Create Task</button>
</form>
```

### 2. CSS Development

```css
/* web/static/css/style.css */
:root {
    --primary-color: #007bff;
    --secondary-color: #6c757d;
    --success-color: #28a745;
    --danger-color: #dc3545;
    --warning-color: #ffc107;
    --info-color: #17a2b8;
    --light-color: #f8f9fa;
    --dark-color: #343a40;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: var(--dark-color);
    background-color: var(--light-color);
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 15px;
}

.task-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    padding: 1rem;
    margin-bottom: 1rem;
    transition: transform 0.2s ease;
}

.task-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.task-priority-high {
    border-left: 4px solid var(--danger-color);
}

.task-priority-medium {
    border-left: 4px solid var(--warning-color);
}

.task-priority-low {
    border-left: 4px solid var(--success-color);
}
```

## Common Development Tasks

### 1. Adding New API Endpoint

```bash
# 1. Create request/response models
# internal/handler/task_handler.go

# 2. Add handler method
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

# 3. Add route
# cmd/main.go
router.HandleFunc("/api/tasks/{id}", h.UpdateTask).Methods("PUT")

# 4. Add tests
# internal/handler/task_handler_test.go
func TestTaskHandler_UpdateTask(t *testing.T) {
    // Test implementation
}
```

### 2. Adding New Domain Entity

```bash
# 1. Create domain entity
# internal/domains/comment/comment.go

# 2. Create repository interface
# internal/domains/comment/repo/comment.go

# 3. Create repository implementation
# internal/domains/comment/repo/comment_repository.go

# 4. Create service
# internal/app/comment_service.go

# 5. Create handler
# internal/handler/comment_handler.go

# 6. Add tests
# internal/domains/comment/comment_test.go
```

### 3. Database Schema Changes

```bash
# 1. Create migration
goose -dir ./migrations create add_comments_table sql

# 2. Write migration SQL
# migrations/003_add_comments_table.sql

# 3. Run migration
goose -dir ./migrations postgres "connection_string" up

# 4. Update domain models
# internal/domains/comment/comment.go

# 5. Update repository
# internal/domains/comment/repo/comment.go
```

## Troubleshooting

### Common Development Issues

#### 1. Build Errors

```bash
# Check Go version
go version

# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Check for missing imports
go build ./...
```

#### 2. Database Connection Issues

```bash
# Check PostgreSQL status
docker ps | grep postgres

# Test database connection
psql -h localhost -U taskhub -d taskhub

# Check database logs
docker logs taskhub-postgres
```

#### 3. Test Failures

```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestTaskService_CreateTask ./internal/app/

# Check test coverage
go test -cover ./...
```

#### 4. Performance Issues

```bash
# Profile CPU usage
go tool pprof http://localhost:8080/debug/pprof/profile

# Profile memory usage
go tool pprof http://localhost:8080/debug/pprof/heap

# Check goroutine leaks
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

### Development Tips

1. **Use Air for hot reloading** during development
2. **Keep dependencies minimal** and well-maintained
3. **Write tests as you code** rather than after
4. **Use meaningful commit messages**
5. **Regularly update dependencies** with `go get -u ./...`
6. **Use the debugger** instead of print statements
7. **Profile regularly** to catch performance issues early

---

This development guide provides comprehensive information for working with TaskHub. Follow these standards and practices to maintain code quality and ensure smooth collaboration.
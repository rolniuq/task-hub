# AGENTS.md - TaskHub Development Guide for AI Agents

## Build, Lint, and Test Commands

### Go Backend Commands
```bash
# Build commands
go build ./cmd/main.go                    # Build web server
go build ./cmd/desktop/main.go            # Build desktop app

# Test commands
go test ./...                              # Run all tests
go test ./internal/handler/...             # Run specific package tests
go test -v ./...                           # Run tests with verbose output
go test -cover ./...                       # Run tests with coverage
go test -race ./...                        # Run tests with race detection
go test -run TestTaskService_CreateTask    # Run single test by name

# Linting and formatting
golangci-lint run                          # Run linter
go fmt ./...                               # Format Go code
goimports -w .                             # Organize imports

# Development tools
go mod download                            # Download dependencies
go mod tidy                                # Clean up dependencies
air                                        # Hot reload (if installed)
```

### Frontend Commands (Vue.js)
```bash
cd web/frontend
npm install                                # Install dependencies
npm run dev                                # Start development server
npm run build                              # Build for production
npm run preview                            # Preview production build
```

### Docker Commands
```bash
task build                                 # Build Docker image
task run                                   # Run with Docker Compose
task run-desktop                           # Run desktop with Docker services
task down                                  # Stop all services
task clean                                 # Clean Docker resources
```

## Code Style Guidelines

### Go Code Style

#### Imports Organization
```go
// Standard library imports
import (
    "context"
    "encoding/json"
    "net/http"
    "time"
)

// Third-party imports
import (
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "go.uber.org/fx"
)

// Internal imports
import (
    "taskhub/internal/app"
    "taskhub/internal/domains/task"
    "taskhub/pkg/logger"
)
```

#### Naming Conventions
- **Packages**: Short, lowercase (`user`, `task`, `auth`)
- **Exported types**: PascalCase (`UserService`, `TaskRepository`)
- **Local variables**: camelCase (`userService`, `taskRepository`)
- **Constants**: UPPER_SNAKE_CASE (`STATUS_TODO`, `PRIORITY_HIGH`)
- **Interfaces**: Suffix with interface purpose (`Repository`, `Service`, `Handler`)

#### Error Handling
```go
// Always handle errors immediately
user, err := userRepo.GetByID(ctx, userID)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// Wrap errors with context
func (s *TaskService) CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    // ...
}
```

#### Function Structure
```go
// Document all exported functions
// CreateUser creates a new user with the given request data.
// It validates the request, hashes the password, and stores the user.
// Returns the created user without sensitive information.
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Implementation
}
```

#### Interface Design
```go
// Define interfaces in the client package
type TaskRepository interface {
    Create(ctx context.Context, task *Task) error
    GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Frontend Code Style (Vue.js)

#### Component Structure
```vue
<template>
  <!-- Component template -->
</template>

<script setup>
// Component logic with Composition API
import { ref, computed } from 'vue'

// Reactive state
const tasks = ref([])
const loading = ref(false)

// Computed properties
const completedTasks = computed(() => tasks.value.filter(t => t.completed))
</script>

<style scoped>
/* Component-specific styles */
</style>
```

#### Naming Conventions
- **Components**: PascalCase (`TaskList.vue`, `UserProfile.vue`)
- **Composables**: camelCase starting with "use" (`useAuth.js`, `useTasks.js`)
- **Props**: camelCase (`taskId`, `userName`)
- **Events**: kebab-case (`task-created`, `user-updated`)
- **CSS classes**: kebab-case (`task-list`, `user-profile`)

### Database Code Style

#### SQL Migrations
```sql
-- migrations/001_create_users_table.sql
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
```

#### SQLC Configuration
- Use camelCase for JSON tags
- Map UUID types to `github.com/google/uuid.UUID`
- Map timestamp types to `time.Time`
- Use interfaces for generated code

## Project Structure Conventions

### Directory Structure
```
task-hub/
├── cmd/                    # Application entry points
│   ├── main.go            # Web server
│   └── desktop/           # Desktop application
├── internal/              # Private application code
│   ├── app/               # Business logic services
│   ├── domains/           # Domain models and repositories
│   ├── handler/           # HTTP handlers
│   └── desktop/           # Desktop UI components
├── pkg/                   # Public library code
│   ├── logger/            # Logging utilities
│   ├── middleware/        # HTTP middleware
│   └── utils/             # Common utilities
├── web/                   # Frontend assets
│   ├── frontend/          # Vue.js application
│   ├── static/            # Static assets
│   └── templates/         # Server-side templates
└── config/                # Configuration
```

### Module Organization
- Use `fx` dependency injection framework
- Organize by domain (user, task, notification)
- Separate interfaces from implementations
- Keep handlers thin, logic in services

## Testing Patterns

### Unit Tests
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

### Integration Tests
```go
func TestTaskHandler_CreateTask_Integration(t *testing.T) {
    // Setup test database and server
    // Create test user and authenticate
    // Make HTTP request and verify response
}
```

### Test Utilities
- Use `testify/assert` for assertions
- Create test helpers in `tests/` package
- Use table-driven tests for multiple scenarios
- Mock external dependencies

## Common Development Tasks

### Adding New API Endpoint
1. Create request/response models
2. Add handler method
3. Register route in main.go
4. Write tests
5. Update documentation

### Adding New Domain Entity
1. Create domain model in `internal/domains/`
2. Create repository interface
3. Implement repository
4. Create service in `internal/app/`
5. Create handler in `internal/handler/`
6. Write tests

### Database Schema Changes
1. Create migration with Goose
2. Update domain models
3. Update repositories
4. Update services if needed
5. Test with existing data

## Error Handling Best Practices

### HTTP Error Responses
```go
func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": false,
        "error":   message,
    })
}
```

### Service Layer Errors
- Return wrapped errors with context
- Use custom error types for specific cases
- Log errors with appropriate levels
- Don't expose internal errors to clients

## Security Guidelines

### Authentication
- Use JWT tokens for API authentication
- Validate tokens in middleware
- Use bcrypt for password hashing
- Implement token refresh mechanism

### Input Validation
- Validate all user inputs
- Use struct tags for validation
- Sanitize data before storage
- Implement rate limiting

### Database Security
- Use parameterized queries
- Implement proper access controls
- Encrypt sensitive data
- Use connection pooling

## Performance Considerations

### Database Optimization
- Use indexes on frequently queried columns
- Implement connection pooling
- Use prepared statements
- Monitor query performance

### Caching Strategy
- Cache frequently accessed data
- Use appropriate cache expiration
- Implement cache invalidation
- Monitor cache hit rates

### API Performance
- Implement pagination for list endpoints
- Use appropriate HTTP status codes
- Compress responses when possible
- Monitor response times

## Environment Configuration

### Required Environment Variables
```bash
APP_ENV=development
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskhub
DB_PASSWORD=dev_password
DB_NAME=taskhub
JWT_SECRET=your_jwt_secret
NATS_URL=nats://localhost:4222
```

### Development vs Production
- Use different log levels (debug vs info)
- Enable hot reloading in development
- Use different database configurations
- Implement proper error reporting in production
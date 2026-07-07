# CONTEXT.md — Task Hub Deep Architecture & Logic

This file provides the complete architectural context, business logic, data flow, and design decisions of the Task Hub project. Useful for any developer or AI agent who needs to understand the full picture.

---

## 1. Architecture Overview

Task Hub uses a **layered architecture** with **Uber FX** dependency injection. The layers are:

```
┌──────────────┐  ┌──────────────┐
│   Web UI     │  │  Desktop UI  │  ← Presentation
│ (HTMX+HTML)  │  │ (Fyne v2)   │
└──────┬───────┘  └──────┬───────┘
       │                  │
       ▼                  ▼
┌──────────────────────────────────┐
│         Gateway (HTTP)           │  ← Interface Adapters
│  Routing · Middleware · Handlers │
└──────────────┬───────────────────┘
               │
               ▼
┌──────────────────────────────────┐
│       Application Services       │  ← Business Logic
│   Auth · Task                     │
└──────────────┬───────────────────┘
               │
               ▼
┌──────────────────────────────────┐
│         Domain Layer             │  ← Core Domain
│   Models · Repos (interfaces)    │
└──────────────┬───────────────────┘
               │
               ▼
┌──────────────────────────────────┐
│       Infrastructure (pkg/)      │  ← Frameworks & Drivers
│  DB · NATS · Logger · Middleware │
└──────────────────────────────────┘
```

### Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| DI Framework | Uber FX v1.24 | Lifecycle management, clean module composition |
| Database | PostgreSQL 16+ | UUID support, JSON, reliability |
| Messaging | NATS | Lightweight, high-performance, Go-native client |
| Auth | JWT (HS256) | Stateless, easy to implement, no session store needed |
| Web frontend | HTMX | Server-side rendering with dynamic interactions, no JS framework |
| Desktop | Fyne v2 | Cross-platform native UI, pure Go |
| Config | godotenv | Simple, 12-factor app style |
| Password hashing | bcrypt | Industry standard |

---

## 2. Module Dependency Graph (Uber FX Modules)

### Web App
```
ConfigModule ──► LoggerModule
     │                  │
     ├──────────────────┤
     │                  │
     ▼                  ▼
UserRepositoryModule  TaskRepositoryModule
     │                  │
     ▼                  ▼
AuthServiceModule    TaskServiceModule
     │                  │
     └─────┬────────────┘
           │
           ▼
     GatewayModule ◄── NatsModule
           │
           ▼
      startApp (Invoke)
```

### Desktop App
```
ConfigModule ──► LoggerModule ◄── NatsModule
     │                  │
     ├──────────────────┤
     │                  │
     ▼                  ▼
UserRepositoryModule  TaskRepositoryModule
     │                  │
     ▼                  ▼
AuthServiceModule    TaskServiceModule
     │
     ▼
DesktopApp (Provide)
     │
     ▼
RunDesktopApp (Invoke)
```

**All modules listed are actively wired.** `UserServiceModule`, `NotificationServiceModule`, and the entire `internal/domains/notification/` directory were previously removed as unused code.

---

## 3. Detailed Data Flow

### Task Creation (HTMX flow)
```
1. User fills form in dashboard.html
2. HTMX sends POST /api/tasks (form data)
3. AuthMiddleware: extracts Bearer token → validates JWT → injects user_id/email into context
4. handleTasks → switch on Method → TaskHandler.Create
5. TaskHandler.Create:
   a. Gets userID from context
   b. Parses form values (title, description, priority, deadline)
   c. Validates: title required
   d. Calls TaskService.CreateTask
6. TaskService.CreateTask:
   a. Creates Task domain object via task.NewTask()
   b. Sets status = "todo", creates BaseEntity with new UUID
   c. Calls TaskRepository.Create (INSERT INTO tasks ...)
7. Handler returns HTML snippet (task card) via renderTaskCard
8. HTMX swaps the response into the DOM
```

### Login (JSON API flow)
```
1. Client sends POST /api/auth/login with JSON body {email, password}
2. AuthHandler.Login:
   a. Checks HX-Request header → false → JSON path
   b. Decodes JSON body
   c. Calls AuthService.Login
3. AuthService.Login:
   a. UserRepository.FindByEmail → lookup user
   b. bcrypt.CompareHashAndPassword → validate password
   c. GenerateTokenPair → creates access token (15m) + refresh token (7d)
   d. Returns LoginResponse{User, Tokens}
4. Handler returns JSON response
```

### Task List with Filters
```
1. Dashboard loads → HTMX triggers GET /api/tasks
2. TaskHandler.List:
   a. Extracts query params: status, priority, deadline, search
   b. Creates ListTasksRequest
   c. Calls TaskService.ListTasks
3. TaskService.ListTasks:
   a. Creates TaskFilter from request params
   b. Calls TaskRepository.FindAll(filter)
4. TaskRepository.FindAll:
   a. Builds dynamic SQL: WHERE deleted_at IS NULL AND status=$1 AND priority=$2 ...
   b. Returns []*Task
5. Service applies client-side search filter (SQL doesn't have search clause)
6. Handler renders each task as HTML card via renderTaskCard
```

### NATS Event Flow (Removed)
The `NotificationService` (and entire `internal/domains/notification/` package) was removed as unused code.
NATS subjects (`task.created`, `task.updated`, `task.reminder`) were only defined there.
Currently, TaskService does NOT publish any NATS events after create/update.

To add event-driven notifications in the future:
1. Recreate `internal/app/notification_service.go` with NATS publishing
2. Wire `NotificationServiceModule` into `cmd/main.go`
3. Call `PublishXxx()` from `TaskService` after mutations

---

## 4. Configuration System

### Config struct (`config/config.go`)
```go
type Config struct {
    Port         string        // from env PORT
    NatsUrl      string        // from env NATS_URL
    JWTSecret    string        // from env JWT_SECRET
    DB           *DB           // nested struct
    ReadTimeout  time.Duration // hardcoded: 15s
    WriteTimeout time.Duration // hardcoded: 15s
    IdleTimeout  time.Duration // hardcoded: 60s
}
```

### .env.example
```
PORT=
NATS_URL=
JWT_SECRET=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
```

### Important
- `PORT` env var is used (NOT `SERVER_PORT`)
- `.env` file is REQUIRED — `godotenv.Load()` panics if file not found
- No `LOG_LEVEL`, `APP_ENV`, `SERVER_HOST` or similar env vars exist in code

---

## 5. Domain Model Details

### BaseEntity (`pkg/base/entity/entity.go`)
```go
type BaseEntity struct {
    Id        uuid.UUID
    CreatedAt time.Time
    CreatedBy uuid.UUID
    UpdateAt  *time.Time      // Note: "UpdateAt" not "UpdatedAt" (field naming quirk)
    UpdateBy  *uuid.UUID
    DeletedAt *time.Time
    DeletedBy *uuid.UUID
}
```

**Note**: The field is `UpdateAt` (not `UpdatedAt`). This inconsistency exists in the codebase.

### User (`internal/domains/user/user.go`)
```go
type User struct {
    BaseEntity
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"` // bcrypt hash, cleared before JSON response
}
```

### Task (`internal/domains/task/task.go`)
```go
type TaskStatus string   // "todo" | "in_progress" | "done"
type TaskPriority string // "low" | "medium" | "high"

type Task struct {
    BaseEntity
    Title       string
    Description string
    Status      TaskStatus     // defaults to "todo"
    Priority    TaskPriority   // defaults to "medium"
    Deadline    *time.Time
    UserID      uuid.UUID
}
```

### TaskFilter (`internal/domains/task/task.go`)
```go
type TaskFilter struct {
    Status   *TaskStatus
    Priority *TaskPriority
    UserID   *uuid.UUID
    Deadline *time.Time
}
```

---

## 6. Database Schema

### Table: `users`
```sql
id         UUID PRIMARY KEY
name       VARCHAR(255) NOT NULL
email      VARCHAR(255) UNIQUE NOT NULL
password   VARCHAR(255) NOT NULL
created_at TIMESTAMP NOT NULL DEFAULT NOW()
updated_at TIMESTAMP
updated_by UUID
deleted_at TIMESTAMP
deleted_by UUID
```

### Table: `tasks`
```sql
id          UUID PRIMARY KEY
title       VARCHAR(255) NOT NULL
description TEXT
status      VARCHAR(20) NOT NULL DEFAULT 'todo'
priority    VARCHAR(20) NOT NULL DEFAULT 'medium'
deadline    TIMESTAMP
user_id     UUID NOT NULL
created_at  TIMESTAMP NOT NULL DEFAULT NOW()
created_by  UUID NOT NULL
updated_at  TIMESTAMP
updated_by  UUID
deleted_at  TIMESTAMP
deleted_by  UUID
```

### Indexes
- `idx_tasks_user_id`, `idx_tasks_status`, `idx_tasks_priority`
- `idx_tasks_deadline`, `idx_tasks_created_at`
- `idx_users_email`
- `idx_tasks_deleted_at`, `idx_users_deleted_at`

---

## 7. Repository Pattern

### UserRepository (`internal/domains/user/repo/user.go`)
| Method | SQL | Description |
|--------|-----|-------------|
| `Create` | INSERT INTO users | Creates user with UUID, name, email, hashed password |
| `FindByEmail` | SELECT WHERE email=$1 | Looks up user by email |
| `FindById` | SELECT WHERE id=$1 | Looks up user by UUID string |
| `Update` | UPDATE users SET name, email | Updates user profile |
| `Delete` | DELETE FROM users | Hard delete (no soft delete for users) |

### TaskRepository (`internal/domains/task/repo/task.go`)
| Method | SQL | Description |
|--------|-----|-------------|
| `Create` | INSERT INTO tasks | Creates task with all fields |
| `UpdateById` | UPDATE tasks SET ... WHERE id | Updates title, desc, status, priority, deadline |
| `FindById` | SELECT ... WHERE id AND deleted_at IS NULL | Soft-delete aware |
| `FindAll` | SELECT ... WHERE conditions + ORDER BY created_at DESC | Dynamic WHERE builder |
| `FindByUserId` | Delegates to FindAll with user_id filter | Convenience method |
| `DeleteById` | UPDATE SET deleted_at=NOW(), deleted_by | Soft delete |
| `MarkAsCompleted` | UPDATE SET status='done' | Status switch |
| `FindTasksNearDeadline` | SELECT ... WHERE deadline <= NOW() + interval | Reminder logic |

### Key SQL patterns:
- **Soft delete**: `deleted_at IS NULL` on all SELECT queries
- **Dynamic filtering**: String concatenation of WHERE clauses with parameterized args
- **NULL handling**: Uses `sql.NullTime` and `sql.NullString` for nullable fields

---

## 8. Authentication & Authorization

### Token Generation
```
Access Token:
  - Algorithm: HS256
  - Claims: {user_id, email, exp(15m), iat, iss:"taskhub"}
  - Signed with: JWT_SECRET

Refresh Token:
  - Algorithm: HS256
  - Claims: {user_id, exp(7d), iat, iss:"taskhub"}
  - Signed with: JWT_SECRET
```

### Auth Flow
1. **Register**: POST /api/auth/register → bcrypt hash → INSERT user → return user (password blanked)
2. **Login**: POST /api/auth/login → bcrypt compare → generate token pair → return tokens + user
3. **Validate**: AuthMiddleware reads `Authorization: Bearer <token>` → parse JWT → inject user_id/email into context
4. **Refresh**: POST /api/auth/refresh → validate refresh token → lookup user → generate new pair
5. **Logout**: POST /api/auth/logout → clear cookies

### Known Gap
AuthMiddleware only checks `Authorization: Bearer` header. It does NOT extract tokens from cookies (`access_token` cookie). The HTMX web app sets cookies on login but the middleware can't authenticate from cookies. This means:
- Web pages that require auth (like `/dashboard`) will fail if the browser sends the JWT only in cookies
- API calls from the web page include the Bearer header via HTMX configuration

---

## 9. HTMX Integration Pattern

The project uses a dual-response pattern:

```go
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    isHTMX := r.Header.Get("HX-Request") == "true"

    if isHTMX {
        // Read form values instead of JSON
        // Return HTML snippets
        // Use HX-Redirect header for navigation
    } else {
        // Read JSON body
        // Return JSON response
    }
}
```

This allows the same endpoint to:
- Serve HTMX-enhanced web pages (HTML responses)
- Serve API clients (JSON responses)

### Web Templates
- `base.html` — Layout with HTMX loaded from CDN (`unpkg.com/htmx.org@1.9.10`)
- `login.html` — Form posts to `/api/auth/login` via HTMX
- `register.html` — Form posts to `/api/auth/register` via HTMX
- `dashboard.html` — Task list, filters, search, create/edit modal

### Dashboard.js features (inline in dashboard.html)
- Task search with debounce (500ms)
- Status/priority filter dropdowns
- Create/edit task modal (`<dialog>` element)
- Task statistics (total, todo, in_progress, done)
- HTMX event handlers for refresh after mutations

---

## 10. Desktop App Architecture

### App Lifecycle
1. `cmd/desktop/main.go` creates Uber FX app with modules
2. `desktop.NewApp` creates Fyne app with CustomTheme
3. `desktop.RunDesktopApp` creates window, shows login screen, runs Fyne event loop
4. On successful login: stores `currentUser`, switches to dashboard
5. On logout: clears `currentUser`, switches back to login

### Desktop Screens
- **Login**: Email + password fields, login button, register link
- **Register**: Name + email + password + confirm password, register button, back link
- **Dashboard**: Welcome message, user info card, logout button (no task CRUD yet)

### CustomTheme (`internal/desktop/theme.go`)
- Primary: Indigo `#6366F1`
- Background: Dark slate `#111827` / Light slate `#F8FAF8`
- Text: System colors based on variant
- Input: White background
- Custom font sizes (16 body, 24 heading, 20 subheading, 12 caption)

---

## 11. NATS Messaging System

### Connection
- URL from `NATS_URL` env var (defaults to `nats://localhost:4222`)
- Connected in `pkg/nats/nats.go` via `nats.Connect()`
- Returns `nil` if connection fails (no retry)
- Thread-safe via NATS Go client

### Pub/Sub Interface
```go
type Nats struct {
    conn   *nats.Conn
    logger *logger.Logger
}

func (n *Nats) Publish(subject string, data []byte) error
func (n *Nats) Subscribe(subject string, handler func([]byte)) error
func (n *Nats) QueueSubscribe(subject, queue string, handler func([]byte)) error
func (n *Nats) IsConnected() bool
```

### Note
The NotificationService and all NATS event subjects (`task.created`, `task.updated`, `task.reminder`) were removed as unused code. The NATS connection in `pkg/nats` remains for future use.

---

## 12. Error Handling Patterns

### Sentinel Errors
```go
// Auth service
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserAlreadyExists  = errors.New("user already exists")
var ErrInvalidToken       = errors.New("invalid token")
var ErrTokenExpired       = errors.New("token expired")

// Task service
var ErrTaskNotFound = errors.New("task not found")
var ErrUnauthorized = errors.New("unauthorized")
```

### Handler error patterns
- JSON API: `writeError(w, statusCode, message)` → `{"error": "message"}`
- HTMX: `writeHTMXError(w, message)` → `<div class="alert alert-error shake">message</div>`
- Success HTMX: `writeHTMXSuccess(w, message, redirectURL)` → sets `HX-Redirect` header

---

## 13. Testing Structure

### Test files
```
internal/app/auth_service_test.go
internal/app/task_service_test.go
internal/app/notification_service_test.go
internal/handler/auth_handler_test.go
internal/handler/task_handler_test.go
internal/domains/task/task_test.go
internal/domains/user/user_test.go
pkg/base/entity/entity_test.go
pkg/middleware/auth_test.go
pkg/nats/nats_test.go
pkg/utils/utils_test.go
config/config_test.go
```

### Testing framework
- **testify** (`github.com/stretchr/testify`) for assertions
- Standard Go `testing` package
- No mock framework (no gomock/mockery) — tests use real implementations or simple stubs

---

## 14. Docker Deployment

### docker-compose.yml
```
Services:
  - db: postgres:16.3-alpine3.20 (port 5432)
  - nats: nats (port 4222)
  - task-hub: built from Dockerfile (port 8080)
```

### Dockerfile (web)
- Multi-stage: builder (golang:1.25-alpine) → runtime (alpine)
- CGO_ENABLED=0 for static binary
- Exposes port 8080

### Dockerfile.desktop
- Multi-stage with buildx support (`--platform=$BUILDPLATFORM` / `$TARGETPLATFORM`)
- CGO_ENABLED=0
- Note: Desktop in Docker is headless; mainly for CI/testing

### docker-compose.desktop.yml
```
Services:
  - task-hub-desktop: built from Dockerfile.desktop
  - db: postgres
  - nats: nats
```

---

## 15. Known Issues & TODOs

1. **Cookie auth gap**: AuthMiddleware doesn't read `access_token` cookie → HTMX web pages can't authenticate
2. **Field naming inconsistency**: `UpdateAt` vs `UpdatedAt` (field name in BaseEntity is `UpdateAt`)
3. **No validation library**: Simple string checks; no input validation beyond empty checks
4. **No error wrapping**: Services return sentinel errors, handlers check with `==`, no `errors.Is`/`errors.As`
5. **Desktop task CRUD missing**: Desktop app has no task creation/editing UI yet
6. **Dockerfile uses Go 1.25**: Dockerfile says `golang:1.25-alpine` but go.mod says `go 1.24.0` — works because 1.25 is backward compatible but should be consistent
7. **No graceful shutdown**: No OS signal handling (SIGINT/SIGTERM) in web or desktop app
8. **Search filtering in application layer**: Text search is done in-memory after fetching all tasks, not at the database level

---

## 16. Quick Reference

### Common Operations

**Add a new dependency to FX:**
```go
// 1. Create your module
var MyModule = fx.Module("my-module", fx.Provide(NewMyService))

// 2. Add to cmd/main.go's fx.New()
app.MyModule,
```

**Add a new API endpoint:**
```go
// 1. Create handler method in internal/handler/
// 2. Add route in internal/gateway/gateway.go's Start() method
mux.HandleFunc("/api/foo", h.myHandler)
// 3. If auth required: wrap with middleware
mux.Handle("/api/foo", g.authMiddleware.Authenticate(http.HandlerFunc(h.myHandler)))
```

**Add a new domain entity:**
```go
// 1. Create internal/domains/<entity>/<entity>.go (model)
// 2. Create internal/domains/<entity>/repo/<entity>.go (repository)
// 3. Create internal/app/<entity>_service.go (business logic)
// 4. Create internal/handler/<entity>_handler.go (HTTP handlers)
// 5. Wire modules in cmd/main.go
```

### Environment Variables Required
| Variable | Example | Used In |
|----------|---------|---------|
| `PORT` | `8080` | config.go → Gateway |
| `NATS_URL` | `nats://localhost:4222` | config.go → Nats |
| `JWT_SECRET` | `min-32-chars-secret!!` | config.go → AuthService |
| `DB_HOST` | `localhost` | config.go → DB |
| `DB_PORT` | `5432` | config.go → DB |
| `DB_USER` | `taskhub` | config.go → DB |
| `DB_PASSWORD` | `...` | config.go → DB |
| `DB_NAME` | `taskhub` | config.go → DB |

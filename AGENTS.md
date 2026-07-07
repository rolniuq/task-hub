# AGENTS.md — Task Hub Rules for AI Agents

This file defines how AI agents should understand, navigate, and modify the Task Hub project.

---

## 1. Project Identity

| Attribute | Value |
|-----------|-------|
| **Name** | Task Hub |
| **Module** | `taskhub` (Go module) |
| **Go version** | `1.24.0` (`go.mod`) |
| **Repo** | `github.com/rolniuq/task-hub` |
| **Description** | Task management app: web (Go + HTMX) + desktop (Fyne v2) + PostgreSQL + NATS |

## 2. Documentation Map

| File | Purpose | Audience |
|------|---------|----------|
| `AGENTS.md` | Rules & quick reference for AI agents | AI agents |
| `CONTEXT.md` | Deep architecture, data flow, known issues | Developers & AI |
| `README.md` | Project overview, quick start | End users |
| `docs/README.md` | Secondary overview (needs cleanup) | End users |
| `docs/API.md` | API documentation (aspirational — does NOT match actual code) | API consumers |
| `docs/ARCHITECTURE.md` | Architecture docs (aspirational — does NOT match actual code) | Architects |
| `docs/DEVELOPMENT.md` | Development guide (aspirational — does NOT match actual code) | Developers |
| `docs/DEPLOYMENT.md` | Deployment guide (aspirational — does NOT match actual code) | DevOps |
| `docs/CONTRIBUTING.md` | Contribution guide (aspirational — does NOT match actual code) | Contributors |
| `install.sh` | macOS installer script | End users |

> **⚠️ CRITICAL**: The files in `docs/` are **aspirational** — they describe what the project *could become*, not what it *currently is*. Always verify against the actual code in `internal/`, `pkg/`, `cmd/`, and `config/`. The **source of truth** is the codebase itself.

## 3. Architecture Rules

### 3.1 Layer Rules

```
cmd/          → Entry points (thin — only FX wiring)
internal/
  app/        → Business logic services
  desktop/    → Fyne UI (desktop only)
  domains/    → Domain models + repository implementations
  gateway/    → HTTP server, routing, middleware setup
  handler/    → HTTP handlers (JSON + HTMX dual responses)
pkg/          → Reusable infrastructure code
config/       → Config loading from env vars
web/          → Static assets + HTML templates
```

**Rules:**
- `cmd/` files must ONLY wire FX modules and call `fx.Invoke` — no business logic
- `internal/app/` services must NOT import `net/http`, `handler`, `gateway`, or `desktop`
- `internal/domains/` must NOT import `internal/app/` or `internal/handler/` or `internal/gateway/`
- `pkg/` must NOT import any `internal/` packages
- `internal/handler/` must ONLY depend on `internal/app/` services, never on `internal/domains/` directly

### 3.2 FX Module Naming Convention

Every module follows this pattern:
```go
var XxxModule = fx.Module("xxx", fx.Provide(NewXxx))
```

| Module variable | Package | Provides |
|----------------|---------|----------|
| `config.ConfigModule` | `config` | `*config.Config` |
| `logger.LoggerModule` | `pkg/logger` | `*logger.Logger` |
| `userrepo.UserRepositoryModule` | `internal/domains/user/repo` | `*repo.UserRepository` |
| `taskrepo.TaskRepositoryModule` | `internal/domains/task/repo` | `*repo.TaskRepository` |
| `app.AuthServiceModule` | `internal/app` | `*app.AuthService` |
| `app.TaskServiceModule` | `internal/app` | `*app.TaskService` |
| `nats.NatsModule` | `pkg/nats` | `*nats.Nats` |
| `gateway.GatewayModule` | `internal/gateway` | `*gateway.Gateway` |

**All modules listed are actively wired.** `UserServiceModule`, `NotificationServiceModule`, and the entire `internal/domains/notification/` directory were removed as unused code.

### 3.3 Adding a New Feature

To add a new domain entity (e.g. "Comment"):

```
Step 1: internal/domains/comment/comment.go         — entity + interfaces
Step 2: internal/domains/comment/repo/comment.go     — SQL repository
Step 3: internal/app/comment_service.go              — business logic + FX module
Step 4: internal/handler/comment_handler.go          — HTTP handlers
Step 5: internal/gateway/gateway.go                  — register routes
Step 6: cmd/main.go                                  — wire new FX module
Step 7: .init/01-init.sql                            — add DDL
Step 8: internal/domains/comment/comment_test.go     — tests
Step 9: internal/handler/comment_handler_test.go     — handler tests
```

## 4. Coding Rules

### 4.1 HTMX Dual Response Pattern

ALL handlers must check `HX-Request` header and return appropriate responses:

```go
func (h *XxxHandler) DoSomething(w http.ResponseWriter, r *http.Request) {
    isHTMX := r.Header.Get("HX-Request") == "true"

    if isHTMX {
        // 1. Read form values (not JSON)
        // 2. Return HTML snippet
        // 3. Use HX-Redirect for navigation
        // 4. Use HX-Trigger for events
    } else {
        // 1. Read JSON body
        // 2. Return JSON response
        // 3. Return standard HTTP status codes
    }
}
```

Helper functions (in `internal/handler/auth_handler.go`):
- `writeJSON(w, status, data)` — JSON response
- `writeError(w, status, message)` — JSON error
- `writeHTMXError(w, message)` — HTML error alert
- `writeHTMXSuccess(w, message, redirectURL)` — HTML success + optional redirect

### 4.2 Error Handling

Use sentinel errors defined as package-level vars:

```go
// Define in service package
var ErrXxx = errors.New("xxx error")

// Check with == (not errors.Is)
if err == app.ErrXxx { ... }
```

**Do NOT** use `errors.Is()` or `errors.As()` — the codebase uses `==` comparison everywhere.

### 4.3 Database Access Pattern

```go
// Repository reads config and creates its own DB connection
func NewXxxRepository(config *config.Config, logger *logger.Logger) *XxxRepository {
    conn := db.NewDB(config).GetConnection()
    return &XxxRepository{conn: conn, logger: logger}
}
```

**Rules:**
- Repositories create their own DB connections via `db.NewDB(config).GetConnection()`
- Soft-deleted records: always check `WHERE deleted_at IS NULL`
- UUID primary keys, no auto-increment
- Use `sql.NullTime` / `sql.NullString` for nullable columns

### 4.4 Auth Middleware Pattern

```go
// middleware/auth.go
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ONLY checks Authorization: Bearer <token> header
        // Does NOT read cookies!
        // Injects user_id + email into context
    })
}

// Extract user info:
userID := middleware.GetUserIDFromContext(r.Context())
email := middleware.GetEmailFromContext(r.Context())
```

### 4.5 Task Ownership Check

Every task operation must check ownership server-side:

```go
if existingTask.UserID != userID {
    return nil, app.ErrUnauthorized
}
```

### 4.6 Soft Delete

Tasks use soft delete. The `DELETE` endpoint runs:
```sql
UPDATE tasks SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2
```

All SELECT queries on tasks include: `WHERE deleted_at IS NULL`

### 4.7 Password Handling

Passwords are bcrypt-hashed. The `Password` field is cleared (set to `""`) before returning user data in any response.

## 5. Key Gaps & Pitfalls

### Known Issues to NEVER introduce by accident:

1. **Go version**: Must stay `1.24.0` — Dockerfiles may say `1.25-alpine` but go.mod is `1.24.0`
2. **Config key**: `PORT` env var, NOT `SERVER_PORT`
3. **`.env` required**: `godotenv.Load()` panics if `.env` missing
4. **Cookie auth gap**: AuthMiddleware only reads `Authorization: Bearer` header, NOT cookies. HTMX pages set cookies but can't authenticate from them — fix this before adding cookie-based features
5. **Field naming**: `BaseEntity.UpdateAt` (not `UpdatedAt`) — this is inconsistent with SQL column `updated_at`
6. **NotificationService not wired**: The event publishing code exists but is never called. TaskService doesn't publish events
7. **UserServiceModule unused**: Defined but never imported
8. **Docker Go version**: Dockerfile uses `golang:1.25-alpine` while go.mod says `1.24.0` — works (1.25 builds 1.24 code) but inconsistent
9. **Search is in-memory**: `ListTasks` does text search client-side after fetching all tasks, not at DB level
10. **No graceful shutdown**: No OS signal handling (SIGINT/SIGTERM) in web or desktop app

## 6. Testing Rules

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Specific test
go test -v -run TestXxx ./internal/app/

# Race detection
go test -race ./...
```

- Use `testify/assert` for assertions
- No mock framework (gomock/mockery) — use simple stubs or real implementations
- Tests use real services, not mocked interfaces
- Integration tests (build tag `integration`) connect to real PostgreSQL/NATS

## 7. Environment Variables

| Env Var | Required | Used In |
|---------|----------|---------|
| `PORT` | Yes | `config.go` → Gateway HTTP server |
| `NATS_URL` | Yes | `config.go` → NATS connection |
| `JWT_SECRET` | Yes | `config.go` → JWT signing/validation |
| `DB_HOST` | Yes | `config.go` → Postgres DSN |
| `DB_PORT` | Yes | `config.go` → Postgres DSN |
| `DB_USER` | Yes | `config.go` → Postgres DSN |
| `DB_PASSWORD` | Yes | `config.go` → Postgres DSN |
| `DB_NAME` | Yes | `config.go` → Postgres DSN |

Hardcoded config: `ReadTimeout=15s`, `WriteTimeout=15s`, `IdleTimeout=60s`

## 8. Common File Editing Patterns

### Add a route
```go
// internal/gateway/gateway.go — in Start() method
mux.HandleFunc("/api/foo", g.myHandler.HandleFoo)                 // no auth
mux.Handle("/api/foo", g.authMiddleware.Authenticate(...))         // with auth
mux.Handle("/api/foo/", g.authMiddleware.Authenticate(...))        // with sub-paths
```

### Add a new FX module
```go
// 1. Define module + provider
var MyModule = fx.Module("my", fx.Provide(NewMyService))

// 2. Wire in cmd/main.go
app.MyModule,
```

### Add a NATS subject
```go
// In internal/app/notification_service.go
const SubjectXxx = "xxx.event"

// Publish
func (s *NotificationService) PublishXxx(...) error {
    return s.nats.Publish(SubjectXxx, data)
}
```

## 9. Desktop App Specifics

- Uses Fyne v2 (`fyne.io/fyne/v2`)
- CustomTheme with indigo primary (`#6366F1`)
- Icon generated programmatically (see `internal/desktop/icon.go`)
- Runs natively (not in Docker) on macOS — needs GUI
- Services (DB, NATS) run in Docker, app runs natively
- `task run-desktop` → `docker compose up -d db nats` + `go run cmd/desktop/main.go`

## 10. Workflow & CI Rules

- GitHub Actions in `.github/workflows/`
- CI: lint → test → build → docker → security scan → release → deploy
- `golangci-lint` config in `.golangci.yml` (no `version` field — incompatible with v1.64+)
- Go version in workflows must match `go.mod`: `'1.24'`
- Gosec SARIF output goes to `${{ runner.temp }}/results.sarif`

## 11. Quick Reference Commands

```bash
# Build
go build -o task-hub ./cmd/main.go
go build -o task-hub-desktop ./cmd/desktop/main.go

# Test
go test ./...
go test -race ./...
go test -cover ./...
go test -v -run TestAuthService ./internal/app/

# Run (dev)
docker compose up -d db nats   # start dependencies
go run cmd/main.go              # web app
go run cmd/desktop/main.go      # desktop app

# Taskfile
task run-web        # docker compose up -d
task run-desktop    # docker compose up -d db nats + go run desktop
task build          # docker build
task clean          # docker compose down -v --remove-orphans

# Docker
docker compose up -d
docker compose down -v --remove-orphans
```

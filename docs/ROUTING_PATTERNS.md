# Route Registration Patterns

This document explains three different patterns to replace manual `mux.Handle()` calls in `gateway.go`.

## Current Problem

Your `gateway.go` has many manual route registrations:
```go
mux.HandleFunc("/api/auth/register", g.authHandler.Register)
mux.HandleFunc("/api/auth/login", g.authHandler.Login)
mux.Handle("/api/auth/refresh", g.authMiddleware.Authenticate(http.HandlerFunc(g.authHandler.RefreshToken)))
// ... 15+ more lines
```

## Pattern 1: Interface-based Auto-Registration (Recommended) ✅

**How it works:**
- Handlers implement a `RouteRegistrar` interface with a `Routes()` method
- fx (dependency injection) automatically collects all handlers that implement the interface
- Gateway iterates through registered routes and sets them up

**Pros:**
- ✅ No external tools needed
- ✅ Works seamlessly with existing fx setup
- ✅ Type-safe at compile time
- ✅ Easy to test and mock
- ✅ No magic - explicit route definitions

**Cons:**
- ❌ Need to add `Routes()` method to each handler
- ❌ Slightly more boilerplate than code generation

**To implement:**
1. Add `Routes()` method to AuthHandler, TaskHandler, WebHandler (see routes_example.go)
2. Replace gateway.go with gateway_refactored.go
3. That's it - routes auto-register via fx

---

## Pattern 2: File-based Routing (Next.js style)

**How it works:**
- Routes are defined by file structure
- `internal/routes/api/tasks/route.go` → handles `/api/tasks`
- `internal/routes/api/tasks/[id]/route.go` → handles `/api/tasks/{id}`
- A scanner discovers files and generates route registrations

**Pros:**
- ✅ Very intuitive - route matches file path
- ✅ No manual registration needed
- ✅ Easy to see all routes by looking at folder structure

**Cons:**
- ❌ Requires a file scanner/code generator
- ❌ Less explicit - magic happens at build time
- ❌ Harder to navigate in IDE
- ❌ File structure can become deep

**Example structure:**
```
internal/routes/
├── api/
│   ├── auth/
│   │   ├── login.post.go
│   │   ├── register.post.go
│   │   └── logout.post.go
│   └── tasks/
│       ├── route.go (GET, POST)
│       └── [id]/
│           ├── route.go (GET, PUT, DELETE)
│           └── complete.post.go
```

---

## Pattern 3: Annotation/Code Generation

**How it works:**
- Use special comments like `// @Route GET /api/tasks`
- A code generator parses comments and creates routes.go
- Gateway imports the generated file

**Pros:**
- ✅ Clean handler code - just add a comment
- ✅ Can generate OpenAPI/Swagger docs simultaneously
- ✅ One source of truth

**Cons:**
- ❌ Requires build step/code generator
- ❌ Comments can get out of sync
- ❌ More complex setup
- ❌ Hidden magic - routes not visible in code

**Example:**
```go
// @Route POST /api/auth/login
// @Public
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

---

## Recommendation

**Start with Pattern 1 (Interface-based)** because:
1. It's already compatible with your fx setup
2. No external tools or build steps needed
3. Explicit and easy to understand
4. Can migrate to Pattern 3 later if you want code generation

**When to use others:**
- **Pattern 2**: If you have many simple CRUD routes and want file-based organization
- **Pattern 3**: If you also want OpenAPI documentation generation

---

## Quick Start: Pattern 1

1. **Add to auth_handler.go:**
```go
func (h *AuthHandler) Routes() []handler.Route {
    return []handler.Route{
        {Method: handler.POST, Path: "/api/auth/register", Handler: h.Register},
        {Method: handler.POST, Path: "/api/auth/login", Handler: h.Login},
        {Method: handler.POST, Path: "/api/auth/refresh", Handler: h.RefreshToken},
        {Method: handler.POST, Path: "/api/auth/logout", Handler: h.Logout},
    }
}
```

2. **Add to task_handler.go:**
```go
func (h *TaskHandler) Routes() []handler.Route {
    return []handler.Route{
        {Method: handler.GET, Path: "/api/tasks", Handler: h.List},
        {Method: handler.POST, Path: "/api/tasks", Handler: h.Create},
        {Method: handler.GET, Path: "/api/tasks/", Handler: h.Get},     // Note: /{id} pattern
        {Method: handler.PUT, Path: "/api/tasks/", Handler: h.Update},
        {Method: handler.DELETE, Path: "/api/tasks/", Handler: h.Delete},
    }
}
```

3. **Replace gateway.go** with `gateway_refactored.go` (rename it)

4. **Done!** fx will automatically wire everything up.

# Architecture Documentation

## Table of Contents

1. [Overview](#overview)
2. [Design Principles](#design-principles)
3. [System Architecture](#system-architecture)
4. [Domain Model](#domain-model)
5. [Component Architecture](#component-architecture)
6. [Data Flow](#data-flow)
7. [Security Architecture](#security-architecture)
8. [Event-Driven Architecture](#event-driven-architecture)
9. [Scalability Considerations](#scalability-considerations)
10. [Technology Decisions](#technology-decisions)

## Overview

TaskHub is built using **Clean Architecture** principles combined with **Domain-Driven Design (DDD)**. This approach ensures:

- **Independence from frameworks**: The business logic doesn't depend on external frameworks
- **Testability**: The system can be tested without UI, database, or external services
- **Independence from UI**: The business logic doesn't know about the UI
- **Independence from database**: Business rules are not bound to the database
- **Independence from external agents**: Business logic doesn't know about the outside world

## Design Principles

### 1. Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                        Presentation                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Web UI    │  │   REST API │  │   GraphQL API       │  │
│  │   (HTMX)    │  │  Handlers  │  │   (Future)          │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                     Application Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Services  │  │   Use Cases │  │   Application      │  │
│  │             │  │             │  │   Coordinators     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                       Domain Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Entities   │  │  Value Objs │  │   Domain Events     │  │
│  │             │  │             │  │                     │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────┐
│                    Infrastructure Layer                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ Database    │  │   External  │  │   Message Broker    │  │
│  │ Repositories │  │   APIs      │  │   (NATS)           │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 2. SOLID Principles

- **Single Responsibility**: Each class has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Derived types must be substitutable for base types
- **Interface Segregation**: Clients shouldn't depend on unused interfaces
- **Dependency Inversion**: Depend on abstractions, not concretions

### 3. Domain-Driven Design

- **Ubiquitous Language**: Shared language between developers and domain experts
- **Bounded Contexts**: Clear boundaries for different domains
- **Aggregates**: Consistency boundaries around domain entities
- **Repositories**: Abstractions for data access
- **Domain Services**: Business logic that doesn't naturally fit in entities

## System Architecture

### High-Level Architecture

```
                    ┌─────────────────┐
                    │   Load Balancer │
                    │    (Nginx)      │
                    └─────────────────┘
                            │
                            ▼
                    ┌─────────────────┐
                    │  TaskHub App    │
                    │  (Go + FX)      │
                    └─────────────────┘
                            │
            ┌───────────────┼───────────────┐
            │               │               │
            ▼               ▼               ▼
    ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
    │ PostgreSQL  │ │    NATS     │ │   Static    │
    │ (Database)  │ │ (Messaging) │ │   Files     │
    └─────────────┘ └─────────────┘ └─────────────┘
```

### Module Architecture

TaskHub uses **Uber FX** for dependency injection and module composition:

```go
// Main application composition
app := fx.New(
    // Infrastructure modules
    pkg.DatabaseModule,
    pkg.LoggerModule,
    pkg.NATSModule,
    
    // Domain modules
    domains.UserModule,
    domains.TaskModule,
    domains.NotificationModule,
    
    // Application modules
    app.AuthModule,
    app.TaskModule,
    app.UserModule,
    app.NotificationModule,
    
    // Interface modules
    handler.WebModule,
    handler.APIModule,
    gateway.GatewayModule,
)
```

## Domain Model

### Core Entities

#### User Entity
```go
type User struct {
    BaseEntity
    Name     string    `json:"name" db:"name"`
    Email    string    `json:"email" db:"email"`
    Password string    `json:"-" db:"password"` // bcrypt hashed
}

// User business rules
func (u *User) Validate() error {
    if u.Name == "" {
        return errors.New("name is required")
    }
    if !isValidEmail(u.Email) {
        return errors.New("invalid email format")
    }
    return nil
}
```

#### Task Entity
```go
type Task struct {
    BaseEntity
    Title       string     `json:"title" db:"title"`
    Description string     `json:"description" db:"description"`
    Status      TaskStatus `json:"status" db:"status"`
    Priority    TaskPriority `json:"priority" db:"priority"`
    Deadline    *time.Time `json:"deadline" db:"deadline"`
    UserID      uuid.UUID  `json:"user_id" db:"user_id"`
}

// Task business rules
func (t *Task) CanTransitionTo(newStatus TaskStatus) bool {
    validTransitions := map[TaskStatus][]TaskStatus{
        StatusTodo:        {StatusInProgress},
        StatusInProgress:  {StatusDone, StatusTodo},
        StatusDone:        {}, // Terminal state
    }
    
    for _, valid := range validTransitions[t.Status] {
        if valid == newStatus {
            return true
        }
    }
    return false
}
```

### Value Objects

#### TaskStatus
```go
type TaskStatus string

const (
    StatusTodo        TaskStatus = "todo"
    StatusInProgress  TaskStatus = "in_progress"
    StatusDone        TaskStatus = "done"
)
```

#### TaskPriority
```go
type TaskPriority string

const (
    PriorityLow    TaskPriority = "low"
    PriorityMedium TaskPriority = "medium"
    PriorityHigh   TaskPriority = "high"
)
```

### Domain Events

```go
type TaskEvent struct {
    ID        uuid.UUID `json:"id"`
    TaskID    uuid.UUID `json:"task_id"`
    UserID    uuid.UUID `json:"user_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
    Data      interface{} `json:"data"`
}

// Event types
const (
    EventTaskCreated   = "task.created"
    EventTaskUpdated   = "task.updated"
    EventTaskCompleted = "task.completed"
    EventTaskDeleted   = "task.deleted"
    EventTaskReminder  = "task.reminder"
)
```

## Component Architecture

### 1. Gateway Layer

```go
type Gateway struct {
    db     *sql.DB
    nats   *nats.Conn
    logger *slog.Logger
}

// Database operations
func (g *Gateway) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
func (g *Gateway) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

// Message publishing
func (g *Gateway) Publish(subject string, data interface{}) error
func (g *Gateway) Subscribe(subject string, handler func(*nats.Msg)) error
```

### 2. Repository Pattern

```go
// Interface definition
type TaskRepository interface {
    Create(ctx context.Context, task *Task) error
    GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
    GetByUserID(ctx context.Context, userID uuid.UUID, filters TaskFilters) ([]*Task, error)
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// Implementation
type taskRepository struct {
    gateway *gateway.Gateway
}

func (r *taskRepository) Create(ctx context.Context, task *Task) error {
    query := `
        INSERT INTO tasks (title, description, status, priority, deadline, user_id, created_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at
    `
    
    return r.gateway.QueryRow(ctx, query, 
        task.Title, task.Description, task.Status, 
        task.Priority, task.Deadline, task.UserID, task.UserID).
        Scan(&task.ID, &task.CreatedAt)
}
```

### 3. Service Layer

```go
type TaskService struct {
    taskRepo    domains.TaskRepository
    userRepo    domains.UserRepository
    gateway     *gateway.Gateway
    eventBus    *EventBus
}

func (s *TaskService) CreateTask(ctx context.Context, req *CreateTaskRequest, userID uuid.UUID) (*Task, error) {
    // Validate user exists
    if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
        return nil, errors.New("user not found")
    }
    
    // Create task
    task := &Task{
        Title:       req.Title,
        Description: req.Description,
        Status:      StatusTodo,
        Priority:    req.Priority,
        Deadline:    req.Deadline,
        UserID:      userID,
    }
    
    if err := s.taskRepo.Create(ctx, task); err != nil {
        return nil, err
    }
    
    // Publish domain event
    event := &TaskEvent{
        ID:        uuid.New(),
        TaskID:    task.ID,
        UserID:    userID,
        EventType: EventTaskCreated,
        Timestamp: time.Now(),
        Data:      task,
    }
    
    s.eventBus.Publish(event)
    
    return task, nil
}
```

## Data Flow

### Request Flow

```
1. HTTP Request → Handler
2. Handler → Service (Business Logic)
3. Service → Repository (Data Access)
4. Repository → Gateway (Database/External)
5. Gateway → Database/NATS
6. Response flows back through layers
```

### Event Flow

```
1. Domain Event → EventBus
2. EventBus → NATS Publisher
3. NATS → Subscribers (Notification Service)
4. Subscribers → Background Jobs
5. Jobs → External Services (Email, SMS, etc.)
```

### Authentication Flow

```
1. Login Request → AuthHandler
2. AuthHandler → AuthService
3. AuthService → UserRepository
4. User Validation → JWT Token Generation
5. Token → HTTP Response (Cookie/Header)
6. Subsequent Requests → Middleware Validation
7. Middleware → Context User Injection
```

## Security Architecture

### 1. Authentication

```go
type JWTClaims struct {
    UserID uuid.UUID `json:"user_id"`
    Email  string    `json:"email"`
    jwt.RegisteredClaims
}

type AuthService struct {
    userRepo       domains.UserRepository
    tokenSecret    string
    accessDuration time.Duration
    refreshDuration time.Duration
}

func (s *AuthService) GenerateTokenPair(user *User) (*TokenPair, error) {
    // Access token (15 minutes)
    accessClaims := JWTClaims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    // Refresh token (7 days)
    refreshClaims := JWTClaims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    // Generate tokens...
}
```

### 2. Authorization Middleware

```go
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from cookie or header
            token := extractToken(r)
            if token == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            // Validate token
            claims, err := validateToken(token, secret)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            // Add user context
            ctx := context.WithValue(r.Context(), "user", claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 3. Resource-Based Authorization

```go
func (s *TaskService) GetTask(ctx context.Context, taskID uuid.UUID, userID uuid.UUID) (*Task, error) {
    task, err := s.taskRepo.GetByID(ctx, taskID)
    if err != nil {
        return nil, err
    }
    
    // Authorization check: user can only access their own tasks
    if task.UserID != userID {
        return nil, errors.New("unauthorized access to task")
    }
    
    return task, nil
}
```

## Event-Driven Architecture

### 1. Event Bus Implementation

```go
type EventBus struct {
    conn *nats.Conn
    logger *slog.Logger
}

func (eb *EventBus) Publish(event interface{}) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    subject := getSubjectForEvent(event)
    return eb.conn.Publish(subject, data)
}

func (eb *EventBus) Subscribe(eventType string, handler func(interface{})) error {
    subject := fmt.Sprintf("taskhub.%s.*", eventType)
    
    _, err := eb.conn.Subscribe(subject, func(msg *nats.Msg) {
        var event interface{}
        if err := json.Unmarshal(msg.Data, &event); err != nil {
            eb.logger.Error("failed to unmarshal event", "error", err)
            return
        }
        
        handler(event)
    })
    
    return err
}
```

### 2. Event Types and Subjects

```
Task Events:
- taskhub.task.created     → Notification Service
- taskhub.task.updated     → Search Index Service
- taskhub.task.completed   → Analytics Service
- taskhub.task.reminder    → Notification Service

User Events:
- taskhub.user.registered  → Welcome Email Service
- taskhub.user.login       → Security Audit Service
```

### 3. Notification Service

```go
type NotificationService struct {
    eventBus *EventBus
    emailClient *EmailClient
    smsClient   *SMSClient
}

func (ns *NotificationService) Start() {
    // Subscribe to task events
    ns.eventBus.Subscribe("task.created", ns.handleTaskCreated)
    ns.eventBus.Subscribe("task.reminder", ns.handleTaskReminder)
}

func (ns *NotificationService) handleTaskReminder(event *TaskEvent) error {
    // Send reminder notification
    task := event.Data.(*Task)
    
    notification := &Notification{
        UserID:  task.UserID,
        Title:   "Task Reminder",
        Message: fmt.Sprintf("Task '%s' is due soon!", task.Title),
        Type:    NotificationTypeReminder,
    }
    
    return ns.sendNotification(notification)
}
```

## Scalability Considerations

### 1. Horizontal Scaling

- **Stateless Application**: All state stored in database or message broker
- **Load Balancing**: Multiple instances behind load balancer
- **Database Connection Pooling**: Efficient database resource usage
- **Caching Layer**: Redis for session and frequently accessed data

### 2. Database Scaling

```go
// Read replica configuration
type DatabaseConfig struct {
    Master *sql.DB
    Slaves []*sql.DB
}

func (dc *DatabaseConfig) Read(query string, args ...interface{}) (*sql.Rows, error) {
    // Use read replica for read operations
    slave := dc.getRandomSlave()
    return slave.Query(query, args...)
}

func (dc *DatabaseConfig) Write(query string, args ...interface{}) (sql.Result, error) {
    // Use master for write operations
    return dc.Master.Exec(query, args...)
}
```

### 3. Message Queue Scaling

- **Partitioning**: NATS subject partitioning for high throughput
- **Consumer Groups**: Multiple consumers for event processing
- **Backpressure Handling**: Flow control for message processing

### 4. Caching Strategy

```go
type CacheService struct {
    redis *redis.Client
    ttl   time.Duration
}

func (cs *CacheService) GetTasks(userID uuid.UUID) ([]*Task, error) {
    // Try cache first
    key := fmt.Sprintf("tasks:user:%s", userID)
    cached, err := cs.redis.Get(key).Result()
    if err == nil {
        var tasks []*Task
        json.Unmarshal([]byte(cached), &tasks)
        return tasks, nil
    }
    
    // Fallback to database
    tasks, err := cs.taskRepo.GetByUserID(context.Background(), userID)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    data, _ := json.Marshal(tasks)
    cs.redis.Set(key, data, cs.ttl)
    
    return tasks, nil
}
```

## Technology Decisions

### 1. Go Language

**Chosen for:**
- Performance and efficiency
- Strong typing and compile-time safety
- Excellent concurrency support
- Rich standard library
- Fast compilation and deployment

### 2. PostgreSQL

**Chosen for:**
- ACID compliance and data integrity
- Advanced features (JSONB, indexes, transactions)
- Strong consistency guarantees
- Excellent tooling and ecosystem
- Proven scalability

### 3. NATS Messaging

**Chosen for:**
- Lightweight and high-performance
- Simple protocol and deployment
- Excellent Go client library
- Support for various messaging patterns
- Cloud-native design

### 4. Uber FX Dependency Injection

**Chosen for:**
- Clean dependency management
- Module-based architecture
- Lifecycle management
- Excellent Go integration
- Minimal runtime overhead

### 5. HTMX for Frontend

**Chosen for:**
- Simplicity and progressive enhancement
- No complex JavaScript framework required
- Server-side rendering with modern interactions
- Small learning curve
- Excellent performance

### 6. JWT Authentication

**Chosen for:**
- Stateless authentication
- Cross-platform compatibility
- Built-in expiration handling
- Security best practices
- Easy integration with clients

---

This architecture documentation provides a comprehensive view of TaskHub's design decisions, patterns, and implementation details. It serves as a guide for developers, architects, and system administrators working with the system.
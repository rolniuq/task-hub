# 🧩 TaskHub – Personal Task Management System

## 📘 Overview
**TaskHub** is a backend system for managing personal or team tasks.
It allows users to create, track, and manage tasks efficiently while demonstrating professional backend engineering practices.
This project is designed to showcase expertise in **Golang**, **system design**, and **clean architecture**.

---

## 🚀 Key Features
- User authentication and authorization (JWT + refresh tokens)
- CRUD operations for tasks (create, update, delete, complete)
- Task prioritization and filtering (by status, priority, deadline)
- Background job system for reminders and notifications
- Caching for recent task lists and session management
- Modular, scalable architecture (microservice-ready)
- Well-documented REST API with Swagger

---

## 🧱 Tech Stack
| Component | Technology | Purpose |
|------------|-------------|----------|
| Language | **Golang (Fiber or Echo)** | High-performance backend |
| Database | **PostgreSQL** | Transactional data storage |
| Cache | **Redis** | Session, caching, rate limiting |
| Message Broker | **NATS / Asynq** | Background jobs & event-driven processing |
| ORM | **GORM or MikroORM** | Data mapping |
| Auth | **JWT + Refresh Token** | Stateless authentication |
| Config | **Viper / Env** | Environment-based configuration |
| Logging | **Zap / Zerolog** | Structured logging |
| Testing | **Go test + Testify** | Unit & integration testing |
| Container | **Docker Compose** | Local orchestration |
| Docs | **Swagger / Redoc** | API documentation |

---

## 🧩 System Architecture

```
                  +----------------------+
                  |      API Gateway     |
                  +----------------------+
                            |
          ------------------------------------------
          |                    |                   |
     User Service         Task Service        Notification Service
 (auth, profile)       (CRUD, filter)        (async reminder jobs)
          |                    |                   |
     PostgreSQL           PostgreSQL / Redis        Redis / NATS
```

### Core Flow
1. User logs in → receives JWT token
2. Sends request `POST /tasks` → Task Service saves to PostgreSQL
3. Task Service publishes event → NATS → Notification Service consumes
4. Notification Service runs background jobs (Asynq) for reminders
5. Redis caches user session and recent tasks

---

## 🧠 Business Logic

### Auth Flow
- Register → Hash password with bcrypt → store user
- Login → Validate password → issue JWT + refresh token
- Middleware verifies JWT for all protected routes

### Task Flow
- CRUD: create, update, delete, mark as completed
- Filtering: by `status`, `priority`, `deadline`
- Background jobs send reminders before deadlines
- (Future) Assign tasks between users

---

## 🏗️ Deployment
**Docker Compose** setup includes:
- `taskhub-api` – main Golang service
- `postgres` – database
- `redis` – caching and job queue
- `nats` – message broker
- `asynq-worker` – background job processor

Command:
```bash
docker compose up
```

---

## 📦 Folder Structure (proposed)
```
taskhub/
├── cmd/
│   └── main.go
├── internal/
│   ├── user/
│   ├── task/
│   ├── notification/
│   └── pkg/
│       ├── db/
│       ├── logger/
│       └── middleware/
├── config/
├── docs/
│   └── README.md
├── docker-compose.yml
└── go.mod
```

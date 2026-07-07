# TaskHub - Modern Task Management System

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-blue.svg)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

TaskHub is a modern, high-performance task management system built with Go, featuring a clean architecture, event-driven design, and responsive web interface powered by HTMX.

## ✨ Key Features

- **🚀 High Performance**: Built with Go 1.24+ and optimized for speed
- **🏗️ Clean Architecture**: Domain-driven design with clear separation of concerns
- **🔐 Secure Authentication**: JWT-based auth with refresh tokens
- **📱 Responsive UI**: Modern web interface with HTMX for seamless interactions
- **🔄 Event-Driven**: NATS-powered messaging for real-time notifications
- **🐳 Containerized**: Docker support with multi-stage builds
- **🧪 Well Tested**: Comprehensive test coverage with unit and integration tests

## 🏛️ Architecture Overview

TaskHub follows **Clean Architecture** principles with **Domain-Driven Design**:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Layer     │    │   API Layer     │    │   Gateway       │
│                 │    │                 │    │                 │
│ • HTMX UI       │◄──►│ • HTTP Handlers │◄──►│ • Database      │
│ • Static Files  │    │ • Middleware    │    │ • NATS          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  Application    │
                       │                 │
                       │ • Services      │
                       │ • Business Logic│
                       └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │    Domain       │
                       │                 │
                       │ • Entities      │
                       │ • Repositories  │
                       │ • Use Cases     │
                       └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 16+
- Docker & Docker Compose (optional)
- NATS Server (included in Docker Compose)

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-org/task-hub.git
cd task-hub

# Start all services
docker compose up -d

# The application will be available at http://localhost:8080
```

### Manual Installation

```bash
# Clone the repository
git clone https://github.com/your-org/task-hub.git
cd task-hub

# Install dependencies
go mod download

# Create .env file (required)
cat > .env << EOF
PORT=8080
NATS_URL=nats://localhost:4222
JWT_SECRET=your_jwt_secret_key_at_least_32_characters
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskhub
DB_PASSWORD=dev_password
DB_NAME=taskhub
EOF

# Start the application
go run cmd/main.go
```

## 📊 Technology Stack

| Component | Technology | Version |
|-----------|------------|---------|
| **Backend** | Go | 1.24+ |
| **Database** | PostgreSQL | 16+ |
| **Messaging** | NATS | Latest |
| **Authentication** | JWT | v5 |
| **DI Framework** | Uber FX | v1.24 |
| **Frontend** | HTMX + HTML/CSS | Latest |
| **Deployment** | Docker | Latest |

## 📖 Documentation

- [**Architecture**](ARCHITECTURE.md) - Detailed system design and patterns
- [**API Reference**](API.md) - Complete API documentation
- [**Development**](DEVELOPMENT.md) - Development setup and guidelines
- [**Deployment**](DEPLOYMENT.md) - Production deployment guide
- [**Contributing**](CONTRIBUTING.md) - How to contribute

## 🔧 Configuration

TaskHub uses environment variables for configuration. Key settings:

```bash
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskhub
DB_PASSWORD=your_password
DB_NAME=taskhub

# Authentication
JWT_SECRET=your_jwt_secret_key_at_least_32_characters

# NATS
NATS_URL=nats://localhost:4222
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/domains/task/...

# Run tests with race detection
go test -race ./...
```

## 📝 API Usage Examples

### Authentication

```bash
# Register user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secure123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secure123"}'
```

### Task Management

```bash
# Create task (requires auth)
curl -X POST http://localhost:8080/api/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Complete project","description":"Finish the TaskHub project","priority":"high"}'

# List tasks
curl -X GET http://localhost:8080/api/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🤝 Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code of Conduct
- Development Process
- Pull Request Guidelines
- Coding Standards

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Uber FX](https://github.com/uber-go/fx) for dependency injection
- [NATS](https://nats.io/) for messaging system
- [HTMX](https://htmx.org/) for modern web interactions
- [PostgreSQL](https://www.postgresql.org/) for robust data storage

## 📞 Support

- 📧 Email: support@taskhub.dev
- 💬 Discord: [Join our community](https://discord.gg/taskhub)
- 🐛 Issues: [GitHub Issues](https://github.com/your-org/task-hub/issues)

---

**Built with ❤️ by the TaskHub Team**

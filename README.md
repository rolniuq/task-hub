# Task Hub

[![CI/CD Pipeline](https://github.com/${{ github.repository }}/actions/workflows/ci.yml/badge.svg)](https://github.com/${{ github.repository }}/actions/workflows/ci.yml)
[![Dependencies](https://github.com/${{ github.repository }}/actions/workflows/dependencies.yml/badge.svg)](https://github.com/${{ github.repository }}/actions/workflows/dependencies.yml)
[![Performance](https://github.com/${{ github.repository }}/actions/workflows/performance.yml/badge.svg)](https://github.com/${{ github.repository }}/actions/workflows/performance.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/${{ github.repository }})](https://goreportcard.com/report/github.com/${{ github.repository }})
[![Coverage](https://codecov.io/gh/${{ github.repository }}/branch/main/graph/badge.svg)](https://codecov.io/gh/${{ github.repository }})

A modern task management application with both web and desktop interfaces built with Go, PostgreSQL, and NATS messaging.

## 🚀 Features

- **Dual Interface**: Web application and native desktop app
- **Modern UI**: Clean, responsive design with Fyne desktop framework
- **Real-time Updates**: NATS-powered messaging for instant sync
- **Secure Authentication**: JWT-based authentication system
- **Task Management**: Create, update, and organize tasks
- **Docker Support**: Full containerization for easy deployment
- **Database**: PostgreSQL for reliable data storage

## 📋 Prerequisites

- **Go**: 1.25 or higher
- **PostgreSQL**: 16 or higher
- **NATS Server**: 2.9 or higher
- **Docker**: 20.10 or higher (for containerized deployment)
- **Task**: Task runner for development tasks

## 🛠️ Quick Start

### Using Task Runner (Recommended)

```bash
# Install Task (if not already installed)
go install github.com/go-task/task/v3/cmd/task@latest

# Clone and set up
git clone <repository-url>
cd task-hub

# Copy environment file
cp .env.example .env
# Edit .env with your configuration

# Start development services (PostgreSQL + NATS)
task run

# Run web version locally
task run-web

# Run desktop version locally
task run-desktop
```

### Using Docker

```bash
# Start all services with Docker Compose
task run

# Or directly with docker-compose
docker compose up -d
```

### Manual Setup

```bash
# Install dependencies
go mod download

# Start PostgreSQL and NATS
docker run --name taskhub-postgres -e POSTGRES_USER=taskhub -e POSTGRES_PASSWORD=dev_password -e POSTGRES_DB=taskhub -p 5432:5432 -d postgres:16-alpine
docker run --name taskhub-nats -p 4222:4222 -d nats:2.9-alpine

# Run web application
go run cmd/main.go

# Run desktop application
go run cmd/desktop/main.go
```

## 📁 Project Structure

```
task-hub/
├── cmd/                    # Application entry points
│   ├── main.go            # Web application
│   └── desktop/           # Desktop application
│       └── main.go        # Desktop app entry point
├── internal/              # Private application code
│   ├── app/               # Application services (auth, tasks, etc.)
│   ├── desktop/           # Desktop UI components and themes
│   ├── domains/           # Domain models and business logic
│   ├── gateway/           # External service integration
│   └── handler/           # HTTP handlers
├── pkg/                   # Public library code
│   ├── base/              # Base entities and repositories
│   ├── db/                # Database configuration
│   ├── logger/            # Logging utilities
│   ├── middleware/        # HTTP middleware
│   ├── nats/              # NATS messaging
│   └── utils/             # Utility functions
├── web/                   # Web assets
│   ├── static/            # CSS, JS, images
│   └── templates/         # HTML templates
├── config/                # Configuration management
├── docs/                  # Documentation
├── .env.example           # Environment template
├── Taskfile.yml           # Task runner configuration
├── docker-compose.yml     # Docker services
├── Dockerfile             # Web app container
├── Dockerfile.desktop     # Desktop app container
└── README.md
```

## 🎨 UI Components

### Web Interface
- Built with standard Go net/http and html/template
- HTMX for dynamic interactions
- Responsive CSS design
- RESTful API endpoints

### Desktop Interface
- Built with Fyne cross-platform framework
- Modern Material Design-inspired theme
- Native desktop experience
- Same backend services as web version

## 🔧 Available Tasks

The project uses Task for development automation. Available commands:

```bash
# Build and deployment
task build                # Build web Docker image
task build-desktop        # Build desktop Docker image

# Running applications
task run                  # Run web version with Docker
task run-web              # Run web version locally
task run-desktop          # Run desktop version locally

# Docker management
task down                 # Stop all services
task clean                # Clean Docker images, volumes, containers
```

## 🌐 API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Tasks
- `GET /api/tasks` - List user tasks
- `POST /api/tasks` - Create new task
- `PUT /api/tasks/{id}` - Update task
- `DELETE /api/tasks/{id}` - Delete task

### Web Pages
- `GET /` - Dashboard (authenticated)
- `GET /login` - Login page
- `GET /register` - Registration page

## 🔐 Configuration

Copy `.env.example` to `.env` and configure:

```bash
# Application
APP_ENV=development
SERVER_PORT=8080

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

# Logging
LOG_LEVEL=debug
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./tests/...
```

## 📱 Desktop UI Features

The desktop application includes:

- **Modern Theme**: Indigo-based color scheme with proper contrast
- **Enhanced Forms**: Better layout, icons, and placeholders
- **Responsive Design**: Proper window sizing and scaling
- **Professional Layout**: Card-based components with visual hierarchy
- **Native Experience**: Platform-specific optimizations

## 🚀 Deployment

### Production Docker

```bash
# Build production images
docker build -t task-hub:latest .
docker build -f Dockerfile.desktop -t task-hub-desktop:latest .

# Deploy with Docker Compose
docker compose -f docker-compose.yml up -d
```

### Development

```bash
# Start development environment
task run

# Run desktop app (native)
task run-desktop
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## 📚 Documentation

- [Development Guide](docs/DEVELOPMENT.md) - Comprehensive development documentation
- [API Documentation](docs/API.md) - API endpoint reference
- [Architecture](docs/ARCHITECTURE.md) - System architecture overview
- [Deployment](docs/DEPLOYMENT.md) - Deployment guide

## 🐛 Troubleshooting

### Common Issues

1. **Desktop app won't start on Mac M3/M4**
   - The desktop app runs natively (not in Docker) on Mac for proper GUI support
   - Use `task run-desktop` which starts services in Docker but runs the app natively

2. **Database connection errors**
   - Ensure PostgreSQL is running: `docker ps | grep postgres`
   - Check `.env` configuration matches database settings

3. **NATS connection issues**
   - Verify NATS is running: `docker ps | grep nats`
   - Check NATS URL configuration

For more troubleshooting, see the [Development Guide](docs/DEVELOPMENT.md).

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Task Hub** - Modern task management for teams and individuals.
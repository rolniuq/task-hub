# API Documentation

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Base URL](#base-url)
4. [Response Format](#response-format)
5. [Error Handling](#error-handling)
6. [Rate Limiting](#rate-limiting)
7. [Endpoints](#endpoints)
   - [Authentication](#authentication-endpoints)
   - [Tasks](#task-endpoints)
   - [Users](#user-endpoints)
   - [Health](#health-endpoints)
8. [Webhooks](#webhooks)
9. [SDK Examples](#sdk-examples)

## Overview

TaskHub provides a RESTful API for managing tasks and users. The API supports both traditional HTTP requests and HTMX-enhanced web interactions.

### Key Features

- **JWT-based Authentication**: Secure token-based authentication
- **RESTful Design**: Clean, predictable endpoint structure
- **HTMX Support**: Enhanced web interactions for frontend
- **Comprehensive Filtering**: Advanced task filtering and sorting
- **Real-time Events**: Event-driven notifications via NATS
- **Soft Deletes**: Data integrity with soft delete patterns

## Authentication

TaskHub uses JWT (JSON Web Tokens) for authentication. The API supports two authentication methods:

### 1. Cookie-based Authentication (Web)

```http
Cookie: access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 2. Bearer Token Authentication (API)

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Token Types

| Token Type | Duration | Purpose |
|------------|----------|---------|
| Access Token | 15 minutes | API requests |
| Refresh Token | 7 days | Token renewal |

## Base URL

```
Development: http://localhost:8080
Production:  https://api.taskhub.dev
```

## Response Format

All API responses follow a consistent format:

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data
  },
  "message": "Operation completed successfully",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Error Handling

### HTTP Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid request data |
| 401 | Unauthorized | Authentication required |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource conflict |
| 422 | Unprocessable Entity | Validation errors |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |

### Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Input validation failed |
| `AUTHENTICATION_FAILED` | Invalid credentials |
| `AUTHORIZATION_FAILED` | Insufficient permissions |
| `RESOURCE_NOT_FOUND` | Resource does not exist |
| `RESOURCE_CONFLICT` | Resource already exists |
| `RATE_LIMIT_EXCEEDED` | Too many requests |
| `INTERNAL_ERROR` | Server error |

## Rate Limiting

API requests are limited to prevent abuse:

- **Anonymous requests**: 100 requests per hour
- **Authenticated requests**: 1000 requests per hour
- **Burst limit**: 10 requests per second

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1642248600
```

## Endpoints

### Authentication Endpoints

#### Register User

```http
POST /api/auth/register
```

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-15T10:30:00Z"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
      "expires_in": 900
    }
  }
}
```

**Validation Rules:**
- `name`: Required, 2-100 characters
- `email`: Required, valid email format
- `password`: Required, 8-100 characters

#### Login User

```http
POST /api/auth/login
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "secure123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
      "expires_in": 900
    }
  }
}
```

#### Refresh Token

```http
POST /api/auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

#### Logout

```http
POST /api/auth/logout
```

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

### Task Endpoints

#### List Tasks

```http
GET /api/tasks
```

**Query Parameters:**
- `status`: Filter by status (`todo`, `in_progress`, `done`)
- `priority`: Filter by priority (`low`, `medium`, `high`)
- `deadline_before`: Filter tasks with deadline before date (ISO 8601)
- `deadline_after`: Filter tasks with deadline after date (ISO 8601)
- `sort`: Sort field (`title`, `created_at`, `deadline`, `priority`)
- `order`: Sort order (`asc`, `desc`)
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

**Example:**
```http
GET /api/tasks?status=todo&priority=high&sort=created_at&order=desc&page=1&limit=10
```

**Response:**
```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "Complete project documentation",
        "description": "Write comprehensive API documentation",
        "status": "todo",
        "priority": "high",
        "deadline": "2024-01-20T23:59:59Z",
        "user_id": "550e8400-e29b-41d4-a716-446655440001",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "total_pages": 3
    }
  }
}
```

#### Create Task

```http
POST /api/tasks
```

**Request Body:**
```json
{
  "title": "Complete project documentation",
  "description": "Write comprehensive API documentation",
  "priority": "high",
  "deadline": "2024-01-20T23:59:59Z"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Complete project documentation",
    "description": "Write comprehensive API documentation",
    "status": "todo",
    "priority": "high",
    "deadline": "2024-01-20T23:59:59Z",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Validation Rules:**
- `title`: Required, 1-200 characters
- `description`: Optional, max 1000 characters
- `priority`: Optional, `low`, `medium`, `high` (default: `medium`)
- `deadline`: Optional, ISO 8601 datetime

#### Get Task

```http
GET /api/tasks/{id}
```

**Path Parameters:**
- `id`: Task UUID

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Complete project documentation",
    "description": "Write comprehensive API documentation",
    "status": "todo",
    "priority": "high",
    "deadline": "2024-01-20T23:59:59Z",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Update Task

```http
PUT /api/tasks/{id}
```

**Request Body:**
```json
{
  "title": "Updated task title",
  "description": "Updated description",
  "status": "in_progress",
  "priority": "medium",
  "deadline": "2024-01-25T23:59:59Z"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Updated task title",
    "description": "Updated description",
    "status": "in_progress",
    "priority": "medium",
    "deadline": "2024-01-25T23:59:59Z",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

#### Complete Task

```http
POST /api/tasks/{id}/complete
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Complete project documentation",
    "description": "Write comprehensive API documentation",
    "status": "done",
    "priority": "high",
    "deadline": "2024-01-20T23:59:59Z",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z",
    "completed_at": "2024-01-15T11:00:00Z"
  }
}
```

#### Delete Task

```http
DELETE /api/tasks/{id}
```

**Response:**
```json
{
  "success": true,
  "message": "Task deleted successfully"
}
```

**Note:** This performs a soft delete. The task is marked as deleted but not removed from the database.

### User Endpoints

#### Get Current User

```http
GET /api/users/me
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Update User Profile

```http
PUT /api/users/me
```

**Request Body:**
```json
{
  "name": "John Smith"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "John Smith",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### Health Endpoints

#### Health Check

```http
GET /health
```

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:00Z",
    "version": "1.0.0",
    "services": {
      "database": "healthy",
      "nats": "healthy"
    }
  }
}
```

## Webhooks

TaskHub supports webhooks for real-time notifications about task events.

### Configure Webhook

Contact support to configure webhook endpoints for your account.

### Supported Events

- `task.created` - New task created
- `task.updated` - Task updated
- `task.completed` - Task marked as done
- `task.deleted` - Task deleted

### Webhook Payload

```json
{
  "event": "task.created",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "task": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Complete project documentation",
      "status": "todo",
      "priority": "high",
      "user_id": "550e8400-e29b-41d4-a716-446655440001"
    },
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

## SDK Examples

### Go Client

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/taskhub/go-client"
)

func main() {
    // Initialize client
    client := taskhub.NewClient("http://localhost:8080")
    
    // Login
    tokens, err := client.Auth.Login(context.Background(), "john@example.com", "secure123")
    if err != nil {
        panic(err)
    }
    
    // Set authentication
    client.SetAccessToken(tokens.AccessToken)
    
    // Create task
    task, err := client.Tasks.Create(context.Background(), &taskhub.CreateTaskRequest{
        Title:       "Complete project documentation",
        Description: "Write comprehensive API documentation",
        Priority:    taskhub.PriorityHigh,
        Deadline:    time.Now().Add(5 * 24 * time.Hour),
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Created task: %s\n", task.Title)
}
```

### JavaScript Client

```javascript
import { TaskHubClient } from '@taskhub/js-client';

const client = new TaskHubClient('http://localhost:8080');

// Login
const tokens = await client.auth.login('john@example.com', 'secure123');
client.setAccessToken(tokens.access_token);

// Create task
const task = await client.tasks.create({
    title: 'Complete project documentation',
    description: 'Write comprehensive API documentation',
    priority: 'high',
    deadline: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000)
});

console.log('Created task:', task.title);
```

### Python Client

```python
from taskhub_client import TaskHubClient
from datetime import datetime, timedelta

# Initialize client
client = TaskHubClient('http://localhost:8080')

# Login
tokens = client.auth.login('john@example.com', 'secure123')
client.set_access_token(tokens.access_token)

# Create task
task = client.tasks.create({
    'title': 'Complete project documentation',
    'description': 'Write comprehensive API documentation',
    'priority': 'high',
    'deadline': datetime.now() + timedelta(days=5)
})

print(f'Created task: {task.title}')
```

### cURL Examples

```bash
# Register user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secure123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secure123"}'

# Create task
curl -X POST http://localhost:8080/api/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Complete project","description":"Finish the TaskHub project","priority":"high"}'

# List tasks
curl -X GET http://localhost:8080/api/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get specific task
curl -X GET http://localhost:8080/api/tasks/TASK_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Update task
curl -X PUT http://localhost:8080/api/tasks/TASK_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated title","status":"in_progress"}'

# Complete task
curl -X POST http://localhost:8080/api/tasks/TASK_ID/complete \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Delete task
curl -X DELETE http://localhost:8080/api/tasks/TASK_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Testing

### Postman Collection

Import the provided Postman collection to test all API endpoints:

```json
{
  "info": {
    "name": "TaskHub API",
    "description": "Complete API collection for TaskHub"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080"
    },
    {
      "key": "accessToken",
      "value": ""
    }
  ]
}
```

### OpenAPI/Swagger

The API is also documented with OpenAPI 3.0 specification:

```
http://localhost:8080/swagger/index.html
```

---

This API documentation provides comprehensive information for integrating with TaskHub. For additional support, contact our development team or refer to the SDK documentation.
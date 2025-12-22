# Deployment Guide

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Environment Configuration](#environment-configuration)
4. [Deployment Options](#deployment-options)
   - [Docker Compose](#docker-compose-deployment)
   - [Kubernetes](#kubernetes-deployment)
   - [Manual Deployment](#manual-deployment)
5. [Database Setup](#database-setup)
6. [SSL/TLS Configuration](#ssltls-configuration)
7. [Monitoring and Logging](#monitoring-and-logging)
8. [Backup and Recovery](#backup-and-recovery)
9. [Security Considerations](#security-considerations)
10. [Troubleshooting](#troubleshooting)

## Overview

This guide covers various deployment strategies for TaskHub, from local development to production environments. TaskHub is designed to be cloud-native and can be deployed in containerized, virtualized, or bare-metal environments.

## Prerequisites

### System Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| CPU | 2 cores | 4 cores |
| Memory | 4 GB | 8 GB |
| Storage | 20 GB | 50 GB SSD |
| Network | 100 Mbps | 1 Gbps |

### Software Requirements

- **Docker**: 20.10+ (for containerized deployment)
- **Docker Compose**: 2.0+ (for local development)
- **Kubernetes**: 1.24+ (for K8s deployment)
- **PostgreSQL**: 16+ (for manual deployment)
- **NATS Server**: 2.9+ (for manual deployment)

## Environment Configuration

TaskHub uses environment variables for configuration. Create a `.env` file based on the template:

```bash
# Copy the example environment file
cp .env.example .env
```

### Required Environment Variables

```bash
# Application
APP_NAME=TaskHub
APP_VERSION=1.0.0
APP_ENV=production
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskhub
DB_PASSWORD=your_secure_password
DB_NAME=taskhub
DB_SSLMODE=require
DB_MAX_CONNECTIONS=20
DB_MAX_IDLE_CONNECTIONS=5

# JWT
JWT_SECRET=your_jwt_secret_key_at_least_32_characters
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h

# NATS
NATS_URL=nats://localhost:4222
NATS_SUBJECT_PREFIX=taskhub

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# CORS (if needed)
CORS_ALLOWED_ORIGINS=https://yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
```

### Optional Environment Variables

```bash
# Redis (for caching)
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=your_redis_password
REDIS_DB=0

# Email (for notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Monitoring
PROMETHEUS_ENABLED=true
PROMETHEUS_PORT=9090

# Security
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=1h
```

## Deployment Options

### Docker Compose Deployment

This is the recommended approach for development, staging, and small production environments.

#### 1. Create Docker Compose File

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=taskhub
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=taskhub
      - NATS_URL=nats://nats:4222
    depends_on:
      postgres:
        condition: service_healthy
      nats:
        condition: service_started
    restart: unless-stopped
    networks:
      - taskhub-network

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=taskhub
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=taskhub
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U taskhub"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    networks:
      - taskhub-network

  nats:
    image: nats:2.9-alpine
    ports:
      - "4222:4222"
    command: ["-js", "-m", "8222"]
    restart: unless-stopped
    networks:
      - taskhub-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - taskhub-network

volumes:
  postgres_data:

networks:
  taskhub-network:
    driver: bridge
```

#### 2. Create Nginx Configuration

```nginx
# nginx/nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream taskhub {
        server app:8080;
    }

    server {
        listen 80;
        server_name your-domain.com;
        
        # Redirect to HTTPS
        return 301 https://$server_name$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name your-domain.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/nginx/ssl/key.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;

        location / {
            proxy_pass http://taskhub;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        location /health {
            proxy_pass http://taskhub/health;
            access_log off;
        }
    }
}
```

#### 3. Deploy with Docker Compose

```bash
# Create environment file
echo "DB_PASSWORD=your_secure_password" > .env

# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f app
```

### Kubernetes Deployment

For production environments, Kubernetes provides better scalability and management.

#### 1. Create Namespace

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: taskhub
```

#### 2. Create ConfigMap

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: taskhub-config
  namespace: taskhub
data:
  APP_ENV: "production"
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_USER: "taskhub"
  DB_NAME: "taskhub"
  NATS_URL: "nats://nats-service:4222"
  LOG_LEVEL: "info"
```

#### 3. Create Secret

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: taskhub-secret
  namespace: taskhub
type: Opaque
data:
  DB_PASSWORD: <base64-encoded-password>
  JWT_SECRET: <base64-encoded-jwt-secret>
```

#### 4. Create PostgreSQL Deployment

```yaml
# k8s/postgres.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: taskhub
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:16-alpine
        env:
        - name: POSTGRES_USER
          value: taskhub
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: taskhub-secret
              key: DB_PASSWORD
        - name: POSTGRES_DB
          value: taskhub
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: postgres-service
  namespace: taskhub
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
```

#### 5. Create NATS Deployment

```yaml
# k8s/nats.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nats
  namespace: taskhub
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nats
  template:
    metadata:
      labels:
        app: nats
    spec:
      containers:
      - name: nats
        image: nats:2.9-alpine
        command: ["nats-server", "-js", "-m", "8222"]
        ports:
        - containerPort: 4222
        - containerPort: 8222

---
apiVersion: v1
kind: Service
metadata:
  name: nats-service
  namespace: taskhub
spec:
  selector:
    app: nats
  ports:
  - name: client
    port: 4222
    targetPort: 4222
  - name: monitor
    port: 8222
    targetPort: 8222
```

#### 6. Create TaskHub Deployment

```yaml
# k8s/taskhub.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: taskhub
  namespace: taskhub
spec:
  replicas: 3
  selector:
    matchLabels:
      app: taskhub
  template:
    metadata:
      labels:
        app: taskhub
    spec:
      containers:
      - name: taskhub
        image: taskhub:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: taskhub-config
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: taskhub-secret
              key: DB_PASSWORD
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: taskhub-secret
              key: JWT_SECRET
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"

---
apiVersion: v1
kind: Service
metadata:
  name: taskhub-service
  namespace: taskhub
spec:
  selector:
    app: taskhub
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP
```

#### 7. Create Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: taskhub-ingress
  namespace: taskhub
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - your-domain.com
    secretName: taskhub-tls
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: taskhub-service
            port:
              number: 8080
```

#### 8. Deploy to Kubernetes

```bash
# Apply all configurations
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n taskhub
kubectl get services -n taskhub
kubectl get ingress -n taskhub

# View logs
kubectl logs -f deployment/taskhub -n taskhub
```

### Manual Deployment

For environments where containers are not suitable, you can deploy TaskHub manually.

#### 1. Build Binary

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o taskhub-linux ./cmd/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o taskhub-darwin ./cmd/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o taskhub-windows.exe ./cmd/main.go
```

#### 2. Install Dependencies

```bash
# PostgreSQL (Ubuntu/Debian)
sudo apt update
sudo apt install postgresql postgresql-contrib

# NATS Server
wget https://github.com/nats-io/nats-server/releases/download/v2.9.0/nats-server-v2.9.0-linux-amd64.tar.gz
tar xzf nats-server-v2.9.0-linux-amd64.tar.gz
sudo cp nats-server-v2.9.0-linux-amd64/nats-server /usr/local/bin/
```

#### 3. Setup Database

```bash
# Create database user
sudo -u postgres createuser taskhub

# Create database
sudo -u postgres createdb taskhub

# Set password
sudo -u postgres psql -c "ALTER USER taskhub PASSWORD 'your_secure_password';"

# Grant privileges
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE taskhub TO taskhub;"
```

#### 4. Run Database Migrations

```bash
# Run migrations (if you have migration scripts)
./taskhub-linux migrate
```

#### 5. Create Systemd Service

```ini
# /etc/systemd/system/taskhub.service
[Unit]
Description=TaskHub Application
After=network.target postgresql.service nats.service

[Service]
Type=simple
User=taskhub
Group=taskhub
WorkingDirectory=/opt/taskhub
ExecStart=/opt/taskhub/taskhub-linux
Restart=always
RestartSec=5
Environment=APP_ENV=production
Environment=DB_HOST=localhost
Environment=DB_PORT=5432
Environment=DB_USER=taskhub
Environment=DB_PASSWORD=your_secure_password
Environment=DB_NAME=taskhub
Environment=NATS_URL=nats://localhost:4222
Environment=JWT_SECRET=your_jwt_secret_key

[Install]
WantedBy=multi-user.target
```

#### 6. Start Service

```bash
# Enable and start service
sudo systemctl enable taskhub
sudo systemctl start taskhub

# Check status
sudo systemctl status taskhub

# View logs
sudo journalctl -u taskhub -f
```

## Database Setup

### Initialize Database Schema

```sql
-- Create tables
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

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_deadline ON tasks(deadline);
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at);

-- Create triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### Database Backup Script

```bash
#!/bin/bash
# backup.sh

DB_NAME="taskhub"
DB_USER="taskhub"
DB_PASSWORD="your_password"
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p $BACKUP_DIR

# Perform backup
PGPASSWORD=$DB_PASSWORD pg_dump -h localhost -U $DB_USER -d $DB_NAME > $BACKUP_DIR/taskhub_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/taskhub_$DATE.sql

# Remove old backups (keep last 7 days)
find $BACKUP_DIR -name "taskhub_*.sql.gz" -mtime +7 -delete

echo "Backup completed: taskhub_$DATE.sql.gz"
```

## SSL/TLS Configuration

### Let's Encrypt with Certbot

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Manual SSL Certificate

```bash
# Generate private key
openssl genrsa -out private.key 2048

# Generate certificate signing request
openssl req -new -key private.key -out certificate.csr

# Generate self-signed certificate (for development)
openssl x509 -req -days 365 -in certificate.csr -signkey private.key -out certificate.crt
```

## Monitoring and Logging

### Prometheus Monitoring

```yaml
# monitoring/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'taskhub'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s
```

### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "TaskHub Monitoring",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))"
          }
        ]
      }
    ]
  }
}
```

### Log Aggregation with ELK Stack

```yaml
# logging/logstash.conf
input {
  beats {
    port => 5044
  }
}

filter {
  if [fields][service] == "taskhub" {
    json {
      source => "message"
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "taskhub-%{+YYYY.MM.dd}"
  }
}
```

## Backup and Recovery

### Automated Backup Script

```bash
#!/bin/bash
# scripts/backup.sh

set -e

BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Create backup directory
mkdir -p $BACKUP_DIR

# Database backup
echo "Starting database backup..."
PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/db_$DATE.sql.gz

# Application files backup
echo "Starting application files backup..."
tar -czf $BACKUP_DIR/files_$DATE.tar.gz /opt/taskhub

# Clean old backups
echo "Cleaning old backups..."
find $BACKUP_DIR -name "*.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed successfully"
```

### Recovery Script

```bash
#!/bin/bash
# scripts/recover.sh

set -e

BACKUP_FILE=$1
DB_NAME="taskhub"
DB_USER="taskhub"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

# Stop application
sudo systemctl stop taskhub

# Restore database
echo "Restoring database..."
gunzip -c $BACKUP_FILE | PGPASSWORD=$DB_PASSWORD psql -h localhost -U $DB_USER -d $DB_NAME

# Start application
sudo systemctl start taskhub

echo "Recovery completed successfully"
```

## Security Considerations

### 1. Network Security

```bash
# Firewall rules (UFW)
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw deny 5432/tcp   # PostgreSQL (internal only)
sudo ufw deny 4222/tcp   # NATS (internal only)
sudo ufw enable
```

### 2. Application Security

```bash
# Secure file permissions
sudo chmod 600 /etc/taskhub/.env
sudo chmod 755 /opt/taskhub/taskhub-linux
sudo chown taskhub:taskhub /opt/taskhub/taskhub-linux

# SELinux context (if enabled)
sudo semanage fcontext -a -t bin_t "/opt/taskhub/taskhub-linux"
sudo restorecon -v /opt/taskhub/taskhub-linux
```

### 3. Database Security

```sql
-- Restrict database connections
-- In postgresql.conf:
listen_addresses = 'localhost'
port = 5432

-- In pg_hba.conf:
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             postgres                                peer
local   all             taskhub                                 md5
host    all             taskhub         127.0.0.1/32            md5
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Failed

```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Test connection
psql -h localhost -U taskhub -d taskhub

# Check logs
sudo tail -f /var/log/postgresql/postgresql-16-main.log
```

#### 2. NATS Connection Failed

```bash
# Test NATS connection
telnet localhost 4222

# Check NATS logs
journalctl -u nats -f
```

#### 3. Application Won't Start

```bash
# Check application logs
sudo journalctl -u taskhub -f

# Test configuration
./taskhub-linux -config-test

# Check port availability
netstat -tlnp | grep 8080
```

#### 4. High Memory Usage

```bash
# Monitor memory usage
top -p $(pgrep taskhub)

# Check for memory leaks
valgrind --tool=memcheck --leak-check=full ./taskhub-linux
```

### Performance Tuning

#### Database Optimization

```sql
-- Analyze query performance
EXPLAIN ANALYZE SELECT * FROM tasks WHERE user_id = 'uuid';

-- Create missing indexes
CREATE INDEX CONCURRENTLY idx_tasks_user_status ON tasks(user_id, status);

-- Update statistics
ANALYZE tasks;
```

#### Application Optimization

```bash
# Enable Go profiling
export GODEBUG=gctrace=1
./taskhub-linux

# Monitor goroutines
curl http://localhost:8080/debug/pprof/goroutine?debug=1
```

---

This deployment guide provides comprehensive instructions for deploying TaskHub in various environments. For additional support, refer to the monitoring and troubleshooting sections or contact our support team.
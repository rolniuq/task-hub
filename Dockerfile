FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

COPY web/frontend/package*.json ./

RUN npm ci

COPY web/frontend/ .

RUN npm run build

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o task-hub ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/task-hub .
COPY --from=frontend-builder /frontend/dist ./web/dist

EXPOSE 8080

CMD ["./task-hub"]

# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Cài CA certificates để tránh lỗi TLS
RUN apk update && apk add --no-cache ca-certificates

# Copy go mod and sum, download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o server main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Cài CA certificates ở runtime
RUN apk update && apk add --no-cache ca-certificates

# Copy binary từ builder
COPY --from=builder /app/server .

# Copy config nếu cần
COPY --from=builder /app/config ./config

# Expose port (Fly.io sẽ tự map)
EXPOSE 8080

# Set env cho Fly.io
ENV PORT=8080

# Run app
CMD ["./server"]
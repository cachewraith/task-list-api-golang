# Multi-stage build for production optimization
# Stage 1: Build stage
FROM golang:1.22-alpine AS builder

# Install git and ca-certificates for building
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary, -ldflags for smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/api \
    ./cmd/api/main.go

# Stage 2: Production stage - minimal image
FROM scratch AS production

# Copy CA certificates from builder (needed for HTTPS)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder
COPY --from=builder /app/bin/api /api

# Copy migrations (if your app needs them at runtime)
COPY --from=builder /app/migrations /migrations

# Use non-root user (create in scratch via numeric ID)
USER 65534:65534

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/api"]

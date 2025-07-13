# Multi-stage build for NVS CLI
# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nvs .

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata git

# Create non-root user
RUN addgroup -g 1001 -S nvs && \
    adduser -u 1001 -S nvs -G nvs

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/nvs .

# Change ownership to non-root user
RUN chown nvs:nvs /app/nvs

# Switch to non-root user
USER nvs

# Expose port (if needed for web interface)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./nvs"]

# Default command
CMD ["--help"] 
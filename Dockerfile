# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for getting Go modules)
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Create a user to run the application
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a user to run the application
RUN adduser -D -g '' appuser

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy the user info from builder stage
COPY --from=builder /etc/passwd /etc/passwd

# Use the created user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/data || exit 1

# Run the binary
CMD ["./main"]

# Stage 1: Build the Go application
FROM golang:alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hardlinker .

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Create non-root user for security
RUN adduser -D -u 10001 appuser

# Create directory for YAML file storage
RUN mkdir /app-data
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/hardlinker .

# Change ownership to non-root user
RUN chown appuser:appuser hardlinker
RUN chmod +x hardlinker

# Expose port 5070
EXPOSE 5070

# Switch to non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["./hardlinker"]
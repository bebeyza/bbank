# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install sqlite (for runtime)
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy any config files if needed
COPY --from=builder /app/.env* ./

# Create directory for database
RUN mkdir -p /root/data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
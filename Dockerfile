# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install required packages for MySQL
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate SQLC code
RUN go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Build the server
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install MySQL client for connectivity
RUN apk --no-cache add ca-certificates mysql-client

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Expose port 9000
EXPOSE 9000

# Run the binary
CMD ["./server"]
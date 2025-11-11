# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./server/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./server"]
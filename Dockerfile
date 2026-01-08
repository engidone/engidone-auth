# Build stage
FROM golang:1.25.5-alpine AS builder


# SSH completo
RUN apk add --no-cache git openssh-client ca-certificates && \
    mkdir -p /root/.ssh && chmod 700 /root/.ssh

COPY keys/id_rsa_docker /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts

RUN mkdir -p /app/sql
COPY sqlc.yaml  /app
COPY sql/*.sql /app/sql

WORKDIR /app

# GOPRIVATE para repos privados
RUN go env -w GOPRIVATE=github.com/engidone/* && \
    git config --global url."git@github.com:engidone/".insteadOf "https://github.com/engidone/"

# Install required packages for PostgreSQL
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# 3. LIMPIA cache
RUN go clean -modcache

# Download dependencies
RUN go mod download

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && sqlc generate

# Copy source code
COPY . .

# Build the server
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server/main.go
RUN CGO_ENABLED=1 GOOS=linux go build -o client ./cmd/client


# Final stage
FROM alpine:latest

# Install PostgreSQL client for connectivity
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/client .
RUN mkdir -p ./cmd/config
RUN mkdir -p ./keys
COPY --from=builder /app/cmd/config ./cmd/config
COPY --from=builder /app/keys/*_key.pem ./keys
RUN ls

# Expose port 9000
EXPOSE 9000

# Run the binary
CMD ["./server"]
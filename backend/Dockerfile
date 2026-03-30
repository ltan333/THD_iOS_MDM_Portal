# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/bin/server ./cmd/server

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app \
	&& adduser -S app -G app

# Copy binary from builder
COPY --from=builder /app/bin/server ./server
# Copy necessary runtime files
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/internal/ent/migrate/migrations ./internal/ent/migrate/migrations

# Expose port
EXPOSE 8000

# Run as non-root
USER app

# Command to run
CMD ["./server"]

# ======================================
# Build Stage
# ======================================
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# ======================================
# Production Stage
# ======================================
FROM alpine:3.19 AS production

WORKDIR /app

# Install runtime dependencies including PDF, video processing and HLS encoding
RUN apk add --no-cache ca-certificates tzdata poppler-utils yt-dlp ffmpeg

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /build/main .

# Copy migrations and seeds for initialization
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/seeds ./seeds

# Create uploads directory
RUN mkdir -p /app/uploads && chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

EXPOSE 8080

CMD ["./main"]

# ======================================
# Development Stage
# ======================================
FROM golang:1.21-alpine AS development

WORKDIR /app

# Install dev tools, PDF/video processing and HLS encoding
RUN apk add --no-cache git poppler-utils yt-dlp ffmpeg

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

EXPOSE 8080

CMD ["go", "run", "main.go"]

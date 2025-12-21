# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (for go modules)
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

# ---------- Runtime Stage ----------
FROM alpine:3.19

WORKDIR /app

# Create non-root user (security best practice)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary from build stage
COPY --from=builder /app/server .

# Change ownership
RUN chown -R appuser:appgroup /app

USER appuser

# Render sets PORT automatically
EXPOSE 8080

# Start server
CMD ["./server"]

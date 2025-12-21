# --build stage----
FROM golang:1.22-alphine AS builder

WORKDIR /app

# Installing git
RUN apk add --no-cache git

#Copy the golang files -> for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the remaining files
COPY . .

# Build the binary file

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

# -------- Runtime Stage --------
FROM alpine:3.19

WORKDIR /app

# Create a non-root user (security best practice)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/server .

# Change ownership
RUN chown -R appuser:appgroup /app

USER appuser

# Expose API port
EXPOSE 8080

# Start server
CMD ["./server"]

# ---------- Build Stage ----------
FROM golang:1.23-alpine AS builder

# Enable Go modules and disable CGO for static binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Copy go.mod and go.sum first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the app
RUN go build -o amialive ./cmd

# ---------- Runtime Stage ----------
FROM alpine:3.20

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/amialive .

# Optional: run as non-root user
RUN adduser -D -g '' apps
USER apps

# Set environment defaults
ENV JOKE_CHANCE=30
ENV PORT=8080

# Start the app
CMD ["./amialive"]

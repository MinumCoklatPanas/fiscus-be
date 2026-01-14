# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

# Install git (needed for go mod download in some cases)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o fiscus .

# ---------- Runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/fiscus /app/fiscus

# Expose Gin port
EXPOSE 8080

# Run binary
CMD ["/app/fiscus"]

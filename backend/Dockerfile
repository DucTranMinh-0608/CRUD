# =========================
# Stage 1: Build
# =========================

FROM golang:1.27.1-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build Go application
RUN go build -o server .


# =========================
# Stage 2: Run
# =========================

FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Backend port
EXPOSE 8080

# Start application
CMD ["./server"]
# --- Build Stage ---
FROM golang:1.22-alpine AS builder

# Install git and build tools (some deps need it)
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary statically (no CGO)
RUN CGO_ENABLED=0 go build -o queue-backend main.go


# --- Run Stage ---
FROM alpine:latest

# Add certificate authorities (needed for https requests)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/queue-backend .

EXPOSE 8080

ENTRYPOINT ["./queue-backend"]

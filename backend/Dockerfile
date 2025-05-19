# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o queue-backend main.go

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/queue-backend .

EXPOSE 8080

CMD ["./queue-backend"]

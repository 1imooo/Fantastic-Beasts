# Build stage
FROM golang:1.25.11-alpine AS builder

WORKDIR /app

COPY src/go.mod ./
RUN go mod download 2>/dev/null || true

COPY src/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server .

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY assets/ ./assets/

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]

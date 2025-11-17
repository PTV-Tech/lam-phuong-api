# syntax=docker/dockerfile:1

### BUILD STAGE
FROM golang:1.24.5-alpine AS builder

# Cần git để lấy commit
RUN apk add --no-cache git

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Cache dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary + inject build info
RUN go build -o server ./cmd/server

### RUNTIME STAGE
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server .

EXPOSE 8080

USER nonroot:nonroot

CMD ["./server"]

# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (kalau ada private repo dependency)
RUN apk add --no-cache git

# Copy dependency file dulu (biar cache optimal)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# ---------- Run Stage ----------
FROM alpine:latest

WORKDIR /app

# Tambah CA certificate (biar bisa HTTPS request)
RUN apk add --no-cache ca-certificates

# Copy binary dari builder
COPY --from=builder /app/app .

EXPOSE 8081

CMD ["./app"]
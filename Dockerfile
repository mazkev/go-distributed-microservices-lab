# ==========================================
# STAGE 1: Build Environment
# ==========================================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency definition
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary yang optimal untuk production
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api-server ./cmd/api

# ==========================================
# STAGE 2: Lightweight Production Runner
# ==========================================
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates & tzdata untuk timezone dan HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary dari builder stage (Hanya binary, tanpa Go compiler)
COPY --from=builder /app/api-server .

EXPOSE 8080

CMD ["./api-server"]

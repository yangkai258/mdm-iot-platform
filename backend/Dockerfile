# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mdm-server .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install certificates and curl for healthcheck
RUN apk --no-cache add ca-certificates curl

# Copy binary from builder
COPY --from=builder /app/mdm-server .

# Create non-root user
RUN adduser -D -u 1000 mdm

USER mdm

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./mdm-server"]

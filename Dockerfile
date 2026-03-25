# ===================================================================
# MDM Backend - Multi-stage Docker Build
# ===================================================================
# Stage 1: Builder - compile Go binary
FROM golang:1.23-alpine AS builder

# Enable automatic Go toolchain download for newer Go version
# Use Chinese Go proxy mirror for network accessibility
ENV GOTOOLCHAIN=auto
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
# CGO_ENABLED=0 for static binary, suitable for alpine
# GOTOOLCHAIN=auto downloads correct Go version as required by go.mod
RUN GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" \
    -o mdm-server \
    main.go

# ===================================================================
# Stage 2: Runtime - minimal runtime image
# ===================================================================
FROM alpine:3.19 AS runtime

LABEL maintainer="MDM Team"
LABEL description="MDM IoT Platform Backend"

# Use Aliyun mirror for China network access
RUN sed -i 's/https:\/\/dl-cdn.alpinelinux.org/https:\/\/mirrors.aliyun.com/g' /etc/apk/repositories

# Install ca-certificates and curl for health checks
RUN apk add --no-cache ca-certificates curl

# Create non-root user for security
RUN adduser -D -u 1000 -s /bin/sh mdm

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/mdm-server /app/mdm-server

# Create directory for logs
RUN mkdir -p /app/logs && chown -R mdm:mdm /app

USER mdm

EXPOSE 8080

# Health check - verify the service responds
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the binary
CMD ["/app/mdm-server"]

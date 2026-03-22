# MDM Backend - Pre-built binary deployment
FROM alpine:latest

WORKDIR /app

# Use Aliyun mirror for China network access
RUN sed -i 's/https:\/\/dl-cdn.alpinelinux.org/https:\/\/mirrors.aliyun.com/g' /etc/apk/repositories

# Install certificates and curl for healthcheck
RUN apk --no-cache add ca-certificates curl

# Copy pre-built binary
COPY mdm-server-linux /app/mdm-server

# Create non-root user
RUN adduser -D -u 1000 mdm

USER mdm

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["/app/mdm-server"]

# ==============================================
# Frontend Dockerfile
# ==============================================
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Install pnpm for faster package management
RUN npm install -g pnpm

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN pnpm build

# ==============================================
# Stage 2: Nginx Runtime
# ==============================================
FROM nginx:alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built static files
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/health || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]

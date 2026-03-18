# Dockerfile for BAC Messaging Service (Node.js)
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY .opencode/package*.json .opencode/

# Install dependencies
RUN npm ci --only=production

# Copy source code
COPY . .

# Production stage
FROM node:18-alpine

WORKDIR /app

# Install dependencies
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.opencode/node_modules ./.opencode/node_modules
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.opencode/package*.json ./.opencode/

# Copy source code
COPY --from=builder /app/src ./src
COPY --from=builder /app/scripts ./scripts

# Create a non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Change ownership
RUN chown -R nodejs:nodejs /app

# Switch to non-root user
USER nodejs

# Expose port
EXPOSE 3001

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:3001/health', (r) => {if(r.statusCode!==200)throw new Error(r.statusCode)})"

# Start command
CMD ["node", "src/messaging/server.ts"]
# Stage 1: Build the frontend
# Use a slim, Debian-based image for better compatibility with npm native modules
FROM node:18-slim AS frontend-builder
WORKDIR /app/web

# Create and switch to a non-root user
RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser
USER appuser

# Copy package files and install dependencies
COPY --chown=appuser:appgroup web/package*.json ./
RUN npm install

# Copy the rest of the source code
COPY --chown=appuser:appgroup web/ .

# Run the build script
RUN npm run build

# Stage 2: Build the backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

# Stage 3: Create the final image
FROM alpine:latest
WORKDIR /app

# Create a non-root user for the final image
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy built assets
COPY --from=frontend-builder /app/web/dist ./web/dist
COPY --from=backend-builder /app/server ./server

# Set ownership for the entire app directory
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

EXPOSE 8080
CMD ["./server"]
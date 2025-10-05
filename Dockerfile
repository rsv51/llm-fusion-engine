# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web

# Create and switch to a non-root user to avoid permission issues
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy package files and install dependencies as the non-root user
# Using --chown ensures the files are owned by the correct user
COPY --chown=appuser:appgroup web/package*.json ./
RUN npm ci

# Copy the rest of the source code
COPY --chown=appuser:appgroup web/ .

# Run the build script as the non-root user. This is the most reliable way.
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
RUN apk --no-cache add ca-certificates

# Create a non-root user for the final image for security
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
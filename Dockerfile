# Stage 1: Build the frontend
# Use a slim, Debian-based image for better compatibility with npm native modules
FROM node:18-slim AS frontend-builder

# Create a non-root user with a valid home directory
RUN groupadd -r appgroup && useradd -r -g appgroup -d /home/appuser -m appuser

# Set environment variables for the non-root user
ENV HOME=/home/appuser
ENV PATH=$HOME/.npm-global/bin:$PATH

# Set the working directory
WORKDIR /app/web

# Ensure workspace directory is owned by appuser for npm writes
RUN mkdir -p /app/web && chown -R appuser:appgroup /app

# Switch to the non-root user
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
# Copy go.mod and go.sum first to leverage Docker's layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/llm-fusion-engine ./cmd/server/main.go

# Stage 3: Create the final image
FROM alpine:latest
WORKDIR /app

# Create a non-root user for the final image
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy built assets
COPY --from=frontend-builder /app/web/dist ./web/dist
COPY --from=backend-builder /app/llm-fusion-engine .

# Create an empty database file and set ownership for the entire app directory
# This ensures the app can write to the database file.
RUN touch fusion.db && chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

EXPOSE 8080
CMD ["./llm-fusion-engine"]
# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder

# Create a non-root user first
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create app directory and set permissions
RUN mkdir -p /app/web && chown -R appuser:appgroup /app

# Set the working directory
WORKDIR /app/web

# Switch to the non-root user
USER appuser

# Copy package files and install dependencies
COPY web/package*.json ./
RUN npm install

# Copy the rest of the source code
COPY web/ .

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

# Create a non-root user for the final image
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create app directory and set permissions
WORKDIR /app
RUN chown -R appuser:appgroup /app

# Copy built assets
COPY --from=frontend-builder /app/web/dist ./web/dist
COPY --from=backend-builder /app/server ./server

# Switch to the non-root user
USER appuser

EXPOSE 8080
CMD ["./server"]
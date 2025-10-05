# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci --only=production=false
COPY web/ .
# 使用 npx 确保可以找到并执行 tsc 和 vite
RUN npx tsc && npx vite build

# Stage 2: Build the backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
# Copy all source code first, including go.mod
COPY . .

# Now that all source code is present, run go mod tidy to ensure
# go.mod and go.sum are complete and in sync with the code.
RUN go mod tidy

# Build the application
# The build command will now use the correctly generated go.sum
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

# Stage 3: Create the final image
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates

# Copy the built frontend from the frontend-builder stage
COPY --from=frontend-builder /app/web/dist ./web/dist

# Copy the built backend from the backend-builder stage
COPY --from=backend-builder /app/server ./server

# The database file will be created by the application on first run.
# For production, you should mount a volume to /app/fusion.db to persist data.

EXPOSE 8080

CMD ["./server"]
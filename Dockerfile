# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ .
RUN npm run build

# Stage 2: Build the backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
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
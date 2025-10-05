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
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server /app/cmd/server/main.go

# Stage 3: Create the final image
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates

# Copy the built frontend from the frontend-builder stage
COPY --from=frontend-builder /app/web/dist ./web/dist

# Copy the built backend from the backend-builder stage
COPY --from=backend-builder /app/server ./server

# Copy the database file (assuming it's pre-populated or will be created)
# For production, you'd likely mount a volume for the database.
COPY fusion.db ./fusion.db

EXPOSE 8080

CMD ["./server"]
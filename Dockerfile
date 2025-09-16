# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy source code
COPY . .

# Build arguments for environment variables
ARG WEATHER_API_KEY
ARG PORT=8080
ARG HOST=0.0.0.0

# Set build-time environment variables
ENV WEATHER_API_KEY=$WEATHER_API_KEY
ENV PORT=$PORT
ENV HOST=$HOST

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy configuration files
COPY configs/ ./configs/

# Expose port
EXPOSE 8080

# Set default environment variables (can be overridden at runtime)
ENV PORT=8080

# Run the application
CMD ["./main"]

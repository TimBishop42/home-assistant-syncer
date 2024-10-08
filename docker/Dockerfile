# Start from the official Go image
FROM golang:1.23-alpine AS builder

# Set environment variables for Go to compile for x86_64 (AMD64) architecture
ENV GOOS=linux
ENV GOARCH=amd64

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o home-assistant-syncer ./cmd/home-assistant-syncer

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/home-assistant-syncer .

# Expose port 8080 (if you have an HTTP server or metrics server)
EXPOSE 8085

# Command to run the executable
CMD ["./home-assistant-syncer"]

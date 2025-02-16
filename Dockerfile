# Use official Golang image as the builder stage
FROM golang:1.23 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy Go modules manifests first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Ensure the working directory is set to the project root
WORKDIR /app

# Build the application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/main.go

# Verify that the binary was built successfully
RUN ls -lah /app/server

# Use a minimal Alpine Linux image for final execution
FROM alpine:latest

# Set working directory inside the runtime container
WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /app/server .

# Ensure the binary is executable
RUN chmod +x server

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]

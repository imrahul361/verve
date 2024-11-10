# First stage: build the application
FROM golang:1.22-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app binary
RUN go build -o verve .

# Second stage: create a smaller image for running the app
FROM alpine:latest

# Install certificates for HTTPS (if needed by the app)
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/verve .

# Make the binary executable (in case permissions are missing)
RUN chmod +x /root/verve

# Expose the port the app will use
EXPOSE 8080

# Command to run the executable
CMD ["./verve"]

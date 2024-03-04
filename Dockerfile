# Use the official golang image as the base
FROM golang:latest AS builder

# Set the working directory for the build stage
WORKDIR /app

# Copy the application source code and additional files to the container
COPY . .

# Build the Go application
RUN go mod download && go build -o main cmd/main.go

# Use a slimmer alpine image for the final container
FROM alpine:latest

# Set the working directory for the final container
WORKDIR /app

# Copy the built Go binary and additional files to the final container
COPY --from=builder /app/main /app/main

# Expose port 3000
EXPOSE 3000

# Set the entrypoint to run the Go binary
CMD ["main"]

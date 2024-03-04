# Use the official golang image as the base
FROM golang:latest

# Set the working directory for the build stage
WORKDIR /app

# Copy the application source code and additional files to the container
ADD . .

# Build the Go application
RUN go build cmd/main.go

# Expose port 3000
EXPOSE 3000

# Set the entrypoint to run the Go binary
CMD ["/app/main"]

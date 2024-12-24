# Use an official Golang runtime as the base image
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application files to the container
COPY . .

# Get the dependencies
# RUN go mod tidy

# Build the Go application
RUN go build -o app

# Use a smaller base image for the final image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Copy the html views to the container
COPY views views

# Copy the sql migrations to the container
COPY init.sql init.sql

# Expose the port the application runs on
EXPOSE 3000

# Command to run the application
CMD ["./app"]

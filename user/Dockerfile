# Start from a base Go image
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY ./main.go .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o user

# Create a minimal container to run the Go application
FROM scratch

# Copy the built binary from the builder stage
COPY --from=builder /app/user /

# Expose the port that the application listens on
EXPOSE 8080

# Set the entry point to run the built Go application
CMD ["./user"]
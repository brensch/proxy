# First stage: build the application
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Second stage: create the runtime image
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]
# Start from the official golang image
FROM golang:alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

ENV OPENAI_API_KEY=""
ENV MONGODB_CONNECTION_STRING_SRV=""

# Expose port 3000 to the outside world
EXPOSE 80

# Command to run the executable
ENTRYPOINT ["/app/main"]

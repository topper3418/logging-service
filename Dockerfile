# -----------------------------------------------------
# Stage 1: Build the Go binary
# -----------------------------------------------------
FROM golang:1.22-alpine AS builder

# Install build dependencies for CGO (i.e. a C compiler, etc.)
RUN apk add --no-cache build-base

# Enable CGO
ENV CGO_ENABLED=1

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for caching modules)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go binary
RUN go build -o main .

# -----------------------------------------------------
# Stage 2: Create a minimal runtime image
# -----------------------------------------------------
FROM alpine:3.17

# Create a directory for the app
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# (Optional) Copy your .env file if you want it inside the image
# Typically you might pass env vars in at runtime or via Docker Compose
COPY .env /app/.env

# set and expose the port
ENV PORT=8080
EXPOSE ${PORT}

# Run the Go binary
CMD ["/app/main"]

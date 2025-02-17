# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /jixifood

# Copy go.mod and go.sum from the project root to the working directory
COPY ../../../go.mod ../../../go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project (including the specific microservice code)
COPY . .

# Build the microservice binary
RUN go build -o menu_service ./services/menu/cmd

# Stage 2: Run
FROM alpine:3.20.1

# Set the working directory
WORKDIR /jixifood

# Copy the binary from the builder stage
COPY --from=builder /jixifood/menu_service .

COPY ../../../k8s.env ./.env
RUN ls -al /jixifood/.env

# Command to run the microservice
CMD ["./menu_service"]
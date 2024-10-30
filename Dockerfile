# Start from the golang base image
FROM golang:1.22-bullseye as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY main.go .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

# Start from the distroless image
FROM gcr.io/distroless/static-debian11

# Copy the binary from the builder stage
COPY --from=builder /main /main

# Command to run the executable
CMD ["/main"]

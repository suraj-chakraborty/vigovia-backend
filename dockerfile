FROM golang:1.24-alpine

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum separately (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project source
COPY . .

# Build the Go app
RUN go build -o main .
RUN ls -la /app

# Expose port (match your app port)
EXPOSE 8080


# Run the app
CMD ["./main"]

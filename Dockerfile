FROM golang:1.20

# Set the current working directory inside the container
WORKDIR /app/go/src

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go binary
RUN go build -o main .

# Expose port 8080 for the HTTP server
EXPOSE 8080

# Start the HTTP server
CMD ["./main"]

FROM golang:1.20

# Set the current working directory inside the container
WORKDIR /app/go/src

# Install necessary packages for Postgres
RUN apt-get update && \
    apt-get install -y postgresql postgresql-contrib

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go binary
RUN go build -o main .

# Expose port 8080 for the HTTP server
EXPOSE 8080

# Expose port 5432 for PostgreSQL
EXPOSE 5432

# Set up Postgres configuration
USER postgres
RUN /etc/init.d/postgresql start && \
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" && \
    createdb -O docker my_database

# Start the HTTP server
CMD service postgresql start && ./main

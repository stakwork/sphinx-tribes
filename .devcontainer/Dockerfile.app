# .devcontainer/Dockerfile.app
FROM mcr.microsoft.com/devcontainers/go:1-ubuntu

# Set the working directory inside the container
WORKDIR /workspace

# Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Expose the port your app uses
EXPOSE 8080

# Default command (can be overridden in docker-compose)
CMD ["go", "run", "main.go"]


# Use the official Go dev container image as the base
FROM mcr.microsoft.com/devcontainers/go:0-1

# Set the working directory to where your repository is mounted
WORKDIR /workspaces/sphinx-tribes

# (Optional) Copy dependency files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your repository into the container
COPY . .

# Expose the port your app uses
EXPOSE 8080

# Set the default command to run your app automatically
CMD ["go", "run", "main.go"]


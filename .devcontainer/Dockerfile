# Use the official Go image as the base
FROM golang:1.20-bullseye

# Set the working directory inside the container
WORKDIR /workspace

# Copy go.mod and go.sum to leverage Docker caching (adjust paths if necessary)
COPY go.mod .
COPY go.sum .

# Installing xgo
RUN go install github.com/xhd2015/xgo/cmd/xgo@latest
# ENV PATH="/usr/local/xgo/bin:${PATH}"
ENV PATH="${GOPATH}/bin:${PATH}"

RUN go mod download

# Copy the rest of the source code into the container
COPY . .

# Expose the port your app listens on (adjust if needed)
EXPOSE 8080

# The command to run your app is defined in docker-compose.yml,
# so you can leave CMD empty or provide a default.
CMD [ "go", "run", "main.go" ]

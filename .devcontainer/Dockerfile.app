FROM golang:1.20-buster

WORKDIR /workspaces/sphinx-tribes

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["bash", "-c", "echo 'Container started'; go version; ls -l; tail -f /dev/null"]



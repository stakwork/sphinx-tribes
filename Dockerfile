# build stage
FROM golang:1.20.13 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

# Installing xgo
RUN go install github.com/xhd2015/xgo/cmd/xgo@latest
# ENV PATH="/usr/local/xgo/bin:${PATH}"
ENV PATH="${GOPATH}/bin:${PATH}"


RUN go mod download

COPY . .

# RUN CGO_ENABLED=0 xgo build

RUN CGO_ENABLED=0 xgo build -o sphinx-tribes


# final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates openssl

COPY --from=builder /app/sphinx-tribes /app/

RUN ls app

EXPOSE 5002
ENTRYPOINT ["/app/sphinx-tribes"]

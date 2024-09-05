# build stage
FROM golang:1.20.13 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build

# final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates openssl

COPY --from=builder /app/sphinx-tribes /app/

RUN ls app

EXPOSE 5002
ENTRYPOINT ["/app/sphinx-tribes"]
# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config.json .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build

# final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates openssl

COPY --from=builder /app/sphinx-tribes /app/
COPY --from=builder /app/config.json /app/

EXPOSE 5002
ENTRYPOINT ["/app/sphinx-tribes"]
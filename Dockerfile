FROM golang:1.21 AS base
RUN apt update && apt install -y git make ca-certificates

ENV CGO_ENABLED=1
RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.14.0

FROM base AS builder
WORKDIR /app
ENV CGO_ENABLED=0

COPY go.mod /app/
COPY go.sum /app/

RUN go mod download && go mod verify

COPY . /app
RUN go build -o bin/library -a .

## Actual runtime
FROM alpine:latest AS prod

RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    curl \
    tzdata \
    && update-ca-certificates

COPY --from=builder /app/bin/library /usr/local/bin/library
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/library"]

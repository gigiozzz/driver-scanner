FROM golang:1.24.12-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o /driver-scanner ./cmd

FROM alpine:3.22.3

RUN apk add --no-cache util-linux

COPY --from=builder /driver-scanner /usr/local/bin/driver-scanner

ENTRYPOINT ["driver-scanner"]

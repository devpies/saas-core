FROM golang:1.21-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /app

RUN mkdir log

COPY go.* ./
COPY cmd/user ./cmd/user
COPY internal/user ./internal/user

RUN go mod download && go build ./cmd/user

CMD ["/app/user"]

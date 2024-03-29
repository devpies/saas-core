FROM golang:1.21-alpine as base

LABEL org.opencontainers.image.authors="devpie"

ENV CGO_ENABLED=0

WORKDIR /app

RUN mkdir log

COPY go.* ./
COPY cmd/tenant ./cmd/tenant
COPY internal/tenant ./internal/tenant

RUN go mod download && go build ./cmd/tenant

CMD ["/app/tenant"]

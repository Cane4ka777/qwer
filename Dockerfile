# Multi-stage build for QWER Band API

# 1) Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates && update-ca-certificates

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./main.go

# 2) Runtime stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates && adduser -D -H appuser
WORKDIR /srv

COPY --from=builder /app/server /srv/server
COPY public /srv/public

ENV PORT=8080
EXPOSE 8080
USER appuser
CMD ["/srv/server"]

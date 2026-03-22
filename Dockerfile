FROM golang:1.25-alpine@sha256:8e02eb337d9e0ea459e041f1ee5eece41cbb61f1d83e7d883a3e2fb4862063fa AS deps

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.25-alpine@sha256:8e02eb337d9e0ea459e041f1ee5eece41cbb61f1d83e7d883a3e2fb4862063fa AS builder

WORKDIR /app
COPY --from=deps /app/go.mod /app/go.sum ./
COPY --from=deps /go/pkg /go/pkg
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o rest_api_app

FROM alpine:3.23.3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659

RUN apk update && apk upgrade

RUN adduser -D -u 10001 -h /home/appuser appuser
WORKDIR /home/appuser

COPY --from=builder /app/rest_api_app .
USER appuser:appuser

ENTRYPOINT ["/home/appuser/rest_api_app"]

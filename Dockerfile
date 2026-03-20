FROM golang:1.25-alpine AS deps

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY --from=deps /app/go.mod /app/go.sum ./
COPY --from=deps /go/pkg /go/pkg
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o rest_api_app

FROM alpine:3.23

COPY --from=builder /app/rest_api_app /rest_api_app
ENTRYPOINT ["/rest_api_app"]

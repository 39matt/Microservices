# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server
RUN chmod +x ./main

FROM builder AS test
RUN go test -v ./...


FROM alpine:latest AS release
WORKDIR /release

COPY --from=builder ./app/main .

EXPOSE 8110

CMD ["./main"]
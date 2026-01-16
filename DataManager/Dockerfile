# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server


FROM builder AS test
RUN go test -v ./...


FROM alpine:latest AS release
WORKDIR /release

COPY --from=builder ./app/main .

EXPOSE 8080

CMD ["./main"]
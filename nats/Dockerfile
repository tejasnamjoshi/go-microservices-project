FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./nats .
RUN go build -o /nats-service

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /nats-service /app/nats-service
CMD ["./nats-service"]
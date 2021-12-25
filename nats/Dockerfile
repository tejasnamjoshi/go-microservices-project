FROM golang:1.16-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /nats-service

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /nats-service /app/nats-service
CMD ["./nats-service"]
FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./auth ./auth
WORKDIR /app/auth
RUN go build -o /auth-service

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /auth-service /app/auth/auth-service
CMD ["./auth/auth-service"]
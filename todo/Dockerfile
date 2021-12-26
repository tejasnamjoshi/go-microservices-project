FROM golang:1.16-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./todo ./todo
WORKDIR /app/todo
RUN go build -o /todo-service

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /todo-service /app/todo/todo-service
CMD ["./todo/todo-service"]
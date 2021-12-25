buildAll:
	go build -o build/auth auth/main.go
	go build -o build/todo todo/main.go
	go build -o build/nats nats/main.go
	cp .env build/.env
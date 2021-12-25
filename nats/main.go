package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	uri := os.Getenv("NATS_URI")
	var nc *nats.Conn

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(nats.DefaultURL)
		if err == nil {
			break
		}

		fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}
	fmt.Println("Connected to NATS at:", nc.ConnectedUrl())

	fmt.Println("Worker subscribed for processing requests...")
	fmt.Println("Server listening on port 8181...")

	http.HandleFunc("/health", health)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		log.Fatal(err)
	}
}

func health(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Success"))
}

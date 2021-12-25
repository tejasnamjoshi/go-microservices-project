package main

import (
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func main() {
	l := getLogger()
	err := godotenv.Load(".env")
	if err != nil {
		l.Error("Error loading .env file")
	}
	uri := os.Getenv("NATS_URI")
	var nc *nats.Conn

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(nats.DefaultURL)
		if err == nil {
			break
		}

		l.Warn("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		l.Fatal("Error establishing connection to NATS:", err)
	}
	l.Info("Connected to NATS at:", nc.ConnectedUrl())

	l.Info("Worker subscribed for processing requests...")
	l.Info("Server listening on port 8181...")

	http.HandleFunc("/health", health)
	if err := http.ListenAndServe(":8181", nil); err != nil {
		l.Fatal(err)
	}
}

func health(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Success"))
}

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return sugar

	// l := log.New(os.Stdout, "go-todo", log.LstdFlags)
}

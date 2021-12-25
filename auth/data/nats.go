package data

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func GetNats() (*nats.Conn, error) {
	var nc *nats.Conn
	uri := os.Getenv("NATS_URI")
	var err error

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}

		fmt.Println("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("Error establishing connection to NATS: %s", err)
	}
	fmt.Println("Connected to NATS at:", nc.ConnectedUrl())

	return nc, nil
}

package data

import (
	"fmt"
	"go-todo/auth/logging"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func GetNats(logger logging.Logger) (*nats.Conn, error) {
	var nc *nats.Conn
	uri := os.Getenv("NATS_URI")
	var err error

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}

		logger.Info("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		err := fmt.Errorf("error establishing connection to NATS: %s", err)
		logger.Error(err.Error())
		return nil, err
	}
	logger.Info("Connected to NATS at:", nc.ConnectedUrl())

	return nc, nil
}

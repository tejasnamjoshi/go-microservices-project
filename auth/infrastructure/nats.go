package infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func (n *Infrastructure) GetNats() (*nats.Conn, error) {
	var nc *nats.Conn
	uri := os.Getenv("NATS_URI")
	var err error

	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}

		n.Logger.Info("Waiting before connecting to NATS at:", uri)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		err := fmt.Errorf("error establishing connection to NATS: %s", err)
		n.Logger.Error(err.Error())
		return nil, err
	}
	n.Logger.Info("Connected to NATS at:", nc.ConnectedUrl())

	return nc, nil
}

func (n *Infrastructure) InitNats() {
	nc, err := n.GetNats()
	if err != nil {
		n.Logger.Fatal(err.Error())
	}

	n.Logger.Info("Listening for auth", nc.ConnectedUrl())
	nc.Subscribe("authenticate", func(msg *nats.Msg) {
		authHeader := string(msg.Data)
		jwtParts := strings.Split(authHeader, " ")
		if len(jwtParts) <= 1 {
			nc.Publish(msg.Reply, []byte("No token available"))
			return
		}

		claims, err := n.JwtService.GetAuthorizationData(jwtParts[1])
		if err != nil {
			n.Logger.Error("Error getting authorization status")
			return
		}
		resp, err := json.Marshal(claims)
		if err != nil {
			n.Logger.Error("Error sending the data")
			return
		}
		nc.Publish(msg.Reply, []byte(resp))
	})
}

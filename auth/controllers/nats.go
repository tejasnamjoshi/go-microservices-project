package controllers

import (
	"encoding/json"
	"go-todo/auth/data"
	"strings"

	"github.com/nats-io/nats.go"
)

func (a Auth) InitNats() {
	nc, err := data.GetNats()
	if err != nil {
		a.Logger.Fatal(err)
	}

	a.Logger.Info("Listening for auth", nc.ConnectedUrl())
	nc.Subscribe("authenticate", func(msg *nats.Msg) {
		authHeader := string(msg.Data)
		jwtParts := strings.Split(authHeader, " ")
		if len(jwtParts) <= 1 {
			nc.Publish(msg.Reply, []byte("No token available"))
			return
		}

		claims, err := a.JwtService.GetAuthorizationData(jwtParts[1])
		if err != nil {
			a.Logger.Error("Error getting authorization status")
			return
		}
		resp, err := json.Marshal(claims)
		if err != nil {
			a.Logger.Error("Error sending the data")
			return
		}
		nc.Publish(msg.Reply, []byte(resp))
	})
}
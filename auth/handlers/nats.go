package handlers

import (
	"encoding/json"
	"go-todo/auth/data"
	"strings"

	"github.com/nats-io/nats.go"
)

func (a Auth) InitNats() {
	nc, err := data.GetNats()
	if err != nil {
		a.l.Fatal(err)
	}

	a.l.Info("Listening for auth", nc.ConnectedUrl())
	nc.Subscribe("authenticate", func(msg *nats.Msg) {
		authHeader := string(msg.Data)
		jwtParts := strings.Split(authHeader, " ")
		if len(jwtParts) <= 1 {
			nc.Publish(msg.Reply, []byte("No token available"))
			return
		}

		var claims = &data.CustomClaims{}
		ok, err := a.GetAuthorizationStatus(authHeader)
		if ok {
			claims, err = data.ParseJWT(jwtParts[1])
			if err != nil {
				a.l.Error("Error parsing token")
			}
		}
		resp, err := json.Marshal(*claims)
		nc.Publish(msg.Reply, []byte(resp))
	})
}

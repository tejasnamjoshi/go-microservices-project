package handlers

import "log"

type Auth struct {
	l *log.Logger
}

func NewAuth(l *log.Logger) *Auth {
	return &Auth{l}
}

package logging

type Logger interface {
	Error(err string)
	Warn(warn string)
	Info(info ...string)
	Fatal(info string)
}

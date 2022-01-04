package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type chiRouter struct{}

var (
	chiDispatcher = chi.NewRouter()
)

func NewChiRouter() Router {
	return &chiRouter{}
}

func (*chiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Get(uri, f)
}

func (*chiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Post(uri, f)
}

func (*chiRouter) PATCH(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Patch(uri, f)
}

func (*chiRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Put(uri, f)
}

func (*chiRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Delete(uri, f)
}

func (*chiRouter) USE(next func(http.Handler) http.Handler) {
	chiDispatcher.Use(next)
}

func (c *chiRouter) WITH(next func(http.Handler) http.Handler) Router {
	chiDispatcher.With(next)
	return c
}

func (*chiRouter) SERVE(port string) {
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), chiDispatcher)
	if err != nil {
		fmt.Printf("Error starting server: %v", err)
		os.Exit(1)
	}
}

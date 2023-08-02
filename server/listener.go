package server

import (
	"fmt"
	"net/http"
	"os"
)

type HttpListener struct {
	serverMux *http.ServeMux
}

func (listener *HttpListener) Init() {
	listener.serverMux = http.NewServeMux()
}

func (listener *HttpListener) RegisterHandler(path string, callback http.HandlerFunc) {
	listener.serverMux.Handle(path, callback)
}

func (listener *HttpListener) Listen(defaultPort int16) error {
	envPort := os.Getenv("PORT")
	if envPort != "" {
		return http.ListenAndServe(":"+envPort, listener.serverMux)
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), listener.serverMux)
}

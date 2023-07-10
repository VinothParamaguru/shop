package server

import (
	"fmt"
	"net/http"
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

func (listener *HttpListener) Listen(port int16) error {
	return http.ListenAndServe(":"+fmt.Sprintf("%d", port), listener.serverMux)
}

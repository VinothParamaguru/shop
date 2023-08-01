package server

import (
	"shop/core"
)

// Start  the http server
func Start() {

	// Create a http listener object to handle incoming http connections
	// init and register the handler functions
	listener := &HttpListener{}
	listener.Init()
	listener.RegisterHandler("/Register", core.RegisterUser)
	listener.RegisterHandler("/Login", core.LoginUser)
	err := listener.Listen(8080)
	if err != nil {
		panic(err)
	}
}

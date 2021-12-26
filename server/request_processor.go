package server

import (
	"net/http"
	"workspace/shop/core"
)

// start the http server
func Start() {

	http.HandleFunc("/Register", core.RegisterUser)
	http.HandleFunc("/Login", core.Login)

	http.ListenAndServe(":5000", nil)

}

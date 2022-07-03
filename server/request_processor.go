package server

import (
	"net/http"
	"workspace/shop/core"
)

// Start  the http server
func Start() {

	http.HandleFunc("/Register", core.RegisterUser)
	http.HandleFunc("/Login", core.Login)

	http.ListenAndServe(":8000", nil)

}

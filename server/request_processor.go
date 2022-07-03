package server

import (
	"net/http"
	"workspace/shop/core"
)

// Start  the http server
func Start() {

	http.HandleFunc("/Register", core.RegisterUser)
	http.HandleFunc("/Login", core.Login)

	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		panic(err)
	}
}

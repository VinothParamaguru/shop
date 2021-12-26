package core

import (
	"fmt"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Registration... called")
}

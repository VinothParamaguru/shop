package core

import (
	"fmt"
	"net/http"
)

func Login(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Login... called")
}

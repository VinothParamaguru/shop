package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"workspace/shop/security"

	"workspace/shop/utilities"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginUser(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Login... called")

	payload, err := ioutil.ReadAll(req.Body)
	utilities.HandlePanic(err)

	var loginParams Login
	err = json.Unmarshal(payload, &loginParams)
	utilities.HandlePanic(err)

	status, code := security.ValidateRequiredFields([]string{loginParams.Username,
		loginParams.Password})
	if !status {
		utilities.HandleError(w, status, code)
	}

	status, code = security.ValidateInput("email", loginParams.Username)
	if !status {
		utilities.HandleError(w, status, code)
	}

	status, code = security.ValidateInput("password", loginParams.Username)
	if !status {
		utilities.HandleError(w, status, code)
	}

	// Get the stored password
	
}

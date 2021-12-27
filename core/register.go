package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"workspace/shop/security"
)

// define the struct for registration

type Register struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

func RegisterUser(http_response_writer http.ResponseWriter, http_request *http.Request) {

	fmt.Println("Registration... called")
	payload, error := ioutil.ReadAll(http_request.Body)
	if error != nil {
		panic(error)
	}
	var register_params Register
	error = json.Unmarshal(payload, &register_params)
	if error != nil {
		panic(error)
	}
	status := security.Validate("email", register_params.Email)
	fmt.Println(status)

}

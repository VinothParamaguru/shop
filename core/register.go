package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	//app_db "workspace/shop/database"
	"workspace/shop/security"
	"workspace/shop/utilities"
)

// define the struct for registration

type Register struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

// Register the user
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

	// Validate all the fields

	// required field validation
	status, code := security.ValidateRequiredFields([]string{register_params.Email,
		register_params.Firstname, register_params.Lastname})
	if !status {
		utilities.HandleError(http_response_writer, status, code)
		return
	}
	database_config := utilities.GetDataBaseConfig()

	fmt.Println(database_config.Database)
	// input fields validation
	status, code = security.ValidateInput("email", register_params.Email)
	if !status {
		utilities.HandleError(http_response_writer, status, code)
		fmt.Println("email")
		return
	}
	status, code = security.ValidateInput("name", register_params.Firstname)
	if !status {
		utilities.HandleError(http_response_writer, status, code)
		fmt.Println("fname")
		return
	}
	status, code = security.ValidateInput("name", register_params.Lastname)
	if !status {
		utilities.HandleError(http_response_writer, status, code)
		fmt.Println("lname")
		return
	}

}

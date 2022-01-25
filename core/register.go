package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	app_db "workspace/shop/database"
	"workspace/shop/security"
	"workspace/shop/utilities"
)

// define the struct for registration

type Register struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DateOfBirth string `json:"date_of_birth"`
}

// Register the user
func RegisterUser(http_response_writer http.ResponseWriter,
	http_request *http.Request) {

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

	// required field validation
	status, code := security.ValidateRequiredFields([]string{register_params.Email,
		register_params.Password})
	if !status {
		utilities.HandleError(http_response_writer, status, code)
		return
	}

	// input fields validation
	status, code = security.ValidateInput("email", register_params.Email)
	if !status {
		utilities.HandleError(http_response_writer, status, code)
		return
	}

	status, code = security.ValidateInput("password", register_params.Password)
	if !status {
		utilities.HandleError(http_response_writer, status, code)
		return
	}

	if register_params.Firstname != "" {
		status, code = security.ValidateInput("name", register_params.Firstname)
		if !status {
			utilities.HandleError(http_response_writer, status, code)
			return
		}
	}

	if register_params.Lastname != "" {
		status, code = security.ValidateInput("name", register_params.Lastname)
		if !status {
			utilities.HandleError(http_response_writer, status, code)
			return
		}
	}

	// get database configuration
	database_config := utilities.GetDataBaseConfig()
	fmt.Println(database_config.Schema)
	db := app_db.DataBase{Connector: nil, Config: database_config}

	password_seed := utilities.GetRandomString(len(register_params.Password))

	fmt.Println(password_seed)

	temporary_hash := utilities.GenerateHash(register_params.Password + password_seed)

	final_hash := utilities.GenerateHash(temporary_hash)

	// open the database
	if status, code = db.Open(); !status {
		utilities.HandleError(http_response_writer, status, code)
	}

	// perform insert
	fields := []app_db.Field{
		{Name: "username", Value: "Vinoth1"},
		{Name: "password", Value: final_hash},
		{Name: "password_seed", Value: password_seed},
	}

	db.Insert("users", fields)
	fmt.Println(db.Config.Schema)
}

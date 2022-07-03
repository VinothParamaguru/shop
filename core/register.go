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

// RegisterUser Register the user
func RegisterUser(httpResponseWriter http.ResponseWriter,
	httpRequest *http.Request) {

	fmt.Println("Registration... called")
	payload, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		panic(err)
	}
	var registerParams Register
	err = json.Unmarshal(payload, &registerParams)
	if err != nil {
		panic(err)
	}

	// required field validation
	status, code := security.ValidateRequiredFields([]string{registerParams.Email,
		registerParams.Password})
	if !status {
		utilities.HandleError(httpResponseWriter, status, code)
		return
	}

	// input fields validation
	status, code = security.ValidateInput("email", registerParams.Email)
	if !status {
		utilities.HandleError(httpResponseWriter, status, code)
		return
	}

	status, code = security.ValidateInput("password", registerParams.Password)
	if !status {
		utilities.HandleError(httpResponseWriter, status, code)
		return
	}

	if registerParams.Firstname != "" {
		status, code = security.ValidateInput("name", registerParams.Firstname)
		if !status {
			utilities.HandleError(httpResponseWriter, status, code)
			return
		}
	}

	if registerParams.Lastname != "" {
		status, code = security.ValidateInput("name", registerParams.Lastname)
		if !status {
			utilities.HandleError(httpResponseWriter, status, code)
			return
		}
	}

	// get database configuration
	databaseConfig := utilities.GetDataBaseConfig()
	fmt.Println(databaseConfig.Schema)
	db := app_db.DataBase{Connector: nil, Config: databaseConfig}

	passwordSeed := utilities.GetRandomString(len(registerParams.Password))

	fmt.Println(passwordSeed)

	temporaryHash := utilities.GenerateHash(registerParams.Password + passwordSeed)

	finalHash := utilities.GenerateHash(temporaryHash)

	// open the database
	if status, code = db.Open(); !status {
		utilities.HandleError(httpResponseWriter, status, code)
	}

	// perform insert
	fields := []app_db.Field{
		{Name: "username", Value: "Vinoth1"},
		{Name: "password", Value: finalHash},
		{Name: "password_seed", Value: passwordSeed},
	}

	db.Insert("users", fields)
	fmt.Println(db.Config.Schema)
}

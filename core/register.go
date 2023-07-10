package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"workspace/shop/request"
	"workspace/shop/response"

	app_db "workspace/shop/database"
	"workspace/shop/security"
	"workspace/shop/utilities"
)

// define the struct for registration

type Register struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DateOfBirth string `json:"date_of_birth"`
}

// RegisterUser Register the user
func RegisterUser(httpResponseWriter http.ResponseWriter,
	httpRequest *http.Request) {

	requestProcessor := request.Processor{}
	responseProcessor := response.Processor{}

	// extract the json payload from the request
	// with some basic checks
	payload, err := requestProcessor.ReadRequest(httpRequest)
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}
	var registerParams Register
	err = json.Unmarshal(payload, &registerParams)
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	// required field validation
	err = security.ValidateRequiredFields([]string{registerParams.Username,
		registerParams.Password})
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	// input fields validation
	err = security.ValidateInput("email", registerParams.Username)
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	err = security.ValidateInput("password", registerParams.Password)
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	if registerParams.Firstname != "" {
		err = security.ValidateInput("name", registerParams.Firstname)
		if err != nil {
			responseProcessor.SendError(err, httpResponseWriter)
			return
		}
	}

	if registerParams.Lastname != "" {
		err = security.ValidateInput("name", registerParams.Lastname)
		if err != nil {
			responseProcessor.SendError(err, httpResponseWriter)
			return
		}
	}

	// get database configuration
	databaseConfig := app_db.GetDataBaseConfig()
	fmt.Println(databaseConfig.Schema)
	db := app_db.DataBase{Connector: nil, Config: databaseConfig}

	passwordSeed := utilities.GetRandomString(len(registerParams.Password))

	fmt.Println(passwordSeed)

	temporaryHash := utilities.GenerateHash(registerParams.Password + passwordSeed)

	hashedPassword := utilities.GenerateHash(temporaryHash)

	// open the database
	err = db.Open()
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	// perform insert
	fields := []app_db.Field{
		{Name: "username", Value: registerParams.Username},
		{Name: "password", Value: hashedPassword},
		{Name: "password_seed", Value: passwordSeed},
	}

	if registerParams.Firstname != "" {
		fields = append(fields, app_db.Field{Name: "firstname", Value: registerParams.Firstname})
	}
	if registerParams.Lastname != "" {
		fields = append(fields, app_db.Field{Name: "lastname", Value: registerParams.Lastname})
	}

	db.Insert("users", fields)
	fmt.Println(db.Config.Schema)
}

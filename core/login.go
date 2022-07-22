package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	app_db "workspace/shop/database"
	"workspace/shop/errors"
	"workspace/shop/security"
	"workspace/shop/utilities"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginUser(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	fmt.Println("Login... called")

	payload, err := ioutil.ReadAll(httpRequest.Body)
	utilities.HandlePanic(err)

	var loginParams Login
	err = json.Unmarshal(payload, &loginParams)
	utilities.HandlePanic(err)

	status, code := security.ValidateRequiredFields([]string{loginParams.Username,
		loginParams.Password})
	if !status {
		utilities.HandleSecurityError(httpResponseWriter, status, code)
	}

	status, code = security.ValidateInput("email", loginParams.Username)
	if !status {
		utilities.HandleSecurityError(httpResponseWriter, status, code)
	}

	status, code = security.ValidateInput("password", loginParams.Password)
	if !status {
		utilities.HandleSecurityError(httpResponseWriter, status, code)
	}

	// get database configuration
	databaseConfig := app_db.GetDataBaseConfig()
	fmt.Println(databaseConfig.Schema)
	db := app_db.DataBase{Connector: nil, Config: databaseConfig}

	// open the database
	if status, code = db.Open(); !status {
		utilities.HandleDataBaseError(httpResponseWriter, status, code)
	}

	db.InitSql("SELECT username, password as hashed_password, " +
		"password_seed  FROM users where username = @username")
	db.BindParam("@username", loginParams.Username)
	status, code, results := db.Select()

	var (
		username             string
		hashedPasswordStored string
		passwordSeed         string
	)

	if status && results.Next() {
		errorInfo := results.Scan(&username, &hashedPasswordStored, &passwordSeed)
		if errorInfo != nil {
			utilities.HandleDataBaseError(httpResponseWriter, false,
				errors.DbErrorQueryExecution)
		}
		temporaryHash := utilities.GenerateHash(loginParams.Password + passwordSeed)
		hashedPassword := utilities.GenerateHash(temporaryHash)
		if hashedPasswordStored != hashedPassword {
			utilities.HandleApplicationError(httpResponseWriter, false,
				errors.AppInvalidUserNameOrPassword)
		} else {
			type LoginResponse struct {
				Id int
			}
			loginResponse := LoginResponse{Id: 1}
			utilities.SendResponse(httpResponseWriter, loginResponse)
		}
	} else {
		utilities.HandleDataBaseError(httpResponseWriter, status, code)
	}

}

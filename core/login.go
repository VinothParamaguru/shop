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

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Token string `json:"token"`
	Id    int    `json:"id"`
}

func LoginUser(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	fmt.Println("Login... called")

	payload, err := ioutil.ReadAll(httpRequest.Body)
	utilities.HandlePanic(err)

	var loginParams Request
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

	db.InitSql("SELECT id, password as hashed_password, " +
		"password_seed  FROM users where username = @username")
	db.BindParam("@username", loginParams.Username)
	status, code, results := db.Select()

	var (
		id                   int
		hashedPasswordStored string
		passwordSeed         string
	)

	if status && results.Next() {
		errorInfo := results.Scan(&id, &hashedPasswordStored, &passwordSeed)
		if errorInfo != nil {
			utilities.HandleDataBaseError(httpResponseWriter, false,
				errors.DbErrorQueryExecution)
		}
		temporaryHash := utilities.GenerateHash(loginParams.Password + passwordSeed)
		hashedPassword := utilities.GenerateHash(temporaryHash)
		if hashedPasswordStored != hashedPassword {

			// login failed
			utilities.HandleApplicationError(httpResponseWriter, false,
				errors.AppInvalidUserNameOrPassword)
		} else {

			// login successful
			sessionToken := utilities.GenerateRandomToken()
			fields := []app_db.Field{
				{Name: "id", Value: id},
				{Name: "token", Value: sessionToken},
			}
			db.Insert("session", fields)
			loginResponse := Response{Id: id, Token: sessionToken}
			utilities.SendResponse(httpResponseWriter, loginResponse)
		}
	} else {
		utilities.HandleDataBaseError(httpResponseWriter, status, code)
	}

}

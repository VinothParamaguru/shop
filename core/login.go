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
		utilities.HandleError(httpResponseWriter, status, code)
	}

	status, code = security.ValidateInput("email", loginParams.Username)
	if !status {
		utilities.HandleError(httpResponseWriter, status, code)
	}

	status, code = security.ValidateInput("password", loginParams.Username)
	if !status {
		utilities.HandleError(httpResponseWriter, status, code)
	}

	// get database configuration
	databaseConfig := app_db.GetDataBaseConfig()
	fmt.Println(databaseConfig.Schema)
	db := app_db.DataBase{Connector: nil, Config: databaseConfig}

	// open the database
	if status, code = db.Open(); !status {
		utilities.HandleError(httpResponseWriter, status, code)
	}

	db.InitSql("SELECT * FROM users where username = @username")

}

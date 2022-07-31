package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	app_db "workspace/shop/database"
	"workspace/shop/errors"
	"workspace/shop/security"
	"workspace/shop/static"
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
		return
	}

	status, code = security.ValidateInput("email", loginParams.Username)
	if !status {
		utilities.HandleSecurityError(httpResponseWriter, status, code)
		return
	}

	status, code = security.ValidateInput("password", loginParams.Password)
	if !status {
		utilities.HandleSecurityError(httpResponseWriter, status, code)
		return
	}

	// get database configuration
	databaseConfig := app_db.GetDataBaseConfig()
	fmt.Println(databaseConfig.Schema)
	db := app_db.DataBase{Connector: nil, Config: databaseConfig}

	// open the database
	if status, code = db.Open(); !status {
		utilities.HandleDataBaseError(httpResponseWriter, status, code)
		return
	}

	db.InitSql("SELECT id, password as hashed_password, " +
		"password_seed  FROM users where username = @username")
	db.BindParam("@username", loginParams.Username)
	status, code, userResults := db.Select()

	var (
		userId               int
		hashedPasswordStored string
		passwordSeed         string
	)

	if status && userResults.Next() {

		errorInfo := userResults.Scan(&userId, &hashedPasswordStored, &passwordSeed)
		if errorInfo != nil {
			utilities.HandleDataBaseError(httpResponseWriter, false,
				errors.DbErrorScanning)
			return
		}
		temporaryHash := utilities.GenerateHash(loginParams.Password + passwordSeed)
		hashedPassword := utilities.GenerateHash(temporaryHash)
		if hashedPasswordStored != hashedPassword {

			// login failed
			utilities.HandleApplicationError(httpResponseWriter, false,
				errors.AppInvalidUserNameOrPassword)

			// delete the session token associated with the user here
			return
		} else {

			// login successful

			// check if there is an ongoing session for the current user
			// if there is one, delete the session by deleting the session token

			db.InitSql("SELECT id FROM session where fkid_users = @id")
			db.BindParam("@id", userId)
			status, code, sessionResults := db.Select()

			if status && sessionResults.Next() {

				// login successful but there is a session token
				// already
				// Active session existing update the session token
				sessionToken := utilities.GenerateRandomToken()
				updateStatus, updateCode := updateSession(&db, userId, sessionToken)
				if !updateStatus {
					utilities.HandleDataBaseError(httpResponseWriter, updateStatus, updateCode)
					return
				}
				loginResponse := Response{Id: userId, Token: sessionToken}
				utilities.SendResponse(httpResponseWriter, loginResponse)

			} else if status && !sessionResults.Next() {

				// login successful for the first time
				sessionToken := utilities.GenerateRandomToken()
				createSession(&db, userId, sessionToken)
				loginResponse := Response{Id: userId, Token: sessionToken}
				utilities.SendResponse(httpResponseWriter, loginResponse)

			} else if !status {
				utilities.HandleDataBaseError(httpResponseWriter, status, code)
				return
			}
		}
	} else {
		utilities.HandleDataBaseError(httpResponseWriter, status, code)
		return
	}

}

func createSession(db *app_db.DataBase, userId int, token string) {

	fields := []app_db.Field{
		{Name: "fkid_users", Value: userId},
		{Name: "fkid_token_types", Value: static.GetSessionTokenType()},
		{Name: "token", Value: token},
		{Name: "session_start_time",
			Value: utilities.GetCurrentTimeStampString()},
	}
	db.Insert("session", fields)
}

func updateSession(db *app_db.DataBase, userId int, token string) (status bool, code int) {

	db.InitSql("UPDATE session SET token = @token, " +
		"session_start_time = @session_start_time " +
		"WHERE fkid_users = @user_id")
	db.BindParam("@token", token)
	db.BindParam("@session_start_time",
		utilities.GetCurrentTimeStampString())
	db.BindParam("@user_id", userId)
	return db.Execute()
}

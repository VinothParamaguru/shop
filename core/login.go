package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	appdb "shop/database"
	apperrors "shop/errors"
	"shop/request"
	"shop/response"
	"shop/security"
	"shop/static"
	"shop/utilities"
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

	requestProcessor := request.Processor{}
	responseProcessor := response.Processor{}

	payload, err := requestProcessor.ReadRequest(httpRequest)
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}
	var loginParams Request
	err = json.Unmarshal(payload, &loginParams)
	utilities.HandlePanic(err)

	err = security.ValidateRequiredFields([]string{loginParams.Username,
		loginParams.Password})
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	err = security.ValidateInput("email", loginParams.Username)
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	err = security.ValidateInput("password", loginParams.Password)
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	// get database configuration
	databaseConfig := appdb.GetDataBaseConfig()
	fmt.Println(databaseConfig.Schema)
	db := appdb.DataBase{Connector: nil, Config: databaseConfig}

	// open the database
	err = db.Open()
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

	db.InitSql("SELECT id, password as hashed_password, " +
		"password_seed  FROM users where username = @username")
	db.BindParam("@username", loginParams.Username)
	userResults, err := db.Select()

	var (
		userId               int
		hashedPasswordStored string
		passwordSeed         string
	)

	if err == nil && userResults.Next() {
		err = userResults.Scan(&userId, &hashedPasswordStored, &passwordSeed)
		if err != nil {
			responseProcessor.SendError(err, httpResponseWriter)
			return
		}
		temporaryHash := utilities.GenerateHash(loginParams.Password + passwordSeed)
		hashedPassword := utilities.GenerateHash(temporaryHash)
		if hashedPasswordStored != hashedPassword {

			// login failed
			responseProcessor.SendError(
				errors.New(apperrors.ApplicationErrorDescriptions[apperrors.AppInvalidUserNameOrPassword]),
				httpResponseWriter)
			// delete the session token associated with the user here
			return
		} else {

			// login successful

			// check if there is an ongoing session for the current user
			// if there is one, delete the session by deleting the session token

			db.InitSql("SELECT id FROM session where fkid_users = @id")
			db.BindParam("@id", userId)
			sessionResults, err := db.Select()

			if err == nil && sessionResults.Next() {

				// login successful but there is a session token
				// already
				// Active session existing update the session token
				sessionToken := utilities.GenerateRandomToken()
				err = updateSession(&db, userId, sessionToken)
				if err != nil {
					responseProcessor.SendError(err, httpResponseWriter)
					return
				}
				loginResponse := Response{Id: userId, Token: sessionToken}
				utilities.SendResponse(httpResponseWriter, loginResponse)

			} else if err == nil && !sessionResults.Next() {

				// login successful for the first time
				sessionToken := utilities.GenerateRandomToken()
				createSession(&db, userId, sessionToken)
				loginResponse := Response{Id: userId, Token: sessionToken}
				utilities.SendResponse(httpResponseWriter, loginResponse)

			} else if err != nil {
				responseProcessor.SendError(err, httpResponseWriter)
				return
			}
		}
	} else {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}

}

func createSession(db *appdb.DataBase, userId int, token string) {

	fields := []appdb.Field{
		{Name: "fkid_users", Value: userId},
		{Name: "fkid_token_types", Value: static.GetSessionTokenType()},
		{Name: "token", Value: utilities.GenerateHash(token)},
		{Name: "session_start_time",
			Value: utilities.GetCurrentTimeStampString()},
	}
	db.Insert("session", fields)
}

func updateSession(db *appdb.DataBase, userId int, token string) error {

	db.InitSql("UPDATE session SET token = @token, " +
		"session_start_time = @session_start_time " +
		"WHERE fkid_users = @user_id")
	db.BindParam("@token", utilities.GenerateHash(token))
	db.BindParam("@session_start_time",
		utilities.GetCurrentTimeStampString())
	db.BindParam("@user_id", userId)
	return db.Execute()
}

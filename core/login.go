package core

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
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
	Otp      string `json:"otp"`
}

type Response struct {
	Token     any `json:"token,omitempty"`
	SecretKey any `json:"secret_key,omitempty"`
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
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}
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
	db := appdb.DataBase{Connector: nil, Config: databaseConfig}
	// open the database
	err = db.Open()
	if err != nil {
		responseProcessor.SendError(err, httpResponseWriter)
		return
	}
	defer db.Close()
	db.InitSql("SELECT id, password as hashed_password, " +
		"password_seed, secret_key  FROM users where username = @username")
	_ = db.BindParam("@username", loginParams.Username)
	userResults, err := db.Select()
	var (
		userId               int
		hashedPasswordStored string
		passwordSeed         string
		secretKey            sql.NullString
	)
	if err == nil && userResults.Next() {
		err = userResults.Scan(&userId, &hashedPasswordStored, &passwordSeed, &secretKey)
		if err != nil {
			responseProcessor.SendError(err, httpResponseWriter)
			return
		}
		temporaryHash := utilities.GenerateHash(loginParams.Password + passwordSeed)
		hashedPassword := utilities.GenerateHash(temporaryHash)
		if hashedPasswordStored != hashedPassword {
			// wrong password, login fails
			responseProcessor.SendError(
				apperrors.GetError(apperrors.AppInvalidUserNameOrPassword),
				httpResponseWriter)
			// delete the session token associated with the user here
			return
		} else {
			// login successful
			// do otp validation if otp is used in the request
			// all the successful Login requests doesn't contain otp are used
			// to generate secrets and treated as if the user is re-initialising the
			// secret key creation
			sessionToken := utilities.GenerateRandomToken()
			if loginParams.Otp == "" { // user didn't provide otp, secret init call
				key, _ := totp.Generate(totp.GenerateOpts{
					Issuer:      "bitssimplified.com",
					AccountName: loginParams.Username,
					Algorithm:   otp.AlgorithmSHA1,
				})
				_ = updateSecretKey(&db, userId, key.Secret())
				_ = sendQRCodeToUser(httpResponseWriter, key)
				return
			} else { // user provided otp
				if secretKey.Valid && totp.Validate(loginParams.Otp, secretKey.String) {
					loginResponse := Response{SecretKey: nil, Token: sessionToken}
					utilities.SendResponse(httpResponseWriter, loginResponse)
				} else {
					// otp wrong, login fails
					responseProcessor.SendError(
						apperrors.GetError(apperrors.AppInvalidUserNameOrPassword),
						httpResponseWriter)
					return
				}
			}
			// check if there is an ongoing session for the current user
			// if there is one, delete the session by deleting the session token
			db.InitSql("SELECT id FROM session where fkid_users = @id")
			_ = db.BindParam("@id", userId)
			sessionResults, err := db.Select()
			if err == nil && sessionResults.Next() {
				// login successful but there is a session token already for the user
				// Active session existing update the session token
				_ = updateSession(&db, userId, sessionToken)
			} else if err == nil && !sessionResults.Next() {
				// login successful for the first time, no session exists for the user
				// create new session
				_ = createSession(&db, userId, sessionToken)
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

func createSession(db *appdb.DataBase, userId int, token string) error {

	fields := []appdb.Field{
		{Name: "fkid_users", Value: userId},
		{Name: "fkid_token_types", Value: static.GetSessionTokenType()},
		{Name: "token", Value: utilities.GenerateHash(token)},
		{Name: "session_start_time",
			Value: utilities.GetCurrentTimeStampString()},
	}
	return db.Insert("session", fields)
}

func updateSession(db *appdb.DataBase, userId int, token string) error {

	db.InitSql("UPDATE session SET token = @token, " +
		"session_start_time = @session_start_time " +
		"WHERE fkid_users = @user_id")
	_ = db.BindParam("@token", utilities.GenerateHash(token))
	_ = db.BindParam("@session_start_time",
		utilities.GetCurrentTimeStampString())
	_ = db.BindParam("@user_id", userId)
	return db.Execute()
}

func updateSecretKey(db *appdb.DataBase, userId int, secretKey string) error {

	db.InitSql("UPDATE users SET secret_key = @secret_key WHERE id = @user_id")
	_ = db.BindParam("@secret_key", secretKey)
	_ = db.BindParam("@user_id", userId)
	return db.Execute()
}

func sendQRCodeToUser(httpResponseWriter http.ResponseWriter, key *otp.Key) error {
	// Convert TOTP key into a PNG
	var buffer bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return err
	}
	_ = png.Encode(&buffer, img)
	httpResponseWriter.WriteHeader(http.StatusOK)
	httpResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	_, _ = httpResponseWriter.Write(buffer.Bytes())
	return nil
}

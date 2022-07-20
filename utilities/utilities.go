package utilities

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"workspace/shop/errors"
)

// ExtractPayloadAsString extract the http payload as string
func ExtractPayloadAsString(httpRequest *http.Request) string {

	payload, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		panic(err)
	}
	// restore the content of body if needed
	return string(payload)
}

// WriteMessage Uses the response writer object to write the message to the client
func WriteMessage() {

}

// HandleSecurityError . Handle the security error using the api return status and error code returned
func HandleSecurityError(httpResponseWriter http.ResponseWriter, status bool, code int) {

	if !status && code != errors.Success {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusBadRequest)
		responseParams := errors.ErrorResponse{Code: code, Description: errors.SecurityErrorDescriptions[code]}
		jsonResponse, err := json.Marshal(responseParams)
		if err != nil {
			panic(err)
		}
		httpResponseWriter.Write(jsonResponse)
	}
}

// HandleDataBaseError . Handle the database error using the api return status and error code returned
func HandleDataBaseError(httpResponseWriter http.ResponseWriter, status bool, code int) {

	if !status && code != errors.Success {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusBadRequest)
		responseParams := errors.ErrorResponse{Code: code, Description: errors.DataBaseErrorDescriptions[code]}
		jsonResponse, err := json.Marshal(responseParams)
		if err != nil {
			panic(err)
		}
		httpResponseWriter.Write(jsonResponse)
	}
}

// HandleApplicationError . Handle the application error using the api return status and error code returned
func HandleApplicationError(httpResponseWriter http.ResponseWriter, status bool, code int) {

	if !status && code != errors.Success {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusBadRequest)
		responseParams := errors.ErrorResponse{Code: code, Description: errors.SecurityErrorDescriptions[code]}
		jsonResponse, err := json.Marshal(responseParams)
		if err != nil {
			panic(err)
		}
		httpResponseWriter.Write(jsonResponse)
	}
}

// generate the hash using the provided input string
// this function uses sha256 algorithm to generate hashes

func GenerateHash(input string) string {

	sha256Hash := sha256.New()

	if _, errorInfo := sha256Hash.Write([]byte(input)); errorInfo != nil {
		panic(errorInfo)
	}

	byteSlice := sha256Hash.Sum(nil)

	return fmt.Sprintf("%x", byteSlice)

}

func GetRandomNumber(minimum int, maximum int) int {
	// copied from Go cook book
	return rand.Intn(maximum-minimum) + minimum
}

func GetRandomString(length int) string {

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		// Get character equivalent of random numbers
		// in the range 33 - 126 ascii
		// this makes a good mix of letters, alphabets
		// and special characters
		bytes[i] = byte(GetRandomNumber(33, 126))
	}
	return string(bytes)
}

func HandlePanic(err any) {
	if err != nil {
		panic(err)
	}
}

func ParamInAcceptableCharRange(paramChar byte) bool {
	if paramChar >= '0' || paramChar <= '8' ||
		paramChar >= 'A' || paramChar <= 'Z' ||
		paramChar >= 'a' || paramChar <= 'z' {
		return true
	}
	return false
}

func GetSqlParams(sql string) []string {

	var sqlParams []string

	sqlLength := len(sql)
	sqlParsed := false
	for pos, char := range sql {
		if char == '@' {
			param := "@"
			for i := pos + 1; i < sqlLength && ParamInAcceptableCharRange(sql[i]) && sql[i] != ' '; i++ {
				param += string(sql[i])
				if i+1 == sqlLength {
					sqlParsed = true
				}
			}
			sqlParams = append(sqlParams, param)
			if sqlParsed == true {
				break
			}
		}
	}

	return sqlParams

}

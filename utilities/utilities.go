package utilities

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	app_db "workspace/shop/database"
	app_error "workspace/shop/error"
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

// HandleError Handle the error using the api return status and error code returned
func HandleError(httpResponseWriter http.ResponseWriter, status bool, code int) {

	if !status && code != app_error.Success {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusBadRequest)
		responseParams := app_error.ErrorResponse{Code: code, Description: app_error.ErrorDescriptions[code]}
		jsonResponse, err := json.Marshal(responseParams)
		if err != nil {
			panic(err)
		}
		httpResponseWriter.Write(jsonResponse)
	}
}

func GetDataBaseConfig() app_db.DataBaseConfig {

	// read the database configuration file
	// using relative path for now. This should be replaced by absolute path of the file
	configFileRelativePath := "config/database_config.json"
	configFileAbsolutePath, errorInfo := filepath.Abs(configFileRelativePath)

	if errorInfo != nil {
		panic(errorInfo)
	}

	// read
	data, errorInfo := ioutil.ReadFile(configFileAbsolutePath)
	if errorInfo != nil {
		panic(errorInfo)
	}

	// unmarshal
	var databaseConfig app_db.DataBaseConfig
	errorInfo = json.Unmarshal(data, &databaseConfig)
	if errorInfo != nil {
		panic(errorInfo)
	}

	return databaseConfig
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

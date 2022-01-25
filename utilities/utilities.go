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

// extract the http payload as string
func ExtractPayloadAsString(http_request *http.Request) string {

	payload, error := ioutil.ReadAll(http_request.Body)
	if error != nil {
		panic(error)
	}
	// restore the content of body if needed
	return string(payload)
}

// Uses the response writer object to write the message to the client
func WriteMessage(http_response_writer *http.ResponseWriter, message string) {

}

// Handle the error using the api return status and error code returned
func HandleError(http_response_writer http.ResponseWriter, status bool, code int) {

	if !status && code != app_error.Success {
		http_response_writer.Header().Set("Content-Type", "application/json")
		http_response_writer.WriteHeader(http.StatusBadRequest)
		response_params := app_error.ErrorResponse{Code: code, Description: app_error.ErrorDescriptions[code]}
		json_response, error := json.Marshal(response_params)
		if error != nil {
			panic(error)
		}
		http_response_writer.Write(json_response)
	}
}

func GetDataBaseConfig() app_db.DataBaseConfig {

	// read the database configuration file
	// using relative path for now. This should be replaced by absolute path of the file
	config_file_relative_path := "config/database_config.json"
	config_file_absolute_path, error_info := filepath.Abs(config_file_relative_path)

	if error_info != nil {
		panic(error_info)
	}

	// read
	data, error_info := ioutil.ReadFile(config_file_absolute_path)
	if error_info != nil {
		panic(error_info)
	}

	// unmarshal
	var database_config app_db.DataBaseConfig
	error_info = json.Unmarshal(data, &database_config)
	if error_info != nil {
		panic(error_info)
	}

	return database_config
}

// generate the hash using the provided input string
// this function uses sha256 algorithm to generate hashes

func GenerateHash(input string) string {

	sha256_hash := sha256.New()

	if _, error_info := sha256_hash.Write([]byte(input)); error_info != nil {
		panic(error_info)
	}

	byte_slice := sha256_hash.Sum(nil)

	return fmt.Sprintf("%x", byte_slice)

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

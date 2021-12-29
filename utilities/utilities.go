package utilities

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

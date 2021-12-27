package utilities

import (
	"io/ioutil"
	"net/http"
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

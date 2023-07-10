package response

import (
	"encoding/json"
	"log"
	"net/http"
	"workspace/shop/errors"
)

type Processor struct {
}

func (p Processor) SendError(errorInfo error, httpResponseWriter http.ResponseWriter) {
	if errorInfo != nil {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusBadRequest)
		responseParams := errors.ErrorResponseV2{Description: errorInfo.Error()}
		jsonResponse, err := json.Marshal(responseParams)
		if err != nil {
			log.Printf("Error occurred: %s", errorInfo.Error())
		}
		httpResponseWriter.Write(jsonResponse)
	}
}

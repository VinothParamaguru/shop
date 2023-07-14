package response

import (
	"encoding/json"
	"log"
	"net/http"
	"shop/errors"
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

func (p Processor) SendAck(httpResponseWriter http.ResponseWriter) {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusOK)
	responseParams := AckResponse{Value: "ack"}
	jsonResponse, err := json.Marshal(responseParams)
	if err != nil {
		log.Printf("Error occurred: %s", err.Error())
	}
	httpResponseWriter.Write(jsonResponse)
}

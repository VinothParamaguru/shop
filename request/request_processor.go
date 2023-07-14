package request

import (
	"errors"
	"io"
	"net/http"
	apperrors "shop/errors"
)

type Processor struct {
}

// ReadRequest
// reads the http.Request, checks the content type for the incoming request to be application/json
// add additional checks wherever necessary
// can include checks like size, etc
func (p Processor) ReadRequest(httpRequest *http.Request) ([]byte, error) {
	bytes, err := io.ReadAll(httpRequest.Body)
	if err != nil {
		return nil, errors.New(apperrors.HttpErrorDescriptions[apperrors.HttpReadFailed])
	}
	if httpRequest.Header.Get("Content-Type") != "application/json" {
		return nil, errors.New(apperrors.HttpErrorDescriptions[apperrors.HttpUnknownContentType])
	}
	return bytes, nil
}

package errors

// ErrorResponse Error description structure for the client
type ErrorResponse struct {
	Code        int
	Description string
}

// ErrorResponseV2 Error description structure for the client
type ErrorResponseV2 struct {
	Description string `json:"Error"`
}

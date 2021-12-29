package error

// Error codes - application specific
const (
	Success = 0
	// Starts from 100 just not to mingle with the http standard error codes
	InvalidInput = 1000
	FieldMissing = 1001
)

// Error code and descriptions map
var ErrorDescriptions = map[int]string{
	InvalidInput: "Input is invalid",
	FieldMissing: "Input field missing",
}

// Error description structure for the client
type ErrorResponse struct {
	Code        int
	Description string
}

package error

// Error codes - application specific
const (
	Success = 0

	// Starts from 1000 just not to mingle with the http standard error codes
	// validation errors
	InvalidInput = 1000
	FieldMissing = 1001

	// database errors
	DbOpenFailed = 2000
)

// Error code and descriptions map
var ErrorDescriptions = map[int]string{
	InvalidInput: "Input is invalid",
	FieldMissing: "Input field missing",
	DbOpenFailed: "Problem with database connection",
}

// Error description structure for the client
type ErrorResponse struct {
	Code        int
	Description string
}

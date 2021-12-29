package error

// Error codes - application specific
const (
	Success      = 0
	InvalidInput = 100
)

// Error code and descriptions map
var ErrorDescriptions = map[int]string{
	InvalidInput: "Input is invalid. Please check",
}

// Error description structure for the client
type ErrorResponse struct {
	Code        int
	Description string
}

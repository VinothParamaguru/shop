package errors

// Error codes - security specific
const (
	SecInvalidInput = 1000
	SecFieldMissing = 1001
)

// SecurityErrorDescriptions and mappings
var SecurityErrorDescriptions = map[int]string{
	SecInvalidInput: "Input is invalid",
	SecFieldMissing: "Input field missing",
}

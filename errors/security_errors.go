package errors

// Error codes - security specific
const (
	InvalidInput = 1000
	FieldMissing = 1001
)

// SecurityErrorDescriptions and mappings
var SecurityErrorDescriptions = map[int]string{
	InvalidInput: "Input is invalid",
	FieldMissing: "Input field missing",
}

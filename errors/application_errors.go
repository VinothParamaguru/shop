package errors

// Error codes - database specific
const (
	AppInvalidUserNameOrPassword = 3000
	AppUserNameExistsAlready     = 3001
)

// ApplicationErrorDescriptions and mappings
var ApplicationErrorDescriptions = map[int]string{
	AppInvalidUserNameOrPassword: "Error in Login. Invalid Username or password",
	AppUserNameExistsAlready:     "Username not available",
}

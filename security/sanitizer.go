package security

import "net/mail"

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var validator_mappings = map[string]func(string) bool{
	"email": ValidateEmail,
}

// treat all the values as strings for now, for simplicity
func Validate(name string, value string) bool {
	return validator_mappings[name](value)
}

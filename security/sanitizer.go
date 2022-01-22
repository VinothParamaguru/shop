package security

import (
	"net/mail"
	"regexp"
	app_error "workspace/shop/error"
)

type ValidatorParams struct {
	ValidatorCustom func(string) (bool, int)
	ValidatorRegex  func(string, string) (bool, int)
	Expression      string
}

func ValidateEmail(email string) (bool, int) {
	_, error := mail.ParseAddress(email)
	if error == nil {
		return true, app_error.Success
	}
	return false, app_error.InvalidInput
}

func ValidatePassword(password string) (bool, int) {

	// should be atleast 8 characters long
	if len(password) < 8 {
		return false, app_error.InvalidInput
	}

	status_upper := false
	status_lower := false
	status_special := false

	for _, char := range password {
		// should contain one uppercase
		if char >= 'A' && char <= 'Z' {
			status_upper = true
		}
		// should contain one lowercase
		if char >= 'a' && char <= 'z' {
			status_lower = true
		}
		// should contain one special character
		if char == '@' || char == '$' || char == '£' {
			status_special = true
		}
	}

	if status_upper && status_lower && status_special {
		return true, app_error.Success
	}

	return false, app_error.InvalidInput
}

func ValidateField(value string, expression string) (bool, int) {
	match, _ := regexp.MatchString(expression, value)
	if match {
		return true, app_error.Success
	}
	return false, app_error.InvalidInput
}

var validator_mappings = map[string]ValidatorParams{
	"email":    {ValidateEmail, nil, ""},
	"name":     {nil, ValidateField, "^[A-Z][-'a-zA-Z]+$"},
	"password": {ValidatePassword, nil, ""},
}

// treat all the values as strings for now, for simplicity
func ValidateInput(name string, value string) (bool, int) {
	expression := validator_mappings[name].Expression
	if expression == "" {
		return validator_mappings[name].ValidatorCustom(value)
	} else {
		return validator_mappings[name].ValidatorRegex(value, expression)
	}
}

func ValidateRequiredFields(fields []string) (bool, int) {

	for _, field := range fields {
		if field == "" {
			return false, app_error.FieldMissing
		}
	}
	return true, app_error.Success
}

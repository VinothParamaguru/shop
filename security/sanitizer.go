package security

import (
	"net/mail"
	"regexp"
	"workspace/shop/errors"
)

type ValidatorParams struct {
	ValidatorCustom func(string) (bool, int)
	ValidatorRegex  func(string, string) (bool, int)
	Expression      string
}

func ValidateEmail(email string) (bool, int) {
	_, err := mail.ParseAddress(email)
	if err == nil {
		return true, errors.Success
	}
	return false, errors.InvalidInput
}

func ValidatePassword(password string) (bool, int) {

	// should be at least 8 characters long
	if len(password) < 8 {
		return false, errors.InvalidInput
	}

	statusUpper := false
	statusLower := false
	statusSpecial := false

	for _, char := range password {
		// should contain one uppercase
		if char >= 'A' && char <= 'Z' {
			statusUpper = true
		}
		// should contain one lowercase
		if char >= 'a' && char <= 'z' {
			statusLower = true
		}
		// should contain one special character
		if char == '@' || char == '$' || char == '£' {
			statusSpecial = true
		}
	}

	if statusUpper && statusLower && statusSpecial {
		return true, errors.Success
	}

	return false, errors.InvalidInput
}

func ValidateField(value string, expression string) (bool, int) {
	match, _ := regexp.MatchString(expression, value)
	if match {
		return true, errors.Success
	}
	return false, errors.InvalidInput
}

var validatorMappings = map[string]ValidatorParams{
	"email":    {ValidateEmail, nil, ""},
	"name":     {nil, ValidateField, "^[A-Z][-'a-zA-Z]+$"},
	"password": {ValidatePassword, nil, ""},
}

// ValidateInput treat all the values as strings for now, for simplicity
func ValidateInput(name string, value string) (bool, int) {
	expression := validatorMappings[name].Expression
	if expression == "" {
		return validatorMappings[name].ValidatorCustom(value)
	} else {
		return validatorMappings[name].ValidatorRegex(value, expression)
	}
}

func ValidateRequiredFields(fields []string) (bool, int) {

	for _, field := range fields {
		if field == "" {
			return false, errors.FieldMissing
		}
	}
	return true, errors.Success
}

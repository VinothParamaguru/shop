package security

import (
	"errors"
	"net/mail"
	"regexp"
	apperrors "shop/errors"
)

type ValidatorParams struct {
	ValidatorCustom func(string) error
	ValidatorRegex  func(string, string) error
	Expression      string
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err == nil {
		return nil
	}
	return errors.New(apperrors.SecurityErrorDescriptions[apperrors.SecInvalidInput])
}

func ValidatePassword(password string) error {

	// should be at least 8 characters long
	if len(password) < 8 {
		return errors.New(apperrors.SecurityErrorDescriptions[apperrors.SecInvalidInput])
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
		return nil
	}

	return errors.New(apperrors.SecurityErrorDescriptions[apperrors.SecInvalidInput])
}

func ValidateField(value string, expression string) error {
	match, _ := regexp.MatchString(expression, value)
	if match {
		return nil
	}
	return errors.New(apperrors.SecurityErrorDescriptions[apperrors.SecInvalidInput])
}

var validatorMappings = map[string]ValidatorParams{
	"email":    {ValidateEmail, nil, ""},
	"name":     {nil, ValidateField, "^[A-Z][-'a-zA-Z]+$"},
	"password": {ValidatePassword, nil, ""},
}

// ValidateInput treat all the values as strings for now, for simplicity
func ValidateInput(name string, value string) error {
	expression := validatorMappings[name].Expression
	if expression == "" {
		return validatorMappings[name].ValidatorCustom(value)
	} else {
		return validatorMappings[name].ValidatorRegex(value, expression)
	}
}

func ValidateRequiredFields(fields []string) error {
	for _, field := range fields {
		if field == "" {
			return errors.New(apperrors.SecurityErrorDescriptions[apperrors.SecFieldMissing])
		}
	}
	return nil
}

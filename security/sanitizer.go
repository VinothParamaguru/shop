package security

import (
	"net/mail"
	"regexp"
	app_error "workspace/shop/error"
)

type ValidatorParams struct {
	ValidatorAuto   func(string) (bool, int)
	ValidatorManual func(string, string) (bool, int)
	Expression      string
}

func ValidateEmail(email string) (bool, int) {
	_, error := mail.ParseAddress(email)
	if error == nil {
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
	"email": {ValidateEmail, nil, ""},
	"name":  {nil, ValidateField, "^([a-zA-Z]{2,}\\s[a-zA-Z]{1,}'?-?[a-zA-Z]{2,}\\s?([a-zA-Z]{1,})?)"},
}

// treat all the values as strings for now, for simplicity
func Validate(name string, value string) (bool, int) {
	expression := validator_mappings[name].Expression
	if expression == "" {
		return validator_mappings[name].ValidatorAuto(value)
	} else {
		return validator_mappings[name].ValidatorManual(value, expression)
	}
}

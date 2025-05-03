package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	//temporarily only supports id phone number
	matched, _ := regexp.MatchString(`^\+62[0-9]{8,13}$`, phone)
	return matched
}
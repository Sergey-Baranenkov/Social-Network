package validators

import "github.com/go-playground/validator/v10"

func ValidateSex(fl validator.FieldLevel) bool {
	str:= fl.Field().String()
	return str == "М" || str == "Ж"
}

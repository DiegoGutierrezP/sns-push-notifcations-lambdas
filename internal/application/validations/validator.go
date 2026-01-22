package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = func() *validator.Validate {
	v := validator.New()
	RegisterCustomValidations(v)
	return v
}()

func MapValidationErrors(err error) map[string]string {
	errors := map[string]string{}
	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = msgForTag(e)
	}
	return errors
}

func StringValidationErrors(err error) string {
	var msgs []string

	for _, e := range err.(validator.ValidationErrors) {
		msgs = append(msgs, formatField(e.Field())+": "+msgForTag(e))
	}

	return strings.Join(msgs, ", ")
}

func msgForTag(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "es requerido"
	case "oneof":
		return "valor inválido"
	case "min":
		return "longitud mínima inválida"
	case "subscription_filters":
		return "estructura de filtros inválida"
	default:
		return "campo inválido"
	}
}

func formatField(f string) string {
	return strings.ToLower(f[:1]) + f[1:]
}

package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(v *validator.Validate) {
	_ = v.RegisterValidation("subscription_filters", validateSubsciptionFilters)
}

func validateSubsciptionFilters(fl validator.FieldLevel) bool {
	filters, ok := fl.Field().Interface().(map[string][]string)
	if !ok || len(filters) == 0 {
		return false
	}

	for k, values := range filters {
		if strings.TrimSpace(k) == "" {
			return false
		}
		if len(values) == 0 {
			return false
		}
		for _, v := range values {
			if strings.TrimSpace(v) == "" {
				return false
			}
		}
	}
	return true
}

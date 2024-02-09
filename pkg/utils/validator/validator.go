package route_validator

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string
	Tag   string
	Value string
}

var validate = validator.New()

type routeValidator[T any] struct {
}

func Validate[T any](t *T) []*ValidationError {
	var errors []*ValidationError
	err := validate.Struct(t)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			element := ValidationError{
				Field: e.Field(),
				Tag:   e.Tag(),
				Value: e.Param(),
			}
			errors = append(errors, &element)
		}
	}
	return errors
}

func IsDate(fl validator.FieldLevel) bool {
	t, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	} else if t.Before(time.Now().UTC()) {
		fmt.Println(time.Now().UTC())
		return false
	}

	return true
}

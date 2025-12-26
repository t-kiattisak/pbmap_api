package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Wrapper struct {
	validator *validator.Validate
}

func New() *Wrapper {
	return &Wrapper{
		validator: validator.New(),
	}
}

func (v *Wrapper) Validate(i interface{}) map[string]string {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[strings.ToLower(err.Field())] = fmt.Sprintf("failed on the '%s' tag", err.Tag())
	}
	return errors
}

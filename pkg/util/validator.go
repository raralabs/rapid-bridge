package util

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()

	validate.RegisterValidation("date", validateDate)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{
		Validator: validate,
	}
}

func validateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	_, err := time.Parse(time.DateOnly, date)
	return err == nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

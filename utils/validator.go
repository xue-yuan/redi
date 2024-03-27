package utils

import (
	"fmt"
	"os"
	"redi/config"
	"redi/constants"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()

		index := strings.Index(field, "-")
		if index == -1 {
			return false
		}

		u := field[index+1:]
		if _, err := uuid.Parse(u); err != nil {
			return false
		}

		return true
	})

	_ = validate.RegisterValidation("exist", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()

		if field == "" {
			return true
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s", config.Config.ImageFolder, field)); err != nil {
			return false
		}

		return true
	})

	_ = validate.RegisterValidation("order", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()

		if field != string(constants.ASC) && field != string(constants.DESC) {
			return false
		}

		return true
	})

	return validate
}

func ValidatorErrors(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}

// func

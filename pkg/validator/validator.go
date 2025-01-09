package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var (
	name     = regexp.MustCompile(`^[\p{L}\p{N}_.-]{3,32}$`)
	password = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+={}\[\]:;"'<>,.?/-]{8,16}$`)
	mobile   = regexp.MustCompile(`^1([3-9][0-9]{1}|4[5|7][0-9]{1}|5[0-2|7-9][0-9]{1}|6[2-5][0-9]{1}|8[0-9]{1})\d{8}$`)
	validate *validator.Validate
)

func init() {
	validate = validator.New()
	validate.RegisterValidation("name", func(fl validator.FieldLevel) bool {
		return IsName(fl.Field().String())
	})
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return IsPassword(fl.Field().String())
	})
	validate.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		return IsMobile(fl.Field().String())
	})
}

func IsName(s string) bool {
	return name.MatchString(s)
}

func IsPassword(s string) bool {
	return password.MatchString(s)
}

func IsMobile(s string) bool {
	return mobile.MatchString(s)
}

func Struct(s interface{}) error {
	return validate.Struct(s)
}

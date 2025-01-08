package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var (
	name     = regexp.MustCompile(`^[\p{L}\p{N}_.-]{3,32}$`)
	password = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+={}\[\]:;"'<>,.?/-]{8,16}$`)
	mobile   = regexp.MustCompile(`^1([3-9][0-9]{1}|4[5|7][0-9]{1}|5[0-2|7-9][0-9]{1}|6[2-5][0-9]{1}|8[0-9]{1})\d{8}$`)
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{}
	vv := validator.New()
	vv.RegisterValidation("name", func(fl validator.FieldLevel) bool {
		return v.IsName(fl.Field().String())
	})
	vv.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return v.IsPassword(fl.Field().String())
	})
	vv.RegisterValidation("mobile", func(fl validator.FieldLevel) bool {
		return v.IsMobile(fl.Field().String())
	})
	v.validate = vv
	return v
}

func (v *Validator) IsName(s string) bool {
	return name.MatchString(s)
}

func (v *Validator) IsPassword(s string) bool {
	return password.MatchString(s)
}

func (v *Validator) IsMobile(s string) bool {
	return mobile.MatchString(s)
}

func (v *Validator) Struct(s interface{}) error {
	return v.validate.Struct(s)
}

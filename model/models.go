package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type OnboardPostBody struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Email     string `json:"email" form:"email"`
}

func (r OnboardPostBody) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
	)
}

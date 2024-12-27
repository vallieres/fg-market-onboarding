package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type OnboardPostBody struct {
	FirstName           string `json:"first_name" form:"first_name"`
	LastName            string `json:"last_name" form:"last_name"`
	Email               string `json:"email" form:"email"`
	ZipCode             string `json:"zip_code" form:"zip_code"`
	DogName             string `json:"dog_name" form:"dog_name"`
	DogBreed            string `json:"dog_breed" form:"dog_breed"`
	DogAge              int    `json:"dog_age" form:"dog_age"`
	DogWeight           int    `json:"dog_weight" form:"dog_weight"`
	DogHealthConditions string `json:"dog_health_conditions" form:"dog_health_conditions"`
	MailingList         bool   `json:"mailing_list" form:"mailing_list"`
}

func (r OnboardPostBody) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.LastName, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.DogName, validation.Required),
		validation.Field(&r.DogBreed, validation.Required),
		validation.Field(&r.DogAge, validation.Required),
	)
}

type ZipCodeEntry struct {
	ZipCode   string     `json:"zipcode" db:"zipcode"`
	City      string     `json:"city" db:"city"`
	State     string     `json:"state" db:"state"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

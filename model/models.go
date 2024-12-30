package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type OnboardPostBody struct {
	FirstName           string  `json:"first_name" form:"first_name"`
	LastName            string  `json:"last_name" form:"last_name"`
	Email               string  `json:"email" form:"email"`
	ZipCode             string  `json:"zip_code" form:"zip_code"`
	Country             string  `json:"country" form:"country"`
	PetName             string  `json:"pet_name" form:"pet_name"`
	PetSpecies          string  `json:"pet_species" form:"pet_species"`
	PetBreed            string  `json:"pet_breed" form:"pet_breed"`
	PetAge              float64 `json:"pet_age" form:"pet_age"`
	PetWeight           int     `json:"pet_weight" form:"pet_weight"`
	PetHealthConditions string  `json:"pet_health_conditions" form:"pet_health_conditions"`
	MailingList         bool    `json:"mailing_list" form:"mailing_list"`
}

func (r OnboardPostBody) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FirstName, validation.Required),
		validation.Field(&r.LastName, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.PetName, validation.Required),
		validation.Field(&r.PetSpecies, validation.Required),
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

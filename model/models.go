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

type Plan struct {
	ID                             int64      `json:"id" db:"id"`
	Name                           string     `json:"name" db:"name"`
	Status                         string     `json:"status" db:"status"`
	PetName                        string     `json:"pet_name" db:"pet_name"`
	PetSpecies                     string     `json:"pet_species" db:"pet_species"`
	PetBreed                       string     `json:"pet_breed" db:"pet_breed"`
	PetWeightLbs                   string     `json:"pet_weight_lbs" db:"pet_weight_lbs"`
	PetActivityLevel               string     `json:"pet_activity_level" db:"pet_activity_level"`
	DailyTotalCalories             float64    `json:"daily_total_calories" db:"daily_total_calories"`
	DailyTotalProtein              float64    `json:"daily_total_protein" db:"daily_total_protein"`
	DailyTotalFat                  float64    `json:"daily_total_fat" db:"daily_total_fat"`
	DailyTotalCarbohydrates        float64    `json:"daily_total_carbohydrates" db:"daily_total_carbohydrates"`
	ProteinPercentOfCalories       float64    `json:"protein_percent_of_calories" db:"protein_percent_of_calories"`
	FatPercentOfCalories           float64    `json:"fat_percent_of_calories" db:"fat_percent_of_calories"`
	CarbohydratesPercentOfCalories float64    `json:"carbohydrates_percent_of_calories" db:"carbohydrates_percent_of_calories"` //nolint:lll
	CreatedAt                      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                      *time.Time `json:"deleted_at" db:"deleted_at"` // Pointer to handle NULL values
}

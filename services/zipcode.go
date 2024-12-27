package services

import (
	"github.com/vallieres/fg-market-onboarding/repository"
)

type ZipCodeService struct {
	*repository.ZipCodeRepository
}

func NewZipCodeService(repository *repository.ZipCodeRepository) *ZipCodeService {
	return &ZipCodeService{repository}
}

func (z *ZipCodeService) GetCityByZipCode(zipCode string) ([]string, error) {
	var entries []string

	zipCodes, errGetZipCodes := z.GetCitiesByZipCode(zipCode)
	if errGetZipCodes != nil {
		return entries, errGetZipCodes
	}

	for _, zip := range zipCodes {
		entries = append(entries, zip.City)
	}

	return entries, nil
}

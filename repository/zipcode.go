// Package repository takes care of working with different data sources.
// For the Zip Codes, the source for the data is https://postalpro.usps.com/ZIP_Locale_Detail
package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/vallieres/fg-market-onboarding/model"
)

type ZipCodeRepository struct {
	db *sqlx.DB
}

func NewZipCodeRepository(database *sqlx.DB) *ZipCodeRepository {
	return &ZipCodeRepository{db: database}
}

func (z *ZipCodeRepository) GetCitiesByZipCode(zipCode string) ([]model.ZipCodeEntry, error) {
	var entries []model.ZipCodeEntry
	query := `
SELECT zipcode, city, state
  FROM zipcodes
 WHERE zipcode = ?
`
	rows, errQueryRow := z.db.Queryx(query, zipCode) //nolint:sqlclosecheck // closing rows just below
	if errQueryRow != nil {
		return entries, fmt.Errorf("unable to select to get zip code entries: %w", errQueryRow)
	}
	_ = rows.Close()

	for rows.Next() {
		var entry model.ZipCodeEntry
		if errScan := rows.StructScan(&entry); errScan != nil {
			return entries, fmt.Errorf("unable to structscan into zip code entries: %w", errScan)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// ZipCode NOT found
	if len(entries) == 0 {
		return entries, errors.New("no zip code entries found")
	}

	// ZipCodes found
	return entries, nil
}

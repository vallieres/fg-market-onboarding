package repository

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // or the appropriate driver for your database
	"github.com/jmoiron/sqlx"

	"github.com/vallieres/fg-market-onboarding/model"
)

type PlanRepository struct {
	db *sqlx.DB
}

func NewPlanRepository(database *sqlx.DB) *PlanRepository {
	return &PlanRepository{db: database}
}

//nolint:lll
func (u *PlanRepository) CreateMealPlan(userEmail string, planDetails model.Plan) error {
	query := `
INSERT INTO plans(
    	name, pet_name, pet_species, pet_breed, pet_weight_lbs, pet_activity_level
) VALUES (?, ?, ?, ?, ?, ?)
 `
	_, errExec := u.db.Exec(query,
		userEmail, planDetails.PetName, planDetails.PetSpecies, planDetails.PetBreed, planDetails.PetWeightLbs, planDetails.PetActivityLevel,
	)
	if errExec != nil {
		return fmt.Errorf("unable to upsert customer subscription: %w", errExec)
	}

	return nil
}

func (u *PlanRepository) GetPlan(planID int64) (model.Plan, error) {
	var plan model.Plan
	query := `
SELECT id, name, status, pet_name, pet_species, pet_breed, pet_weight_lbs, pet_activity_level,
       created_at, updated_at, deleted_at
  FROM plans
 WHERE id = ?
`
	rows, errQueryRow := u.db.Queryx(query, planID) //nolint:sqlclosecheck // closing rows just below
	if errQueryRow != nil {
		return plan, fmt.Errorf("unable to select plan: %w", errQueryRow)
	}
	defer rows.Close()

	for rows.Next() {
		errScan := rows.StructScan(&plan)
		if errScan != nil {
			return plan, fmt.Errorf("unable to structscan into plan: %w", errScan)
		}
	}
	if err := rows.Err(); err != nil {
		return plan, err
	}

	// Plan NOT found
	if plan.ID == 0 {
		return plan, nil
	}

	// Plan found
	return plan, nil
}

package services

import (
	"context"
	"errors"

	"github.com/vallieres/fg-market-onboarding/model"
	"github.com/vallieres/fg-market-onboarding/repository"
)

type PlanService struct {
	shopifyAppToken        string
	shopifyStorefrontToken string
	planRepository         *repository.PlanRepository
}

func NewPlanService(
	planRepository *repository.PlanRepository,
	adminAPIToken string,
	storefrontToken string) *PlanService {
	return &PlanService{
		shopifyAppToken:        adminAPIToken,
		shopifyStorefrontToken: storefrontToken,
		planRepository:         planRepository,
	}
}

func (c *PlanService) CreateBasicPlan(_ context.Context, _ model.OnboardPostBody) (int64, error) {
	return 0, nil
}

func (c *PlanService) IsPlanReady(planID int64) (bool, error) {
	plan, errGetPlan := c.planRepository.GetPlan(planID)
	if errGetPlan != nil {
		return false, errGetPlan
	}

	if plan.Status == "READY" {
		return true, nil
	} else if plan.ID == 0 {
		return false, errors.New("plan is not found")
	}
	return false, nil
}

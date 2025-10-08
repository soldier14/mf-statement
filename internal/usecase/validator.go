package usecase

import (
	"mf-statement/internal/domain"
)

type Validator interface {
	ValidatePeriod(year, month int) error
	ValidateDateRange(startDate, endDate interface{}) error
}

type PeriodValidator struct{}

func NewPeriodValidator() Validator {
	return &PeriodValidator{}
}

func (v *PeriodValidator) ValidatePeriod(year, month int) error {
	if month < 1 || month > 12 {
		return domain.NewValidationError("invalid month", map[string]interface{}{
			"month": month,
		})
	}
	return nil
}

func (v *PeriodValidator) ValidateDateRange(startDate, endDate interface{}) error {
	// For now, we'll keep this simple since we removed year validation
	// In the future, we could add more sophisticated date range validation
	return nil
}

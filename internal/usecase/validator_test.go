package usecase_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/domain"
	"mf-statement/internal/usecase"
)

var _ = Describe("PeriodValidator", func() {
	var validator usecase.Validator

	BeforeEach(func() {
		validator = usecase.NewPeriodValidator()
	})

	Context("ValidatePeriod", func() {
		Context("with valid inputs", func() {
			It("should accept valid month", func() {
				err := validator.ValidatePeriod(2025, 1)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should accept any year", func() {
				err := validator.ValidatePeriod(1950, 6)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should accept future year", func() {
				err := validator.ValidatePeriod(2050, 12)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid inputs", func() {
			It("should reject month 0", func() {
				err := validator.ValidatePeriod(2025, 0)
				Expect(err).To(HaveOccurred())
				Expect(domain.IsValidationError(err)).To(BeTrue())
			})

			It("should reject month 13", func() {
				err := validator.ValidatePeriod(2025, 13)
				Expect(err).To(HaveOccurred())
				Expect(domain.IsValidationError(err)).To(BeTrue())
			})

			It("should reject negative month", func() {
				err := validator.ValidatePeriod(2025, -1)
				Expect(err).To(HaveOccurred())
				Expect(domain.IsValidationError(err)).To(BeTrue())
			})
		})
	})

	Context("ValidateDateRange", func() {
		It("should always pass for now", func() {
			err := validator.ValidateDateRange(nil, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

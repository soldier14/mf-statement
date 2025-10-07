package domain_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/domain"
)

var _ = Describe("Domain", func() {
	Context("Transaction", func() {
		It("should create valid transaction", func() {
			date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			transaction, err := domain.NewTransaction(date, 1000, "Salary")

			Expect(err).NotTo(HaveOccurred())
			Expect(transaction.Date).To(Equal(date))
			Expect(transaction.Amount).To(Equal(int64(1000)))
			Expect(transaction.Content).To(Equal("Salary"))
		})

		It("should return error for empty content", func() {
			date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			_, err := domain.NewTransaction(date, 1000, "")

			Expect(err).To(HaveOccurred())
		})

		It("should return error for whitespace-only content", func() {
			date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			_, err := domain.NewTransaction(date, 1000, "   ")

			Expect(err).To(HaveOccurred())
		})

		It("should identify income correctly", func() {
			transaction := domain.Transaction{Amount: 1000}
			Expect(transaction.IsIncome()).To(BeTrue())
		})

		It("should identify expense correctly", func() {
			transaction := domain.Transaction{Amount: -1000}
			Expect(transaction.IsExpense()).To(BeTrue())
		})

		It("should calculate absolute amount", func() {
			transaction := domain.Transaction{Amount: -1000}
			Expect(transaction.AbsAmount()).To(Equal(int64(1000)))
		})

		It("should calculate absolute amount for positive values", func() {
			transaction := domain.Transaction{Amount: 1000}
			Expect(transaction.AbsAmount()).To(Equal(int64(1000)))
		})

		It("should calculate absolute amount for zero", func() {
			transaction := domain.Transaction{Amount: 0}
			Expect(transaction.AbsAmount()).To(Equal(int64(0)))
		})
	})

	Context("Statement", func() {
		It("should create statement with transactions", func() {
			date := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			transactions := []domain.Transaction{
				{Date: date, Amount: 1000, Content: "Salary"},
				{Date: date, Amount: -200, Content: "Groceries"},
			}

			statement := domain.NewStatement("2025/01", transactions, 1000, -200)

			Expect(statement.Period).To(Equal("2025/01"))
			Expect(statement.TotalIncome).To(Equal(int64(1000)))
			Expect(statement.TotalExpenditure).To(Equal(int64(-200)))
			Expect(statement.NetAmount).To(Equal(int64(800)))
			Expect(statement.TransactionCount).To(Equal(2))
			Expect(statement.Transactions).To(HaveLen(2))
		})
	})

	Context("DomainError", func() {
		It("should create validation error", func() {
			err := domain.NewValidationError("test error", map[string]interface{}{"field": "value"})
			Expect(err.Type).To(Equal(domain.ErrorTypeValidation))
			Expect(err.Message).To(Equal("test error"))
		})

		It("should create parse error", func() {
			cause := errors.New("parse failed")
			err := domain.NewParseError("test error", cause)
			Expect(err.Type).To(Equal(domain.ErrorTypeParse))
			Expect(err.Cause).To(Equal(cause))
		})

		It("should create IO error", func() {
			cause := errors.New("file not found")
			err := domain.NewIOError("test error", cause)
			Expect(err.Type).To(Equal(domain.ErrorTypeIO))
			Expect(err.Cause).To(Equal(cause))
		})

		It("should check error types", func() {
			validationErr := domain.NewValidationError("test", nil)
			parseErr := domain.NewParseError("test", nil)
			ioErr := domain.NewIOError("test", nil)

			Expect(domain.IsValidationError(validationErr)).To(BeTrue())
			Expect(domain.IsParseError(parseErr)).To(BeTrue())
			Expect(domain.IsIOError(ioErr)).To(BeTrue())
		})

		It("should return false for non-domain errors", func() {
			otherErr := errors.New("other error")
			Expect(domain.IsValidationError(otherErr)).To(BeFalse())
			Expect(domain.IsParseError(otherErr)).To(BeFalse())
			Expect(domain.IsIOError(otherErr)).To(BeFalse())
		})

		It("should return false for nil errors", func() {
			Expect(domain.IsValidationError(nil)).To(BeFalse())
			Expect(domain.IsParseError(nil)).To(BeFalse())
			Expect(domain.IsIOError(nil)).To(BeFalse())
		})

		It("should create not found error", func() {
			err := domain.NewNotFoundError("test error")
			Expect(err.Type).To(Equal(domain.ErrorTypeNotFound))
			Expect(err.Message).To(ContainSubstring("test error"))
		})

		It("should create internal error", func() {
			cause := errors.New("internal error")
			err := domain.NewInternalError("test error", cause)
			Expect(err.Type).To(Equal(domain.ErrorTypeInternal))
			Expect(err.Cause).To(Equal(cause))
		})

		It("should handle error unwrapping", func() {
			cause := errors.New("original error")
			err := domain.NewParseError("wrapped error", cause)
			Expect(err.Unwrap()).To(Equal(cause))
		})

		It("should format error string", func() {
			err := domain.NewValidationError("test error", map[string]interface{}{"field": "value"})
			Expect(err.Error()).To(ContainSubstring("test error"))
		})
	})

	Context("ErrorSummary", func() {
		It("should add and check errors", func() {
			errorSummary := &domain.ErrorSummary{}
			Expect(errorSummary.HasErrors()).To(BeFalse())

			errorSummary.AddError(domain.NewValidationError("test", nil))
			Expect(errorSummary.HasErrors()).To(BeTrue())
		})

		It("should format error string with multiple errors", func() {
			errorSummary := &domain.ErrorSummary{}
			errorSummary.AddError(domain.NewValidationError("error 1", nil))
			errorSummary.AddError(domain.NewParseError("error 2", nil))

			errorStr := errorSummary.Error()
			Expect(errorStr).To(ContainSubstring("multiple errors"))
			Expect(errorStr).To(ContainSubstring("error 1"))
			Expect(errorStr).To(ContainSubstring("error 2"))
		})

		It("should return no errors message when empty", func() {
			errorSummary := &domain.ErrorSummary{}
			Expect(errorSummary.Error()).To(Equal("no errors"))
		})

		It("should handle single error", func() {
			errorSummary := &domain.ErrorSummary{}
			errorSummary.AddError(domain.NewValidationError("single error", nil))

			errorStr := errorSummary.Error()
			Expect(errorStr).To(ContainSubstring("multiple errors"))
			Expect(errorStr).To(ContainSubstring("single error"))
		})
	})
})

package usecase_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/domain"
	"mf-statement/internal/usecase"
	"mf-statement/internal/util"
)

type mockSource struct {
	reader io.ReadCloser
	err    error
}

func (m mockSource) Open(_ context.Context, _ string) (io.ReadCloser, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.reader != nil {
		return m.reader, nil
	}
	return io.NopCloser(strings.NewReader("")), nil
}

type mockParser struct {
	transactions []domain.Transaction
	err          error
}

func (m mockParser) Parse(_ context.Context, _ io.Reader) ([]domain.Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.transactions, nil
}

var _ = Describe("TransactionService", func() {
	var (
		service usecase.TransactionService
		ctx     context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("when getting all transactions", func() {
		It("should return all transactions successfully", func() {
			// Given
			expectedTransactions := []domain.Transaction{
				{Date: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC), Amount: 2000, Content: "Salary"},
				{Date: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), Amount: -300, Content: "Grocery"},
			}
			service = usecase.NewTransactionService(mockSource{}, mockParser{transactions: expectedTransactions})

			// When
			transactions, err := service.GetAllTransactions(ctx, "test.csv")

			// Then
			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(len(expectedTransactions)))
		})
	})

	Context("when getting transactions by period", func() {
		It("should filter transactions by year and month", func() {
			// Given
			allTransactions := []domain.Transaction{
				{Date: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC), Amount: 2000, Content: "Salary"},
				{Date: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), Amount: -300, Content: "Grocery"},
				{Date: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC), Amount: 999, Content: "Next Month"},
			}
			service = usecase.NewTransactionService(mockSource{}, mockParser{transactions: allTransactions})

			// When
			transactions, err := service.GetTransactionsByPeriod(ctx, "test.csv", 2025, 1)

			// Then
			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2), "Should only return January transactions")
			Expect(transactions[0].Date.Day()).To(Equal(9), "Should be sorted by date (newest first)")
		})
	})

	Context("when calculating totals", func() {
		It("should calculate income and expenditure correctly", func() {
			// Given
			transactions := []domain.Transaction{
				{Date: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC), Amount: 2000, Content: "Salary"},
				{Date: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), Amount: -300, Content: "Grocery"},
				{Date: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC), Amount: 100, Content: "Gift"},
			}
			service = usecase.NewTransactionService(mockSource{}, mockParser{})

			// When
			totalIncome, totalExpenditure := service.CalculateTotals(transactions)

			// Then
			Expect(totalIncome).To(Equal(int64(2100)), "Should calculate total income correctly")
			Expect(totalExpenditure).To(Equal(int64(-300)), "Should calculate total expenditure correctly")
		})
	})

	Context("error handling", func() {
		It("should return error when source fails", func() {
			// Given
			service = usecase.NewTransactionService(mockSource{err: errors.New("source error")}, mockParser{})

			// When
			_, err := service.GetAllTransactions(ctx, "test.csv")

			// Then
			Expect(err).To(HaveOccurred())
		})

		It("should return error when parser fails", func() {
			// Given
			service = usecase.NewTransactionService(mockSource{}, mockParser{err: errors.New("parse error")})

			// When
			_, err := service.GetAllTransactions(ctx, "test.csv")

			// Then
			Expect(err).To(HaveOccurred())
		})
	})

	Context("validation", func() {
		BeforeEach(func() {
			service = usecase.NewTransactionService(mockSource{}, mockParser{})
		})

		It("should return error for invalid month", func() {
			// When
			_, err := service.GetTransactionsByPeriod(ctx, "test.csv", 2025, 13)

			// Then
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetTransactionsByDateRange", func() {
		It("should get transactions within date range", func() {
			allTransactions := []domain.Transaction{
				{Date: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC), Amount: 2000, Content: "Salary"},
				{Date: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC), Amount: -300, Content: "Grocery"},
				{Date: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC), Amount: 999, Content: "Next Month"},
			}
			service = usecase.NewTransactionService(mockSource{}, mockParser{transactions: allTransactions})

			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

			transactions, err := service.GetTransactionsByDateRange(ctx, "test.csv", startDate, endDate)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
		})

		It("should handle context cancellation", func() {
			service = usecase.NewTransactionService(mockSource{}, mockParser{})
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

			transactions, err := service.GetTransactionsByDateRange(cancelledCtx, "test.csv", startDate, endDate)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(BeEmpty())
		})
	})

	Context("Between", func() {
		It("should include dates within range", func() {
			start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

			// Test dates within range
			Expect(util.Between(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue())  // start date
			Expect(util.Between(time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue()) // middle
			Expect(util.Between(time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue()) // end date
		})

		It("should exclude dates outside range", func() {
			start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

			// Test dates outside range
			Expect(util.Between(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), start, end)).To(BeFalse()) // before start
			Expect(util.Between(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC), start, end)).To(BeFalse())   // after end
		})

	})
})

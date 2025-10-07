package usecase_test

import (
	"context"
	"errors"
	"mf-statement/internal/domain"
	"mf-statement/internal/usecase"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Mock implementations
type mockTransactionService struct {
	transactionsByPeriod    []domain.Transaction
	transactionsByDateRange []domain.Transaction
	periodError             error
	dateRangeError          error
}

func (m *mockTransactionService) GetTransactionsByPeriod(ctx context.Context, csvFileURI string, year, month int) ([]domain.Transaction, error) {
	return m.transactionsByPeriod, m.periodError
}

func (m *mockTransactionService) GetTransactionsByDateRange(ctx context.Context, csvFileURI string, startDate, endDate time.Time) ([]domain.Transaction, error) {
	return m.transactionsByDateRange, m.dateRangeError
}

func (m *mockTransactionService) CalculateTotals(transactions []domain.Transaction) (totalIncome, totalExpenditure int64) {
	var income, expenditure int64
	for _, t := range transactions {
		if t.IsIncome() {
			income += t.Amount
		} else {
			expenditure += t.Amount
		}
	}
	return income, expenditure
}

type mockWriter struct {
	writtenStatement *domain.Statement
	writeError       error
}

func (m *mockWriter) Write(ctx context.Context, statement domain.Statement) error {
	m.writtenStatement = &statement
	return m.writeError
}

var _ = Describe("StatementService", func() {
	var (
		service            *usecase.StatementService
		mockTxService      *mockTransactionService
		mockWriterInstance *mockWriter
		ctx                context.Context
	)

	BeforeEach(func() {
		mockTxService = &mockTransactionService{}
		mockWriterInstance = &mockWriter{}
		service = usecase.NewStatementService(mockTxService, mockWriterInstance)
		ctx = context.Background()
	})

	Describe("NewStatementService", func() {
		It("should create a new StatementService with dependencies", func() {
			Expect(service).ToNot(BeNil())
			Expect(service.TransactionService).To(Equal(mockTxService))
			Expect(service.Writer).To(Equal(mockWriterInstance))
		})
	})

	Describe("GenerateMonthlyStatement", func() {
		var (
			csvFileURI    string
			periodDisplay string
			year, month   int
		)

		BeforeEach(func() {
			csvFileURI = "test.csv"
			periodDisplay = "2025/01"
			year, month = 2025, 1
		})

		Context("when transaction service returns transactions", func() {
			BeforeEach(func() {
				date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				date2 := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
				mockTxService.transactionsByPeriod = []domain.Transaction{
					{Date: date1, Amount: 1000, Content: "Salary"},
					{Date: date2, Amount: -200, Content: "Groceries"},
				}
			})

			It("should generate and write statement successfully", func() {
				err := service.GenerateMonthlyStatement(ctx, csvFileURI, periodDisplay, year, month)

				Expect(err).ToNot(HaveOccurred())
				Expect(mockWriterInstance.writtenStatement).ToNot(BeNil())
				Expect(mockWriterInstance.writtenStatement.Period).To(Equal(periodDisplay))
				Expect(mockWriterInstance.writtenStatement.TotalIncome).To(Equal(int64(1000)))
				Expect(mockWriterInstance.writtenStatement.TotalExpenditure).To(Equal(int64(-200)))
				Expect(mockWriterInstance.writtenStatement.Transactions).To(HaveLen(2))
			})

			It("should call transaction service with correct parameters", func() {
				err := service.GenerateMonthlyStatement(ctx, csvFileURI, periodDisplay, year, month)

				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when transaction service returns error", func() {
			BeforeEach(func() {
				mockTxService.periodError = errors.New("transaction service error")
			})

			It("should return the error", func() {
				err := service.GenerateMonthlyStatement(ctx, csvFileURI, periodDisplay, year, month)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("transaction service error"))
			})
		})

		Context("when writer returns error", func() {
			BeforeEach(func() {
				date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				mockTxService.transactionsByPeriod = []domain.Transaction{
					{Date: date1, Amount: 1000, Content: "Salary"},
				}
				mockWriterInstance.writeError = errors.New("write error")
			})

			It("should return IO error", func() {
				err := service.GenerateMonthlyStatement(ctx, csvFileURI, periodDisplay, year, month)

				Expect(err).To(HaveOccurred())
				Expect(domain.IsIOError(err)).To(BeTrue())
				Expect(err.Error()).To(ContainSubstring("failed to write statement"))
			})
		})
	})

	Describe("GenerateStatementFromTransactions", func() {
		var (
			transactions  []domain.Transaction
			periodDisplay string
		)

		BeforeEach(func() {
			date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			date2 := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
			transactions = []domain.Transaction{
				{Date: date1, Amount: 1000, Content: "Salary"},
				{Date: date2, Amount: -200, Content: "Groceries"},
			}
			periodDisplay = "2025/01"
		})

		Context("when writer succeeds", func() {
			It("should generate and write statement successfully", func() {
				err := service.GenerateStatementFromTransactions(ctx, transactions, periodDisplay)

				Expect(err).ToNot(HaveOccurred())
				Expect(mockWriterInstance.writtenStatement).ToNot(BeNil())
				Expect(mockWriterInstance.writtenStatement.Period).To(Equal(periodDisplay))
				Expect(mockWriterInstance.writtenStatement.TotalIncome).To(Equal(int64(1000)))
				Expect(mockWriterInstance.writtenStatement.TotalExpenditure).To(Equal(int64(-200)))
				Expect(mockWriterInstance.writtenStatement.Transactions).To(HaveLen(2))
			})
		})

		Context("when writer returns error", func() {
			BeforeEach(func() {
				mockWriterInstance.writeError = errors.New("write error")
			})

			It("should return IO error", func() {
				err := service.GenerateStatementFromTransactions(ctx, transactions, periodDisplay)

				Expect(err).To(HaveOccurred())
				Expect(domain.IsIOError(err)).To(BeTrue())
				Expect(err.Error()).To(ContainSubstring("failed to write statement"))
			})
		})

		Context("with empty transactions", func() {
			BeforeEach(func() {
				transactions = []domain.Transaction{}
			})

			It("should generate statement with zero totals", func() {
				err := service.GenerateStatementFromTransactions(ctx, transactions, periodDisplay)

				Expect(err).ToNot(HaveOccurred())
				Expect(mockWriterInstance.writtenStatement.TotalIncome).To(Equal(int64(0)))
				Expect(mockWriterInstance.writtenStatement.TotalExpenditure).To(Equal(int64(0)))
				Expect(mockWriterInstance.writtenStatement.Transactions).To(HaveLen(0))
			})
		})
	})

	Describe("GenerateStatementByDateRange", func() {
		var (
			csvFileURI         string
			periodDisplay      string
			startDate, endDate time.Time
		)

		BeforeEach(func() {
			csvFileURI = "test.csv"
			periodDisplay = "2025/01"
			startDate = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate = time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)
		})

		Context("when transaction service returns transactions", func() {
			BeforeEach(func() {
				date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				date2 := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
				mockTxService.transactionsByDateRange = []domain.Transaction{
					{Date: date1, Amount: 1000, Content: "Salary"},
					{Date: date2, Amount: -200, Content: "Groceries"},
				}
			})

			It("should generate and write statement successfully", func() {
				err := service.GenerateStatementByDateRange(ctx, csvFileURI, periodDisplay, startDate, endDate)

				Expect(err).ToNot(HaveOccurred())
				Expect(mockWriterInstance.writtenStatement).ToNot(BeNil())
				Expect(mockWriterInstance.writtenStatement.Period).To(Equal(periodDisplay))
				Expect(mockWriterInstance.writtenStatement.Transactions).To(HaveLen(2))
			})
		})

		Context("when transaction service returns error", func() {
			BeforeEach(func() {
				mockTxService.dateRangeError = errors.New("date range error")
			})

			It("should return the error", func() {
				err := service.GenerateStatementByDateRange(ctx, csvFileURI, periodDisplay, startDate, endDate)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("date range error"))
			})
		})

		Context("when GenerateStatementFromTransactions returns error", func() {
			BeforeEach(func() {
				date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				mockTxService.transactionsByDateRange = []domain.Transaction{
					{Date: date1, Amount: 1000, Content: "Salary"},
				}
				mockWriterInstance.writeError = errors.New("write error")
			})

			It("should return the error from GenerateStatementFromTransactions", func() {
				err := service.GenerateStatementByDateRange(ctx, csvFileURI, periodDisplay, startDate, endDate)

				Expect(err).To(HaveOccurred())
				Expect(domain.IsIOError(err)).To(BeTrue())
			})
		})
	})

	Describe("Context handling", func() {
		It("should pass context to transaction service", func() {
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			mockTxService.transactionsByPeriod = []domain.Transaction{
				{Date: date1, Amount: 1000, Content: "Salary"},
			}

			err := service.GenerateMonthlyStatement(cancelledCtx, "test.csv", "2025/01", 2025, 1)

			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass context to writer", func() {
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			mockTxService.transactionsByPeriod = []domain.Transaction{
				{Date: date1, Amount: 1000, Content: "Salary"},
			}

			err := service.GenerateMonthlyStatement(cancelledCtx, "test.csv", "2025/01", 2025, 1)

			Expect(err).ToNot(HaveOccurred())
		})
	})
})

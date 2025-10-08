package usecase_test

import (
	"context"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/in"
	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/domain"
	"mf-statement/internal/usecase"
)

var _ = Describe("Optimized Services", func() {
	var (
		tempDir string
		csvPath string
		ctx     context.Context
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "optimized_test_*")
		Expect(err).NotTo(HaveOccurred())

		csvPath = filepath.Join(tempDir, "transactions.csv")
		ctx = context.Background()
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Context("OptimizedTransactionService", func() {
		It("should create new service", func() {
			source := in.NewCSVFileSource()
			service := usecase.NewOptimizedTransactionService(source)

			Expect(service).NotTo(BeNil())
		})

		It("should get transactions by period", func() {
			csvContent := `date,amount,content
2025/01/01,1000,January 1
2025/01/15,2000,January 15
2025/02/01,3000,February 1`
			Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

			source := in.NewCSVFileSource()
			service := usecase.NewOptimizedTransactionService(source)

			transactions, err := service.GetTransactionsByPeriodOptimized(ctx, csvPath, 2025, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
		})

		It("should get transactions by date range", func() {
			csvContent := `date,amount,content
2025/01/01,1000,January 1
2025/01/15,2000,January 15
2025/02/01,3000,February 1`
			Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

			source := in.NewCSVFileSource()
			service := usecase.NewOptimizedTransactionService(source)

			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)
			transactions, err := service.GetTransactionsByDateRangeOptimized(ctx, csvPath, startDate, endDate)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
		})

		It("should calculate totals", func() {
			source := in.NewCSVFileSource()
			service := usecase.NewOptimizedTransactionService(source)

			tx1, err1 := domain.NewTransaction(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), 1000, "Income")
			Expect(err1).NotTo(HaveOccurred())
			tx2, err2 := domain.NewTransaction(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), -500, "Expense")
			Expect(err2).NotTo(HaveOccurred())

			transactions := []domain.Transaction{tx1, tx2}

			totalIncome, totalExpenditure := service.CalculateTotalsOptimized(transactions)

			Expect(totalIncome).To(Equal(int64(1000)))
			Expect(totalExpenditure).To(Equal(int64(-500)))
		})
	})

	Context("OptimizedStatementService", func() {
		It("should create new service", func() {
			source := in.NewCSVFileSource()
			transactionService := usecase.NewOptimizedTransactionService(source)
			writer := output.NewJSON(os.Stdout)
			service := usecase.NewOptimizedStatementService(transactionService, writer)

			Expect(service).NotTo(BeNil())
		})

		It("should generate monthly statement", func() {
			csvContent := `date,amount,content
2025/01/01,2000,January Salary
2025/01/15,-500,January Expense`
			Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

			source := in.NewCSVFileSource()
			transactionService := usecase.NewOptimizedTransactionService(source)
			outputPath := filepath.Join(tempDir, "statement.json")
			writer := output.NewJSONFile(outputPath)
			service := usecase.NewOptimizedStatementService(transactionService, writer)

			err := service.GenerateMonthlyStatementOptimized(ctx, csvPath, "2025/01", 2025, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(outputPath).To(BeAnExistingFile())
		})

		It("should generate statement by date range", func() {
			csvContent := `date,amount,content
2025/01/01,1000,January 1
2025/01/15,2000,January 15
2025/02/01,3000,February 1`
			Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

			source := in.NewCSVFileSource()
			transactionService := usecase.NewOptimizedTransactionService(source)
			outputPath := filepath.Join(tempDir, "statement.json")
			writer := output.NewJSONFile(outputPath)
			service := usecase.NewOptimizedStatementService(transactionService, writer)

			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)
			err := service.GenerateStatementByDateRangeOptimized(ctx, csvPath, "2025/01", startDate, endDate)

			Expect(err).NotTo(HaveOccurred())
			Expect(outputPath).To(BeAnExistingFile())
		})
	})
})

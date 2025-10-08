package cli_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/in"
	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/usecase"
)

var _ = Describe("GenerateOptimizedCommand", func() {
	var (
		tempDir    string
		csvPath    string
		outputPath string
		ctx        context.Context
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "generate_optimized_test_*")
		Expect(err).NotTo(HaveOccurred())

		csvPath = filepath.Join(tempDir, "transactions.csv")
		outputPath = filepath.Join(tempDir, "statement.json")
		ctx = context.Background()
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Context("when generating optimized statements", func() {
		It("should generate statement for specific period", func() {
			// Create test CSV
			csvContent := `date,amount,content
2025/01/01,2000,January Salary
2025/01/15,-500,January Expense
2025/02/01,2500,February Salary`
			Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

			// Create optimized services
			source := in.NewCSVFileSource()
			optimizedTransactionService := usecase.NewOptimizedTransactionService(source)
			writer := output.NewJSONFile(outputPath)
			optimizedStatementService := usecase.NewOptimizedStatementService(optimizedTransactionService, writer)

			// Generate statement
			err := optimizedStatementService.GenerateMonthlyStatementOptimized(ctx, csvPath, "2025/01", 2025, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(outputPath).To(BeAnExistingFile())

			// Verify output content
			content, err := os.ReadFile(outputPath)
			Expect(err).NotTo(HaveOccurred())
			output := string(content)

			Expect(output).To(ContainSubstring(`"period": "2025/01"`))
			Expect(output).To(ContainSubstring(`"total_income": 2000`))
			Expect(output).To(ContainSubstring(`"total_expenditure": -500`))
			Expect(output).To(ContainSubstring(`"January Salary"`))
			Expect(output).To(ContainSubstring(`"January Expense"`))
		})

		It("should filter transactions by period during parsing", func() {
			// Create CSV with mixed periods
			csvContent := `date,amount,content
2025/01/01,1000,January 1
2025/01/15,2000,January 15
2025/02/01,3000,February 1
2025/02/15,4000,February 15
2025/03/01,5000,March 1`
			Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

			// Create optimized services
			source := in.NewCSVFileSource()
			optimizedTransactionService := usecase.NewOptimizedTransactionService(source)
			writer := output.NewJSONFile(outputPath)
			optimizedStatementService := usecase.NewOptimizedStatementService(optimizedTransactionService, writer)

			// Generate statement for February only
			err := optimizedStatementService.GenerateMonthlyStatementOptimized(ctx, csvPath, "2025/02", 2025, 2)

			Expect(err).NotTo(HaveOccurred())

			// Verify only February transactions are included
			content, err := os.ReadFile(outputPath)
			Expect(err).NotTo(HaveOccurred())
			output := string(content)

			Expect(output).To(ContainSubstring(`"period": "2025/02"`))
			Expect(output).To(ContainSubstring(`"total_income": 7000`)) // 3000 + 4000
			Expect(output).To(ContainSubstring(`"February 1"`))
			Expect(output).To(ContainSubstring(`"February 15"`))
			Expect(output).NotTo(ContainSubstring(`"January"`))
			Expect(output).NotTo(ContainSubstring(`"March"`))
		})

		It("should handle empty period gracefully", func() {
			// Create CSV with no transactions for the period
			csvContent := `date,amount,content
2025/02/01,2000,February Salary
2025/03/01,3000,March Salary`
			Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

			// Create optimized services
			source := in.NewCSVFileSource()
			optimizedTransactionService := usecase.NewOptimizedTransactionService(source)
			writer := output.NewJSONFile(outputPath)
			optimizedStatementService := usecase.NewOptimizedStatementService(optimizedTransactionService, writer)

			// Generate statement for January (no transactions)
			err := optimizedStatementService.GenerateMonthlyStatementOptimized(ctx, csvPath, "2025/01", 2025, 1)

			Expect(err).NotTo(HaveOccurred())

			// Verify empty statement
			content, err := os.ReadFile(outputPath)
			Expect(err).NotTo(HaveOccurred())
			output := string(content)

			Expect(output).To(ContainSubstring(`"period": "2025/01"`))
			Expect(output).To(ContainSubstring(`"total_income": 0`))
			Expect(output).To(ContainSubstring(`"total_expenditure": 0`))
			Expect(output).To(ContainSubstring(`"transactions": []`))
		})
	})

	Context("when generating statements by date range", func() {
		It("should filter transactions by date range", func() {
			// Create CSV with transactions spanning multiple months
			csvContent := `date,amount,content
2025/01/15,1000,January Mid
2025/01/31,2000,January End
2025/02/01,3000,February Start
2025/02/15,4000,February Mid
2025/03/01,5000,March Start`
			Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

			// Create optimized services
			source := in.NewCSVFileSource()
			optimizedTransactionService := usecase.NewOptimizedTransactionService(source)
			writer := output.NewJSONFile(outputPath)
			optimizedStatementService := usecase.NewOptimizedStatementService(optimizedTransactionService, writer)

			// Generate statement for January 15 to February 15
			startDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC)
			err := optimizedStatementService.GenerateStatementByDateRangeOptimized(ctx, csvPath, "2025/01-02", startDate, endDate)

			Expect(err).NotTo(HaveOccurred())

			// Verify filtered transactions
			content, err := os.ReadFile(outputPath)
			Expect(err).NotTo(HaveOccurred())
			output := string(content)

			Expect(output).To(ContainSubstring(`"period": "2025/01-02"`))
			Expect(output).To(ContainSubstring(`"total_income": 10000`)) // 1000 + 2000 + 3000 + 4000
			Expect(output).To(ContainSubstring(`"January Mid"`))
			Expect(output).To(ContainSubstring(`"January End"`))
			Expect(output).To(ContainSubstring(`"February Start"`))
			Expect(output).To(ContainSubstring(`"February Mid"`))
			Expect(output).NotTo(ContainSubstring(`"March Start"`))
		})
	})

	Context("when handling errors", func() {
		It("should return error for non-existent CSV file", func() {
			// Create optimized services
			source := in.NewCSVFileSource()
			optimizedTransactionService := usecase.NewOptimizedTransactionService(source)
			writer := output.NewJSONFile(outputPath)
			optimizedStatementService := usecase.NewOptimizedStatementService(optimizedTransactionService, writer)

			// Try to generate statement from non-existent file
			err := optimizedStatementService.GenerateMonthlyStatementOptimized(ctx, "non-existent.csv", "2025/01", 2025, 1)

			Expect(err).To(HaveOccurred())
			Expect(outputPath).NotTo(BeAnExistingFile())
		})

		It("should return error for invalid CSV format", func() {
			// Create invalid CSV
			invalidCSV := `invalid,header
2025/01/01,1000,Test`
			Expect(os.WriteFile(csvPath, []byte(invalidCSV), 0644)).To(Succeed())

			// Create optimized services
			source := in.NewCSVFileSource()
			optimizedTransactionService := usecase.NewOptimizedTransactionService(source)
			writer := output.NewJSONFile(outputPath)
			optimizedStatementService := usecase.NewOptimizedStatementService(optimizedTransactionService, writer)

			// Try to generate statement
			err := optimizedStatementService.GenerateMonthlyStatementOptimized(ctx, csvPath, "2025/01", 2025, 1)

			Expect(err).To(HaveOccurred())
			Expect(outputPath).NotTo(BeAnExistingFile())
		})
	})

	Context("when testing memory efficiency", func() {
		It("should handle large CSV files efficiently", func() {
			// Create a large CSV with 1000 transactions
			var csvBuilder strings.Builder
			csvBuilder.WriteString("date,amount,content\n")
			for i := 0; i < 1000; i++ {
				csvBuilder.WriteString("2025/01/01,1000,Transaction\n")
			}
			Expect(os.WriteFile(csvPath, []byte(csvBuilder.String()), 0644)).To(Succeed())

			// Create optimized services
			source := in.NewCSVFileSource()
			optimizedTransactionService := usecase.NewOptimizedTransactionService(source)
			writer := output.NewJSONFile(outputPath)
			optimizedStatementService := usecase.NewOptimizedStatementService(optimizedTransactionService, writer)

			// Generate statement
			err := optimizedStatementService.GenerateMonthlyStatementOptimized(ctx, csvPath, "2025/01", 2025, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(outputPath).To(BeAnExistingFile())

			// Verify all transactions are included
			content, err := os.ReadFile(outputPath)
			Expect(err).NotTo(HaveOccurred())
			output := string(content)

			Expect(output).To(ContainSubstring(`"total_income": 1000000`)) // 1000 * 1000
			Expect(output).To(ContainSubstring(`"Transaction"`))
		})
	})
})

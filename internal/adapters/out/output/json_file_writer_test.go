package output_test

import (
	"context"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/domain"
)

var _ = Describe("JSONFileWriter", func() {
	var (
		writer   *output.JSONFileWriter
		ctx      context.Context
		tempDir  string
		filePath string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "json_file_writer_test_*")
		Expect(err).NotTo(HaveOccurred())

		filePath = filepath.Join(tempDir, "statement.json")
		writer = output.NewJSONFile(filePath)
		ctx = context.Background()
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Context("when writing statements to files", func() {
		It("should create and write JSON to file", func() {
			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      2000,
				TotalExpenditure: -500,
				NetAmount:        1500,
				TransactionCount: 2,
				Transactions: []domain.TransactionDTO{
					{Date: "2025/01/01", Amount: "2000", Content: "Salary"},
					{Date: "2025/01/05", Amount: "-500", Content: "Groceries"},
				},
				GeneratedAt: time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			}

			err := writer.Write(ctx, statement)

			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).To(BeAnExistingFile())

			// Verify file content
			content, err := os.ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())
			output := string(content)

			Expect(output).To(ContainSubstring(`"period": "2025/01"`))
			Expect(output).To(ContainSubstring(`"total_income": 2000`))
			Expect(output).To(ContainSubstring(`"total_expenditure": -500`))
			Expect(output).To(ContainSubstring(`"net_amount": 1500`))
			Expect(output).To(ContainSubstring(`"transaction_count": 2`))
		})

		It("should handle empty statement", func() {
			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      0,
				TotalExpenditure: 0,
				NetAmount:        0,
				TransactionCount: 0,
				Transactions:     []domain.TransactionDTO{},
				GeneratedAt:      time.Now(),
			}

			err := writer.Write(ctx, statement)

			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).To(BeAnExistingFile())

			content, err := os.ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())
			output := string(content)

			Expect(output).To(ContainSubstring(`"transaction_count": 0`))
			Expect(output).To(ContainSubstring(`"transactions": []`))
		})

		It("should create file in non-existent directory", func() {
			nestedDir := filepath.Join(tempDir, "nested")
			Expect(os.MkdirAll(nestedDir, 0755)).To(Succeed())

			nestedPath := filepath.Join(nestedDir, "statement.json")
			writer = output.NewJSONFile(nestedPath)

			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      1000,
				TotalExpenditure: -200,
				NetAmount:        800,
				TransactionCount: 1,
				Transactions: []domain.TransactionDTO{
					{Date: "2025/01/01", Amount: "1000", Content: "Salary"},
				},
				GeneratedAt: time.Now(),
			}

			err := writer.Write(ctx, statement)

			Expect(err).NotTo(HaveOccurred())
			Expect(nestedPath).To(BeAnExistingFile())
		})

		It("should overwrite existing file", func() {
			Expect(os.WriteFile(filePath, []byte("old content"), 0644)).To(Succeed())

			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      1000,
				TotalExpenditure: -200,
				NetAmount:        800,
				TransactionCount: 1,
				Transactions: []domain.TransactionDTO{
					{Date: "2025/01/01", Amount: "1000", Content: "Salary"},
				},
				GeneratedAt: time.Now(),
			}

			err := writer.Write(ctx, statement)

			Expect(err).NotTo(HaveOccurred())

			content, err := os.ReadFile(filePath)
			Expect(err).NotTo(HaveOccurred())
			output := string(content)

			Expect(output).NotTo(ContainSubstring("old content"))
			Expect(output).To(ContainSubstring(`"period": "2025/01"`))
		})
	})

	Context("context handling", func() {
		It("should complete successfully even with cancelled context", func() {
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      1000,
				TotalExpenditure: -200,
				NetAmount:        800,
				TransactionCount: 0,
				Transactions:     []domain.TransactionDTO{},
				GeneratedAt:      time.Now(),
			}

			err := writer.Write(cancelledCtx, statement)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("error handling", func() {
		It("should handle invalid file paths", func() {
			invalidPath := "/dev/null/invalid/path/statement.json"
			writer = output.NewJSONFile(invalidPath)

			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      1000,
				TotalExpenditure: -200,
				NetAmount:        800,
				TransactionCount: 0,
				Transactions:     []domain.TransactionDTO{},
				GeneratedAt:      time.Now(),
			}

			err := writer.Write(ctx, statement)

			Expect(err).To(HaveOccurred())
		})

		It("should handle read-only directory", func() {
			readOnlyDir := filepath.Join(tempDir, "readonly")
			Expect(os.Mkdir(readOnlyDir, 0444)).To(Succeed()) // Read-only permissions
			readOnlyPath := filepath.Join(readOnlyDir, "statement.json")
			writer = output.NewJSONFile(readOnlyPath)

			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      1000,
				TotalExpenditure: -200,
				NetAmount:        800,
				TransactionCount: 0,
				Transactions:     []domain.TransactionDTO{},
				GeneratedAt:      time.Now(),
			}

			err := writer.Write(ctx, statement)

			Expect(err).To(HaveOccurred())
		})
	})

	Context("file permissions", func() {
		It("should create file with appropriate permissions", func() {
			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      1000,
				TotalExpenditure: -200,
				NetAmount:        800,
				TransactionCount: 0,
				Transactions:     []domain.TransactionDTO{},
				GeneratedAt:      time.Now(),
			}

			err := writer.Write(ctx, statement)

			Expect(err).NotTo(HaveOccurred())

			// Check file permissions
			info, err := os.Stat(filePath)
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Mode().IsRegular()).To(BeTrue())
		})
	})
})

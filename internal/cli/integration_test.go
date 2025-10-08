package cli_test

import (
	"bytes"
	"context"
	"io"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/adapters/out/parser"
	"mf-statement/internal/domain"
	"mf-statement/internal/usecase"
)

type mockSource struct {
	content string
}

func (m mockSource) Open(ctx context.Context, uri string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(m.content)), nil
}

type mockWriter struct {
	written domain.Statement
	called  bool
}

func (m *mockWriter) Write(ctx context.Context, s domain.Statement) error {
	m.called = true
	m.written = s
	return nil
}

var _ = Describe("Statement Generation Integration", func() {
	var (
		statementService   usecase.StatementService
		transactionService usecase.TransactionService
		writer             *mockWriter
		ctx                context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		writer = &mockWriter{}
	})

	Context("given valid CSV data with mixed months", func() {
		It("should generate correct monthly statement", func() {
			// Given
			csvContent := `date,amount,content
2025/01/05,2000,Salary
2025/01/09,-300,Grocery
2025/01/01,100,Gift
2025/02/01,999,Next Month (ignored)
2025/01/15,-150,Transport`

			source := mockSource{content: csvContent}
			csvParser := parser.NewCSV()
			transactionService = usecase.NewTransactionService(source, csvParser)
			statementService = usecase.NewStatementService(transactionService, writer)

			// When
			err := statementService.GenerateMonthlyStatement(ctx, "test.csv", "2025/01", 2025, 1)

			// Then
			Expect(err).NotTo(HaveOccurred())
			Expect(writer.called).To(BeTrue(), "Writer should be called")

			statement := writer.written
			Expect(statement.Period).To(Equal("2025/01"))
			Expect(statement.TotalIncome).To(Equal(int64(2100)))
			Expect(statement.TotalExpenditure).To(Equal(int64(-450)))
			Expect(statement.Transactions).To(HaveLen(4))

			// Verify transactions are sorted by date (newest first)
			expectedDates := []string{"2025/01/15", "2025/01/09", "2025/01/05", "2025/01/01"}
			for i, expectedDate := range expectedDates {
				Expect(statement.Transactions[i].Date).To(Equal(expectedDate), "Transaction %d should have correct date", i)
			}
		})
	})

	Context("error handling", func() {
		It("should handle invalid CSV data gracefully", func() {
			// Given
			invalidCSV := `date,amount,content
invalid-date,not-a-number,Empty Content`

			source := mockSource{content: invalidCSV}
			csvParser := parser.NewCSV()
			transactionService = usecase.NewTransactionService(source, csvParser)
			statementService = usecase.NewStatementService(transactionService, writer)

			// When
			err := statementService.GenerateMonthlyStatement(ctx, "test.csv", "2025/01", 2025, 1)

			// Then
			Expect(err).To(HaveOccurred(), "Should return error for invalid CSV")
			Expect(err).To(BeAssignableToTypeOf(domain.DomainError{}), "Should return domain error")
		})
	})

	Context("JSON output", func() {
		It("should generate valid JSON with all required fields", func() {
			// Given
			csvContent := `date,amount,content
2025/01/05,2000,Salary
2025/01/09,-300,Grocery`

			var buf bytes.Buffer
			jsonWriter := output.NewJSON(&buf)
			source := mockSource{content: csvContent}
			csvParser := parser.NewCSV()
			transactionService = usecase.NewTransactionService(source, csvParser)
			statementService = usecase.NewStatementService(transactionService, jsonWriter)

			// When
			err := statementService.GenerateMonthlyStatement(ctx, "test.csv", "2025/01", 2025, 1)

			// Then
			Expect(err).NotTo(HaveOccurred())

			output := buf.String()
			expectedFields := []string{
				`"period": "2025/01"`,
				`"total_income": 2000`,
				`"total_expenditure": -300`,
			}

			for _, field := range expectedFields {
				Expect(output).To(ContainSubstring(field), "JSON should contain %s", field)
			}
		})
	})
})

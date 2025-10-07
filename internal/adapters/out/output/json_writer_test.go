package output_test

import (
	"bytes"
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/domain"
)

var _ = Describe("JSONWriter", func() {
	var (
		writer *output.JSONWriter
		buf    *bytes.Buffer
		ctx    context.Context
	)

	BeforeEach(func() {
		buf = new(bytes.Buffer)
		writer = output.NewJSON(buf)
		ctx = context.Background()
	})

	Context("when writing statements", func() {
		It("should write valid JSON to the writer", func() {
			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      2000,
				TotalExpenditure: -500,
				Transactions: []domain.TransactionDTO{
					{Date: "2025/01/01", Amount: "2000", Content: "Salary"},
					{Date: "2025/01/05", Amount: "-500", Content: "Groceries"},
				},
			}

			err := writer.Write(ctx, statement)

			Expect(err).NotTo(HaveOccurred())
			output := buf.String()

			Expect(output).To(ContainSubstring(`"period": "2025/01"`))
			Expect(output).To(ContainSubstring(`"total_income": 2000`))
			Expect(output).To(ContainSubstring(`"total_expenditure": -500`))

			Expect(output).To(ContainSubstring(`"transactions": [`))
			Expect(output).To(ContainSubstring(`"date": "2025/01/01"`))
			Expect(output).To(ContainSubstring(`"amount": "2000"`))
			Expect(output).To(ContainSubstring(`"content": "Salary"`))
		})

		It("should handle empty statement", func() {
			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      0,
				TotalExpenditure: 0,
				Transactions:     []domain.TransactionDTO{},
			}

			err := writer.Write(ctx, statement)

			Expect(err).NotTo(HaveOccurred())
			output := buf.String()
			Expect(output).To(ContainSubstring(`"transactions": []`))
		})

		It("should format JSON with proper indentation", func() {
			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      1000,
				TotalExpenditure: -200,
				Transactions: []domain.TransactionDTO{
					{Date: "2025/01/01", Amount: "1000", Content: "Salary"},
				},
			}

			err := writer.Write(ctx, statement)

			Expect(err).NotTo(HaveOccurred())
			output := buf.String()

			lines := strings.Split(output, "\n")
			Expect(len(lines)).To(BeNumerically(">", 5)) // Should have multiple lines due to indentation
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
				Transactions:     []domain.TransactionDTO{},
			}

			err := writer.Write(cancelledCtx, statement)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("error handling", func() {
		It("should handle writer errors gracefully", func() {
			failingWriter := &failingWriter{}
			writer = output.NewJSON(failingWriter)

			statement := domain.Statement{
				Period:           "2025/01",
				TotalIncome:      1000,
				TotalExpenditure: -200,
				Transactions:     []domain.TransactionDTO{},
			}

			err := writer.Write(ctx, statement)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("write failed"))
		})
	})
})

// failingWriter is a test helper that always fails to write
type failingWriter struct{}

func (f *failingWriter) Write(p []byte) (n int, err error) {
	return 0, &mockError{message: "write failed"}
}

type mockError struct {
	message string
}

func (m *mockError) Error() string {
	return m.message
}

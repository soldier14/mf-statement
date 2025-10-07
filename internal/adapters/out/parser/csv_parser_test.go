package parser_test

import (
	"context"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/out/parser"
)

var _ = Describe("CSVParser", func() {
	var (
		csvParser *parser.CSVParser
		ctx       context.Context
	)

	BeforeEach(func() {
		csvParser = parser.NewCSV()
		ctx = context.Background()
	})

	Context("when parsing valid CSV data", func() {
		It("should parse transactions with correct headers", func() {
			csvContent := `date,amount,content
2025/01/01,1000,Salary
2025/01/05,-200,Groceries
2025/01/10,500,Bonus`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(3))

			Expect(transactions[0].Date).To(Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
			Expect(transactions[0].Amount).To(Equal(int64(1000)))
			Expect(transactions[0].Content).To(Equal("Salary"))

			Expect(transactions[1].Date).To(Equal(time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)))
			Expect(transactions[1].Amount).To(Equal(int64(-200)))
			Expect(transactions[1].Content).To(Equal("Groceries"))

			Expect(transactions[2].Date).To(Equal(time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)))
			Expect(transactions[2].Amount).To(Equal(int64(500)))
			Expect(transactions[2].Content).To(Equal("Bonus"))
		})

		It("should handle headers with different cases", func() {
			csvContent := `DATE,AMOUNT,CONTENT
2025/01/01,1000,Salary
2025/01/05,-200,Groceries`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
		})

		It("should handle headers with extra whitespace", func() {
			csvContent := `  date  ,  amount  ,  content  
2025/01/01,1000,Salary
2025/01/05,-200,Groceries`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
		})

		It("should handle BOM (Byte Order Mark) in headers", func() {
			csvContent := "\uFEFFdate,amount,content\n2025/01/01,1000,Salary"
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(1))
		})

		It("should handle empty CSV with only headers", func() {
			csvContent := `date,amount,content`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(0))
		})

		It("should handle transactions with whitespace", func() {
			csvContent := `date,amount,content
  2025/01/01  ,  1000  ,  Salary  
  2025/01/05  ,  -200  ,  Groceries  `
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
			Expect(transactions[0].Content).To(Equal("Salary"))
			Expect(transactions[1].Content).To(Equal("Groceries"))
		})
	})

	Context("when parsing invalid CSV data", func() {
		It("should return error for missing headers", func() {
			csvContent := `2025/01/01,1000,Salary`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
		})

		It("should return error for wrong number of columns in header", func() {
			csvContent := `date,amount
2025/01/01,1000`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
		})

		It("should return error for wrong header names", func() {
			csvContent := `wrong,header,names
2025/01/01,1000,Salary`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
		})

		It("should return error for invalid date format", func() {
			csvContent := `date,amount,content
invalid-date,1000,Salary`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("failed to parse date"))
		})

		It("should return error for invalid amount format", func() {
			csvContent := `date,amount,content
2025/01/01,not-a-number,Salary`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("failed to parse amount"))
		})

		It("should return error for empty columns", func() {
			csvContent := `date,amount,content
,1000,Salary`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("empty column"))
		})

		It("should return error for wrong number of columns in data row", func() {
			csvContent := `date,amount,content
2025/01/01,1000`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("expected 3 columns"))
		})

		It("should return error for empty content", func() {
			csvContent := `date,amount,content
2025/01/01,1000,`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("empty column"))
		})
	})

	Context("context handling", func() {
		It("should handle context cancellation", func() {
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			csvContent := `date,amount,content
2025/01/01,1000,Salary`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(cancelledCtx, reader)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(context.Canceled))
			Expect(transactions).To(BeNil())
		})
	})

	Context("edge cases", func() {
		It("should handle very large amounts", func() {
			csvContent := `date,amount,content
2025/01/01,9223372036854775807,Max Amount
2025/01/02,-9223372036854775808,Min Amount`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
			Expect(transactions[0].Amount).To(Equal(int64(9223372036854775807)))
			Expect(transactions[1].Amount).To(Equal(int64(-9223372036854775808)))
		})

		It("should handle zero amounts", func() {
			csvContent := `date,amount,content
2025/01/01,0,Zero Amount`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(1))
			Expect(transactions[0].Amount).To(Equal(int64(0)))
		})

		It("should handle special characters in content", func() {
			csvContent := `date,amount,content
2025/01/01,1000,"Salary with ""quotes"" and, commas"`
			reader := strings.NewReader(csvContent)

			transactions, err := csvParser.Parse(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(1))
			Expect(transactions[0].Content).To(Equal(`Salary with "quotes" and, commas`))
		})
	})
})

package parser_test

import (
	"context"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/out/parser"
	"mf-statement/internal/domain"
)

var _ = Describe("FilteredCSVParser", func() {
	var (
		filteredParser *parser.FilteredCSVParser
		ctx            context.Context
	)

	BeforeEach(func() {
		filteredParser = parser.NewFilteredCSV()
		ctx = context.Background()
	})

	Context("when parsing with period filter", func() {
		It("should filter transactions by year and month", func() {
			csvContent := `date,amount,content
2025/01/01,1000,January Salary
2025/02/01,2000,February Salary
2025/01/15,-500,January Expense
2025/02/15,-300,February Expense`
			reader := strings.NewReader(csvContent)

			transactions, err := filteredParser.ParseWithPeriodFilter(ctx, reader, 2025, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
			Expect(transactions[0].Content).To(Equal("January Salary"))
			Expect(transactions[1].Content).To(Equal("January Expense"))
		})

		It("should return empty slice when no transactions match period", func() {
			csvContent := `date,amount,content
2025/02/01,2000,February Salary
2025/03/01,3000,March Salary`
			reader := strings.NewReader(csvContent)

			transactions, err := filteredParser.ParseWithPeriodFilter(ctx, reader, 2025, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(0))
		})
	})

	Context("when parsing with date range filter", func() {
		It("should filter transactions by date range", func() {
			csvContent := `date,amount,content
2025/01/01,1000,New Year
2025/01/15,2000,Mid January
2025/02/01,3000,February
2025/01/31,4000,End January`
			reader := strings.NewReader(csvContent)

			startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)
			transactions, err := filteredParser.ParseWithDateRangeFilter(ctx, reader, startDate, endDate)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(3))
			// Transactions should be in chronological order (newest first)
			Expect(transactions[0].Content).To(Equal("New Year"))
			Expect(transactions[1].Content).To(Equal("Mid January"))
			Expect(transactions[2].Content).To(Equal("End January"))
		})

		It("should handle single day range", func() {
			csvContent := `date,amount,content
2025/01/15,1000,Single Day
2025/01/16,2000,Next Day`
			reader := strings.NewReader(csvContent)

			startDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
			transactions, err := filteredParser.ParseWithDateRangeFilter(ctx, reader, startDate, endDate)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(1))
			Expect(transactions[0].Content).To(Equal("Single Day"))
		})
	})

	Context("when parsing with custom filter", func() {
		It("should filter transactions by custom criteria", func() {
			csvContent := `date,amount,content
2025/01/01,1000,Salary
2025/01/02,-500,Expense
2025/01/03,2000,Bonus`
			reader := strings.NewReader(csvContent)

			// Filter for positive amounts only
			transactions, err := filteredParser.ParseWithFilter(ctx, reader, func(transaction domain.Transaction) bool {
				return transaction.Amount > 0
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
			// Transactions should be in chronological order (newest first)
			Expect(transactions[0].Content).To(Equal("Salary"))
			Expect(transactions[1].Content).To(Equal("Bonus"))
		})

		It("should filter by content pattern", func() {
			csvContent := `date,amount,content
2025/01/01,1000,Salary Payment
2025/01/02,-500,Grocery Shopping
2025/01/03,2000,Salary Bonus`
			reader := strings.NewReader(csvContent)

			// Filter for salary-related transactions
			transactions, err := filteredParser.ParseWithFilter(ctx, reader, func(transaction domain.Transaction) bool {
				return strings.Contains(strings.ToLower(transaction.Content), "salary")
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(2))
			// Transactions should be in chronological order (newest first)
			Expect(transactions[0].Content).To(Equal("Salary Payment"))
			Expect(transactions[1].Content).To(Equal("Salary Bonus"))
		})
	})

	Context("when handling invalid CSV data", func() {
		It("should return error for invalid headers", func() {
			csvContent := `invalid,header,format
2025/01/01,1000,Test`
			reader := strings.NewReader(csvContent)

			transactions, err := filteredParser.ParseWithPeriodFilter(ctx, reader, 2025, 1)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
		})

		It("should return error for malformed records", func() {
			csvContent := `date,amount,content
2025/01/01,1000,Valid
2025/01/02,Invalid Amount,Invalid
2025/01/03,2000,Valid`
			reader := strings.NewReader(csvContent)

			transactions, err := filteredParser.ParseWithPeriodFilter(ctx, reader, 2025, 1)

			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
		})
	})

	Context("when handling context cancellation", func() {
		It("should respect context cancellation", func() {
			csvContent := `date,amount,content
2025/01/01,1000,Test1
2025/01/02,2000,Test2
2025/01/03,3000,Test3`
			reader := strings.NewReader(csvContent)

			cancelledCtx, cancel := context.WithCancel(ctx)
			cancel() // Cancel immediately

			transactions, err := filteredParser.ParseWithPeriodFilter(cancelledCtx, reader, 2025, 1)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(context.Canceled))
			Expect(transactions).To(BeNil())
		})
	})

	Context("when handling large datasets", func() {
		It("should process large CSV efficiently", func() {
			// Generate a large CSV with 1000 transactions
			var csvBuilder strings.Builder
			csvBuilder.WriteString("date,amount,content\n")
			for i := 0; i < 1000; i++ {
				csvBuilder.WriteString("2025/01/01,1000,Transaction\n")
			}
			reader := strings.NewReader(csvBuilder.String())

			transactions, err := filteredParser.ParseWithPeriodFilter(ctx, reader, 2025, 1)

			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(1000))
		})
	})
})

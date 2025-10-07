package in_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/in"
	"mf-statement/internal/adapters/out/parser"
	"mf-statement/internal/domain"
)

var _ = Describe("CSVFileSource", func() {
	var (
		source *in.CSVFileSource
		ctx    context.Context
	)

	BeforeEach(func() {
		source = in.NewCSVFileSource()
		ctx = context.Background()
	})

	Context("when opening files", func() {
		It("should open direct file paths", func() {
			tempFile, err := os.CreateTemp("", "test_*.csv")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(tempFile.Name())

			_, err = tempFile.WriteString("date,amount,content\n2025/01/01,1000,Test")
			Expect(err).NotTo(HaveOccurred())
			tempFile.Close()

			reader, err := source.Open(ctx, tempFile.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(reader).NotTo(BeNil())

			content, err := io.ReadAll(reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("date,amount,content"))
		})

		It("should open file:// URIs", func() {
			tempFile, err := os.CreateTemp("", "test_*.csv")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(tempFile.Name())

			_, err = tempFile.WriteString("date,amount,content\n2025/01/01,1000,Test")
			Expect(err).NotTo(HaveOccurred())
			tempFile.Close()

			// Test opening with file:// URI
			fileURI := "file://" + tempFile.Name()
			reader, err := source.Open(ctx, fileURI)
			Expect(err).NotTo(HaveOccurred())
			Expect(reader).NotTo(BeNil())

			content, err := io.ReadAll(reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("date,amount,content"))
		})

		It("should return error for non-existent files", func() {
			reader, err := source.Open(ctx, "non-existent-file.csv")
			Expect(err).To(HaveOccurred())
			Expect(reader).To(BeNil())
		})

		It("should handle context cancellation", func() {
			// Create a context that's already cancelled
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			reader, err := source.Open(cancelledCtx, "some-file.csv")
			Expect(err).To(HaveOccurred())
			Expect(reader).To(BeNil())
		})
	})
})

var _ = Describe("CSVReaderService", func() {
	var (
		service *in.CSVReaderService
		ctx     context.Context
		tempDir string
		csvPath string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "csv_reader_test_*")
		Expect(err).NotTo(HaveOccurred())

		csvPath = filepath.Join(tempDir, "transactions.csv")
		csvContent := `date,amount,content
2025/01/01,1000,Salary
2025/01/05,-200,Groceries
2025/01/10,500,Bonus`
		Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

		source := in.NewCSVFileSource()
		parser := parser.NewCSV()
		service = in.NewCSVReaderService(source, parser)
		ctx = context.Background()
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Context("when reading valid CSV data", func() {
		It("should parse all transactions correctly", func() {
			transactions, err := service.ReadTransactions(ctx, csvPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(3))

			// Verify first transaction
			Expect(transactions[0].Date).To(Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
			Expect(transactions[0].Amount).To(Equal(int64(1000)))
			Expect(transactions[0].Content).To(Equal("Salary"))

			// Verify second transaction
			Expect(transactions[1].Date).To(Equal(time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)))
			Expect(transactions[1].Amount).To(Equal(int64(-200)))
			Expect(transactions[1].Content).To(Equal("Groceries"))

			// Verify third transaction
			Expect(transactions[2].Date).To(Equal(time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)))
			Expect(transactions[2].Amount).To(Equal(int64(500)))
			Expect(transactions[2].Content).To(Equal("Bonus"))
		})

		It("should handle file:// URIs", func() {
			fileURI := "file://" + csvPath
			transactions, err := service.ReadTransactions(ctx, fileURI)
			Expect(err).NotTo(HaveOccurred())
			Expect(transactions).To(HaveLen(3))
		})
	})

	Context("when reading invalid CSV data", func() {
		It("should return error for non-existent file", func() {
			transactions, err := service.ReadTransactions(ctx, "non-existent.csv")
			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(domain.DomainError{}))
		})

		It("should return error for malformed CSV", func() {
			// Create malformed CSV
			malformedPath := filepath.Join(tempDir, "malformed.csv")
			malformedContent := `date,amount,content
invalid-date,not-a-number,Empty Content`
			Expect(os.WriteFile(malformedPath, []byte(malformedContent), 0644)).To(Succeed())

			transactions, err := service.ReadTransactions(ctx, malformedPath)
			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(domain.DomainError{}))
		})
	})

	Context("context handling", func() {
		It("should handle context cancellation", func() {
			// Create a context that's already cancelled
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel()

			transactions, err := service.ReadTransactions(cancelledCtx, csvPath)
			Expect(err).To(HaveOccurred())
			Expect(transactions).To(BeNil())
		})
	})
})

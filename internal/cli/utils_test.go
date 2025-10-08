package cli_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/cli"
)

var _ = Describe("CLI Utils", func() {
	Context("ParsePeriod", func() {
		It("should parse valid period", func() {
			year, month, display, err := cli.ParsePeriod("202501")

			Expect(err).NotTo(HaveOccurred())
			Expect(year).To(Equal(2025))
			Expect(month).To(Equal(1))
			Expect(display).To(Equal("2025/01"))
		})

		It("should return error for invalid format", func() {
			_, _, _, err := cli.ParsePeriod("invalid")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("period must be in YYYYMM format"))
		})

		It("should return error for invalid month", func() {
			_, _, _, err := cli.ParsePeriod("202513")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("month must be between 01 and 12"))
		})
	})

	Context("CreateWriter", func() {
		It("should create JSON writer for stdout", func() {
			writer := cli.CreateWriter("")

			Expect(writer).NotTo(BeNil())
			// Should be a JSON writer (stdout)
			_, ok := writer.(*output.JSONWriter)
			Expect(ok).To(BeTrue())
		})

		It("should create JSON file writer for file path", func() {
			tempDir, err := os.MkdirTemp("", "cli_test_*")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tempDir)

			filePath := filepath.Join(tempDir, "test.json")
			writer := cli.CreateWriter(filePath)

			Expect(writer).NotTo(BeNil())
			// Should be a JSON file writer
			_, ok := writer.(*output.JSONFileWriter)
			Expect(ok).To(BeTrue())
		})
	})
})

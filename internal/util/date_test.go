package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"mf-statement/internal/util"
)

var _ = Describe("ParseYYYYMM", func() {
	Context("with valid inputs", func() {
		It("parses 202501", func() {
			year, month, display, err := util.ParseYYYYMM("202501")
			Expect(err).NotTo(HaveOccurred())
			Expect(year).To(Equal(2025))
			Expect(month).To(Equal(1))
			Expect(display).To(Equal("2025/01"))
		})

		It("parses 202412", func() {
			year, month, display, err := util.ParseYYYYMM("202412")
			Expect(err).NotTo(HaveOccurred())
			Expect(year).To(Equal(2024))
			Expect(month).To(Equal(12))
			Expect(display).To(Equal("2024/12"))
		})
	})

	Context("with invalid inputs", func() {
		It("rejects too short value", func() {
			_, _, _, err := util.ParseYYYYMM("2025")
			Expect(err).To(HaveOccurred())
		})

		It("rejects invalid month", func() {
			_, _, _, err := util.ParseYYYYMM("202513")
			Expect(err).To(HaveOccurred())
		})

		It("rejects non-numeric", func() {
			_, _, _, err := util.ParseYYYYMM("abcdxx")
			Expect(err).To(HaveOccurred())
		})
	})
})

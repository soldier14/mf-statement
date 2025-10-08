package util_test

import (
	"time"

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

var _ = Describe("Between", func() {
	Context("with valid date ranges", func() {
		It("should include dates within range", func() {
			start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

			// Test dates within range
			Expect(util.Between(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue())  // start date
			Expect(util.Between(time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue()) // middle
			Expect(util.Between(time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue()) // end date
		})

		It("should exclude dates outside range", func() {
			start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

			// Test dates outside range
			Expect(util.Between(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), start, end)).To(BeFalse()) // before start
			Expect(util.Between(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC), start, end)).To(BeFalse())   // after end
		})

		It("should handle 3-day range correctly (1/1/2025 to 1/3/2025)", func() {
			start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)

			// Should include all 3 days
			Expect(util.Between(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue()) // Day 1
			Expect(util.Between(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue()) // Day 2
			Expect(util.Between(time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), start, end)).To(BeTrue()) // Day 3

			// Should exclude day 4
			Expect(util.Between(time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC), start, end)).To(BeFalse()) // Day 4
		})
	})
})

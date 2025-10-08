package main

import (
	"fmt"
	"mf-statement/internal/util"
	"time"
)

func main() {
	// Test dates
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)

	// Test cases
	testDates := []time.Time{
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),   // Day 1 - should be included
		time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),   // Day 2 - should be included
		time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),   // Day 3 - should be included
		time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),   // Day 4 - should be excluded
		time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), // Before range - should be excluded
	}

	fmt.Printf("Testing date range: %s to %s (inclusive)\n",
		startDate.Format("2006/01/02"), endDate.Format("2006/01/02"))
	fmt.Println("==================================================")

	for i, testDate := range testDates {
		isBetween := util.Between(testDate, startDate, endDate)
		status := "❌ EXCLUDED"
		if isBetween {
			status = "✅ INCLUDED"
		}

		fmt.Printf("Day %d: %s - %s\n", i+1, testDate.Format("2006/01/02"), status)
	}
}

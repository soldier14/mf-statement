package cli

import (
	"fmt"
	"os"
	"strconv"

	"mf-statement/internal/adapters/out/output"
)

// ParsePeriod parses a period string in YYYYMM format
func ParsePeriod(period string) (year, month int, display string, err error) {
	if len(period) != 6 {
		return 0, 0, "", fmt.Errorf("period must be in YYYYMM format, got %s", period)
	}

	year, err = strconv.Atoi(period[:4])
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid year in period: %w", err)
	}

	month, err = strconv.Atoi(period[4:])
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid month in period: %w", err)
	}

	if month < 1 || month > 12 {
		return 0, 0, "", fmt.Errorf("month must be between 01 and 12, got %02d", month)
	}

	display = fmt.Sprintf("%d/%02d", year, month)
	return year, month, display, nil
}

// CreateWriter creates an appropriate writer based on output path
func CreateWriter(outputPath string) output.Writer {
	if outputPath == "" {
		return output.NewJSON(os.Stdout)
	}
	return output.NewJSONFile(outputPath)
}

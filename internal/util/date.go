package util

import (
	"fmt"
	"time"
)

func ParseYYYYMM(s string) (year, month int, display string, err error) {
	t, err := time.Parse("200601", s)
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid period format (expected YYYYMM): %w", err)
	}
	return t.Year(), int(t.Month()), t.Format("2006/01"), nil
}

package util

import (
	"fmt"
	"strconv"
)

func ParseYYYYMM(s string) (year, month int, display string, err error) {
	if len(s) != 6 {
		return 0, 0, "", fmt.Errorf("expected YYYYMM")
	}
	y, err := strconv.Atoi(s[:4])
	if err != nil {
		return 0, 0, "", err
	}
	m, err := strconv.Atoi(s[4:])
	if err != nil {
		return 0, 0, "", err
	}
	if m < 1 || m > 12 {
		return 0, 0, "", fmt.Errorf("month out of range")
	}
	return y, m, fmt.Sprintf("%04d/%02d", y, m), nil
}

package parser

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"mf-statement/internal/domain"
)

type CSVParser struct{}

func NewCSV() *CSVParser { return &CSVParser{} }

const (
	colDate    = "date"
	colAmount  = "amount"
	colContent = "content"
)

func (p *CSVParser) Parse(ctx context.Context, r io.Reader) ([]domain.Transaction, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if err := validateHeader(header); err != nil {
		return nil, err
	}

	var (
		out      []domain.Transaction
		rowIndex = 2
	)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read record at line %d: %w", rowIndex, err)
		}
		tx, err := parseRecord(record)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", rowIndex, err)
		}
		out = append(out, tx)
		rowIndex++
	}
	return out, nil
}

func validateHeader(header []string) error {
	if len(header) < 3 {
		return fmt.Errorf("invalid header: expected 3 columns, got %d", len(header))
	}

	header[0] = strings.TrimPrefix(header[0], "\uFEFF")

	if !eq(header[0], colDate) || !eq(header[1], colAmount) || !eq(header[2], colContent) {
		return fmt.Errorf("unexpected header: %v (expected: %s,%s,%s)", header, colDate, colAmount, colContent)
	}
	return nil
}

func parseRecord(record []string) (domain.Transaction, error) {
	if len(record) != 3 {
		return domain.Transaction{}, domain.NewParseError(
			fmt.Sprintf("invalid record: expected 3 columns, got %d", len(record)),
			fmt.Errorf("record: %v", record),
		)
	}

	dateStr := strings.TrimSpace(record[0])
	amountStr := strings.TrimSpace(record[1])
	content := strings.TrimSpace(record[2])

	if dateStr == "" || amountStr == "" || content == "" {
		return domain.Transaction{}, domain.NewValidationError(
			"empty column in record",
			map[string]interface{}{
				"record":  record,
				"date":    dateStr,
				"amount":  amountStr,
				"content": content,
			},
		)
	}

	date, err := time.Parse(domain.CSVDateLayout, dateStr)
	if err != nil {
		return domain.Transaction{}, domain.NewParseError(
			fmt.Sprintf("failed to parse date: %s", dateStr),
			err,
		)
	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		return domain.Transaction{}, domain.NewParseError(
			fmt.Sprintf("failed to parse amount: %s", amountStr),
			err,
		)
	}

	return domain.NewTransaction(date, amount, content)
}

func eq(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

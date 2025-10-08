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
	"mf-statement/internal/util"
)

// FilteredCSVParser provides memory-efficient CSV parsing with early filtering
type FilteredCSVParser struct{}

func NewFilteredCSV() *FilteredCSVParser {
	return &FilteredCSVParser{}
}

const (
	streamingColDate    = "date"
	streamingColAmount  = "amount"
	streamingColContent = "content"
)

// ParseWithFilter parses CSV and filters transactions during parsing to reduce memory usage
func (p *FilteredCSVParser) ParseWithFilter(ctx context.Context, r io.Reader, filterFunc func(domain.Transaction) bool) ([]domain.Transaction, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	// Optimize buffer size for large files
	reader.ReuseRecord = true

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if err := streamingValidateHeader(header); err != nil {
		return nil, err
	}

	var (
		transactions []domain.Transaction
		rowIndex     = 2
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

		transaction, err := streamingParseRecord(record)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", rowIndex, err)
		}

		// Early filtering - only add to result if it matches the filter
		if filterFunc(transaction) {
			transactions = append(transactions, transaction)
		}

		rowIndex++
	}
	return transactions, nil
}

// ParseWithPeriodFilter parses CSV and filters by year/month during parsing
func (p *FilteredCSVParser) ParseWithPeriodFilter(ctx context.Context, r io.Reader, year, month int) ([]domain.Transaction, error) {
	return p.ParseWithFilter(ctx, r, func(transaction domain.Transaction) bool {
		return transaction.Date.Year() == year && int(transaction.Date.Month()) == month
	})
}

// ParseWithDateRangeFilter parses CSV and filters by date range during parsing
func (p *FilteredCSVParser) ParseWithDateRangeFilter(ctx context.Context, r io.Reader, startDate, endDate time.Time) ([]domain.Transaction, error) {
	return p.ParseWithFilter(ctx, r, func(transaction domain.Transaction) bool {
		return util.Between(transaction.Date, startDate, endDate)
	})
}

func streamingValidateHeader(header []string) error {
	if len(header) < 3 {
		return fmt.Errorf("invalid header: expected 3 columns, got %d", len(header))
	}

	header[0] = strings.TrimPrefix(header[0], "\uFEFF")

	if !streamingEq(header[0], streamingColDate) || !streamingEq(header[1], streamingColAmount) || !streamingEq(header[2], streamingColContent) {
		return fmt.Errorf("unexpected header: %v (expected: %s,%s,%s)", header, streamingColDate, streamingColAmount, streamingColContent)
	}
	return nil
}

func streamingParseRecord(record []string) (domain.Transaction, error) {
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

func streamingEq(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

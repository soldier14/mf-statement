package usecase

import (
	"context"
	"mf-statement/internal/adapters/out/parser"
	"mf-statement/internal/domain"
	"sort"
	"time"
)

// OptimizedTransactionService provides memory-efficient transaction processing
type OptimizedTransactionService struct {
	Source         Source
	FilteredParser *parser.FilteredCSVParser
	Validator      Validator
}

func NewOptimizedTransactionService(source Source) *OptimizedTransactionService {
	return &OptimizedTransactionService{
		Source:         source,
		FilteredParser: parser.NewFilteredCSV(),
		Validator:      NewPeriodValidator(),
	}
}

// GetTransactionsByPeriodOptimized uses streaming parser with early filtering
func (s *OptimizedTransactionService) GetTransactionsByPeriodOptimized(ctx context.Context, csvFileURI string, year, month int) ([]domain.Transaction, error) {
	if err := s.Validator.ValidatePeriod(year, month); err != nil {
		return nil, err
	}

	csvReader, err := s.Source.Open(ctx, csvFileURI)
	if err != nil {
		return nil, domain.NewIOError("failed to open CSV source", err)
	}
	defer csvReader.Close()

	// Use filtered parser with period filter - only loads relevant transactions
	transactions, err := s.FilteredParser.ParseWithPeriodFilter(ctx, csvReader, year, month)
	if err != nil {
		return nil, domain.NewParseError("failed to parse CSV", err)
	}

	// Sort by date (newest first)
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].Date.After(transactions[j].Date)
	})

	return transactions, nil
}

// GetTransactionsByDateRangeOptimized uses streaming parser with date range filtering
func (s *OptimizedTransactionService) GetTransactionsByDateRangeOptimized(ctx context.Context, csvFileURI string, startDate, endDate time.Time) ([]domain.Transaction, error) {
	csvReader, err := s.Source.Open(ctx, csvFileURI)
	if err != nil {
		return nil, domain.NewIOError("failed to open CSV source", err)
	}
	defer csvReader.Close()

	// Use filtered parser with date range filter - only loads relevant transactions
	transactions, err := s.FilteredParser.ParseWithDateRangeFilter(ctx, csvReader, startDate, endDate)
	if err != nil {
		return nil, domain.NewParseError("failed to parse CSV", err)
	}

	// Sort by date (newest first)
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].Date.After(transactions[j].Date)
	})

	return transactions, nil
}

// CalculateTotalsOptimized calculates totals with early exit for large datasets
func (s *OptimizedTransactionService) CalculateTotalsOptimized(transactions []domain.Transaction) (totalIncome, totalExpenditure int64) {
	for _, transaction := range transactions {
		if transaction.IsIncome() {
			totalIncome += transaction.Amount
		} else if transaction.IsExpense() {
			totalExpenditure += transaction.Amount
		}
	}
	return totalIncome, totalExpenditure
}

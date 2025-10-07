package usecase

import (
	"context"
	"io"
	"mf-statement/internal/domain"
	"sort"
	"time"
)

type Source interface {
	Open(ctx context.Context, uri string) (io.ReadCloser, error)
}

type Parser interface {
	Parse(ctx context.Context, reader io.Reader) ([]domain.Transaction, error)
}

type TransactionService struct {
	Source Source
	Parser Parser
}

func NewTransactionService(source Source, parser Parser) *TransactionService {
	return &TransactionService{
		Source: source,
		Parser: parser,
	}
}

func (s *TransactionService) GetAllTransactions(ctx context.Context, csvFileURI string) ([]domain.Transaction, error) {
	csvReader, err := s.Source.Open(ctx, csvFileURI)
	if err != nil {
		return nil, domain.NewIOError("failed to open CSV source", err)
	}
	defer csvReader.Close()

	transactions, err := s.Parser.Parse(ctx, csvReader)
	if err != nil {
		return nil, domain.NewParseError("failed to parse CSV", err)
	}

	return transactions, nil
}

func (s *TransactionService) GetTransactionsByPeriod(ctx context.Context, csvFileURI string, year, month int) ([]domain.Transaction, error) {
	if year < 1900 || year > 2100 {
		return nil, domain.NewValidationError("invalid year", map[string]interface{}{
			"year": year,
		})
	}

	if month < 1 || month > 12 {
		return nil, domain.NewValidationError("invalid month", map[string]interface{}{
			"month": month,
		})
	}

	allTransactions, err := s.GetAllTransactions(ctx, csvFileURI)
	if err != nil {
		return nil, err
	}

	var filteredTransactions []domain.Transaction
	for _, transaction := range allTransactions {
		if transaction.Date.Year() == year && int(transaction.Date.Month()) == month {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	sort.SliceStable(filteredTransactions, func(i, j int) bool {
		return filteredTransactions[i].Date.After(filteredTransactions[j].Date)
	})

	return filteredTransactions, nil
}

func (s *TransactionService) GetTransactionsByDateRange(ctx context.Context, csvFileURI string, startDate, endDate time.Time) ([]domain.Transaction, error) {
	allTransactions, err := s.GetAllTransactions(ctx, csvFileURI)
	if err != nil {
		return nil, err
	}

	var filteredTransactions []domain.Transaction
	for _, transaction := range allTransactions {
		if (transaction.Date.Equal(startDate) || transaction.Date.After(startDate)) &&
			(transaction.Date.Equal(endDate) || transaction.Date.Before(endDate)) {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	sort.SliceStable(filteredTransactions, func(i, j int) bool {
		return filteredTransactions[i].Date.After(filteredTransactions[j].Date)
	})

	return filteredTransactions, nil
}

func (s *TransactionService) CalculateTotals(transactions []domain.Transaction) (totalIncome, totalExpenditure int64) {
	for _, transaction := range transactions {
		if transaction.IsIncome() {
			totalIncome += transaction.Amount
		} else if transaction.IsExpense() {
			totalExpenditure += transaction.Amount
		}
	}
	return totalIncome, totalExpenditure
}

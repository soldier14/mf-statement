package usecase

import (
	"context"
	"mf-statement/internal/domain"
	"mf-statement/internal/util"
	"sort"
	"time"
)

type TransactionService interface {
	GetAllTransactions(ctx context.Context, csvFileURI string) ([]domain.Transaction, error)
	GetTransactionsByPeriod(ctx context.Context, csvFileURI string, year, month int) ([]domain.Transaction, error)
	GetTransactionsByDateRange(ctx context.Context, csvFileURI string, startDate, endDate time.Time) ([]domain.Transaction, error)
	CalculateTotals(transactions []domain.Transaction) (totalIncome, totalExpenditure int64)
}

type TransactionServiceImpl struct {
	Source    Source
	Parser    Parser
	Validator Validator
}

func NewTransactionService(source Source, parser Parser) TransactionService {
	return &TransactionServiceImpl{
		Source:    source,
		Parser:    parser,
		Validator: NewPeriodValidator(),
	}
}

func (s *TransactionServiceImpl) GetAllTransactions(ctx context.Context, csvFileURI string) ([]domain.Transaction, error) {
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

func (s *TransactionServiceImpl) GetTransactionsByPeriod(ctx context.Context, csvFileURI string, year, month int) ([]domain.Transaction, error) {
	if err := s.Validator.ValidatePeriod(year, month); err != nil {
		return nil, err
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

func (s *TransactionServiceImpl) GetTransactionsByDateRange(ctx context.Context, csvFileURI string, startDate, endDate time.Time) ([]domain.Transaction, error) {
	allTransactions, err := s.GetAllTransactions(ctx, csvFileURI)
	if err != nil {
		return nil, err
	}

	var filteredTransactions []domain.Transaction
	for _, transaction := range allTransactions {
		if util.Between(transaction.Date, startDate, endDate) {
			filteredTransactions = append(filteredTransactions, transaction)
		}
	}

	sort.SliceStable(filteredTransactions, func(i, j int) bool {
		return filteredTransactions[i].Date.After(filteredTransactions[j].Date)
	})

	return filteredTransactions, nil
}

func (s *TransactionServiceImpl) CalculateTotals(transactions []domain.Transaction) (totalIncome, totalExpenditure int64) {
	for _, transaction := range transactions {
		if transaction.IsIncome() {
			totalIncome += transaction.Amount
		} else if transaction.IsExpense() {
			totalExpenditure += transaction.Amount
		}
	}
	return totalIncome, totalExpenditure
}

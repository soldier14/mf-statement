package usecase

import (
	"context"
	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/domain"
	"time"
)

type TransactionServiceInterface interface {
	GetTransactionsByPeriod(ctx context.Context, csvFileURI string, year, month int) ([]domain.Transaction, error)
	GetTransactionsByDateRange(ctx context.Context, csvFileURI string, startDate, endDate time.Time) ([]domain.Transaction, error)
	CalculateTotals(transactions []domain.Transaction) (totalIncome, totalExpenditure int64)
}

type StatementService struct {
	TransactionService TransactionServiceInterface
	Writer             output.Writer
}

func NewStatementService(transactionService TransactionServiceInterface, writer output.Writer) *StatementService {
	return &StatementService{
		TransactionService: transactionService,
		Writer:             writer,
	}
}

func (s *StatementService) GenerateMonthlyStatement(ctx context.Context, csvFileURI string, periodDisplay string, year, month int) error {
	transactions, err := s.TransactionService.GetTransactionsByPeriod(ctx, csvFileURI, year, month)
	if err != nil {
		return err
	}

	totalIncome, totalExpenditure := s.TransactionService.CalculateTotals(transactions)

	statement := domain.NewStatement(periodDisplay, transactions, totalIncome, totalExpenditure)

	if err := s.Writer.Write(ctx, statement); err != nil {
		return domain.NewIOError("failed to write statement", err)
	}

	return nil
}

func (s *StatementService) GenerateStatementFromTransactions(ctx context.Context, transactions []domain.Transaction, periodDisplay string) error {
	totalIncome, totalExpenditure := s.TransactionService.CalculateTotals(transactions)

	statement := domain.NewStatement(periodDisplay, transactions, totalIncome, totalExpenditure)

	if err := s.Writer.Write(ctx, statement); err != nil {
		return domain.NewIOError("failed to write statement", err)
	}

	return nil
}

func (s *StatementService) GenerateStatementByDateRange(ctx context.Context, csvFileURI string, periodDisplay string, startDate, endDate time.Time) error {
	transactions, err := s.TransactionService.GetTransactionsByDateRange(ctx, csvFileURI, startDate, endDate)
	if err != nil {
		return err
	}

	return s.GenerateStatementFromTransactions(ctx, transactions, periodDisplay)
}

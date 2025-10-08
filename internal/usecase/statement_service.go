package usecase

import (
	"context"
	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/domain"
	"time"
)

type StatementService interface {
	GenerateMonthlyStatement(ctx context.Context, csvFileURI string, periodDisplay string, year, month int) error
	GenerateStatementFromTransactions(ctx context.Context, transactions []domain.Transaction, periodDisplay string) error
	GenerateStatementByDateRange(ctx context.Context, csvFileURI string, periodDisplay string, startDate, endDate time.Time) error
}

type StatementServiceImpl struct {
	TransactionService TransactionService
	Writer             output.Writer
}

func NewStatementService(transactionService TransactionService, writer output.Writer) StatementService {
	return &StatementServiceImpl{
		TransactionService: transactionService,
		Writer:             writer,
	}
}

func (s *StatementServiceImpl) GenerateMonthlyStatement(ctx context.Context, csvFileURI string, periodDisplay string, year, month int) error {
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

func (s *StatementServiceImpl) GenerateStatementFromTransactions(ctx context.Context, transactions []domain.Transaction, periodDisplay string) error {
	totalIncome, totalExpenditure := s.TransactionService.CalculateTotals(transactions)

	statement := domain.NewStatement(periodDisplay, transactions, totalIncome, totalExpenditure)

	if err := s.Writer.Write(ctx, statement); err != nil {
		return domain.NewIOError("failed to write statement", err)
	}

	return nil
}

func (s *StatementServiceImpl) GenerateStatementByDateRange(ctx context.Context, csvFileURI string, periodDisplay string, startDate, endDate time.Time) error {
	transactions, err := s.TransactionService.GetTransactionsByDateRange(ctx, csvFileURI, startDate, endDate)
	if err != nil {
		return err
	}

	return s.GenerateStatementFromTransactions(ctx, transactions, periodDisplay)
}

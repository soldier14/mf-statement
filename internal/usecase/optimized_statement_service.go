package usecase

import (
	"context"
	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/domain"
	"time"
)

// OptimizedStatementService provides memory-efficient statement generation
type OptimizedStatementService struct {
	OptimizedTransactionService *OptimizedTransactionService
	Writer                      output.Writer
}

func NewOptimizedStatementService(optimizedTransactionService *OptimizedTransactionService, writer output.Writer) *OptimizedStatementService {
	return &OptimizedStatementService{
		OptimizedTransactionService: optimizedTransactionService,
		Writer:                      writer,
	}
}

// GenerateMonthlyStatementOptimized uses streaming processing for memory efficiency
func (s *OptimizedStatementService) GenerateMonthlyStatementOptimized(ctx context.Context, csvFileURI string, periodDisplay string, year, month int) error {
	// Use optimized transaction service with streaming parser
	transactions, err := s.OptimizedTransactionService.GetTransactionsByPeriodOptimized(ctx, csvFileURI, year, month)
	if err != nil {
		return err
	}

	// Calculate totals efficiently
	totalIncome, totalExpenditure := s.OptimizedTransactionService.CalculateTotalsOptimized(transactions)

	// Create statement
	statement := domain.NewStatement(periodDisplay, transactions, totalIncome, totalExpenditure)

	// Write statement
	if err := s.Writer.Write(ctx, statement); err != nil {
		return domain.NewIOError("failed to write statement", err)
	}

	return nil
}

// GenerateStatementByDateRangeOptimized uses streaming processing for date range queries
func (s *OptimizedStatementService) GenerateStatementByDateRangeOptimized(ctx context.Context, csvFileURI string, periodDisplay string, startDate, endDate time.Time) error {
	// Use optimized transaction service with streaming parser
	transactions, err := s.OptimizedTransactionService.GetTransactionsByDateRangeOptimized(ctx, csvFileURI, startDate, endDate)
	if err != nil {
		return err
	}

	// Calculate totals efficiently
	totalIncome, totalExpenditure := s.OptimizedTransactionService.CalculateTotalsOptimized(transactions)

	// Create statement
	statement := domain.NewStatement(periodDisplay, transactions, totalIncome, totalExpenditure)

	// Write statement
	if err := s.Writer.Write(ctx, statement); err != nil {
		return domain.NewIOError("failed to write statement", err)
	}

	return nil
}

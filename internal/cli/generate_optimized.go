package cli

import (
	"context"
	"fmt"
	"time"

	"mf-statement/internal/adapters/in"
	"mf-statement/internal/usecase"
	"mf-statement/internal/util"

	"github.com/spf13/cobra"
)

var generateOptimizedCmd = &cobra.Command{
	Use:   "generate-optimized",
	Short: "Generate financial statement with memory optimizations",
	Long: `Generate a financial statement from CSV transaction data with memory optimizations.
This command uses streaming parsing and early filtering to reduce memory usage.`,
	Example: `  # Generate statement for January 2025 with optimizations
  mf-statement generate-optimized --period 202501 --csv transactions.csv
  
  # Generate with custom output file
  mf-statement generate-optimized --period 202501 --csv transactions.csv --out statement.json
  
  # Generate with verbose logging
  mf-statement generate-optimized --period 202501 --csv transactions.csv --verbose
  
  # Generate with custom timeout
  mf-statement generate-optimized --period 202501 --csv transactions.csv --timeout 60`,
	RunE: runGenerateOptimized,
}

var (
	optimizedPeriod  string
	optimizedCSV     string
	optimizedOutput  string
	optimizedVerbose bool
	optimizedTimeout int
)

func init() {
	generateOptimizedCmd.Flags().StringVarP(&optimizedPeriod, "period", "p", "", "Month in YYYYMM format (e.g. 202501)")
	generateOptimizedCmd.Flags().StringVarP(&optimizedCSV, "csv", "c", "", "Path to CSV file or file:// URI")
	generateOptimizedCmd.Flags().StringVarP(&optimizedOutput, "out", "o", "", "Output JSON file path (default: stdout)")
	generateOptimizedCmd.Flags().BoolVarP(&optimizedVerbose, "verbose", "v", false, "Enable verbose logging")
	generateOptimizedCmd.Flags().IntVarP(&optimizedTimeout, "timeout", "t", 30, "Timeout in seconds for processing (default: 30)")

	generateOptimizedCmd.MarkFlagRequired("period")
	generateOptimizedCmd.MarkFlagRequired("csv")
}

func runGenerateOptimized(cmd *cobra.Command, args []string) error {
	logger := util.NewDefaultLogger()
	if optimizedVerbose {
		logger = util.NewDebugLogger()
	}

	logger.Info("Starting MF Statement CLI (Optimized)", "version", "1.0.0")

	year, month, periodDisplay, err := ParsePeriod(optimizedPeriod)
	if err != nil {
		return fmt.Errorf("invalid period: %w", err)
	}

	logger.Info("Generating statement for period", "period", periodDisplay)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(optimizedTimeout)*time.Second)
	defer cancel()

	if optimizedOutput != "" {
		logger.Info("Output will be written to file", "file", optimizedOutput)
	}

	// Create optimized services
	source := in.NewCSVFileSource()
	optimizedTransactionService := usecase.NewOptimizedTransactionService(source)
	writer := CreateWriter(optimizedOutput)
	optimizedStatementService := usecase.NewOptimizedStatementService(optimizedTransactionService, writer)

	// Generate statement with optimizations
	err = optimizedStatementService.GenerateMonthlyStatementOptimized(ctx, optimizedCSV, periodDisplay, year, month)
	if err != nil {
		logger.Error("Failed to generate statement", "error", err)
		return err
	}

	logger.Info("Statement generated successfully")
	return nil
}

package cli

import (
	"context"
	"mf-statement/internal/adapters/in"
	"mf-statement/internal/adapters/out/output"
	"mf-statement/internal/adapters/out/parser"
	"mf-statement/internal/domain"
	"mf-statement/internal/usecase"
	"mf-statement/internal/util"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func NewGenerateCommand() *cobra.Command {
	var (
		periodArg      string
		csvPath        string
		outputFilePath string
		verbose        bool
		timeout        int
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a monthly statement (CSV → JSON)",
		Long: `Reads a CSV of wallet transactions and outputs a monthly statement in JSON format.

The CSV file should have the following format:
  date,amount,content
  2025/01/05,2000,Salary
  2025/01/09,-300,Grocery

Where:
  - date: Date in YYYY/MM/DD format
  - amount: Amount in cents (positive for income, negative for expenses)
  - content: Description of the transaction`,
		Example: `  # Generate statement for January 2025
  mf-statement generate --period 202501 --csv transactions.csv
  
  # Generate with custom output file
  mf-statement generate --period 202501 --csv transactions.csv --out statement.json
  
  # Generate with verbose logging
  mf-statement generate --period 202501 --csv transactions.csv --verbose
  
  # Generate with custom timeout
  mf-statement generate --period 202501 --csv transactions.csv --timeout 60`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if periodArg == "" || csvPath == "" {
				_ = cmd.Help()
				return domain.NewValidationError("missing required arguments", map[string]interface{}{
					"period": periodArg,
					"csv":    csvPath,
				})
			}

			year, month, display, err := util.ParseYYYYMM(periodArg)
			if err != nil {
				return domain.NewValidationError("invalid period format", map[string]interface{}{
					"period": periodArg,
					"error":  err.Error(),
				})
			}

			if verbose {
				logger = util.NewDebugLogger()
			}

			logger.Info("Generating statement for period", "period", display)
			logger.Debug("CSV path", "path", csvPath)
			logger.Debug("Output file", "file", outputFilePath)

			var writer output.Writer
			if outputFilePath != "" {
				writer = output.NewJSONFile(outputFilePath)
				logger.Info("Output will be written to file", "file", outputFilePath)
			} else {
				writer = output.NewJSON(os.Stdout)
				logger.Info("Output will be written to stdout")
			}

			csvSource := in.NewCSVFileSource()
			csvParser := parser.NewCSV()

			transactionService := usecase.NewTransactionService(csvSource, csvParser)

			statementService := usecase.NewStatementService(transactionService, writer)

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
			defer cancel()

			if err := statementService.GenerateMonthlyStatement(ctx, csvPath, display, year, month); err != nil {
				logger.Error("Failed to generate statement", "error", err)
				return err
			}

			logger.Info("Statement generated successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&periodArg, "period", "p", "", "Month in YYYYMM format (e.g. 202501)")
	cmd.Flags().StringVarP(&csvPath, "csv", "c", "", "Path to CSV file or file:// URI")
	cmd.Flags().StringVarP(&outputFilePath, "out", "o", "", "Output JSON file path (default: stdout)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Timeout in seconds for processing (default: 30)")

	_ = cmd.MarkFlagRequired("period")
	_ = cmd.MarkFlagRequired("csv")

	return cmd
}

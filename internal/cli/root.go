package cli

import (
	"mf-statement/internal/util"
	"os"

	"github.com/spf13/cobra"
)

var (
	logger *util.Logger
)

func NewRootCommand() *cobra.Command {
	logger = util.NewDefaultLogger()

	root := &cobra.Command{
		Use:   "mf-statement",
		Short: "Monthly Financial Statement Generator",
		Long: `A command-line tool to process transaction CSVs into structured JSON format statements.
This tool helps you generate monthly financial statements from CSV transaction data,
calculating income, expenditure, and providing detailed transaction summaries.`,
		Version: "1.0.0",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Starting MF Statement CLI")
			return nil
		},
	}

	root.AddCommand(NewVersionCommand())
	root.AddCommand(NewGenerateCommand())
	root.AddCommand(generateOptimizedCmd)

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		logger.Error("Command failed", "error", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"io"
	"mf-statement/internal/output"
	"mf-statement/internal/parser"
	"mf-statement/internal/service"
	"mf-statement/internal/util"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type fileSource struct{}

func (fileSource) Open(ctx context.Context, uri string) (io.ReadCloser, error) {
	if u, err := url.Parse(uri); err == nil && (u.Scheme == "file") {
		return os.Open(u.Path)
	}
	return os.Open(uri)
}

func main() {
	var (
		periodArg      string
		csvPath        string
		timeout        time.Duration
		outputFilePath string
	)

	root := &cobra.Command{
		Use:   "statement",
		Short: "Generate a monthly statement (CSV → JSON)",
		Long:  "Reads a CSV of wallet transactions and prints a monthly statement in JSON.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if periodArg == "" || csvPath == "" {
				_ = cmd.Help()
				return fmt.Errorf("both --period and --csv are required")
			}

			year, month, display, err := util.ParseYYYYMM(periodArg)
			if err != nil {
				return fmt.Errorf("invalid period: %w", err)
			}
			var writer output.Writer
			if outputFilePath != "" {
				writer = output.NewJSONFile(outputFilePath)
			} else {
				writer = output.NewJSON(os.Stdout)
			}

			gen := &service.StatementGenerator{
				Source: fileSource{},
				Parse:  parser.NewCSV(),
				Write:  writer,
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			return gen.GenerateMonthly(ctx, display, year, month, csvPath)
		},
	}

	root.Flags().StringVarP(&periodArg, "period", "p", "", "Month in YYYYMM (e.g. 202201)")
	root.Flags().StringVarP(&csvPath, "csv", "c", "", "Path or file:// URI to CSV")
	root.Flags().DurationVar(&timeout, "timeout", 10*time.Second, "Execution timeout")
	root.Flags().StringVarP(&outputFilePath, "out", "o", "", "Output JSON file path (optional)")

	// Example future subcommands:
	// root.AddCommand(apiCmd)    // REST API entrypoint
	// root.AddCommand(workerCmd) // SQS worker

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

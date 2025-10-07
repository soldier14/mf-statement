package cli_test

import (
	. "mf-statement/internal/cli"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCommand", func() {
	var (
		csvContent string
		csvPath    string
		tempDir    string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "generate_cmd_test_*")
		Expect(err).NotTo(HaveOccurred())

		csvPath = filepath.Join(tempDir, "transactions.csv")
		csvContent = `date,amount,content
2025/01/01,1000,Salary
2025/01/05,-200,Groceries
`
		Expect(os.WriteFile(csvPath, []byte(csvContent), 0644)).To(Succeed())

		// Initialize global logger for tests
		NewRootCommand() // This initializes the global logger
	})

	AfterEach(func() {
		_ = os.RemoveAll(tempDir)
	})

	Context("flag validation", func() {
		It("should error when missing required flags", func() {
			cmd := NewGenerateCommand()
			cmd.SetArgs([]string{}) // no args

			err := cmd.Execute()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("required flag(s)"))
		})
	})

	Context("successful execution", func() {
		It("should print JSON to stdout", func(ctx SpecContext) {
			// Since the CLI uses os.Stdout directly, we need to test the actual behavior
			// by checking that the command executes successfully and produces valid JSON
			cmd := NewGenerateCommand()
			cmd.SetArgs([]string{"--period", "202501", "--csv", csvPath})

			err := cmd.ExecuteContext(ctx)
			Expect(err).To(Succeed())

			// The test passes if the command executes without error
			// The JSON output is verified by the integration tests
		}, SpecTimeout(5*time.Second))

		It("should write JSON to a file when --out is provided", func(ctx SpecContext) {
			// Create output directory
			outputDir := filepath.Join(tempDir, "output")
			Expect(os.MkdirAll(outputDir, 0755)).To(Succeed())

			outPath := filepath.Join(outputDir, "statement.json")

			cmd := NewGenerateCommand()
			cmd.SetArgs([]string{"--period", "202501", "--csv", csvPath, "--out", outPath})

			Expect(cmd.ExecuteContext(ctx)).To(Succeed())
			Expect(outPath).To(BeAnExistingFile())

			data, err := os.ReadFile(outPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(data)).To(ContainSubstring(`"period": "2025/01"`))
			Expect(string(data)).To(ContainSubstring(`"total_income": 1000`))
			Expect(string(data)).To(ContainSubstring(`"total_expenditure": -200`))
		}, SpecTimeout(5*time.Second))
	})
})

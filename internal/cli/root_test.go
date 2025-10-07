package cli_test

import (
	"bytes"
	. "mf-statement/internal/cli"

	"github.com/spf13/cobra"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RootCommand", func() {
	var root *cobra.Command

	BeforeEach(func() {
		root = NewRootCommand()
	})

	It("should have correct use and descriptions", func() {
		Expect(root.Use).To(Equal("mf-statement"))
		Expect(root.Short).To(Equal("Monthly Financial Statement Generator"))
		Expect(root.Long).To(ContainSubstring("transaction CSVs"))
	})

	It("should include the generate and version subcommands", func() {
		commands := root.Commands()
		commandNames := make([]string, len(commands))
		for i, cmd := range commands {
			commandNames[i] = cmd.Use
		}

		Expect(commandNames).To(ContainElements("generate", "version"))
	})

	It("should execute help text by default", func(ctx SpecContext) {
		buf := new(bytes.Buffer)
		root.SetOut(buf)
		root.SetArgs([]string{}) // no args should trigger help

		err := root.ExecuteContext(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(buf.String()).To(ContainSubstring("A command-line tool to process transaction CSVs"))
	})

	It("should execute version command successfully", func(ctx SpecContext) {
		root.SetArgs([]string{"version"})

		err := root.ExecuteContext(ctx)
		Expect(err).NotTo(HaveOccurred())
		// The version command prints to stdout, so we just verify it executes successfully
	})

	It("should handle invalid command", func(ctx SpecContext) {
		buf := new(bytes.Buffer)
		root.SetErr(buf)
		root.SetArgs([]string{"invalid"})

		err := root.ExecuteContext(ctx)
		Expect(err).To(HaveOccurred())
		Expect(buf.String()).To(ContainSubstring("unknown command"))
	})
})

var _ = Describe("Execute", func() {
	It("should execute root command successfully", func() {
		// Test the Execute function by creating a command that will succeed
		cmd := NewRootCommand()
		cmd.SetArgs([]string{"version"})

		// Test that the command can be executed without the Execute() wrapper
		err := cmd.Execute()
		Expect(err).ToNot(HaveOccurred())
		// The version command prints to stdout, so we just verify it executes successfully
	})

	It("should handle command execution errors", func() {
		cmd := NewRootCommand()
		cmd.SetArgs([]string{"nonexistent"})

		var buf bytes.Buffer
		cmd.SetErr(&buf)

		// Test that the command structure handles errors properly
		err := cmd.Execute()
		Expect(err).To(HaveOccurred())
		Expect(buf.String()).To(ContainSubstring("unknown command"))
	})
})

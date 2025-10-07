package cli_test

import (
	"mf-statement/internal/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("VersionCommand", func() {
	var cmd *cobra.Command

	BeforeEach(func() {
		cmd = cli.NewVersionCommand()
	})

	It("should create a version command with correct properties", func() {
		Expect(cmd.Use).To(Equal("version"))
		Expect(cmd.Short).To(Equal("Print version information"))
		Expect(cmd.Long).To(Equal("Print version information for the MF Statement CLI tool"))
	})

	It("should print version information", func() {
		cmd.SetArgs([]string{})

		err := cmd.Execute()

		Expect(err).ToNot(HaveOccurred())
		// The version command prints to stdout, so we just verify it executes successfully
	})

	It("should handle execution with arguments", func() {
		cmd.SetArgs([]string{"extra", "args"})

		err := cmd.Execute()

		Expect(err).ToNot(HaveOccurred())
		// The version command prints to stdout, so we just verify it executes successfully
	})
})

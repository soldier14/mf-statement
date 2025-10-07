package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information for the MF Statement CLI tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("MF Statement CLI\n")
			fmt.Printf("Version: %s\n", "1.0.0")
			fmt.Printf("Go Version: %s\n", runtime.Version())
			fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
}

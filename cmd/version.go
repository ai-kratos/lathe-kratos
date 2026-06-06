package cmd

import (
	"fmt"

	"github.com/devenjarvis/lathe/internal/buildinfo"
	"github.com/spf13/cobra"
)

// versionCmd is a friendly alias for `lathe --version`. It prints the same
// resolved version (plus commit/date when stamped) that the --version flag does.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the lathe version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprintf(cmd.OutOrStdout(), "lathe %s\n", buildinfo.String())
		return err
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

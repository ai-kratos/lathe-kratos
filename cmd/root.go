package cmd

import (
	"os"

	"github.com/devenjarvis/lathe/internal/buildinfo"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "lathe",
	Short:   "Generate and manage hands-on technical tutorials",
	Version: buildinfo.Resolve(),
}

func init() {
	// rootCmd.Version (above) enables cobra's --version flag; the template
	// controls what it prints. buildinfo.String() folds in commit/date when
	// they were stamped at build time, so --version shows the full provenance.
	rootCmd.SetVersionTemplate("lathe " + buildinfo.String() + "\n")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package cmd

import (
	"github.com/spf13/cobra"
)

// floatingCmd is the main command for actions to be done on floating windows.
var floatingCmd = &cobra.Command{
	Use:   "floating",
	Short: "Main command for floating windows. See --help for the sub-commands.",
}

func init() {
	rootCmd.AddCommand(floatingCmd)
}

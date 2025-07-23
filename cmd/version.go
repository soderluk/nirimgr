package cmd

import (
	"fmt"
	"runtime"

	"github.com/soderluk/nirimgr/config"
	"github.com/spf13/cobra"
)

// versionCmd prints out build information about nirimgr
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  "Prints the version number and build information about nirimgr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nirimgr " + buildInfo())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// buildInfo returns the build information about nirimgr
func buildInfo() string {
	info := config.Version
	info += fmt.Sprintf(" (%s %s %s %s)",
		runtime.Version(),
		runtime.GOARCH,
		runtime.GOOS,
		config.Date,
	)

	return info
}

package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

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
	info := fmt.Sprintf("\nVersion:\t%s\nCommit:\t%s\nGo Version:\t%s\nBuild Date:\t%s\nBuild info: \n",
		config.Version,
		config.CommitSHA,
		runtime.Version(),
		config.BuildDate,
	)
	bi, _ := debug.ReadBuildInfo()
	for _, setting := range bi.Settings {
		info += fmt.Sprintf("%s:\t%s\n", setting.Key, setting.Value)
	}

	return info
}

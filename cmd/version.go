package cmd

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/soderluk/nirimgr/config"
	"github.com/spf13/cobra"
)

// versionCmd prints out build information about nirimgr
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  "Prints the version number and build information about nirimgr",
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewTable(os.Stdout)
		table.Header([]string{"nirimgr", ""})
		c := 10
		info := [][]string{
			{"Version", config.Version},
			{"Commit", config.CommitSHA},
			{"Build Date", config.BuildDate},
			{r("-", c), r("-", c)},
			{"Build Info", ""},
			{r("-", c), r("-", c)},
			{"Go version", runtime.Version()},
		}
		bi, _ := debug.ReadBuildInfo()
		for _, s := range bi.Settings {
			info = append(info, []string{s.Key, s.Value})
		}
		for _, data := range info {
			if err := table.Append(data[0], data[1]); err != nil {
				fmt.Printf("Could not append %s to info\n", data[0])
			}
		}
		if err := table.Render(); err != nil {
			fmt.Println("Could not render table.")
		}
	},
}

// r returns the repeated string s, count c times.
func r(s string, c int) string {
	return strings.Repeat(s, c)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

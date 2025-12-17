package cmd

import (
	"github.com/soderluk/nirimgr/events"

	"github.com/spf13/cobra"
)

// eventsCmd listens to the Niri event stream.
var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Listen to niri event stream and act to events.",
	Long: `This command listens to the niri event stream, and when an event is seen,
		acts on it as defined in the configuration. See config.json rules section.`,
	Run: func(cmd *cobra.Command, args []string) {
		events.Run()
	},
}

func init() {
	RootCmd.AddCommand(eventsCmd)
}

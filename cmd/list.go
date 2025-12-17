package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/events"
	"github.com/spf13/cobra"
)

// listCmd lists the available actions and events defined in nirimgr.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available actions or events that nirimgr has defined.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]

		switch command {
		case "actions":
			listActions()
		case "events":
			listEvents()
		default:
			slog.Error("Unknown command", "cmd", command)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

// startTable creates the tablewriter and sets the header.
func startTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Name", "Fields"})

	return table
}

// listActions lists all the defined actions.
func listActions() {
	fmt.Println("nirimgr supports the following actions:")
	actionsSorted := msort(actions.ActionRegistry)

	table := startTable()

	for _, name := range actionsSorted {
		action := actions.ActionRegistry[name]
		model := action()
		fields := extractFields(model)
		err := table.Append([]string{name, fmt.Sprintf("%+v", fields)})
		if err != nil {
			fmt.Printf("could not append %v to table, error: %v", name, err)
			continue
		}
	}
	err := table.Render()
	if err != nil {
		fmt.Printf("could not render table %v", err)
	}
}

// listEvents lists all the defined events.
func listEvents() {
	fmt.Println("nirimgr supports the following events:")
	eventsSorted := msort(events.EventRegistry)

	table := startTable()
	for _, name := range eventsSorted {
		event := events.EventRegistry[name]
		model := event()
		fields := extractFields(model)
		err := table.Append([]string{name, fmt.Sprintf("%+v", fields)})
		if err != nil {
			fmt.Printf("could not append %v to table, error: %v", name, err)
			continue
		}
	}
	err := table.Render()
	if err != nil {
		fmt.Printf("could not render table %v", err)
	}
}

// msort sorts the given map by keys and returns the sorted list as a slice.
func msort[T any](m map[string]T) []string {
	l := make([]string, 0, len(m))
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}

// extractFields extracts the fields in a struct.
//
// Note: The AName and EName embedded structs are discarded for convenience.
func extractFields(s any) map[string]any {
	m := make(map[string]any)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	for i := range v.NumField() {
		name := t.Field(i).Name
		if name == "AName" || name == "EName" {
			continue
		}
		value := v.Field(i).Interface()
		m[name] = value
	}
	return m
}

// Package main handles the commands that nirimgr supports.
//
// Currently supports the following commands:
//
// # events
//
// The events command
//
//	nirimgr events
//
// listens on the niri event-stream, and reacts
// to specified events. You need to configure the rules in the config.json file.
// It supports the same matching on windows as the niri config window-rule. Specify
// matches and excludes containing the title or appId of the window you want to match.
// Then provide which actions you want to do on the matched window, e.g. MoveWindowToFloating
// will move the matched window to floating.
//
// # scratch
//
// The scratch command
//
//	nirimgr scratch [move|show]
//
// takes one argument, either `move` or `show`. Move will move the
// currently focused window to the scratchpad workspace. Show will take the last window
// on the scratchpad workspace, and move it to the currently focused workspace.
//
// # list
//
// The list command
//
//	nirimgr list
//
// lists all the available actions and events defined in nirimgr.
package main

import (
	"github.com/soderluk/nirimgr/cmd"
	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/internal/common"
)

func main() {
	if err := config.Configure("config.json"); err != nil {
		panic(err)
	}
	common.SetupLogger()
	cmd.Execute()
}

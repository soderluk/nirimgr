# nirimgr

A manager for niri, written in Go.

This project is my personal way of trying to learn a bit more of Go.

I've been using [Niri WM](https://github.com/YaLTeR/niri) for a while, and
saw some helper script by YaLTeR in the discussions, and thought, why not try
and port the same in Go.

I've also used i3wm before, and it had scratchpad functionality, which I used quite
a lot. I thought this would be a great learning experience for myself, and I try to
mimic the way i3wm handles the scratchpad.

# Installation

Run `go install github.com/soderluk/nirimgr@latest` to install nirimgr.

# Configuration

The configuration file for nirimgr should be put in ~/.config/nirimgr/config.json

Example configuration:

```json
{
  // The socket type to use. The Niri socket is a unix socket.
  "socketType": "unix",
  // Set the log level of the nirimgr command. Supported levels "DEBUG", "INFO", "WARN", "ERROR"
  "logLevel": "DEBUG",
  // Window rules and actions to do on the matched window.
  "rules": [
    {
      "match": [
        {
          // Match the title Bitwarden
          "title": "Bitwarden",
          // Match the app-id zen
          "appId": "zen"
        }
      ],
      "actions": {
        // Move the matching window to floating
        "MoveWindowToFloating": {},
        // Set the floating window width to a fixed 400
        "SetWindowWidth": {
          "change": {
            "SetFixed": 400
          }
        },
        // Set the floating window height to a fixed 600
        "SetWindowHeight": {
          "change": {
            "SetFixed": 600
          }
        }
      }
    },
    {
      "match": [
        {
          // Match the app-id org.gnome.Calculator
          "appId": "org.gnome.Calculator"
        }
      ],
      "actions": {
        // Move the calculator to floating
        "MoveWindowToFloating": {},
        // Move the floating window to a fixed x, y coordinate of 800, 200
        "MoveFloatingWindow": {
          "x": {
            "SetFixed": 800
          },
          "y": {
            "SetFixed": 200
          }
        },
        // Set the floating window width to a fixed 50.
        "SetWindowWidth": {
          "change": {
            "SetFixed": 50
          }
        },
        // Set the floating window height to a fixed 50.
        "SetWindowHeight": {
          "change": {
            "SetFixed": 50
          }
        }
      }
    }
  ]
}
```

The rules are the same as the `window-rule` in niri configuration. Match the window on a given title or app-id.
Then specify which action you want to do to the matched window. In the example above, the gnome calculator
is matched, then we move the calculator window to floating, move the floating window to a specified x and y coordinate,
set the window width and height to a fixed amount.

Each action needs to be a separate action. The actions are applied sequentially on the window.

The actions you can use can be found in the [niri ipc documentation](https://yalter.github.io/niri/niri_ipc/enum.Action.html)

_NOTE_: Currently only `WindowsChanged`, `WindowOpenedOrChanged` and `WindowClosed` events are watched.

# Usage

To use nirimgr, it provides two CLI-commands:

- events: The events command starts listening on the niri event-stream. `nirimgr events`
- scratch: The scratch command moves a window to the scratchpad workspace, or shows the window (moves the window
  to the currently active workspace) from the scratchpad workspace. This command should be configured
  as a keybind in niri configuration. `nirimgr scratch [move|show]`

To use the scratchpad with Niri, you need to have a named workspace `scratchpad`. Set it up like so:

```kdl
workspace "scratchpad"
spawn-at-startup "niri" "msg" "action" "focus-workspace-down"
```

The above will create a new named workspace, `scratchpad`, and focus immediately on the next
workspace.

Then if you want to use the scratchpad, configure `nirimgr scratch move` and `nirimgr scratch show`
in a keybind, like so:

```kdl
binds {
  ...
  Mod+S {
      spawn "nirimgr" "scratch" "move"
  }
  Mod+Shift+S {
      spawn "nirimgr" "scratch" "show"
  }
  ...
}
```

Press `Mod+S` to move the currently focused window to the scratchpad workspace, and `Mod+Shift+S` to
move it back.

If you have multiple windows in the scratchpad, the command will move the last window to
the current workspace.

If you want to use the events (i.e. listen to the niri event-stream and do actions based on the events)
you need to start the `nirimgr events` command on startup like so:

`spawn-at-startup "nirimgr" "events"`

This will listen on the event stream, and react to the matching windows accordingly. You need to
define the `rules` in the `config.json` and add the actions you want to do to the window when the
event happens.

# Known issues

There seems to be an [issue](https://github.com/YaLTeR/niri/issues/1805) with niri that it doesn't respect the `focus true/false` when moving a window.
And here's the [PR](https://github.com/YaLTeR/niri/pull/1820) to fix it (as of now, it hasn't been yet merged.)

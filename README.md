# nirimgr

A manager for niri, written in Go.

This is a small project mainly directed towards learning a bit more of Go.

I've been using [Niri WM](https://github.com/YaLTeR/niri) for a while, and
saw a [helper script by YaLTeR](https://github.com/YaLTeR/niri/discussions/1599) in the discussions, and thought, why not try
and port the same in Go.

I've used i3wm before, and it had scratchpad functionality, which I used quite
a lot. I thought this would be a great learning experience for myself, and I try to
mimic the way i3wm handles the scratchpad.

The scratchpad is very simple, but works for me. If you have any suggestions on how to improve it, please don't hesitate
to open an [issue](https://github.com/soderluk/nirimgr/issues/new).

---

Q: Why not create this in Rust?

A: Because I need to improve my Golang knowledge, and this was a great opportunity for me to dive a bit deeper into
the language, as well as learn a bit more on the Niri IPC and how it works.

# Installation

## Using go install

Run `go install github.com/soderluk/nirimgr@latest` to install the latest version of nirimgr.

## From GitHub Releases

Get the latest release from the [Releases page](https://github.com/soderluk/nirimgr/releases/).

Choose your arch and download the tarball, and optionally the checksums-file if you want to verify the archive.

To verify the archive, run `sha256sum --ignore-missing -c nirimgr_x.x.x_checksums.txt` (replace x.x.x with the current version of the module)

```bash
$ sha256sum --ignore-missing -c nirimgr_x.x.x_checksums.txt

nirimgr_x.x.x_linux_amd64.tar.gz: OK
```

After verification succeeds, extract the tarball somewhere, e.g. `~/Downloads/nirimgr`.

```bash
$ mkdir -p ~/Downloads/nirimgr
$ tar xf nirimgr_x.x.x_linux_amd64.tar.gz --directory=~/Downloads/nirimgr/
```

The tarball includes the following files:

- CHANGELOG.md: The changelog for this version.
- LICENSE: The license file of the repository.
- nirimgr: The executable file.
- README.md: This readme.

Move the `nirimgr` executable somewhere in your `PATH`, e.g. `/usr/local/bin`, or `~/bin`.

```bash
$ sudo cp nirimgr /usr/local/bin
```

## Build from source

Some prerequisites necessary: nirimgr uses a [justfile](https://github.com/casey/just) to make building and installing easier. Install `just` with `cargo install just`.

1. Clone the repository: `git clone https://github.com/soderluk/nirimgr.git && cd nirimgr`
2. Build the project with `just build`. This will compile the source code and create an executable `nirimgr`.
3. Install the executable

   - Copy the executable anywhere on your path: `sudo cp nirimgr /usr/local/bin`
   - If you have GOPATH set, you can run `just install` to install the executable in `$GOPATH/bin`.

## Post install

Verify the version after installation:

```bash
$ nirimgr --version
nirimgr version vx.x.x (commit: aaaaaaa, built at: 2025-07-25T06:05:16Z)
```

# Configuration

The configuration file for nirimgr should be put in ~/.config/nirimgr/config.json

Example configuration (see: [config.json](./examples/config.json)):

```json
{
  // Define the scratchpad workspace name here.
  "scratchpadWorkspace": "scratchpad",
  // Set the log level of the nirimgr command. Supported levels "DEBUG", "INFO", "WARN", "ERROR"
  "logLevel": "DEBUG",
  // Window and/or Workspace rules and actions to do on the matched window/workspace.
  "rules": [
    {
      // This rule matches workspaces.
      "type": "workspace",
      "match": [
        {
          // Matches on the workspace name "chat"
          "name": "chat"
        }
      ],
      "exclude": [
        // Add any excluded matches here.
        {
          "name": "discord"
        }
      ],
      "actions": {
        // Focus the workspace
        "FocusWorkspace": {},
        // Move the workspace to the monitor eDP-1
        // I.e. we always want our chat workspace on a specific monitor.
        "MoveWorkspaceToMonitor": {
          "output": "eDP-1"
        }
      }
    },
    {
      "type": "workspace",
      "match": [
        {
          // Match the workspace named "work", and if it's being on output "eDP-1".
          "name": "work",
          "output": "eDP-1"
        }
      ],
      "actions": {
        // Focus the workspace.
        "FocusWorkspace": {},
        // Move the workspace to the monitor on the right.
        // I.e. we always want to move our work workspace to the "main" screen, if it exists.
        "MoveWorkspaceToMonitorRight": {}
      }
    },
    {
      // This rule matches a window.
      "type": "window",
      "match": [
        {
          // Match the title Bitwarden
          "title": "Bitwarden",
          // Match the app-id zen
          "appId": "zen"
        }
      ],
      "exclude": [
        // Add any exclude matches here.
        {
          "title": "Firefox"
        }
      ],
      "actions": {
        // Move the window to floating.
        "MoveWindowToFloating": {},
        // Set the window width to a fixed 400 pixels.
        "SetWindowWidth": {
          "change": {
            "SetFixed": 400
          }
        },
        // Set the window height to a fixed 600 pixels.
        "SetWindowHeight": {
          "change": {
            "SetFixed": 600
          }
        }
      }
    },
    {
      "type": "window",
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
        // Set the floating window width to a fixed 50 pixels.
        "SetWindowWidth": {
          "change": {
            "SetFixed": 50
          }
        },
        // Set the floating window height to a fixed 50 pixels.
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

- Added in v0.2.0: Rule type - The new rule type defines if the rule should apply to a window or a workspace.

The rules are the same as the `window-rule` in Niri configuration. Currently we only match the window on a given title or app-id.
Then specify which action you want to do with the matched window. In the example above, the gnome calculator
is matched, then we move the calculator window to floating, move the floating window to a specified x and y coordinate,
set the window width and height to a fixed amount.

In addition to the window matching, we can match workspaces. The workspaces matches on a given name or output. The actions are performed
on the matched workspace.

Each action needs to be a separate action. The actions are applied sequentially on the window.

The actions you can use can be found in the [niri ipc documentation](https://yalter.github.io/niri/niri_ipc/enum.Action.html)

_NOTE_: Currently only `WindowsChanged`, `WindowOpenedOrChanged` and `WindowClosed` window events are watched. For workspaces, the
`WorkspacesChanged` event is watched.

Please feel free to open a PR if you have other thoughts what we could do with nirimgr.

# Usage

To use nirimgr, it provides the following CLI-commands:

- `nirimgr events`: The events command starts listening on the niri event-stream.
- `nirimgr scratch [move|show]`: The scratch command moves a window to the scratchpad workspace, or shows the window (moves the window
  to the currently active workspace) from the scratchpad workspace. This command should be configured
  as a key-bind in niri configuration.
- `nirimgr list [actions|events]`: The list command will list all the available actions or events, so you don't need to remember them all.

To use the scratchpad with Niri, you need to have a named workspace `scratchpad`, or if you want to configure it,
set the scratchpadWorkspace configuration option to something else `"scratchpadWorkspace": "scratch"`.

Set it up in niri config like so:

```kdl
workspace "scratchpad"  // or whatever you configured it to be.
spawn-at-startup "niri" "msg" "action" "focus-workspace-down"
```

The above will create a new named workspace, `scratchpad`, and focus immediately on the next
workspace.

Then if you want to use the scratchpad, configure `nirimgr scratch move` and `nirimgr scratch show`
in a key-bind, like so:

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
the current workspace (this is pretty much how i3wm did the scratchpad functionality).

If you want to use the events (i.e. listen to the niri event-stream and do actions based on the events)
you need to start the `nirimgr events` command on startup like so:

```kdl
spawn-at-startup "nirimgr" "events"
```

This will listen on the event stream, and react to the matching windows accordingly. You need to
define the `rules` in `config.json` and add the actions you want to do to the window when the
event happens.

# Justfile

The following actions can be performed with the `just` command:

- build: Builds the source into an executable
- coverage: Opens up the code coverage in the browser
- fmt: Runs `go fmt` on the project.
- help: Lists the available recipes. This is the default, if you run `just` without arguments.
- install: Installs the executable in your GOPATH/bin.
- run RUNARGS: Runs the module with RUNARGS. If the RUNARGS contains a space, you need to quote them, e.g. `just run "list actions"`
- test: Runs `go test` on the project.
- version: Prints the version.
- vet: Runs `go vet` on the project.

# Acknowledgements

Of course the biggest one goes to [Niri WM](https://github.com/YaLTeR/niri) and YaLTeR for an awesome manager!

Since this is mostly a learning project for me, I had to look a bit more into a few of existing libraries,
the most notable being [niri-float-sticky](https://github.com/probeldev/niri-float-sticky)

The goroutine handling of the event stream felt like a better approach than I had before, so thanks to the author for a great library!

# Known issues

There seems to be an [issue](https://github.com/YaLTeR/niri/issues/1805) with niri that it doesn't respect the `focus true/false` when moving a window.
And here's the [PR](https://github.com/YaLTeR/niri/pull/1820) to fix it (as of now, it hasn't been yet merged.)

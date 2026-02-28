Simple X11 status bar inspired by slstatus.
Linux only.

Edit `config/config.go`, then run `make` to build.
Available TAGS:
  with_pulse   Links with libpulse, required for the volume component.

You can use real-time signals to force a redraw of specific component(s). For example:
`kill -35 $(pidof sb)`
This will trigger any component where .Signal == 35.

Note: instantaneous redraw is not guaranteed -- it may take up to the
`RedrawDelay` value (see `config/config.go`) for a redraw to occur.

The `-s` flag can be used to write the output to stdout instead of WM_NAME.

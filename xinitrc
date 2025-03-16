[ -f "$HOME/.Xresources" ] && xrdb -merge "$HOME/.Xresources"
export XDG_SESSION_TYPE="x11"

exec dwm > ~/.cache/dwm_out.log 2> ~/.cache/dwm_err.log

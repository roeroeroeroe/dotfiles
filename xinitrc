[ -f "$HOME/.Xresources" ] && xrdb -merge "$HOME/.Xresources"
export XDG_SESSION_TYPE="x11"

exec dwm

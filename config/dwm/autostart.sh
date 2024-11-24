picom --config ~/.config/picom/picom_dwm.conf &
feh --bg-fill --no-fehbg ~/Pictures/Wallpapers/gruvbox/landscape.jpg &
xrandr --output VGA-0 --off --output LVDS-0 --mode 1366x768 --pos 0x312 --rotate normal --output HDMI-0 --primary --mode 1920x1080 --rate 75.00 --pos 1366x0 --rotate normal &
xsetroot -cursor_name left_ptr &
xset -dpms s off &
xinput --set-prop 'pointer:Razer Razer DeathAdder Essential' 'libinput Accel Profile Enabled' 0, 1 &
setxkbmap -option caps:escape_shifted_capslock &
xset r rate 250 30 &
xset m 0 0 &
slstatus &
dunst &

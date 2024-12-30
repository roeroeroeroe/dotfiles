xrandr --output VGA-0 --off \
	--output LVDS-0 --mode 1366x768 --pos 0x312 --rotate normal \
	--output HDMI-0 --mode 1920x1080 --pos 1366x0 --rate 75.00 --rotate normal --primary
picom -b --config ~/.config/picom/picom_dwm.conf &
feh --bg-fill --no-fehbg ~/Pictures/Wallpapers/gruvbox/1D2021.png &
xsetroot -cursor_name left_ptr &
xset -dpms s off &
xinput --set-prop 'pointer:Razer Razer DeathAdder Essential' 'libinput Accel Profile Enabled' 0, 1 &
setxkbmap -option caps:escape_shifted_capslock &
xset r rate 250 30 &
slstatus &
dunst &

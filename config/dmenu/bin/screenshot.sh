#!/usr/bin/env bash

source "$HOME/.config/dmenu/theme"

option_1="desktop"
option_2="area"
option_3="window"
option_4="capture in 5s"
option_5="capture in 10s"

dmenu_cmd() {
	dmenu -l 5 "${dmenu_options[@]}" 2>/dev/null
}

run_dmenu() {
	echo -e "$option_1\n$option_2\n$option_3\n$option_4\n$option_5" | dmenu_cmd
}

time=$(date +%Y-%m-%d-%H-%M-%S)
dir="$(xdg-user-dir PICTURES)/Screenshots"
file="${time}.png"

if [[ ! -d "$dir" ]]; then
	mkdir -p "$dir"
fi

notify_view() {
	notify_cmd_shot='dunstify -u low --replace=699'
	${notify_cmd_shot} "Copied to clipboard."
	viewnior "${dir}/$file"
	if [[ -e "${dir}/$file" ]]; then
		${notify_cmd_shot} "Screenshot Saved."
	else
		${notify_cmd_shot} "Screenshot Deleted."
	fi
}

copy_shot() {
	tee "$file" | xclip -selection clipboard -t image/png
}

countdown() {
	for sec in $(seq $1 -1 1); do
		dunstify -t 1000 --replace=699 "Taking shot in : $sec"
		sleep 1
	done
}

shotnow() {
	cd "${dir}" && sleep 0.5 && maim -u -f png | copy_shot
	notify_view
}

shot5() {
	countdown '5'
	sleep 1 && cd "${dir}" && maim -u -f png | copy_shot
	notify_view
}

shot10() {
	countdown '10'
	sleep 1 && cd "${dir}" && maim -u -f png | copy_shot
	notify_view
}

shotwin() {
	cd "${dir}" && maim -u -f png -i "$(xdotool getactivewindow)" | copy_shot
	notify_view
}

shotarea() {
	cd "${dir}" && maim -u -f png -s -b 2 -c 0.35,0.55,0.85,0.25 -l | copy_shot
	notify_view
}

run_cmd() {
	case "$1" in
		'--opt1')
			shotnow
			;;
		'--opt2')
			shotarea
			;;
		'--opt3')
			shotwin
			;;
		'--opt4')
			shot5
			;;
		'--opt5')
			shot10
			;;
	esac
}

chosen="$(run_dmenu)"
case ${chosen} in
	$option_1)
		run_cmd --opt1
		;;
	$option_2)
		run_cmd --opt2
		;;
	$option_3)
		run_cmd --opt3
		;;
	$option_4)
		run_cmd --opt4
		;;
	$option_5)
		run_cmd --opt5
		;;
esac

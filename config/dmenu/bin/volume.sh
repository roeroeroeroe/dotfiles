#!/usr/bin/env bash

source "$HOME/.config/dmenu/theme"

mixer="$(amixer info Master | grep 'Mixer name' | cut -d':' -f2 | tr -d ',' )"
speaker="$(amixer get Master | tail -n1 | awk -F ' ' '{print $5}' | tr -d '[]')"
mic="$(amixer get Capture | tail -n1 | awk -F ' ' '{print $5}' | tr -d '[]')"

active=""
urgent=""

amixer get Master | grep '\[on\]' &>/dev/null
if [[ "$?" == 0 ]]; then
	active="-a 1"
	stext='deafen'
	sicon=''
else
	urgent="-u 1"
	stext='undeafen'
	sicon=''
fi

amixer get Capture | grep '\[on\]' &>/dev/null
if [[ "$?" == 0 ]]; then
	[ -n "$active" ] && active+=",2" || active="-a 2"
	mtext='mute'
	micon='󰍬'
else
	[ -n "$urgent" ] && urgent+=",2" || urgent="-u 2"
	mtext='unmute'
	micon='󰍭'
fi

mesg=" $speaker, 󰍬 $mic"

option_1="settings"
option_2="$sicon $stext"
option_3="$micon $mtext"

dmenu_cmd() {
	dmenu -l 3 "${dmenu_options[@]}" 2>/dev/null
}

run_dmenu() {
	echo -e "$option_1\n$option_2\n$option_3" | dmenu_cmd
}

run_cmd() {
	case "$1" in
		'--opt1')
			pavucontrol
			;;
		'--opt2')
			amixer set Master toggle
			;;
		'--opt3')
			amixer set Capture toggle
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
esac

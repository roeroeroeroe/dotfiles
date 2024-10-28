#!/usr/bin/env bash

source "$HOME/.config/dmenu/theme"

option_1="lock"
option_2="logout"
option_3="reboot"
option_4="shutdown"

dmenu_cmd() {
	dmenu -l 4 "${dmenu_options[@]}" 2>/dev/null
}

run_dmenu() {
	echo -e "$option_1\n$option_2\n$option_3\n$option_4" | dmenu_cmd
}

confirm_cmd() {
	echo -e "yes\nno" | dmenu "${dmenu_options[@]}" -l 2 2>/dev/null
}

confirm_exit() {
	confirm_cmd
}

confirm_run () {
	selected="$(confirm_exit)"
	if [[ "$selected" == "yes" ]]; then
		${1} && ${2} && ${3}
	else
		exit
	fi	
}

run_cmd() {
	if [[ "$1" == '--opt1' ]]; then
		slock
	elif [[ "$1" == '--opt2' ]]; then
		confirm_run 'kill -9 -1'
	elif [[ "$1" == '--opt3' ]]; then
		confirm_run 'systemctl reboot'
	elif [[ "$1" == '--opt4' ]]; then
		confirm_run 'systemctl poweroff'
	fi
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
esac

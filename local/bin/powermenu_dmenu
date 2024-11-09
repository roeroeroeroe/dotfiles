#!/usr/bin/env bash

ordered=("lock" "logout" "reboot" "shutdown")
declare -A commands=(
	["lock"]="slock"
	["logout"]="confirm_run 'kill -9 -1'"
	["reboot"]="confirm_run 'systemctl reboot'"
	["shutdown"]="confirm_run 'systemctl poweroff'"
)

confirm_run () {
	selected="$(echo -e "yes\nno" | dmenu -l 2 2>/dev/null)"
	if [[ "$selected" == "yes" ]]; then
		${1} && ${2} && ${3}
	else
		exit
	fi
}

choice=$(printf "%s\n" "${ordered[@]}" | dmenu -l 4 2>/dev/null)
[[ -n "${commands[$choice]}" ]] && eval "${commands[$choice]}"

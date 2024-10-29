#!/usr/bin/env bash

ordered=("desktop" "area" "window" "capture in 5s")
declare -A commands=(
	["desktop"]="shotnow"
	["area"]="shotarea"
	["window"]="shotwin"
	["capture in 5s"]="shot5"
)

dir="$(xdg-user-dir PICTURES)/Screenshots"
file="$(date +%Y-%m-%d-%H-%M-%S).png"
[[ ! -d "$dir" ]] && mkdir -p "$dir"

copy_shot() {
	tee "$file" | xclip -selection clipboard -t image/png
}

shotnow() {
	sleep 0.5; cd "${dir}" && maim -u -f png | copy_shot
}

shot5() {
	for i in {5..1}; do
		dunstify -r 5009 -t 1025 "Taking shot in "$i"s"
		sleep 1
	done
	sleep 1.25; cd "${dir}" && maim -u -f png | copy_shot
}

shotwin() {
	sleep 0.5; cd "${dir}" && maim -u -f png -i "$(xdotool getactivewindow)" | copy_shot
}

shotarea() {
	sleep 0.5; cd "${dir}" && maim -u -f png -s -b 2 -c 0.35,0.55,0.85,0.25 -l | copy_shot
}

choice=$(printf "%s\n" "${ordered[@]}" | dmenu -l 4 2>/dev/null)
if [[ -n "${commands[$choice]}" ]]; then
	eval "${commands[$choice]}"
	dunstify -r 5009 -t 2000 "Copied to clipboard."
	viewnior "${dir}/$file"
fi

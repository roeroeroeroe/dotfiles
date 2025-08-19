ansi() {
	for c in {1..${1:-255}}; do echo -n "\033[38;5;${c}m$c "; done
}

newsh() {
	local file=${1:-"$(LC_ALL=C tr -dc A-Za-z < /dev/urandom | head -c 5).sh"}
	[[ "$file" == */* ]] && { echo "invalid filename"; return; }
	[ -f "$file" ] && { echo "file already exists"; return; }
	echo "#!/usr/bin/env bash\n\n" > "$file"; chmod +x "$file"; ${EDITOR:-vim} "$file"
}

ipapi() {
	local stdin=""; [ -p /dev/stdin ] && stdin=$(</dev/stdin)
	local input="${*}${stdin:+${*:+ }${stdin}}"
	[ -z "$input" ] && return
	echo "$input" | tr " " "\n" | parallel -j 4 'curl -s "ip-api.com/json/{}" | jq'
}

hex() {
	[ -z "$1" ] && return
	local color="${1#\#}"
	! [[ "$color" =~ ^[A-Fa-f0-9]{6}$ ]] && { echo "invalid hex code"; return; }
	magick -size 600x600 xc:#"$color" -strip "$color".png
}

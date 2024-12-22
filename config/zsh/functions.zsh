ansi() {
	for COLOR in {1..255}; do
		echo -en "\033[38;5;${COLOR}m"
		echo -n "${COLOR} "
	done
}

streamlink() {
	for arg in "$@"; do
		command streamlink \
			--twitch-low-latency \
			--hls-live-edge 1 \
			--stream-segment-threads 10 \
			--stream-timeout 20 \
			--player mpv \
			"twitch.tv/$arg" best &
	done
}

yt() {
	pc yt-dlp -f 'bestvideo[height=1080][fps=60]+bestaudio/bestvideo[height<=1440][fps<=30]+bestaudio/best' "$@"
}

http() {
	local addr=$(ip -4 -o a s enp3s0 | cut -d ' ' -f7 | cut -d '/' -f1)
	python -m http.server --bind "$addr" 8000
}

new() {
	local FILE=${1:-"test.sh"}
	[[ "$FILE" == */* ]] && { echo "invalid filename"; return; }
	[ -f "$FILE" ] && { echo "file already exists"; return; }
	echo '#!/usr/bin/env bash\n\n' > "$FILE"; chmod +x "$FILE"; ${EDITOR:-vim} "$FILE"
}

ipapi() {
	[ -z "$1" ] && return
	echo "$@" | tr " " "\n" | parallel -j 4 'curl -s "ip-api.com/json/{}" | jq'
}

hex() {
	[ -z "$1" ] && return
	local color="${1#\#}"
	! [[ "$color" =~ ^[A-Fa-f0-9]{6}$ ]] && { echo "invalid hex code"; return; }
	magick -size 600x600 xc:#"$color" "$color".png
}

cropmon() {
	[ -z "$1" ] && return
	[ ! -f "$1" ] && { echo "file does not exist"; return; }
	identify "$1" &> /dev/null || { echo "not an image"; return; }
	local ext="${1##*.}"
	local new_img="$(head /dev/urandom | LC_ALL=C tr -dc A-Za-z0-9 | head -c 5)"
	local size="1920x1080"
	local offset="+1366+0"
	magick "$1" -crop "$size$offset" +repage "$new_img.$ext"
}

undocc() {
	local cmd="shred -zu"
	[ "$1" = "--dry-run" ] && cmd="echo"
	find "$XDG_STATE_HOME/nvim/undo" -type f -exec sh -c '
		cmd="$1"; shift
		for f; do
			target=$(basename "$f" | sed "s/%/\\//g")
			[ ! -e "$target" ] && $cmd "$f"
		done
	' sh "$cmd" {} +
}

lower() {
	local stdin=""
	[ -p /dev/stdin ] && stdin=$(</dev/stdin)
	input="${*}${stdin:+${*:+ }${stdin}}"
	[ -z "$input" ] && return
	echo "$input" | tr "[:upper:]" "[:lower:]" | tr " " "_"
}

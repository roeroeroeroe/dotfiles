ansi() {
	for c in {1..255}; do echo -n "\033[38;5;${c}m$c "; done
}

sl() {
	for arg in "$@"; do
		streamlink \
			--twitch-low-latency \
			--hls-live-edge 1 \
			--stream-segment-threads 10 \
			--stream-timeout 20 \
			--player mpv \
			"twitch.tv/$arg" best &
	done
}

slrec() {
	[ -z "$1" ] && return 1
	streamlink \
		--output "${1}_{time:%Y%m%d%H%M%S}.ts" \
		"twitch.tv/$1" best
}

yt() {
	proxychains yt-dlp -f 'bestvideo[height=1080][fps>=48]+bestaudio/bestvideo[height<=1440][fps<=30]+bestaudio/best' "$@"
}

http() {
	local addr=$(ip -4 -o a s enp3s0 | cut -d ' ' -f7 | cut -d '/' -f1)
	python -m http.server --bind "$addr" 8000
}

new() {
	local file=${1:-"$(LC_ALL=C tr -dc A-Za-z < /dev/urandom | head -c 5).sh"}
	[[ "$file" == */* ]] && { echo "invalid filename"; return; }
	[ -f "$file" ] && { echo "file already exists"; return; }
	echo '#!/usr/bin/env bash\n\n' > "$file"; chmod +x "$file"; ${EDITOR:-vim} "$file"
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
	magick -size 600x600 xc:#"$color" "$color".png
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

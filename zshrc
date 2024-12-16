# path
export PATH=$HOME/.local/bin:/usr/local/bin:$PATH
ZSH_PATH="$HOME/.config/zsh"

# autoload
autoload -U vcs_info select-word-style compinit; compinit -d "$ZSH_PATH/.zcompdump-$HOST"

# lscolors
. <(dircolors -b 2>/dev/null || :)

# prompt
NEWLINE=$'\n'
setopt prompt_subst
zstyle ":vcs_info:*" check-for-changes true
zstyle ":vcs_info:*" unstagedstr "%F{yellow}*%f"
zstyle ":vcs_info:*" stagedstr "%F{yellow}+%f"
zstyle ":vcs_info:git:*" formats "%b%u%c" # branch, unstaged, staged
precmd() {
	vcs_info
	local git_branch=""
	[ -n "$vcs_info_msg_0_" ] && git_branch="%F{white} [${vcs_info_msg_0_}%f%F{white}]%f"
	PS1="${NEWLINE}%F{white}%~%f${git_branch}%f%F{white} [%n@%M] [%T]%(1j. [%j].)%f%F{yellow}${NEWLINE}>%f "
	RPROMPT="%(?..%F{white}%?%f) "
}
PS2="%F{white}>%f "
PS3="%F{white}?>%f "

# binds
select-word-style bash
bindkey -e
bindkey "^[[1;5D" backward-word
bindkey "^[[1;5C" forward-word
bindkey "^W" backward-kill-word
bindkey "^[w" kill-whole-line
bindkey "^[[P" delete-char
bindkey -s "^[l" "ls^M"

# completion
zstyle ":completion:*" list-colors ${(s.:.)LS_COLORS}
zstyle ":completion:*" menu select
zstyle ":completion:*" matcher-list "m:{a-zA-Z}={A-Za-z}" "r:|=*" "l:|=* r:|=*"

# history
HISTFILE="$HOME/.zsh_history"
HISTSIZE=10000
SAVEHIST=10000
setopt extended_history
setopt hist_expire_dups_first
setopt hist_ignore_dups
setopt hist_ignore_space
setopt hist_verify
setopt share_history
setopt inc_append_history

# aliases
alias ls="ls --color=auto"
alias grep="grep --color=auto"
alias la="ls -A"
alias lh="ls -lh"
alias lah="la -lh"
alias mv="mv -iv"
alias cp="cp -iv"
alias mkdir="mkdir -vp"
alias less="less -igmj .5"
alias copy="xclip -se c"
alias ru="setxkbmap 'ru'"
alias у="setxkbmap 'us'"
alias yay="PKGEXT=.pkg.tar yay" # skip compression
alias whois="whois -H"
alias history="history 1"
alias diff="diff --color=auto -u"
alias shred="shred -vzu"
alias pc="proxychains"
alias temp="awk '{print \$1/1000 \"°C\"}' /sys/class/thermal/thermal_zone0/temp"
alias ttvc="ttvu -pp -query 'channel{name chatters{count broadcasters{login}moderators{login}vips{login}viewers{login}}}'"
alias mirrors="sudo reflector --country DE,GB,US --protocol https --sort rate --latest 20 --save /etc/pacman.d/mirrorlist"

# func
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

ech() {
	[ -z "$1" ] && return
	curl -s "https://dns.google/resolve?name=$1&type=HTTPS" | jq -r ".Answer[0].data" | grep -q "ech=" && echo true || echo false
}

new() {
	[ -z "$1" ] && local FILE="test.sh" || local FILE="$1"
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
	identify "$1" >/dev/null 2>&1 || { echo "not an image"; return; }
	local ext="${1##*.}"
	local new_img="$(head /dev/urandom | LC_ALL=C tr -dc A-Za-z0-9 | head -c 5)"
	local size="1920x1080"
	local offset="+1366+0"
	magick "$1" -crop "$size$offset" +repage "$new_img.$ext"
}

# plugins
. "$ZSH_PATH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
. "$ZSH_PATH/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
ZSH_HIGHLIGHT_HIGHLIGHTERS+=(brackets)

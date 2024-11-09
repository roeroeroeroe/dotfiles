# path
export PATH=$HOME/.local/bin:/usr/local/bin:$PATH
ZSH_PATH="$HOME/.config/zsh"

# autoload
autoload -U vcs_info select-word-style compinit; compinit -d "$ZSH_PATH/.zcompdump-$HOST"

# lscolors
[ -x /usr/bin/dircolors ] && test -r ~/.dircolors && source <(dircolors -b ~/.dircolors) || source <(dircolors -b)

# prompt
NEWLINE=$'\n'
setopt prompt_subst
zstyle ":vcs_info:*" check-for-changes true
zstyle ":vcs_info:*" unstagedstr "%F{green}*%f"
zstyle ":vcs_info:*" stagedstr "%F{green}+%f"
zstyle ":vcs_info:git:*" formats "%b%u%c"
precmd() {
	vcs_info
	local git_branch=""
	[[ -n ${vcs_info_msg_0_} ]] && git_branch="%F{yellow} [git::${vcs_info_msg_0_}%F{yellow}]%f"
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

# plugins
source "$ZSH_PATH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
source "$ZSH_PATH/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
typeset -g ZSH_HIGHLIGHT_HIGHLIGHTERS=(main brackets)

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

# env
export EDITOR="nvim"
export CHATTERINO2_RECENT_MESSAGES_URL="https://recent-messages.zneix.eu/api/v2/recent-messages/%1"

# aliases
alias ls="ls --color=auto"
alias grep="grep --color=auto"
alias la="ls -A"
alias lh="ls -lh"
alias lah="la -lh"
alias mv="mv -iv"
alias cp="cp -iv"
alias df="df -h"
alias free="free -h"
alias mkdir="mkdir -vp"
alias copy="xclip -selection clipboard"
alias fzf="fzf --prompt '$(pwd) '"
alias nfzf='selected=$(fzf) && [ -n "$selected" ] && nvim "$selected"'
alias ru="setxkbmap 'ru'"
alias у="setxkbmap 'us'"
alias yay="PKGEXT=.pkg.tar yay" # skip compression
alias whois="whois -H"
alias history="history 1"
alias diff="diff --color=auto -u"
alias shred="shred -vzu"
alias http="python -m http.server --bind "$(ip -4 -o a s enp3s0 | cut -d ' ' -f7 | cut -d '/' -f1)" 8000"
alias pc="proxychains"
alias temp="awk '{print \$1/1000 \"°C\"}' /sys/class/thermal/thermal_zone0/temp"

# func
ansi() {
	for COLOR in {1..255}; do
		echo -en "\033[38;5;${COLOR}m"
		echo -n "${COLOR} "
	done
}

streamlink() {
	for arg in $*; do
		command streamlink twitch.tv/$arg best &
	done
}

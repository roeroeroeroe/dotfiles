# path
export PATH=$HOME/.local/bin:/usr/local/bin:$PATH

# history
HISTFILE=~/.zsh_history
HISTSIZE=10000
SAVEHIST=10000

# oh-my-zsh
export ZSH="$HOME/.oh-my-zsh"
ZSH_THEME='frisk'
zstyle ':omz:update' mode disabled
DISABLE_UNTRACKED_FILES_DIRTY='true'
plugins=(
	zsh-autosuggestions
	zsh-syntax-highlighting
	timer
)
source $ZSH/oh-my-zsh.sh

# env
export EDITOR='nvim'
export CHATTERINO2_RECENT_MESSAGES_URL='https://recent-messages.zneix.eu/api/v2/recent-messages/%1'

# appimages
alias chatterino="setsid ~/Appimages/Chatterino-x86_64.AppImage"
alias nekoray="setsid ~/Appimages/nekoray-3.26-2023-12-09-linux-x64.AppImage"

# aliases
alias ls='ls --color=auto'
alias la='ls -A'
alias lh='ls -lh'
alias copy='xclip -selection clipboard'
alias fzf='fzf --prompt "$(pwd) "'
alias nfzf='selected_file=$(fzf) && [ -n "$selected_file" ] && nvim "$selected_file"'
alias ru='setxkbmap "ru"'
alias у='setxkbmap "us"'
alias yay='PKGEXT=.pkg.tar yay' # skip compression
alias fixdev='xset -dpms s off; xset r rate 250 30; xset m 0 0; xinput --set-prop "pointer:Razer Razer DeathAdder Essential" "libinput Accel Profile Enabled" 0, 1'

# func
weather() {
	curl "wttr.in/$1?3&F"
}

# misc
setopt histignorespace

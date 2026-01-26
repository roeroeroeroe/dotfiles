export PATH=$HOME/.local/bin:$GOPATH/bin:$PATH
ZSH_CONF_DIR="$HOME/.config/zsh"

autoload -U vcs_info select-word-style edit-command-line compinit
zle -N edit-command-line
compinit -d "$ZSH_CONF_DIR/zcompdump"

. <(dircolors -b 2> /dev/null)

if [ -d "$ZSH_CONF_DIR" ]; then
	for f in "$ZSH_CONF_DIR"/*.zsh(N); do [ -r "$f" ] && . "$f"; done
	unset f
fi

. "$ZSH_CONF_DIR/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
. "$ZSH_CONF_DIR/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
ZSH_HIGHLIGHT_HIGHLIGHTERS+=(brackets)

unset ZSH_CONF_DIR

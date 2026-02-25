export PATH=$HOME/.local/bin:$GOPATH/bin:$PATH
ZSH_CONF_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/zsh"

autoload -U vcs_info select-word-style edit-command-line compinit
zle -N edit-command-line
compinit -d "${XDG_CACHE_HOME:-$HOME/.cache}/zcompdump"

. <(dircolors -b 2> /dev/null)

. "$ZSH_CONF_DIR/aliases.zsh"
. "$ZSH_CONF_DIR/binds.zsh"
. "$ZSH_CONF_DIR/completion_opts.zsh"
. "$ZSH_CONF_DIR/functions.zsh"
. "$ZSH_CONF_DIR/history_opts.zsh"
. "$ZSH_CONF_DIR/prompt.zsh"

. "$ZSH_CONF_DIR/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
. "$ZSH_CONF_DIR/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
ZSH_HIGHLIGHT_HIGHLIGHTERS+=(brackets)

unset ZSH_CONF_DIR

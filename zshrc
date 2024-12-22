export PATH=$HOME/.local/bin:/usr/local/bin:$PATH
ZSH_PATH="$HOME/.config/zsh"

autoload -U vcs_info select-word-style compinit; compinit -d "$ZSH_PATH/zcompdump"

. <(dircolors -b 2>/dev/null || :)

. "$ZSH_PATH/prompt.zsh"

. "$ZSH_PATH/binds.zsh"

. "$ZSH_PATH/completion_opts.zsh"

. "$ZSH_PATH/history_opts.zsh"

. "$ZSH_PATH/aliases.zsh"

. "$ZSH_PATH/functions.zsh"

. "$ZSH_PATH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
. "$ZSH_PATH/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
ZSH_HIGHLIGHT_HIGHLIGHTERS+=(brackets)

export PATH=$HOME/.local/bin:$PATH
ZSH_PATH="$HOME/.config/zsh"

autoload -U vcs_info select-word-style compinit; compinit -d "$ZSH_PATH/zcompdump"

. <(dircolors -b 2>/dev/null || :)

if [ -d "$ZSH_PATH" ]; then
	for f in "$ZSH_PATH"/*.zsh(N); do [ -r "$f" ] && . "$f"; done
	unset f
fi

. "$ZSH_PATH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh"
. "$ZSH_PATH/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh"
ZSH_HIGHLIGHT_HIGHLIGHTERS+=(brackets)

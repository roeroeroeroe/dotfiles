export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_STATE_HOME="$HOME/.local/state"
export XDG_CACHE_HOME="$HOME/.cache"

export CUDA_CACHE_PATH="$XDG_CACHE_HOME"/nv
export __GL_SHADER_DISK_CACHE_PATH="$XDG_CACHE_HOME"/nv
export GOPATH="$XDG_DATA_HOME"/go
export GTK2_RC_FILES="$XDG_CONFIG_HOME"/gtk-2.0/gtkrc
# export RXVT_SOCKET="$XDG_RUNTIME_DIR"/urxvtd
export PARALLEL_HOME="$XDG_CONFIG_HOME"/parallel

export TERMINAL="st"
export EDITOR="nvim"
export PAGER="less -igmj .5"

export NODE_REPL_HISTORY=""

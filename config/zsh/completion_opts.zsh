zstyle ":completion:*" list-colors ${(s.:.)LS_COLORS}
zstyle ":completion:*" menu select
zstyle ":completion:*" matcher-list "m:{a-zA-Z}={A-Za-z}" "r:|=*" "l:|=* r:|=*"

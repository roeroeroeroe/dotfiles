select-word-style bash
bindkey -e
bindkey "^[[1;5D" backward-word
bindkey "^[[1;5C" forward-word
bindkey "^W" backward-kill-word
bindkey "^[w" kill-whole-line
bindkey "^[[P" delete-char
bindkey "^X^E" edit-command-line
bindkey -s "^[l" "ls\n"
bindkey -s "^[e" "nvimf\n"

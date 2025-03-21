_newline=$'\n'
_rainbow() {
	str="$*"
	[ $(tput colors) -lt 256 ] && { print "$str"; return; }

	len=${#str}
	start=$(( (RANDOM % (256 - len)) + 1 ))

	for (( i = 0; i < len; i++ )); do
		color=$(( ((start + i - 1) % 255) + 1 ))
		print -n "%F{$color}${str[i+1]}"
	done

	print "%f"
}
_user=$(_rainbow $USER)

_primary_color="white"
_accent_color="magenta"

setopt prompt_subst
zstyle ":vcs_info:*" check-for-changes true
zstyle ":vcs_info:*" unstagedstr "%F{$_accent_color}*%f"
zstyle ":vcs_info:*" stagedstr "%F{$_accent_color}+%f"
zstyle ":vcs_info:git:*" formats "%b%u%c" # branch, unstaged, staged
precmd() {
	vcs_info
	local git_branch=""
	[ -n "$vcs_info_msg_0_" ] && git_branch="%F{$_primary_color} [$vcs_info_msg_0_%f%F{$_primary_color}]%f"
	PS1="$_newline%F{$_primary_color}%~%f$git_branch%f%F{$_primary_color} [$_user%F{$_primary_color}@%M] [%T]%(1j. [%j].)%f%F{$_accent_color}$_newline>%f "
	RPROMPT="%(?..%F{$_primary_color}%?%f) "
}
PS2="%F{$_primary_color}>%f "
PS3="%F{$_primary_color}?>%f "

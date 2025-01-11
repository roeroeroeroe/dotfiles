_yellow=$([[ $TERM == *256color* ]] && echo "208" || echo "yellow")
_newline=$'\n'
_rainbow() {
	str="$*"
	[[ "$TERM" != *256color* ]] && { print "$str"; return; }

	len=${#str}
	start=$(( (RANDOM % (256 - len)) + 1 ))

	for (( i = 0; i < len; i++ )); do
		color=$(( ((start + i - 1) % 255) + 1 ))
		print -n "%F{$color}${str[i+1]}"
	done

	print "%f"
}
_user=$(_rainbow $USER)

setopt prompt_subst
zstyle ":vcs_info:*" check-for-changes true
zstyle ":vcs_info:*" unstagedstr "%F{$_yellow}*%f"
zstyle ":vcs_info:*" stagedstr "%F{$_yellow}+%f"
zstyle ":vcs_info:git:*" formats "%b%u%c" # branch, unstaged, staged
precmd() {
	vcs_info
	local git_branch=""
	[ -n "$vcs_info_msg_0_" ] && git_branch="%F{white} [$vcs_info_msg_0_%f%F{white}]%f"
	PS1="$_newline%F{white}%~%f$git_branch%f%F{white} [$_user%F{white}@%M] [%T]%(1j. [%j].)%f%F{$_yellow}$_newline>%f "
	RPROMPT="%(?..%F{white}%?%f) "
}
PS2="%F{white}>%f "
PS3="%F{white}?>%f "

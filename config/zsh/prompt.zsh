local yellow=$([[ $TERM == *256color* ]] && echo "208" || echo "yellow")
local newline=$'\n'
setopt prompt_subst
zstyle ":vcs_info:*" check-for-changes true
zstyle ":vcs_info:*" unstagedstr "%F{$yellow}*%f"
zstyle ":vcs_info:*" stagedstr "%F{$yellow}+%f"
zstyle ":vcs_info:git:*" formats "%b%u%c" # branch, unstaged, staged
precmd() {
	vcs_info
	local git_branch=""
	[ -n "$vcs_info_msg_0_" ] && git_branch="%F{white} [$vcs_info_msg_0_%f%F{white}]%f"
	PS1="$newline%F{white}%~%f$git_branch%f%F{white} [%n@%M] [%T]%(1j. [%j].)%f%F{$yellow}$newline>%f "
	RPROMPT="%(?..%F{white}%?%f) "
}
PS2="%F{white}>%f "
PS3="%F{white}?>%f "

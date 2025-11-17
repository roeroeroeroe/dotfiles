_rainbow_saturation=0.2 # 0..1
_rainbow_lightness=0.45 # 0..1
_rainbow_step=15

_colorize_user=2 # 0=false, 1=accent, 2=rainbow
_primary_color="white"
_accent_color=15 # bright white

_lf=$'\n'

_hsl_to_hex() {
	local H=$1 S=$2 L=$3 temp C Hp X R1 G1 B1 m R G B

	temp=$(( 2 * L - 1 ))
	C=$(( (1 - (temp < 0 ? -temp : temp)) * S ))

	Hp=$(( H / 60.0 ))
	temp=$(( Hp % 2 - 1 ))
	X=$(( C * (1 - (temp < 0 ? -temp : temp)) ))

	if (( Hp < 1 )); then
		R1=$C; G1=$X; B1=0
	elif (( Hp < 2 )); then
		R1=$X; G1=$C; B1=0
	elif (( Hp < 3 )); then
		R1=0; G1=$C; B1=$X
	elif (( Hp < 4 )); then
		R1=0; G1=$X; B1=$C
	elif (( Hp < 5 )); then
		R1=$X; G1=0; B1=$C
	else
		R1=$C; G1=0; B1=$X
	fi

	m=$(( L - C / 2.0 ))

	R=$(( (R1 + m) * 255 ))
	G=$(( (G1 + m) * 255 ))
	B=$(( (B1 + m) * 255 ))

	printf "%02x%02x%02x" $R $G $B
}

_rainbow() {
	str="$*"
	[ $(tput colors) -lt 256 ] && { print "$str"; return; }

	start=$(( $(date +%N) % 360 ))
	len=${#str}
	for (( i = 0; i < len; i++ )); do
		hue=$(( (start + i * _rainbow_step) % 360 ))
		hex=$(_hsl_to_hex $hue $_rainbow_saturation $_rainbow_lightness)
		print -n "%F{#$hex}${str[i+1]}"
	done
	print "%f"
}

case $_colorize_user in
	2) _user=$(_rainbow $USER) ;;
	1) _user="%F{$_accent_color}$USER%f" ;;
	*) _user=$USER ;;
esac

setopt prompt_subst
zstyle ":vcs_info:*" check-for-changes true
zstyle ":vcs_info:*" unstagedstr "%F{$_accent_color}*%f"
zstyle ":vcs_info:*" stagedstr "%F{$_accent_color}+%f"
zstyle ":vcs_info:git:*" formats "%b%u%c"
precmd() {
	vcs_info
	local git_branch=""
	[ -n "$vcs_info_msg_0_" ] && git_branch="%F{$_primary_color} [$vcs_info_msg_0_%f%F{$_primary_color}]%f"
	PS1="$_lf%F{$_primary_color}%~%f$git_branch%F{$_primary_color} [$_user%F{$_primary_color}@%M] [%T]%(1j. [%j].)%f%F{$_accent_color}$_lf>%f "
	RPROMPT="%(?..%F{$_primary_color}%?%f) "
}
PS2="%F{$_primary_color}>%f "
PS3="%F{$_primary_color}?>%f "

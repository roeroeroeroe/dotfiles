static int topbar = 1;
static const char *fonts[] = {
	"JetBrainsMono Nerd Font Mono:style=SemiBold:pixelsize=16"
};
static const char *prompt = NULL;
static const char *colors[SchemeLast][2] = {
	/*               fg         bg       */
	[SchemeNorm] = { "#bbbbbb", "#1e1e2e" },
	[SchemeSel]  = { "#eeeeee", "#b4befe" },
	[SchemeOut]  = { "#000000", "#00ffff" },
};
/* -l option; if nonzero, dmenu uses vertical list with given number of lines */
static unsigned int lines = 0;

/*
 * Characters not considered part of a word while deleting words
 * for example: " /?\"&[]"
 */
static const char worddelimiters[] = " ";

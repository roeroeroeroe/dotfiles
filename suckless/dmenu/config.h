static int topbar = 1;
static const char *fonts[] = {
	"JetBrainsMono Nerd Font Mono:style=SemiBold:size=12"
};
static const char *prompt = NULL;
static const char *colors[SchemeLast][2] = {
	/*               fg         bg       */
	[SchemeNorm] = { "#ebdbb2", "#282828" },
	[SchemeSel]  = { "#fbf1c7", "#d65d0e" },
	[SchemeOut]  = { "#000000", "#00ffff" },
};
/* -l option; if nonzero, dmenu uses vertical list with given number of lines */
static unsigned int lines = 0;

/*
 * Characters not considered part of a word while deleting words
 * for example: " /?\"&[]"
 */
static const char worddelimiters[] = " ";

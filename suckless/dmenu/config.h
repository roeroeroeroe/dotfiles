static int topbar = 1;
static const char *fonts[] = {
	"JetBrainsMono NFP:style=SemiBold:pixelsize=14",
	// "Terminus:size=10",
};
static const char *prompt = NULL;
static const char *colors[SchemeLast][2] = {
	/*               fg         bg       */
	[SchemeNorm] = { "#525252", "#040404" },
	[SchemeSel]  = { "#ffa0ff", "#000000" },
	[SchemeOut]  = { "#000000", "#00ffff" },
};
/* -l option; if nonzero, dmenu uses vertical list with given number of lines */
static unsigned int lines = 0;

/*
 * Characters not considered part of a word while deleting words
 * for example: " /?\"&[]"
 */
static const char worddelimiters[] = " ";

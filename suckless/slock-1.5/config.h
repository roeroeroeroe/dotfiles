/* user and group to drop privileges to */
static const char *user  = "roe";
static const char *group = "wheel";

static const char *colorname[NUMCOLS] = {
	[INIT] =   "#282828",   /* after initialization */
	[INPUT] =  "#3c3836",   /* during input */
	[FAILED] = "#d65d0e",   /* wrong password */
};

/* treat a cleared input like a wrong password (color) */
static const int failonclear = 1;

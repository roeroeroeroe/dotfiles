static const char *colorname[NUMCOLS] = {
	[BG] =     "black",     /* background */
	[INIT] =   "#1e1e2e",   /* after initialization */
	[INPUT] =  "#b4befe",   /* during input */
	[FAILED] = "#f38ba8",   /* wrong password */
};

/* treat a cleared input like a wrong password (color) */
static const int failonclear = 1;

/* size of square in px */
static const int squaresize = 50;

/* number of failed password attempts until failcommand is executed.
   Set to 0 to disable */
static const int failcount = 0;

/* command to be executed after [failcount] failed password attempts */
static const char *failcommand = "shutdown";

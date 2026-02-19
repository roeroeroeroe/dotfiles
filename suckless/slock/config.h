static const char *colorname[NUMCOLS] = {
	[BG] =     "#000000",   /* background */
	[INIT] =   "#525252",   /* after initialization */
	[INPUT] =  "#ffa0ff",   /* during input */
	[FAILED] = "#ffffff",   /* wrong password */
};

/* treat a cleared input like a wrong password (color) */
static const int failonclear = 1;

/* size of square in px */
static const int squaresize = 100;

/* number of failed password attempts until failcommand is executed.
   Set to 0 to disable */
static const int failcount = 5;

/* command to be executed after [failcount] failed password attempts */
static const char *failcommand = "systemctl poweroff";

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

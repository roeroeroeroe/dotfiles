static const char *colorname[NUMCOLS] = {
	[BG] =     "#1d2021",   /* background */
	[INIT] =   "#282828",   /* after initialization */
	[INPUT] =  "#3c3836",   /* during input */
	[FAILED] = "#d65d0e",   /* wrong password */
};

/* treat a cleared input like a wrong password (color) */
static const int failonclear = 1;

/* size of square in px */
static const int squaresize = 100;

static const unsigned int borderpx  = 2;
static const unsigned int gappx     = 6;
static const unsigned int snap      = 6;
static const unsigned int systraypinning = 0;   /* 0: sloppy systray follows selected monitor, >0: pin systray to monitor X */
static const unsigned int systrayonleft = 0;    /* 0: systray in the right corner, >0: systray on left of status text */
static const unsigned int systrayspacing = 2;   /* systray spacing */
static const int systraypinningfailfirst = 1;   /* 1: if pinning fails, display systray on the first monitor, False: display systray on the last monitor*/
static const int showsystray        = 1;
static const int showbar            = 1;
static const int topbar             = 1;
static const char *fonts[]          = {
	"JetBrainsMono NFP:style=SemiBold:pixelsize=16",
	// "Terminus:pixelsize=16",
	"Noto Color Emoji:size=12"
};

static const char col_gray1[]       = "#282828";
static const char col_gray2[]       = "#3c3836";
static const char col_gray3[]       = "#ebdbb2";
static const char col_gray4[]       = "#fbf1c7";
static const char col_accent[]      = "#d65d0e";
static const char *colors[][3]      = {
	/*               fg         bg           border   */
	[SchemeNorm] = { col_gray3, col_gray1,   col_gray2 },
	[SchemeSel]  = { col_gray4, col_accent,  col_accent },
};

static const char *tags[] = { ">.<", ":x", "^=^", "^w^", "^o^", "^_^", "^.^", ":w", "c:", ":3" };
// static const char *tags[] = { "1", "2", "3", "4", "5", "6", "7", "8", "9", "10" };

static const Rule rules[] = {
	/* class          instance    title       tags mask     isfloating   monitor */
	{ "nekoray",      NULL,       NULL,       1 << 0,       0,           1 },
	{ "Tor Browser",  NULL,       NULL,       0,            1,           -1 },
	{ "vesktop",      NULL,       NULL,       1 << 3,       0,           1 },
};

static const float mfact     = 0.50; /* factor of master area size [0.05..0.95] */
static const int nmaster     = 1;    /* number of clients in master area */
static const int resizehints = 0;    /* 1 means respect size hints in tiled resizals */
static const int lockfullscreen = 1; /* 1 will force focus on the fullscreen window */

static const Layout layouts[] = {
	/* symbol     arrange function */
	{ "t",      tile },    /* first entry is default */
	{ "f",      NULL },    /* no layout function means floating behavior */
	{ "0",      monocle },
};

#define MODKEY Mod4Mask
#define TAGKEYS(KEY,TAG) \
	{ MODKEY,                       KEY,      view,           {.ui = 1 << TAG} }, \
	{ MODKEY|ControlMask,           KEY,      toggleview,     {.ui = 1 << TAG} }, \
	{ MODKEY|ShiftMask,             KEY,      tag,            {.ui = 1 << TAG} }, \
	{ MODKEY|ControlMask|ShiftMask, KEY,      toggletag,      {.ui = 1 << TAG} },

/* helper for spawning shell commands in the pre dwm-5.0 fashion */
#define SHCMD(cmd) { .v = (const char*[]){ "/bin/sh", "-c", cmd, NULL } }

static char dmenumon[2] = "0";
static const char *termcmd[]  = { "st", NULL };
static const char *dmenucmd[] = {
	"dmenu_run",
	"-m", dmenumon,
	NULL
};

static const Key keys[] = {
	/* modifier                     key        function        argument */
	{ Mod1Mask,                     XK_p,      spawn,          SHCMD("play_pause_notify") },
	{ Mod1Mask,                     XK_m,      spawn,          SHCMD("toggle_mic_notify") },
	{ Mod1Mask|ControlMask,         XK_b,      spawn,          SHCMD("browser_bookmarks_dmenu") },
	{ Mod1Mask|ControlMask,         XK_m,      spawn,          SHCMD("powermenu_dmenu") },
	{ 0,                            XK_Print,  spawn,          SHCMD("screenshot_dmenu") },
	{ Mod1Mask|ControlMask,         XK_n,      spawn,          SHCMD("notes_dmenu") },
	{ Mod1Mask|ControlMask,         XK_p,      spawn,          SHCMD("ttv_live_follows_dmenu") },
	{ Mod1Mask|ControlMask,         XK_l,      spawn,          SHCMD("cmus_lyrics") },
	{ Mod1Mask|ControlMask,         XK_s,      spawn,          SHCMD("pavucontrol") },
	{ Mod1Mask|ControlMask,         XK_t,      spawn,          SHCMD("tor-browser") },
	{ Mod1Mask|ControlMask,         XK_r,      spawn,          SHCMD("simplescreenrecorder --no-systray") },
	{ Mod1Mask|ControlMask,         XK_c,      spawn,          SHCMD("cmus_launcher") },
	{ MODKEY,                       XK_Return, spawn,          {.v = termcmd } },
	{ MODKEY,                       XK_b,      togglebar,      {0} },
	{ MODKEY|ShiftMask,             XK_j,      rotatestack,    {.i = +1 } },
	{ MODKEY|ShiftMask,             XK_k,      rotatestack,    {.i = -1 } },
	{ MODKEY,                       XK_j,      focusstack,     {.i = +1 } },
	{ MODKEY,                       XK_k,      focusstack,     {.i = -1 } },
	{ MODKEY,                       XK_h,      setmfact,       {.f = -0.05} },
	{ MODKEY,                       XK_l,      setmfact,       {.f = +0.05} },
	{ MODKEY|ShiftMask,             XK_Return, zoom,           {0} },
	{ MODKEY,                       XK_Tab,    view,           {0} },
	{ MODKEY,                       XK_w,      killclient,     {0} },
	{ MODKEY,                       XK_t,      setlayout,      {.v = &layouts[0]} },
	{ MODKEY,                       XK_m,      setlayout,      {.v = &layouts[2]} },
	{ MODKEY,                       XK_space,  spawn,          {.v = dmenucmd } },
	{ MODKEY,                       XK_s,      togglefloating, {0} },
	{ MODKEY,                       XK_comma,  focusmon,       {.i = -1 } },
	{ MODKEY,                       XK_period, focusmon,       {.i = +1 } },
	{ MODKEY|ShiftMask,             XK_comma,  tagmon,         {.i = -1 } },
	{ MODKEY|ShiftMask,             XK_period, tagmon,         {.i = +1 } },
	{ MODKEY,                       XK_minus,  setgaps,        {.i = -6 } },
	{ MODKEY,                       XK_equal,  setgaps,        {.i = +6 } },
	{ MODKEY|ShiftMask,             XK_equal,  setgaps,        {.i = 0 } },
	TAGKEYS(                        XK_1,                      0)
	TAGKEYS(                        XK_2,                      1)
	TAGKEYS(                        XK_3,                      2)
	TAGKEYS(                        XK_4,                      3)
	TAGKEYS(                        XK_5,                      4)
	TAGKEYS(                        XK_6,                      5)
	TAGKEYS(                        XK_7,                      6)
	TAGKEYS(                        XK_8,                      7)
	TAGKEYS(                        XK_9,                      8)
	TAGKEYS(                        XK_0,                      9)
	{ MODKEY,                       XK_grave,  view,           {.ui = ~0 } },
	{ MODKEY|ShiftMask,             XK_grave,  tag,            {.ui = ~0 } },
	{ MODKEY|ShiftMask,             XK_q,      quit,           {0} },
};

/* click can be ClkTagBar, ClkLtSymbol, ClkStatusText, ClkWinTitle, ClkClientWin, or ClkRootWin */
static const Button buttons[] = {
	/* click                event mask      button          function        argument */
	{ ClkLtSymbol,          0,              Button1,        setlayout,      {0} },
	// { ClkLtSymbol,          0,              Button3,        setlayout,      {.v = &layouts[2]} },
	{ ClkWinTitle,          0,              Button2,        zoom,           {0} },
	{ ClkStatusText,        0,              Button2,        spawn,          {.v = termcmd } },
	{ ClkClientWin,         MODKEY,         Button1,        movemouse,      {0} },
	{ ClkClientWin,         MODKEY,         Button2,        togglefloating, {0} },
	{ ClkClientWin,         MODKEY,         Button3,        resizemouse,    {0} },
	{ ClkTagBar,            0,              Button1,        view,           {0} },
	{ ClkTagBar,            0,              Button3,        toggleview,     {0} },
	{ ClkTagBar,            MODKEY,         Button1,        tag,            {0} },
	{ ClkTagBar,            MODKEY,         Button3,        toggletag,      {0} },
};


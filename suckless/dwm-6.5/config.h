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
	"JetBrainsMono NFP:style=SemiBold:pixelsize=14",
	// "Terminus:pixelsize=16",
	"Noto Color Emoji:size=10"
};

static const char col_norm_fg[]     = "#525252";
static const char col_norm_bg[]     = "#000000";
static const char col_norm_bd[]     = "#525252";
static const char col_sel_fg[]      = "#ffa0ff";
static const char col_sel_bg[]      = "#000000";
static const char col_sel_bd[]      = "#ffa0ff";
static const char *colors[][3]      = {
	/*               fg           bg           border   */
	[SchemeNorm] = { col_norm_fg, col_norm_bg, col_norm_bd },
	[SchemeSel]  = { col_sel_fg,  col_sel_bg,  col_sel_bd },
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
static const char *dmenucmd[] = { "dmenu_run", "-m", dmenumon, NULL };
static const char *termcmd[]  = { "st", NULL };

static const Key keys[] = {
	/* modifier                     keycode function         argument */
	{ 0,        /* XF86AudioMute */ 121,   spawn,            SHCMD("volume mute") },
	{ 0, /* XF86AudioLowerVolume */ 122,   spawn,            SHCMD("volume -") },
	{ 0, /* XF86AudioRaiseVolume */ 123,   spawn,            SHCMD("volume +") },
	{ Mod1Mask,             /* p */ 33,    spawn,            SHCMD("play_pause") },
	{ Mod1Mask,             /* m */ 58,    spawn,            SHCMD("toggle_mic_notify") },
	{ Mod1Mask|ControlMask, /* b */ 56,    spawn,            SHCMD("browser_bookmarks_dmenu") },
	{ Mod1Mask|ControlMask, /* m */ 58,    spawn,            SHCMD("powermenu_dmenu") },
	{ 0,                /* Print */ 107,   spawn,            SHCMD("screenshot_dmenu") },
	{ Mod1Mask|ControlMask, /* n */ 57,    spawn,            SHCMD("notes_dmenu") },
	{ Mod1Mask|ControlMask, /* p */ 33,    spawn,            SHCMD("ttv_live_follows_dmenu") },
	{ Mod1Mask|ControlMask, /* e */ 26,    spawn,            SHCMD("st -e nvimf") },
	{ Mod1Mask|ControlMask, /* l */ 46,    spawn,            SHCMD("cmus_lyrics") },
	{ Mod1Mask|ControlMask, /* s */ 39,    spawn,            SHCMD("pavucontrol") },
	{ Mod1Mask|ControlMask, /* t */ 28,    spawn,            SHCMD("torbrowser-launcher") },
	{ Mod1Mask|ControlMask, /* r */ 27,    spawn,            SHCMD("simplescreenrecorder --no-systray") },
	{ Mod1Mask|ControlMask, /* c */ 54,    spawn,            SHCMD("cmus_launcher") },
	{ Mod1Mask,             /* b */ 56,    spawn,            SHCMD("playerctl next") },
	{ Mod1Mask,             /* z */ 52,    spawn,            SHCMD("playerctl previous") },
	{ MODKEY,          /* Return */ 36,    spawn,            {.v = termcmd } },
	{ MODKEY,               /* b */ 56,    togglebar,        {0} },
	{ MODKEY|ShiftMask,     /* b */ 56,    toggletopbar,     {0} },
	{ MODKEY|ShiftMask,     /* j */ 44,    rotatestack,      {.i = +1 } },
	{ MODKEY|ShiftMask,     /* k */ 45,    rotatestack,      {.i = -1 } },
	{ MODKEY,               /* j */ 44,    focusstack,       {.i = +1 } },
	{ MODKEY,               /* k */ 45,    focusstack,       {.i = -1 } },
	{ MODKEY,               /* h */ 43,    setmfact,         {.f = -0.05} },
	{ MODKEY,               /* l */ 46,    setmfact,         {.f = +0.05} },
	{ MODKEY|ShiftMask,/* Return */ 36,    zoom,             {0} },
	{ MODKEY,             /* Tab */ 23,    view,             {0} },
	{ MODKEY,               /* w */ 25,    killclient,       {0} },
	{ MODKEY,               /* t */ 28,    setlayout,        {.v = &layouts[0]} },
	{ MODKEY,               /* m */ 58,    setlayout,        {.v = &layouts[2]} },
	{ MODKEY,           /* space */ 65,    spawn,            {.v = dmenucmd } },
	{ MODKEY,               /* s */ 39,    togglefloating,   {0} },
	{ MODKEY,           /* comma */ 59,    focusmon,         {.i = -1 } },
	{ MODKEY,          /* period */ 60,    focusmon,         {.i = +1 } },
	{ MODKEY|ShiftMask, /* comma */ 59,    tagmon,           {.i = -1 } },
	{ MODKEY|ShiftMask,/* period */ 60,    tagmon,           {.i = +1 } },
	{ MODKEY,           /* minus */ 20,    setgaps,          {.i = -6 } },
	{ MODKEY,           /* equal */ 21,    setgaps,          {.i = +6 } },
	{ MODKEY|ShiftMask, /* equal */ 21,    setgaps,          {.i = 0 } },
	{ MODKEY,           /* grave */ 49,    view,             {.ui = ~0 } },
	{ MODKEY|ShiftMask, /* grave */ 49,    tag,              {.ui = ~0 } },
	{ MODKEY|ShiftMask,     /* q */ 24,    quit,             {0} },
	TAGKEYS(                /* 1 */ 10,                      0)
	TAGKEYS(                /* 2 */ 11,                      1)
	TAGKEYS(                /* 3 */ 12,                      2)
	TAGKEYS(                /* 4 */ 13,                      3)
	TAGKEYS(                /* 5 */ 14,                      4)
	TAGKEYS(                /* 6 */ 15,                      5)
	TAGKEYS(                /* 7 */ 16,                      6)
	TAGKEYS(                /* 8 */ 17,                      7)
	TAGKEYS(                /* 9 */ 18,                      8)
	TAGKEYS(                /* 0 */ 19,                      9)
};

/* click can be ClkTagBar, ClkLtSymbol, ClkStatusText, ClkWinTitle, ClkClientWin, or ClkRootWin */
static const Button buttons[] = {
	/* click                event mask      button          function        argument */
	{ ClkLtSymbol,          0,              Button1,        setlayout,      {0} },
	// { ClkLtSymbol,          0,              Button3,        setlayout,      {.v = &layouts[2]} },
	{ ClkWinTitle,          0,              Button2,        killclient,     {0} },
	{ ClkStatusText,        0,              Button2,        spawn,          {.v = termcmd } },
	{ ClkClientWin,         MODKEY,         Button1,        movemouse,      {0} },
	{ ClkClientWin,         MODKEY,         Button2,        togglefloating, {0} },
	{ ClkClientWin,         MODKEY,         Button3,        resizemouse,    {0} },
	{ ClkTagBar,            0,              Button1,        view,           {0} },
	{ ClkTagBar,            0,              Button3,        toggleview,     {0} },
	{ ClkTagBar,            MODKEY,         Button1,        tag,            {0} },
	{ ClkTagBar,            MODKEY,         Button3,        toggletag,      {0} },
};


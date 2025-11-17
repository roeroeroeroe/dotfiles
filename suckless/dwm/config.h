static const unsigned int borderpx = 2;
static const unsigned int gappx    = 6;
static const unsigned int snap     = 6;
static const unsigned int systraypinning = 0;   /* 0: sloppy systray follows selected monitor, >0: pin systray to monitor X */
static const unsigned int systrayonleft  = 0;   /* 0: systray in the right corner, >0: systray on left of status text */
static const unsigned int systrayspacing = 2;   /* systray spacing */
static const int systraypinningfailfirst = 1;   /* 1: if pinning fails, display systray on the first monitor, False: display systray on the last monitor*/
static const int showsystray = 1;
static const int showbar     = 1;
static const int topbar      = 1;
static const char *fonts[]   = {
	"JetBrainsMono NFP:style=SemiBold:pixelsize=14",
	// "Terminus:size=10",
	"Noto Color Emoji:size=10"
};

static const char col_norm_fg[] = "#525252";
static const char col_norm_bg[] = "#000000";
static const char col_norm_bd[] = "#525252";
static const char col_sel_fg[]  = "#cccccc";
static const char col_sel_bg[]  = "#000000";
static const char col_sel_bd[]  = "#cccccc";
static const char *colors[][3]  = {
	/*               fg           bg           border   */
	[SchemeNorm] = { col_norm_fg, col_norm_bg, col_norm_bd },
	[SchemeSel]  = { col_sel_fg,  col_sel_bg,  col_sel_bd },
};

static const char *tags[] = { "*", "*", "*", "*", "*", "*" };

static const Rule rules[] = {
	/* class     instance     title     tags mask     isfloating     monitor */
	{ "Gimp",    NULL,        NULL,     0,            1,             -1 },
};

static const float mfact        = 0.50; /* factor of master area size [0.05..0.95] */
static const int nmaster        = 1;    /* number of clients in master area */
static const int resizehints    = 0;    /* 1 means respect size hints in tiled resizals */
static const int lockfullscreen = 1;    /* 1 will force focus on the fullscreen window */

static const float mfactdelta = 0.05;
static const int gapdelta     = 6;

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

#define CMD(...) { .v = ((const char*[]){ __VA_ARGS__, NULL }) }

static char dmenumon[2] = "0";
static const char *dmenucmd[] = { "dmenu_run", "-m", dmenumon, NULL };
static const char *termcmd[]  = { "st", NULL };

/* use xev(1) */
#define KEY_PRINT  107
#define KEY_RETURN 36
#define KEY_TAB    23
#define KEY_SPACE  65
#define KEY_COMMA  59
#define KEY_MINUS  20
#define KEY_PERIOD 60
#define KEY_EQUAL  21
#define KEY_GRAVE  49
#define KEY_XF86_AUDIO_LOWER_VOLUME 122
#define KEY_XF86_AUDIO_MUTE         121
#define KEY_XF86_AUDIO_NEXT         171
#define KEY_XF86_AUDIO_PLAY         172
#define KEY_XF86_AUDIO_PREV         173
#define KEY_XF86_AUDIO_RAISE_VOLUME 123
#define KEY_B 56
#define KEY_C 54
#define KEY_E 26
#define KEY_H 43
#define KEY_J 44
#define KEY_K 45
#define KEY_L 46
#define KEY_M 58
#define KEY_N 57
#define KEY_Q 24
#define KEY_R 27
#define KEY_S 39
#define KEY_T 28
#define KEY_U 30
#define KEY_W 25
#define KEY_1 10
#define KEY_2 11
#define KEY_3 12
#define KEY_4 13
#define KEY_5 14
#define KEY_6 15
#define KEY_7 16
#define KEY_8 17
#define KEY_9 18
#define KEY_0 19

static const Key keys[] = {
	/* modifier               keycode       function          argument */
	{ 0,                      KEY_XF86_AUDIO_LOWER_VOLUME, spawn, CMD("volume", "-") },
	{ 0,                      KEY_XF86_AUDIO_RAISE_VOLUME, spawn, CMD("volume", "+") },
	{ 0,                      KEY_XF86_AUDIO_MUTE,   spawn,   CMD("volume", "mute") },
	{ 0,                      KEY_XF86_AUDIO_PLAY,   spawn,   SHCMD("playerctl play-pause && kill -35 $(pidof sb)") },
	{ 0,                      KEY_XF86_AUDIO_PREV,   spawn,   SHCMD("playerctl previous && kill -35 $(pidof sb)") },
	{ 0,                      KEY_XF86_AUDIO_NEXT,   spawn,   SHCMD("playerctl next && kill -35 $(pidof sb)") },
	{ Mod1Mask,               KEY_M,        spawn,            CMD("toggle_mic_notify") },
	{ Mod1Mask|ControlMask,   KEY_B,        spawn,            CMD("browser_bookmarks_dmenu") },
	{ Mod1Mask|ControlMask,   KEY_M,        spawn,            CMD("powermenu_dmenu") },
	{ 0,                      KEY_PRINT,    spawn,            CMD("screenshot_dmenu") },
	{ Mod1Mask|ControlMask,   KEY_N,        spawn,            CMD("notes_dmenu") },
	{ Mod1Mask,               KEY_U,        spawn,            CMD("unicode_dmenu") },
	{ Mod1Mask|ControlMask,   KEY_E,        spawn,            CMD("st", "-e", "nvimf", "-i") },
	{ Mod1Mask|ControlMask,   KEY_L,        spawn,            CMD("cmus_lyrics") },
	{ Mod1Mask|ControlMask,   KEY_S,        spawn,            CMD("pavucontrol") },
	{ Mod1Mask|ControlMask,   KEY_R,        spawn,            CMD("simplescreenrecorder", "--no-systray") },
	{ Mod1Mask|ControlMask,   KEY_C,        spawn,            CMD("cmus_launcher") },
	{ MODKEY,                 KEY_RETURN,   spawn,            {.v = termcmd } },
	{ MODKEY,                 KEY_B,        togglebar,        {0} },
	{ MODKEY|ShiftMask,       KEY_B,        toggletopbar,     {0} },
	{ MODKEY|ShiftMask,       KEY_J,        rotatestack,      {.i = +1 } },
	{ MODKEY|ShiftMask,       KEY_K,        rotatestack,      {.i = -1 } },
	{ MODKEY,                 KEY_J,        focusstack,       {.i = +1 } },
	{ MODKEY,                 KEY_K,        focusstack,       {.i = -1 } },
	{ MODKEY,                 KEY_H,        setmfact,         {.f = -mfactdelta} },
	{ MODKEY,                 KEY_L,        setmfact,         {.f = +mfactdelta} },
	{ MODKEY|ShiftMask,       KEY_RETURN,   zoom,             {0} },
	{ MODKEY,                 KEY_TAB,      view,             {0} },
	{ MODKEY,                 KEY_W,        killclient,       {0} },
	{ MODKEY,                 KEY_T,        setlayout,        {.v = &layouts[0]} },
	{ MODKEY,                 KEY_M,        setlayout,        {.v = &layouts[2]} },
	{ MODKEY,                 KEY_SPACE,    spawn,            {.v = dmenucmd } },
	{ MODKEY,                 KEY_S,        togglefloating,   {0} },
	{ MODKEY,                 KEY_COMMA,    focusmon,         {.i = -1 } },
	{ MODKEY,                 KEY_PERIOD,   focusmon,         {.i = +1 } },
	{ MODKEY|ShiftMask,       KEY_COMMA,    tagmon,           {.i = -1 } },
	{ MODKEY|ShiftMask,       KEY_PERIOD,   tagmon,           {.i = +1 } },
	{ MODKEY,                 KEY_MINUS,    setgaps,          {.i = -gapdelta } },
	{ MODKEY,                 KEY_EQUAL,    setgaps,          {.i = +gapdelta } },
	{ MODKEY|ShiftMask,       KEY_EQUAL,    setgaps,          {.i = 0 } },
	{ MODKEY,                 KEY_GRAVE,    view,             {.ui = ~0 } },
	{ MODKEY|ShiftMask,       KEY_GRAVE,    tag,              {.ui = ~0 } },
	{ MODKEY|ShiftMask,       KEY_Q,        quit,             {0} },
	TAGKEYS(                  KEY_1,                          0)
	TAGKEYS(                  KEY_2,                          1)
	TAGKEYS(                  KEY_3,                          2)
	TAGKEYS(                  KEY_4,                          3)
	TAGKEYS(                  KEY_5,                          4)
	TAGKEYS(                  KEY_6,                          5)
	// TAGKEYS(                  KEY_7,                          6)
	// TAGKEYS(                  KEY_8,                          7)
	// TAGKEYS(                  KEY_9,                          8)
	// TAGKEYS(                  KEY_0,                          9)
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

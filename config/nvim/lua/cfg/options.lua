local indent = 4

vim.opt.background = "dark"
vim.opt.colorcolumn = "80"
vim.opt.cursorline = true
vim.opt.cursorlineopt = "line,number"
vim.opt.number = true
vim.opt.relativenumber = true
vim.opt.signcolumn = "yes"
vim.opt.statusline = "%F %h%m%r %=%l,%c/%L — %p%% %y"
vim.opt.termguicolors = true

vim.opt.expandtab = false
vim.opt.shiftwidth = indent
vim.opt.smartindent = true
vim.opt.softtabstop = indent
vim.opt.tabstop = indent

vim.opt.hlsearch = false
vim.opt.ignorecase = true
vim.opt.incsearch = true
vim.opt.smartcase = true

vim.opt.scrolloff = 5
vim.opt.whichwrap:append("<>[]hl")

vim.opt.clipboard = "unnamedplus"
vim.opt.syntax = "on"

vim.opt.inccommand = "split"
vim.opt.shortmess:append("IF")
vim.opt.showcmd = true
vim.opt.showmode = true
vim.opt.wildmenu = true

vim.opt.splitbelow = true
vim.opt.splitright = true

vim.opt.mouse = "a"
vim.opt.ttimeoutlen = 0

vim.opt.hidden = true
vim.opt.undofile = true

vim.g.netrw_banner = 0

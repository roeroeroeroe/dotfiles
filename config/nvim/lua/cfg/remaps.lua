vim.g.mapleader = " "

local function map(mode, lhs, rhs)
	vim.keymap.set(mode, lhs, rhs, { silent = true })
end

map("v", "J", ":m '>+1<CR>gv=gv")
map("v", "K", ":m '<-2<CR>gv=gv")

map("n", "<C-d>", "<C-d>zz")
map("n", "<C-u>", "<C-u>zz")

map("n", "n", "nzz")
map("n", "N", "Nzz")

map("n", "<Tab>", "<Cmd>bnext<CR>")
map("n", "<S-Tab>", "<Cmd>bprevious<CR>")

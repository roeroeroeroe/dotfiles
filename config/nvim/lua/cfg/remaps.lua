vim.g.mapleader = " "

local function map(mode, lhs, rhs)
	vim.keymap.set(mode, lhs, rhs, { silent = true })
end

map('n', '<Tab>', '<Cmd>BufferNext<CR>')
map('n', '<S-Tab>', '<Cmd>BufferPrevious<CR>')

map("n", "<leader>o", "<CMD>vsplit<CR>")
map("n", "<leader>p", "<CMD>split<CR>")

map("n", "<C-h>", "<C-w>h")
map("n", "<C-l>", "<C-w>l")
map("n", "<C-k>", "<C-w>k")
map("n", "<C-j>", "<C-w>j")

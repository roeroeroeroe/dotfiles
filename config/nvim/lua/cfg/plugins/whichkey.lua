return {
	"folke/which-key.nvim",
	event = "VeryLazy",
	opts = { delay = 500 },
	keys = {
		{
			"<leader>?",
			function()
				require("which-key").show({ global = false })
			end,
			desc = "Buffer Local Keymaps (which-key)",
		},
		{
			"<leader>u",
			"<cmd>UndotreeToggle<cr>",
			desc = "Toggle undo tree",
		},
		{
			"<leader>fm",
			function()
				local bufnr = vim.api.nvim_get_current_buf()
				require("conform").format({
					bufnr = bufnr,
					lsp_fallback = true,
					async = true,
					timeout_ms = 1000,
				})
			end,
			desc = "Format",
		},
		{
			"<leader>ff",
			"<cmd>Telescope find_files<cr>",
			desc = "Telescope find files",
		},
		{
			"<leader>fb",
			function()
				require("telescope.builtin").buffers({ initial_mode = "normal" })
			end,
			desc = "Telescope buffers",
		},
		{
			"<leader>fh",
			"<cmd>Telescope help_tags<cr>",
			desc = "Telescope help tags",
		},
		{
			"<leader>fz",
			"<cmd>Telescope current_buffer_fuzzy_find<cr>",
			desc = "Telescope current buffer fuzzy find",
		},
		{
			"<leader>fg",
			"<cmd>Telescope live_grep<cr>",
			desc = "Telescope live grep",
		},
	},
}

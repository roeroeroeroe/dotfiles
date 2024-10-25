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
			"<leader>fm",
			function()
				local bufnr = vim.api.nvim_get_current_buf()
				require("conform").format({ bufnr = bufnr })
			end,
			desc = "Format",
		},
		{
			"<leader>u",
			"<cmd>UndotreeToggle<cr>",
			desc = "Toggle undo tree",
			mode = "n",
		},

	},
}


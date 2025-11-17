return {
	"catppuccin/nvim",
	name = "catppuccin",
	priority = 1000,
	config = function()
		require("catppuccin").setup({
			flavour = "mocha",
			transparent_background = true,
			float = { transparent = true },
			styles = {
				comments = {},
				conditionals = {},
			},
		})

		vim.cmd("colorscheme catppuccin")
		vim.api.nvim_set_hl(0, "CursorLine", { bg = "#111111" })
		vim.api.nvim_set_hl(0, "CursorLineNr", { fg = "#cccccc", bg = "#111111" })
		vim.api.nvim_set_hl(0, "Normal", { bg = "#000000" })
		vim.api.nvim_set_hl(0, "NormalNC", { bg = "#000000" })
		vim.api.nvim_set_hl(0, "VertSplit", { bg = "#000000" })

		vim.api.nvim_set_hl(0, "FloatBorder", { bg = "#000000", fg = "#cccccc" })
	end,
}

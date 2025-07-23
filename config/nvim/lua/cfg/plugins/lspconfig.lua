return {
	"neovim/nvim-lspconfig",
	event = { "BufReadPre", "BufNewFile" },
	dependencies = { "hrsh7th/cmp-nvim-lsp" },
	config = function()
		local capabilities = require("cmp_nvim_lsp").default_capabilities()

		local remap_references = function()
			local wk = require("which-key")
			wk.add({
				{
					"<leader>lfr",
					function()
						require("telescope.builtin").lsp_references({
							initial_mode = "normal",
							layout_strategy = "vertical",
						})
					end,
					desc = "Find references",
				},
				-- code actions: gra
				-- implementation: gri
				-- references: grr
				-- rename: grn
				-- signature help: K | ^s in insert
			})
		end

		vim.lsp.config("*", {
			capabilities = capabilities,
			on_attach = remap_references,
		})
		vim.lsp.config("ts_ls", {
			capabilities = capabilities,
			on_attach = remap_references,
		})

		vim.diagnostic.config({
			virtual_text = {
				prefix = "",
				spacing = 0,
			},
			signs = true,
			underline = true,
			update_in_insert = false,
		})
	end,
}

return {
	"neovim/nvim-lspconfig",
	event = { "BufReadPre", "BufNewFile" },
	dependencies = {
		"hrsh7th/cmp-nvim-lsp",
	},

	config = function()
		local nvim_lsp = require("lspconfig")
		local mason_lspconfig = require("mason-lspconfig")
		local capabilities = require("cmp_nvim_lsp").default_capabilities()

		local register_remaps = function()
			local wk = require("which-key")
			wk.add({
				{
					"<leader>lgd",
					function()
						vim.lsp.buf.definition()
					end,
					desc = "Go to definition",
				},
				{
					"<leader>lh",
					function()
						vim.lsp.buf.hover()
					end,
					desc = "Hover",
				},
				{
					"<leader>lfr",
					function()
						require("telescope.builtin").lsp_references({
							initial_mode = "normal",
							layout_strategy = "vertical",
						})
						-- vim.lsp.buf.references()
					end,
					desc = "Find references",
				},
				{
					"<leader>lrn",
					function()
						vim.lsp.buf.rename()
					end,
					desc = "Rename",
				},
				{
					"<leader>lca",
					function()
						vim.lsp.buf.code_action()
					end,
					desc = "Code actions",
				},
			})
		end

		mason_lspconfig.setup_handlers({
			function(server)
				nvim_lsp[server].setup({
					capabilities = capabilities,
					on_attach = register_remaps,
				})
			end,
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

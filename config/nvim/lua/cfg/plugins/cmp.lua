return {
	"hrsh7th/nvim-cmp",
	dependencies = {
		"hrsh7th/cmp-nvim-lsp",
		"hrsh7th/cmp-buffer",
		"hrsh7th/cmp-path",
		{
			"L3MON4D3/LuaSnip",
			version = "v2.*",
			dependencies = { "rafamadriz/friendly-snippets" },
			build = "make install_jsregexp",
		},
		"saadparwaiz1/cmp_luasnip",
		"onsails/lspkind-nvim",
	},
	config = function()
		local cmp = require("cmp")
		local lspkind = require("lspkind")
		local luasnip = require("luasnip")
		local cmp_select_opts = { behavior = cmp.SelectBehavior.Select }
		require("luasnip.loaders.from_vscode").lazy_load()
		require("luasnip.config").set_config({
			history = true,
			updateEvents = "TextChanged,TextChangedI",
		})
		vim.api.nvim_create_autocmd("InsertLeave", {
			callback = function()
				if
					require("luasnip").session.current_nodes[vim.api.nvim_get_current_buf()]
					and not require("luasnip").session.jump_active
				then
					require("luasnip").unlink_current()
				end
			end,
		})

		cmp.setup({
			snippet = {
				expand = function(args)
					luasnip.lsp_expand(args.body)
				end,
			},
			preselect = cmp.PreselectMode.None,
			completion = {
				completeopt = "menu,menuone,noinsert,noselect",
			},
			mapping = {
				["<CR>"] = cmp.mapping.confirm({
					select = true,
					behavior = cmp.ConfirmBehavior.Replace,
				}),
				["<C-e>"] = cmp.mapping.abort(),
				["<C-d>"] = cmp.mapping.scroll_docs(-4),
				["<C-f>"] = cmp.mapping.scroll_docs(4),
				["<C-Space>"] = cmp.mapping.close(),
				["<Tab>"] = cmp.mapping(function(fallback)
					if cmp.visible() then
						cmp.select_next_item()
					elseif luasnip.expand_or_jumpable() then
						luasnip.expand_or_jump()
					else
						fallback()
					end
				end, { "i", "s" }),
				["<S-Tab>"] = cmp.mapping(function(fallback)
					if cmp.visible() then
						cmp.select_prev_item()
					elseif luasnip.jumpable(-1) then
						luasnip.jump(-1)
					else
						fallback()
					end
				end, { "i", "s" }),
			},
			sources = cmp.config.sources({
				{ name = "nvim_lsp" },
				{ name = "luasnip" },
			}, {
				{ name = "buffer", keyword_length = 3 },
				{ name = "path" },
			}),
			formatting = {
				format = lspkind.cmp_format({
					mode = "symbol",
					menu = {
						nvim_lsp = "[lsp]",
						luasnip = "[snip]",
						path = "[path]",
						nvim_lua = "[api]",
						buffer = "[buf]",
					},
				}),
			},
		})

		cmp.setup.cmdline({ "/" }, {
			mapping = cmp.mapping.preset.cmdline(),
			sources = {
				{ name = "buffer" },
			},
		})
	end,
}

return {
	"nvim-treesitter/nvim-treesitter",
	build = ":TSUpdate",
	config = function()
		local ts = require("nvim-treesitter")
		local parsers = {
			'bash',
			'c',
			'cpp',
			'css',
			'go',
			'html',
			'javascript',
			'json',
			'lua',
			'markdown',
			'nginx',
			'python',
			'sql',
			'typescript',
			'xml',
			'zsh',
		}
		ts.install(parsers):wait(300000)

		local patterns = vim.list_extend(vim.deepcopy(parsers), {
			'sh'
		})
		vim.api.nvim_create_autocmd("FileType", {
			pattern = patterns,
			callback = function()
				vim.bo.indentexpr = "v:lua.require('nvim-treesitter').indentexpr()"
				vim.treesitter.start()
			end,
		})
	end,
}

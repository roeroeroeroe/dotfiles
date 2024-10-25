return {
	"echasnovski/mini.indentscope",
	version = "*",
	config = function()
		require("mini.indentscope").setup({
			symbol = '╎',
			draw = {
				delay = 100
			},
			options = {
				indent_at_cursor = false,
			},
		})
	end,
}

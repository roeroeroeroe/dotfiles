package components

import "roe/sb/statusbar"

type Text struct {
	text string
	statusbar.BaseComponentConfig
}

func NewText(text string) *Text {
	const name = "text"
	if text == "" {
		panic(name + ": empty text")
	}

	base := statusbar.NewBaseComponentConfig(name, 0, 0)
	return &Text{text, *base}
}

func (t *Text) Start(update func(string), _ <-chan struct{}) {
	update(t.text)
}

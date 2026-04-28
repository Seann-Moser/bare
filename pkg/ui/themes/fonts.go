package themes

import (
	"sync"

	"gioui.org/font"
	"gioui.org/font/gofont"
)

var (
	fontsOnce   sync.Once
	fontsCached []font.FontFace
)

func loadFonts() []font.FontFace {
	fontsOnce.Do(func() {
		fontsCached = defaultFonts()
	})

	return fontsCached
}

func defaultFonts() []font.FontFace {
	return gofont.Collection()
}

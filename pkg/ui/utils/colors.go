package utils

import "image/color"

func ReadableOn(bg color.NRGBA) color.NRGBA {
	lum := 0.299*float32(bg.R) + 0.587*float32(bg.G) + 0.114*float32(bg.B)

	if lum > 160 {
		return color.NRGBA{R: 20, G: 24, B: 31, A: 255}
	}

	return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
}

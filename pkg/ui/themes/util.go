package themes

import "image/color"

func readableOn(bg color.NRGBA) color.NRGBA {
	// Simple luminance check.
	lum := 0.299*float32(bg.R) + 0.587*float32(bg.G) + 0.114*float32(bg.B)

	if lum > 160 {
		return color.NRGBA{R: 20, G: 24, B: 31, A: 255}
	}

	return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
}

func lightTokens(p RawPalette) ColorTokens {
	return ColorTokens{
		Primary:   p.Colors[0],
		Secondary: p.Colors[1],
		Tertiary:  p.Colors[2],
		Accent:    p.Colors[3],

		Background: Hex("#FAFAF7"),
		Surface:    Hex("#FFFFFF"),
		SurfaceAlt: Mix(p.Colors[3], Hex("#FFFFFF"), 0.65),
		Border:     Mix(p.Colors[0], Hex("#000000"), 0.18),

		Text:      Hex("#17202A"),
		TextMuted: Hex("#667085"),

		Success: Hex("#3BA55D"),
		Warning: Hex("#D99A21"),
		Error:   Hex("#D64545"),
	}
}

func darkTokens(p RawPalette) ColorTokens {
	return ColorTokens{
		Primary:   p.Colors[0],
		Secondary: p.Colors[1],
		Tertiary:  p.Colors[2],
		Accent:    p.Colors[3],

		Background: Hex("#0E1116"),
		Surface:    Hex("#171B22"),
		SurfaceAlt: Mix(p.Colors[0], Hex("#171B22"), 0.22),
		Border:     Mix(p.Colors[1], Hex("#FFFFFF"), 0.22),

		Text:      Hex("#F5F7FA"),
		TextMuted: Hex("#AAB2C0"),

		Success: Hex("#5CC878"),
		Warning: Hex("#E5B84E"),
		Error:   Hex("#FF6B6B"),
	}
}

func parseHexByte(a, b byte) uint8 {
	return fromHex(a)<<4 | fromHex(b)
}

func fromHex(c byte) uint8 {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}

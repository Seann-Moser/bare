package themes

import (
	"image/color"
)

var Palettes = map[PaletteName]RawPalette{
	PaletteSunset: {
		Name: PaletteSunset,
		Colors: [4]color.NRGBA{
			Hex("#FF9A86"),
			Hex("#FFB399"),
			Hex("#FFD6A6"),
			Hex("#FFF0BE"),
		},
	},
	PaletteCoastal: {
		Name: PaletteCoastal,
		Colors: [4]color.NRGBA{
			Hex("#81A6C6"),
			Hex("#AACDDC"),
			Hex("#F3E3D0"),
			Hex("#D2C4B4"),
		},
	},
	PaletteSky: {
		Name: PaletteSky,
		Colors: [4]color.NRGBA{
			Hex("#355872"),
			Hex("#7AAACE"),
			Hex("#9CD5FF"),
			Hex("#F7F8F0"),
		},
	},
	PaletteBlush: {
		Name: PaletteBlush,
		Colors: [4]color.NRGBA{
			Hex("#9ECAD6"),
			Hex("#748DAE"),
			Hex("#F5CBCB"),
			Hex("#FFEAEA"),
		},
	},
	PaletteOcean: {
		Name: PaletteOcean,
		Colors: [4]color.NRGBA{
			Hex("#0F2854"),
			Hex("#1C4D8D"),
			Hex("#4988C4"),
			Hex("#BDE8F5"),
		},
	},
	PalettePastel: {
		Name: PalettePastel,
		Colors: [4]color.NRGBA{
			Hex("#F5D2D2"),
			Hex("#F8F7BA"),
			Hex("#BDE3C3"),
			Hex("#A3CCDA"),
		},
	},
}

func New(mode Mode, paletteName PaletteName, systemDark bool) Theme {
	effectiveMode := mode
	if effectiveMode == ModeSystem {
		if systemDark {
			effectiveMode = ModeDark
		} else {
			effectiveMode = ModeLight
		}
	}
	raw, ok := Palettes[paletteName]
	if !ok {
		raw = Palettes[PaletteSunset]
	}

	t := Theme{
		Mode:    mode,
		Palette: raw.Name,
		Space: SpaceTokens{
			XS: 4,
			SM: 8,
			MD: 12,
			LG: 16,
			XL: 24,
		},
		Radius: RadiusTokens{
			SM: 6,
			MD: 10,
			LG: 14,
			XL: 20,
		},
	}

	if effectiveMode == ModeDark {
		t.Color = darkTokens(raw)
	} else {
		t.Color = lightTokens(raw)
	}

	return t
}

func Hex(s string) color.NRGBA {
	if len(s) != 7 || s[0] != '#' {
		return color.NRGBA{A: 255}
	}

	var r, g, b uint8
	r = parseHexByte(s[1], s[2])
	g = parseHexByte(s[3], s[4])
	b = parseHexByte(s[5], s[6])

	return color.NRGBA{R: r, G: g, B: b, A: 255}
}

func Mix(a, b color.NRGBA, amount float32) color.NRGBA {
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}

	inv := 1 - amount

	return color.NRGBA{
		R: uint8(float32(a.R)*amount + float32(b.R)*inv),
		G: uint8(float32(a.G)*amount + float32(b.G)*inv),
		B: uint8(float32(a.B)*amount + float32(b.B)*inv),
		A: uint8(float32(a.A)*amount + float32(b.A)*inv),
	}
}

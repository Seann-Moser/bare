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
	PaletteWarmEarth: {
		Name: PaletteWarmEarth,
		Colors: [4]color.NRGBA{
			Hex("#A47251"),
			Hex("#DD9E59"),
			Hex("#F0D8A1"),
			Hex("#DCF0C3"),
		},
	},
	PaletteSoftNeutral: {
		Name: PaletteSoftNeutral,
		Colors: [4]color.NRGBA{
			Hex("#A98B76"),
			Hex("#BFA28C"),
			Hex("#F3E4C9"),
			Hex("#BABF94"),
		},
	},
	PaletteLavender: {
		Name: PaletteLavender,
		Colors: [4]color.NRGBA{
			Hex("#F2EAE0"),
			Hex("#B4D3D9"),
			Hex("#BDA6CE"),
			Hex("#9B8EC7"),
		},
	},
	PaletteHarvest: {
		Name: PaletteHarvest,
		Colors: [4]color.NRGBA{
			Hex("#758A93"),
			Hex("#ECD5BC"),
			Hex("#E9B63B"),
			Hex("#C66E52"),
		},
	},
	PaletteCandy: {
		Name: PaletteCandy,
		Colors: [4]color.NRGBA{
			Hex("#CDF0EA"),
			Hex("#F9F9F9"),
			Hex("#F6C6EA"),
			Hex("#C490E4"),
		},
	},
	PaletteCreamyPop: {
		Name: PaletteCreamyPop,
		Colors: [4]color.NRGBA{
			Hex("#FFFFC1"),
			Hex("#FFD2A5"),
			Hex("#D38CAD"),
			Hex("#8A79AF"),
		},
	},
	PaletteViolet: {
		Name: PaletteViolet,
		Colors: [4]color.NRGBA{
			Hex("#BEDCFA"),
			Hex("#98ACF8"),
			Hex("#B088F9"),
			Hex("#DA9FF9"),
		},
	},
	PaletteForestPop: {
		Name: PaletteForestPop,
		Colors: [4]color.NRGBA{
			Hex("#C0D8C0"),
			Hex("#F5EEDC"),
			Hex("#DD4A48"),
			Hex("#ECB390"),
		},
	},
	PaletteDarkAccent: {
		Name: PaletteDarkAccent,
		Colors: [4]color.NRGBA{
			Hex("#DDDDDD"),
			Hex("#222831"),
			Hex("#30475E"),
			Hex("#F05454"),
		},
	},
	PaletteRetro: {
		Name: PaletteRetro,
		Colors: [4]color.NRGBA{
			Hex("#F6BED6"),
			Hex("#E79CC2"),
			Hex("#1D1B38"),
			Hex("#336D88"),
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

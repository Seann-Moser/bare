package themes

import (
	"image/color"
)

var Palettes = map[PaletteName]RawPalette{
	PaletteMoonlitLibrary: newStudyPalette(
		PaletteMoonlitLibrary,
		"Moonlit Library",
		modeTokens{
			Background: "#F7F4EE",
			Surface:    "#FFFFFF",
			Text:       "#1B1E27",
			TextMuted:  "#596170",
			Primary:    "#285EA8",
			Secondary:  "#7B3FA3",
			Success:    "#257A4F",
			Warning:    "#8A5A00",
			Error:      "#B3261E",
		},
		modeTokens{
			Background: "#10131A",
			Surface:    "#181D27",
			Text:       "#EEF2F7",
			TextMuted:  "#A9B3C1",
			Primary:    "#8FB7FF",
			Secondary:  "#D7A6FF",
			Success:    "#78D8A4",
			Warning:    "#F3C969",
			Error:      "#FF8C8C",
		},
	),
	PaletteSakuraStudy: newStudyPalette(
		PaletteSakuraStudy,
		"Sakura Study",
		modeTokens{
			Background: "#FFF6F8",
			Surface:    "#FFFFFF",
			Text:       "#241820",
			TextMuted:  "#6F5C67",
			Primary:    "#A2345A",
			Secondary:  "#3E5AA8",
			Success:    "#237A54",
			Warning:    "#7B5A00",
			Error:      "#B3261E",
		},
		modeTokens{
			Background: "#151217",
			Surface:    "#211B24",
			Text:       "#F8EEF5",
			TextMuted:  "#C5B5C0",
			Primary:    "#FFB3C7",
			Secondary:  "#B8C7FF",
			Success:    "#8BE0B2",
			Warning:    "#FFD17A",
			Error:      "#FF9A9A",
		},
	),
	PaletteCyberGlossary: newStudyPalette(
		PaletteCyberGlossary,
		"Cyber Glossary",
		modeTokens{
			Background: "#F0FBFD",
			Surface:    "#FFFFFF",
			Text:       "#102227",
			TextMuted:  "#526A72",
			Primary:    "#007C8A",
			Secondary:  "#5E35B1",
			Success:    "#007A4D",
			Warning:    "#7A5800",
			Error:      "#B00020",
		},
		modeTokens{
			Background: "#071217",
			Surface:    "#0E2028",
			Text:       "#E9FBFF",
			TextMuted:  "#A4C3CC",
			Primary:    "#48D7E8",
			Secondary:  "#B388FF",
			Success:    "#67E8A3",
			Warning:    "#FFD166",
			Error:      "#FF7A90",
		},
	),
	PaletteInkPaper: newStudyPalette(
		PaletteInkPaper,
		"Ink & Paper",
		modeTokens{
			Background: "#FAF7EF",
			Surface:    "#FFFFFF",
			Text:       "#24211D",
			TextMuted:  "#605A50",
			Primary:    "#8B5A00",
			Secondary:  "#1E6F62",
			Success:    "#2F7D46",
			Warning:    "#8B5A00",
			Error:      "#A23A2E",
		},
		modeTokens{
			Background: "#111111",
			Surface:    "#1B1A18",
			Text:       "#F2EFE8",
			TextMuted:  "#B9B1A5",
			Primary:    "#E0B15B",
			Secondary:  "#8EC5B5",
			Success:    "#88C999",
			Warning:    "#E0B15B",
			Error:      "#E78A7A",
		},
	),
	PaletteRoyalOtome: newStudyPalette(
		PaletteRoyalOtome,
		"Royal Otome",
		modeTokens{
			Background: "#F8F5FF",
			Surface:    "#FFFFFF",
			Text:       "#211B34",
			TextMuted:  "#615773",
			Primary:    "#5B3FA3",
			Secondary:  "#A03668",
			Success:    "#287A4F",
			Warning:    "#805A00",
			Error:      "#B3261E",
		},
		modeTokens{
			Background: "#11111F",
			Surface:    "#1A1830",
			Text:       "#F2F0FF",
			TextMuted:  "#B9B5D0",
			Primary:    "#C8B6FF",
			Secondary:  "#F4A7C6",
			Success:    "#8CD7A7",
			Warning:    "#F3C76E",
			Error:      "#FF8E9B",
		},
	),
}

type modeTokens struct {
	Background string `yaml:"background"`
	Surface    string `yaml:"surface"`
	Text       string `yaml:"text"`
	TextMuted  string `yaml:"text_muted"`
	Primary    string `yaml:"primary"`
	Secondary  string `yaml:"secondary"`
	Success    string `yaml:"success"`
	Warning    string `yaml:"warning"`
	Error      string `yaml:"error"`
}

func newStudyPalette(name PaletteName, label string, light, dark modeTokens) RawPalette {
	lightColors := colorTokens(light, false)
	darkColors := colorTokens(dark, true)

	return RawPalette{
		Name:  name,
		Label: label,
		Colors: [4]color.NRGBA{
			lightColors.Primary,
			lightColors.Secondary,
			lightColors.Success,
			lightColors.Warning,
		},
		Light: lightColors,
		Dark:  darkColors,
	}
}

func colorTokens(t modeTokens, dark bool) ColorTokens {
	bg := Hex(t.Background)
	surface := Hex(t.Surface)
	primary := Hex(t.Primary)
	secondary := Hex(t.Secondary)

	return ColorTokens{
		Primary:   primary,
		Secondary: secondary,
		Tertiary:  Mix(Hex(t.Success), surface, 0.55),
		Accent:    secondary,

		Background: bg,
		Surface:    surface,
		SurfaceAlt: surfaceAlt(bg, surface, primary, secondary, dark),
		Border:     borderColor(bg, surface, primary, dark),

		Text:      Hex(t.Text),
		TextMuted: Hex(t.TextMuted),

		Success: Hex(t.Success),
		Warning: Hex(t.Warning),
		Error:   Hex(t.Error),
	}
}

func validModeTokens(t modeTokens) bool {
	return validHex(t.Background) &&
		validHex(t.Surface) &&
		validHex(t.Text) &&
		validHex(t.TextMuted) &&
		validHex(t.Primary) &&
		validHex(t.Secondary) &&
		validHex(t.Success) &&
		validHex(t.Warning) &&
		validHex(t.Error)
}

func surfaceAlt(bg, surface, primary, secondary color.NRGBA, dark bool) color.NRGBA {
	if dark {
		return Mix(primary, surface, 0.12)
	}

	return Mix(secondary, bg, 0.08)
}

func borderColor(bg, surface, primary color.NRGBA, dark bool) color.NRGBA {
	if dark {
		return Mix(primary, surface, 0.28)
	}

	return Mix(primary, bg, 0.22)
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
		raw = Palettes[PaletteMoonlitLibrary]
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
		t.Color = raw.Dark
	} else {
		t.Color = raw.Light
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

func validHex(s string) bool {
	if len(s) != 7 || s[0] != '#' {
		return false
	}

	for i := 1; i < len(s); i++ {
		switch {
		case s[i] >= '0' && s[i] <= '9':
		case s[i] >= 'a' && s[i] <= 'f':
		case s[i] >= 'A' && s[i] <= 'F':
		default:
			return false
		}
	}

	return true
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

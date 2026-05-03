package themes

type Mode string

const (
	ModeSystem Mode = "system"
	ModeLight  Mode = "light"
	ModeDark   Mode = "dark"
)

type PaletteName string

const (
	PaletteMoonlitLibrary PaletteName = "moonlit_library"
	PaletteSakuraStudy    PaletteName = "sakura_study"
	PaletteCyberGlossary  PaletteName = "cyber_glossary"
	PaletteInkPaper       PaletteName = "ink_paper"
	PaletteRoyalOtome     PaletteName = "royal_otome"
)

var PaletteOrder = []PaletteName{
	PaletteMoonlitLibrary,
	PaletteSakuraStudy,
	PaletteCyberGlossary,
	PaletteInkPaper,
	PaletteRoyalOtome,
}

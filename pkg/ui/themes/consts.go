package themes

type Mode string

const (
	ModeSystem Mode = "system"
	ModeLight  Mode = "light"
	ModeDark   Mode = "dark"
)

type PaletteName string

const (
	PaletteSunset  PaletteName = "sunset"
	PaletteCoastal PaletteName = "coastal"
	PaletteSky     PaletteName = "sky"
	PaletteBlush   PaletteName = "blush"
	PaletteOcean   PaletteName = "ocean"
	PalettePastel  PaletteName = "pastel"
)

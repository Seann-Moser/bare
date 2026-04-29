package themes

type Mode string

const (
	ModeSystem Mode = "system"
	ModeLight  Mode = "light"
	ModeDark   Mode = "dark"
)

type PaletteName string

const (
	PaletteSunset      PaletteName = "sunset"
	PaletteCoastal     PaletteName = "coastal"
	PaletteSky         PaletteName = "sky"
	PaletteBlush       PaletteName = "blush"
	PaletteOcean       PaletteName = "ocean"
	PalettePastel      PaletteName = "pastel"
	PaletteWarmEarth   PaletteName = "warm_earth"
	PaletteSoftNeutral PaletteName = "soft_neutral"
	PaletteLavender    PaletteName = "lavender"
	PaletteHarvest     PaletteName = "harvest"
	PaletteCandy       PaletteName = "candy"
	PaletteCreamyPop   PaletteName = "creamy_pop"
	PaletteViolet      PaletteName = "violet"
	PaletteForestPop   PaletteName = "forest_pop"
	PaletteDarkAccent  PaletteName = "dark_accent"
	PaletteRetro       PaletteName = "retro"
)

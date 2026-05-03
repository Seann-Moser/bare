package themes

import (
	"image/color"
	"sync"

	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/widget/material"
)

var (
	defaultShaperOnce sync.Once
	defaultShaper     *text.Shaper
)

type RawPalette struct {
	Name  PaletteName
	Label string

	Colors [4]color.NRGBA
	Light  ColorTokens
	Dark   ColorTokens
}

type Theme struct {
	Mode    Mode
	Palette PaletteName

	Color  ColorTokens
	Space  SpaceTokens
	Radius RadiusTokens
}

func (t Theme) Gio() *material.Theme {
	return newMaterialTheme(defaultThemeShaper(), t)
}

func defaultThemeShaper() *text.Shaper {
	defaultShaperOnce.Do(func() {
		defaultShaper = text.NewShaper(text.WithCollection(loadFonts()))
	})

	return defaultShaper
}

func (t Theme) GioFont(fonts []font.FontFace) *material.Theme {
	return newMaterialTheme(text.NewShaper(text.WithCollection(fonts)), t)
}

func newMaterialTheme(shaper *text.Shaper, t Theme) *material.Theme {
	mt := material.NewTheme()
	mt.Shaper = shaper
	mt.Palette = material.Palette{
		Bg:         t.Color.Background,
		Fg:         t.Color.Text,
		ContrastBg: t.Color.Primary,
		ContrastFg: readableOn(t.Color.Primary),
	}

	return mt
}

type ColorTokens struct {
	Primary   color.NRGBA
	Secondary color.NRGBA
	Tertiary  color.NRGBA
	Accent    color.NRGBA

	Background color.NRGBA
	Surface    color.NRGBA
	SurfaceAlt color.NRGBA
	Border     color.NRGBA

	Text      color.NRGBA
	TextMuted color.NRGBA

	Success color.NRGBA
	Warning color.NRGBA
	Error   color.NRGBA
}

type SpaceTokens struct {
	XS int
	SM int
	MD int
	LG int
	XL int
}

type RadiusTokens struct {
	SM int
	MD int
	LG int
	XL int
}

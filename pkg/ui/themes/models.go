package themes

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/widget/material"
)

type RawPalette struct {
	Name   PaletteName
	Colors [4]color.NRGBA
}

type Theme struct {
	Mode    Mode
	Palette PaletteName

	Color  ColorTokens
	Space  SpaceTokens
	Radius RadiusTokens
}

func (t Theme) Gio() *material.Theme {
	shaper := text.NewShaper(text.WithCollection(loadFonts()))
	return &material.Theme{
		Shaper: shaper,
		Palette: material.Palette{
			Bg:         t.Color.Background,
			Fg:         t.Color.Text,
			ContrastBg: t.Color.Primary,
			ContrastFg: readableOn(t.Color.Primary),
		},
	}
}

func (t Theme) GioFont(fonts []font.FontFace) *material.Theme {
	shaper := text.NewShaper(text.WithCollection(fonts))
	return &material.Theme{
		Shaper: shaper,
		Palette: material.Palette{
			Bg:         t.Color.Background,
			Fg:         t.Color.Text,
			ContrastBg: t.Color.Primary,
			ContrastFg: readableOn(t.Color.Primary),
		},
	}
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

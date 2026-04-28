package themes

import (
	"image"
	"sort"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ThemeSelector struct {
	ModeButtons    map[Mode]*widget.Clickable
	PaletteButtons map[PaletteName]*widget.Clickable
}

func NewThemeSelector() *ThemeSelector {
	ts := &ThemeSelector{
		ModeButtons:    map[Mode]*widget.Clickable{},
		PaletteButtons: map[PaletteName]*widget.Clickable{},
	}

	for _, mode := range []Mode{ModeSystem, ModeLight, ModeDark} {
		ts.ModeButtons[mode] = new(widget.Clickable)
	}

	for name := range Palettes {
		ts.PaletteButtons[name] = new(widget.Clickable)
	}

	return ts
}

func (ts *ThemeSelector) Layout(
	gtx layout.Context,
	th Theme,
	systemDark bool,
) layout.Dimensions {
	gioTheme := th.Gio()

	for mode, btn := range ts.ModeButtons {
		for btn.Clicked(gtx) {
			th.Mode = mode
			th = New(mode, th.Palette, systemDark)
		}
	}

	for palette, btn := range ts.PaletteButtons {
		for btn.Clicked(gtx) {
			th.Palette = palette
			th = New(th.Mode, palette, systemDark)
		}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(material.Body1(gioTheme, "Theme").Layout),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(8))
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				modeButton(gioTheme, ts.ModeButtons[ModeSystem], "System", th.Mode == ModeSystem),
				modeButton(gioTheme, ts.ModeButtons[ModeLight], "Light", th.Mode == ModeLight),
				modeButton(gioTheme, ts.ModeButtons[ModeDark], "Dark", th.Mode == ModeDark),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(16))
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),

		layout.Rigid(material.Body1(gioTheme, "Palette").Layout),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(8))
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				ts.genLayout(gioTheme, th)...,
			)
		}),
	)
}

func (ts *ThemeSelector) LayoutThemeSelector(
	gtx layout.Context,
	th Theme,
	systemDark bool,
) (Theme, layout.Dimensions) {
	for mode, btn := range ts.ModeButtons {
		for btn.Clicked(gtx) {
			th = New(mode, th.Palette, systemDark)
		}
	}

	for palette, btn := range ts.PaletteButtons {
		for btn.Clicked(gtx) {
			th = New(th.Mode, palette, systemDark)
		}
	}

	dims := ts.layout(gtx, th)
	return th, dims
}

func (ts *ThemeSelector) layout(gtx layout.Context, th Theme) layout.Dimensions {
	gioTheme := th.Gio()

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(material.Body1(gioTheme, "Theme").Layout),

		layout.Rigid(spacer(8)),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				modeButton(gioTheme, ts.ModeButtons[ModeSystem], "System", th.Mode == ModeSystem),
				modeButton(gioTheme, ts.ModeButtons[ModeLight], "Light", th.Mode == ModeLight),
				modeButton(gioTheme, ts.ModeButtons[ModeDark], "Dark", th.Mode == ModeDark),
			)
		}),

		layout.Rigid(spacer(16)),
		layout.Rigid(material.Body1(gioTheme, "Palette").Layout),
		layout.Rigid(spacer(8)),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				ts.genLayout(gioTheme, th)...,
			)
		}),
	)
}

func (ts *ThemeSelector) genLayout(gioTheme *material.Theme, th Theme) (l []layout.FlexChild) {
	keys := make([]string, 0, len(Palettes))
	for k := range Palettes {
		keys = append(keys, string(k))
	}

	sort.Strings(keys)

	for _, k := range keys {
		j := PaletteName(k)
		l = append(l, paletteButton(gioTheme, ts.PaletteButtons[j], Palettes[j], th.Palette == j))
	}
	return
}

func modeButton(
	th *material.Theme,
	btn *widget.Clickable,
	label string,
	selected bool,
) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		b := material.Button(th, btn, label)
		if selected {
			b.Background = th.Palette.ContrastBg
			b.Color = th.Palette.ContrastFg
		}
		b.Inset = layout.UniformInset(unit.Dp(8))
		return b.Layout(gtx)
	})
}

func paletteButton(
	th *material.Theme,
	btn *widget.Clickable,
	p RawPalette,
	selected bool,
) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		label := string(p.Name)
		if selected {
			label = "✓ " + label
		}

		b := material.Button(th, btn, label)
		b.Background = p.Colors[0]
		b.Color = readableOn(p.Colors[0])
		b.Inset = layout.UniformInset(unit.Dp(8))

		return b.Layout(gtx)
	})
}

func spacer(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := gtx.Dp(dp)
		return layout.Dimensions{
			Size: image.Pt(size, size),
		}
	}
}

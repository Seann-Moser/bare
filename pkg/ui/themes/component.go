package themes

import (
	"image"
	"sort"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ThemeSelector struct {
	ModeButtons    map[Mode]*widget.Clickable
	PaletteButtons map[PaletteName]*widget.Clickable
	PaletteToggle  widget.Clickable
	PaletteOpen    bool
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
	for mode, btn := range ts.ModeButtons {
		for btn.Clicked(gtx) {
			th = New(mode, th.Palette, systemDark)
		}
	}

	for palette, btn := range ts.PaletteButtons {
		for btn.Clicked(gtx) {
			th = New(th.Mode, palette, systemDark)
			ts.PaletteOpen = false
		}
	}

	for ts.PaletteToggle.Clicked(gtx) {
		ts.PaletteOpen = !ts.PaletteOpen
	}

	return ts.layout(gtx, th)
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
			ts.PaletteOpen = false
		}
	}

	for ts.PaletteToggle.Clicked(gtx) {
		ts.PaletteOpen = !ts.PaletteOpen
	}

	dims := ts.layout(gtx, th)
	return th, dims
}

func (ts *ThemeSelector) layout(gtx layout.Context, th Theme) layout.Dimensions {
	gioTheme := th.Gio()

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return settingHeader(gtx, gioTheme, th, "Theme Mode", modeLabel(th.Mode))
		}),

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
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return settingHeader(gtx, gioTheme, th, "Palette", paletteLabel(th.Palette))
		}),
		layout.Rigid(spacer(8)),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ts.paletteDropdown(gtx, gioTheme, th)
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

func (ts *ThemeSelector) paletteDropdown(
	gtx layout.Context,
	gioTheme *material.Theme,
	th Theme,
) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dropdownToggle(gtx, gioTheme, th, &ts.PaletteToggle, ts.PaletteOpen, paletteLabel(th.Palette))
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !ts.PaletteOpen {
				return layout.Dimensions{}
			}

			return layout.Inset{
				Top: unit.Dp(10),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return dropdownSurface(gtx, th, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx, ts.genLayout(gioTheme, th)...)
					})
				})
			})
		}),
	)
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
		label := paletteLabel(p.Name)

		b := material.Button(th, btn, label)
		b.Background = p.Colors[0]
		b.Color = readableOn(p.Colors[0])
		b.Inset = layout.UniformInset(unit.Dp(8))
		if selected {
			b.TextSize = unit.Sp(15)
		}

		return b.Layout(gtx)
	})
}

func settingHeader(
	gtx layout.Context,
	gioTheme *material.Theme,
	th Theme,
	label string,
	current string,
) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body1(gioTheme, label)
			lbl.Color = th.Color.Text
			return lbl.Layout(gtx)
		}),
		layout.Rigid(spacer(4)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(gioTheme, "Current: "+current)
			lbl.Color = th.Color.TextMuted
			return lbl.Layout(gtx)
		}),
	)
}

func dropdownToggle(
	gtx layout.Context,
	gioTheme *material.Theme,
	th Theme,
	btn *widget.Clickable,
	open bool,
	current string,
) layout.Dimensions {
	label := current + "  v"
	if open {
		label = current + "  ^"
	}

	return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return dropdownSurface(gtx, th, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(12),
				Bottom: unit.Dp(12),
				Left:   unit.Dp(14),
				Right:  unit.Dp(14),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body1(gioTheme, label)
						lbl.Color = th.Color.Text
						return lbl.Layout(gtx)
					}),
				)
			})
		})
	})
}

func dropdownSurface(
	gtx layout.Context,
	th Theme,
	child layout.Widget,
) layout.Dimensions {
	defer clip.RRect{
		Rect: image.Rectangle{Max: gtx.Constraints.Max},
		NE:   gtx.Dp(unit.Dp(th.Radius.MD)),
		NW:   gtx.Dp(unit.Dp(th.Radius.MD)),
		SE:   gtx.Dp(unit.Dp(th.Radius.MD)),
		SW:   gtx.Dp(unit.Dp(th.Radius.MD)),
	}.Push(gtx.Ops).Pop()

	paint.Fill(gtx.Ops, th.Color.SurfaceAlt)
	return child(gtx)
}

func modeLabel(mode Mode) string {
	switch mode {
	case ModeDark:
		return "Dark"
	case ModeLight:
		return "Light"
	default:
		return "System"
	}
}

func paletteLabel(name PaletteName) string {
	switch name {
	case PaletteSunset:
		return "Sunset"
	case PaletteCoastal:
		return "Coastal"
	case PaletteSky:
		return "Sky"
	case PaletteBlush:
		return "Blush"
	case PaletteOcean:
		return "Ocean"
	case PalettePastel:
		return "Pastel"
	default:
		return string(name)
	}
}

func spacer(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := gtx.Dp(dp)
		return layout.Dimensions{
			Size: image.Pt(size, size),
		}
	}
}

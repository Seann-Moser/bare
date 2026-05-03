package themes

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	uiutils "github.com/DarlingGoose/bare/pkg/ui/utils"
)

type ThemeSelector struct {
	ModeButtons     map[Mode]*widget.Clickable
	PaletteButtons  map[PaletteName]*widget.Clickable
	PaletteToggle   widget.Clickable
	PaletteOpen     bool
	PaletteList     layout.List
	ReloadCustom    widget.Clickable
	CustomLoadError string
}

func NewThemeSelector() *ThemeSelector {
	loadErr := LoadCustomThemes()
	ts := &ThemeSelector{
		ModeButtons:    map[Mode]*widget.Clickable{},
		PaletteButtons: map[PaletteName]*widget.Clickable{},
		PaletteList: layout.List{
			Axis: layout.Vertical,
		},
	}
	if loadErr != nil {
		ts.CustomLoadError = loadErr.Error()
	}

	for _, mode := range []Mode{ModeSystem, ModeLight, ModeDark} {
		ts.ModeButtons[mode] = new(widget.Clickable)
	}

	ts.ensurePaletteButtons()

	return ts
}

func (ts *ThemeSelector) Layout(
	gtx layout.Context,
	th Theme,
	systemDark bool,
) layout.Dimensions {
	ts.handleCustomReload(gtx)
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
	ts.handleCustomReload(gtx)
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
			return themeGuidance(gtx, gioTheme, th)
		}),

		layout.Rigid(uiutils.Spacer(16)),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return customThemeControls(gtx, gioTheme, th, ts)
		}),

		layout.Rigid(uiutils.Spacer(16)),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return settingHeader(gtx, gioTheme, th, "Theme Mode", "")
		}),

		layout.Rigid(uiutils.Spacer(8)),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				modeButton(gioTheme, ts.ModeButtons[ModeSystem], "System", th.Mode == ModeSystem),
				modeButton(gioTheme, ts.ModeButtons[ModeLight], "Light", th.Mode == ModeLight),
				modeButton(gioTheme, ts.ModeButtons[ModeDark], "Dark", th.Mode == ModeDark),
			)
		}),

		layout.Rigid(uiutils.Spacer(16)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return settingHeader(gtx, gioTheme, th, "Palette", paletteLabel(th.Palette))
		}),
		layout.Rigid(uiutils.Spacer(8)),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ts.paletteDropdown(gtx, gioTheme, th)
		}),
	)
}

func (ts *ThemeSelector) handleCustomReload(gtx layout.Context) {
	for ts.ReloadCustom.Clicked(gtx) {
		if err := LoadCustomThemes(); err != nil {
			ts.CustomLoadError = err.Error()
		} else {
			ts.CustomLoadError = ""
		}
		ts.ensurePaletteButtons()
	}
}

func (ts *ThemeSelector) ensurePaletteButtons() {
	for name := range Palettes {
		if _, ok := ts.PaletteButtons[name]; !ok {
			ts.PaletteButtons[name] = new(widget.Clickable)
		}
	}
}

func themeGuidance(
	gtx layout.Context,
	gioTheme *material.Theme,
	th Theme,
) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body1(gioTheme, "Theme Notes")
			lbl.Color = th.Color.Text
			return lbl.Layout(gtx)
		}),
		layout.Rigid(uiutils.Spacer(6)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(gioTheme, "Moonlit Library is the balanced default. Sakura Study and Royal Otome are warmer VN-style options, Cyber Glossary fits tooling-heavy workflows, and Ink & Paper is best for long reading.")
			lbl.Color = th.Color.TextMuted
			return lbl.Layout(gtx)
		}),
		layout.Rigid(uiutils.Spacer(6)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(gioTheme, "Accessibility: keep contrast high, preserve visible focus states, and pair color-coded states with labels, icons, underlines, or borders.")
			lbl.Color = th.Color.TextMuted
			return lbl.Layout(gtx)
		}),
	)
}

func customThemeControls(
	gtx layout.Context,
	gioTheme *material.Theme,
	th Theme,
	ts *ThemeSelector,
) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return settingHeader(gtx, gioTheme, th, "Custom Themes", "Load from ~/.config/bare/themes.yaml or ~/.config/bare/themes/*.yaml")
		}),
		layout.Rigid(uiutils.Spacer(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			b := material.Button(gioTheme, &ts.ReloadCustom, "Reload Custom Themes")
			b.Background = th.Color.Secondary
			b.Color = readableOn(th.Color.Secondary)
			b.Inset = layout.UniformInset(unit.Dp(8))
			return b.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if ts.CustomLoadError == "" {
				return layout.Dimensions{}
			}

			return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(gioTheme, "Custom theme load error: "+ts.CustomLoadError)
				lbl.Color = th.Color.Error
				return lbl.Layout(gtx)
			})
		}),
	)
}

func (ts *ThemeSelector) genLayout(gioTheme *material.Theme, th Theme) (l []layout.FlexChild) {
	for _, name := range OrderedPalettes() {
		palette, ok := Palettes[name]
		if !ok {
			continue
		}
		l = append(l, paletteButton(gioTheme, ts.PaletteButtons[name], palette, th.Palette == name))
	}
	return
}

func (ts *ThemeSelector) paletteDropdown(
	gtx layout.Context,
	gioTheme *material.Theme,
	th Theme,
) layout.Dimensions {
	dims := dropdownToggle(gtx, gioTheme, th, &ts.PaletteToggle, ts.PaletteOpen, paletteLabel(th.Palette))
	if !ts.PaletteOpen {
		return dims
	}

	macro := op.Record(gtx.Ops)
	offset := op.Offset(image.Pt(0, dims.Size.Y+gtx.Dp(unit.Dp(10)))).Push(gtx.Ops)

	menuGTX := gtx
	menuGTX.Constraints.Min = image.Pt(dims.Size.X, 0)
	menuGTX.Constraints.Max.X = dims.Size.X
	menuGTX.Constraints.Max.Y = gtx.Dp(unit.Dp(240))

	dropdownSurface(menuGTX, th, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return ts.PaletteList.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx, ts.genLayout(gioTheme, th)...)
			})
		})
	})

	offset.Pop()
	op.Defer(gtx.Ops, macro.Stop())

	return dims
}

func modeButton(
	th *material.Theme,
	btn *widget.Clickable,
	label string,
	selected bool,
) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		b := material.Button(th, btn, label)
		b.Background = th.Palette.Bg
		b.Color = th.Palette.Fg
		if selected {
			b.Background = th.Palette.ContrastBg
			b.Color = th.Palette.ContrastFg
			b.TextSize = unit.Sp(15)
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
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
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
	children := []layout.FlexChild{
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body1(gioTheme, label)
			lbl.Color = th.Color.Text
			return lbl.Layout(gtx)
		}),
	}

	if current != "" {
		children = append(children,
			layout.Rigid(uiutils.Spacer(4)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(gioTheme, "Current: "+current)
				lbl.Color = th.Color.TextMuted
				return lbl.Layout(gtx)
			}),
		)
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, children...)
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
	return uiutils.RoundedSurface(gtx, th.Color.SurfaceAlt, unit.Dp(th.Radius.MD), child)
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
	if palette, ok := Palettes[name]; ok && palette.Label != "" {
		return palette.Label
	}

	return string(name)
}

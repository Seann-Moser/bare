package ui

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
	uiutils "github.com/Seann-Moser/bare/pkg/ui/utils"
)

type Modal struct {
	Open bool

	CloseButton widget.Clickable
	ScrimButton widget.Clickable

	CloseOnScrim bool
}

func (m *Modal) Layout(
	gtx layout.Context,
	th themes.Theme,
	title string,
	content layout.Widget,
) layout.Dimensions {
	if !m.Open {
		return layout.Dimensions{}
	}

	if m.CloseOnScrim {
		for m.ScrimButton.Clicked(gtx) {
			m.Open = false
		}
	}

	for m.CloseButton.Clicked(gtx) {
		m.Open = false
	}

	// Full-screen scrim.
	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return m.ScrimButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				paint.FillShape(
					gtx.Ops,
					color.NRGBA{R: 0, G: 0, B: 0, A: 140},
					clip.Rect{Max: gtx.Constraints.Max}.Op(),
				)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),

		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X = min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(520)))

				return uiutils.Card(gtx, unit.Dp(th.Radius.LG), th.Color.Surface, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis:      layout.Horizontal,
									Alignment: layout.Middle,
								}.Layout(gtx,
									layout.Flexed(1, material.H6(th.Gio(), title).Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return material.Button(th.Gio(), &m.CloseButton, "Close").Layout(gtx)
									}),
								)
							}),
							layout.Rigid(uiutils.Spacer(12)),
							layout.Rigid(content),
						)
					})
				})
			})
		}),
	)

	return dims
}

func OpenPopout(title string, th themes.Theme, content func(layout.Context) layout.Dimensions) {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title(title),
			app.Size(unit.Dp(640), unit.Dp(420)),
		)

		var ops op.Ops

		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				content(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()
}

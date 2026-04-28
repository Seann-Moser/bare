package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/icons"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
)

type ButtonVariant string

const (
	ButtonPrimary   ButtonVariant = "primary"
	ButtonSecondary ButtonVariant = "secondary"
	ButtonGhost     ButtonVariant = "ghost"
)

type Button struct {
	Clickable *widget.Clickable
	Icon      bool
	Text      string
	Prefix    string // example: "mdi:plus"
	Suffix    string // example: "mdi:chevron-right"

	Variant ButtonVariant
}

func (b *Button) Clicked(gtx layout.Context) bool {
	if b.Clickable == nil {
		return false
	}
	return b.Clickable.Clicked(gtx)
}

func (b *Button) Layout(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	var bg, fg color.NRGBA

	switch b.Variant {
	case ButtonSecondary:
		bg = th.Color.SurfaceAlt
		fg = th.Color.Text
	case ButtonGhost:
		bg = color.NRGBA{}
		fg = th.Color.Text
	case ButtonPrimary:
		fallthrough
	default:
		bg = th.Color.Primary
		fg = readableOn(th.Color.Primary)
	}

	if b.Clickable != nil && b.Clickable.Hovered() {
		bg = buttonHoverColor(bg, th, b.Variant)
	}

	if b.Clickable == nil {
		return layout.Dimensions{}
	}

	return b.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = 0
		var ce []layout.FlexChild
		if b.Icon {
			ce = []layout.FlexChild{
				iconChild(ic, b.Text, unit.Dp(18), fg),
			}
		} else {
			ce = []layout.FlexChild{
				iconChild(ic, b.Prefix, unit.Dp(18), fg),
				labelChild(th, b.Text, fg),
				iconChild(ic, b.Suffix, unit.Dp(18), fg),
			}
		}

		macro := op.Record(gtx.Ops)
		dims := layout.Inset{
			Top:    unit.Dp(12),
			Bottom: unit.Dp(12),
			Left:   unit.Dp(14),
			Right:  unit.Dp(14),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
				Spacing:   layout.SpaceStart,
			}.Layout(gtx,
				ce...,
			)
		})
		call := macro.Stop()

		if bg.A > 0 {
			paint.FillShape(
				gtx.Ops,
				bg,
				clip.RRect{
					Rect: image.Rectangle{Max: dims.Size},
					SE:   gtx.Dp(unit.Dp(th.Radius.MD)),
					SW:   gtx.Dp(unit.Dp(th.Radius.MD)),
					NW:   gtx.Dp(unit.Dp(th.Radius.MD)),
					NE:   gtx.Dp(unit.Dp(th.Radius.MD)),
				}.Op(gtx.Ops),
			)
		}
		call.Add(gtx.Ops)

		return dims
	})
}

func buttonHoverColor(bg color.NRGBA, th themes.Theme, variant ButtonVariant) color.NRGBA {
	switch variant {
	case ButtonGhost:
		return themes.Mix(th.Color.SurfaceAlt, th.Color.Background, 0.55)
	case ButtonSecondary:
		return themes.Mix(th.Color.Primary, bg, 0.18)
	default:
		return themes.Mix(th.Color.Accent, bg, 0.2)
	}
}

func iconChild(
	ic *icons.Iconify,
	name string,
	size unit.Dp,
	col color.NRGBA,
) layout.FlexChild {
	if name == "" || ic == nil {
		return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		})
	}

	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Right: unit.Dp(8),
			Left:  unit.Dp(8),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return ic.Layout(gtx, name, size, col)
		})
	})
}

func labelChild(
	th themes.Theme,
	text string,
	col color.NRGBA,
) layout.FlexChild {
	if text == "" {
		return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		})
	}

	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		gt := th.Gio()
		lbl := material.Body1(gt, text)
		lbl.Color = col
		return lbl.Layout(gtx)
	})
}

func readableOn(bg color.NRGBA) color.NRGBA {
	lum := 0.299*float32(bg.R) + 0.587*float32(bg.G) + 0.114*float32(bg.B)

	if lum > 160 {
		return color.NRGBA{R: 20, G: 24, B: 31, A: 255}
	}

	return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
}

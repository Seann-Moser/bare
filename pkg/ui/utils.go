package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
)

func card(gtx layout.Context, th themes.Theme, child layout.Widget) layout.Dimensions {
	defer clip.RRect{
		Rect: image.Rectangle{Max: gtx.Constraints.Max},
		NE:   gtx.Dp(unit.Dp(th.Radius.LG)),
		NW:   gtx.Dp(unit.Dp(th.Radius.LG)),
		SE:   gtx.Dp(unit.Dp(th.Radius.LG)),
		SW:   gtx.Dp(unit.Dp(th.Radius.LG)),
	}.Push(gtx.Ops).Pop()

	paint.Fill(gtx.Ops, th.Color.Surface)
	return child(gtx)
}

func spacer(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		v := gtx.Dp(dp)
		return layout.Dimensions{Size: image.Pt(v, v)}
	}
}

func spacerH(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(0, gtx.Dp(dp))}
	}
}

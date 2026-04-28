package utils

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func Panel(
	gtx layout.Context,
	bg color.NRGBA,
	radius unit.Dp,
	child layout.Widget,
) layout.Dimensions {
	defer clip.RRect{
		Rect: image.Rectangle{Max: gtx.Constraints.Max},
		NE:   gtx.Dp(radius),
		NW:   gtx.Dp(radius),
		SE:   gtx.Dp(radius),
		SW:   gtx.Dp(radius),
	}.Push(gtx.Ops).Pop()

	paint.Fill(gtx.Ops, bg)

	return child(gtx)
}

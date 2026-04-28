package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type AppShell struct {
	SidebarWidth unit.Dp
}

func (s AppShell) Layout(
	gtx layout.Context,
	sidebar layout.Widget,
	main layout.Widget,
	overlay layout.Widget,
) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					width := gtx.Dp(s.SidebarWidth)
					gtx.Constraints.Min.X = width
					gtx.Constraints.Max.X = width
					return sidebar(gtx)
				}),
				layout.Flexed(1, main),
			)
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			if overlay == nil {
				return layout.Dimensions{}
			}
			return overlay(gtx)
		}),
	)
}

package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/icons"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
	"github.com/Seann-Moser/bare/pkg/ui/utils"
)

type TopbarAction struct {
	Clickable *widget.Clickable
	Text      string
	Prefix    string
	Variant   ButtonVariant
}

func LayoutTopbar(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	title string,
	actions []TopbarAction,
) layout.Dimensions {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			lbl := material.H5(th.Gio(), title)
			lbl.Color = th.Color.Text
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutActions(gtx, th, ic, actions)
		}),
	)
}

func layoutActions(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	actions []TopbarAction,
) layout.Dimensions {
	children := make([]layout.FlexChild, 0, len(actions)*2)
	for i, action := range actions {
		action := action
		if i > 0 {
			children = append(children, layout.Rigid(utils.SpacerW(8)))
		}
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := Button{
				Clickable: action.Clickable,
				Text:      action.Text,
				Prefix:    action.Prefix,
				Variant:   action.Variant,
			}
			return btn.Layout(gtx, th, ic)
		}))
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx, children...)
}

func LayoutSidebarTitle(
	gtx layout.Context,
	th themes.Theme,
	title string,
) layout.Dimensions {
	return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.H6(th.Gio(), title)
		lbl.Color = th.Color.Text
		return lbl.Layout(gtx)
	})
}

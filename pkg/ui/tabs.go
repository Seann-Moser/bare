package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/icons"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
)

type TabItem struct {
	ID    string
	Label string
	Icon  string // optional, example: "mdi:home"
}

type Tabs struct {
	Items  []TabItem
	Active string

	clicks map[string]*widget.Clickable
}

func NewTabs(items []TabItem, active string) *Tabs {
	t := &Tabs{
		Items:  items,
		Active: active,
		clicks: map[string]*widget.Clickable{},
	}

	for _, item := range items {
		t.clicks[item.ID] = new(widget.Clickable)
	}

	if t.Active == "" && len(items) > 0 {
		t.Active = items[0].ID
	}

	return t
}

func (t *Tabs) Selected() string {
	return t.Active
}

func (t *Tabs) SetItems(items []TabItem) {
	t.Items = items

	for _, item := range items {
		if _, ok := t.clicks[item.ID]; !ok {
			t.clicks[item.ID] = new(widget.Clickable)
		}
	}

	if t.Active == "" && len(items) > 0 {
		t.Active = items[0].ID
	}
}

func (t *Tabs) Layout(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	for _, item := range t.Items {
		btn := t.clicks[item.ID]
		for btn.Clicked(gtx) {
			t.Active = item.ID
		}
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx, t.children(th, ic)...)
}

func (t *Tabs) children(th themes.Theme, ic *icons.Iconify) []layout.FlexChild {
	children := make([]layout.FlexChild, 0, len(t.Items))

	for _, item := range t.Items {
		item := item

		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.layoutTab(gtx, th, ic, item)
		}))
	}

	return children
}

func (t *Tabs) layoutTab(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	item TabItem,
) layout.Dimensions {
	btn := t.clicks[item.ID]
	active := item.ID == t.Active

	bg := th.Color.Surface
	fg := th.Color.TextMuted

	if active {
		bg = th.Color.Primary
		fg = readableOn(th.Color.Primary)
	}

	return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return roundedBackground(gtx, bg, unit.Dp(th.Radius.MD), func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(12),
				Right:  unit.Dp(12),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					tabIcon(ic, item.Icon, fg),
					tabLabel(th, item.Label, fg),
				)
			})
		})
	})
}

func tabIcon(
	ic *icons.Iconify,
	name string,
	col color.NRGBA,
) layout.FlexChild {
	if ic == nil || name == "" {
		return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{}
		})
	}

	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Right: unit.Dp(6),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return ic.Layout(gtx, name, unit.Dp(18), col)
		})
	})
}

func tabLabel(
	th themes.Theme,
	text string,
	col color.NRGBA,
) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		gt := th.Gio()

		lbl := material.Body1(gt, text)
		lbl.Color = col

		return lbl.Layout(gtx)
	})
}

func roundedBackground(
	gtx layout.Context,
	col color.NRGBA,
	radius unit.Dp,
	child layout.Widget,
) layout.Dimensions {
	defer clip.RRect{
		Rect: image.Rectangle{
			Max: gtx.Constraints.Max,
		},
		NE: gtx.Dp(radius),
		NW: gtx.Dp(radius),
		SE: gtx.Dp(radius),
		SW: gtx.Dp(radius),
	}.Push(gtx.Ops).Pop()

	paint.Fill(gtx.Ops, col)

	return child(gtx)
}

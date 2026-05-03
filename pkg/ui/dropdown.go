package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/DarlingGoose/bare/pkg/ui/icons"
	"github.com/DarlingGoose/bare/pkg/ui/themes"
	"github.com/DarlingGoose/bare/pkg/ui/utils"
)

type Dropdown struct {
	Toggle widget.Clickable
	Open   bool
	List   layout.List

	Prefix     string
	Variant    ButtonVariant
	Width      unit.Dp
	MaxHeight  unit.Dp
	OffsetY    unit.Dp
	AlignRight bool
}

func (d *Dropdown) Update(gtx layout.Context) {
	for d.Toggle.Clicked(gtx) {
		d.Open = !d.Open
	}
}

func (d *Dropdown) Close() {
	d.Open = false
}

func (d *Dropdown) Layout(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	label string,
	menu layout.Widget,
) layout.Dimensions {
	suffix := "mdi:chevron-down"
	if d.Open {
		suffix = "mdi:chevron-up"
	}

	btn := Button{
		Clickable: &d.Toggle,
		Text:      label,
		Prefix:    d.Prefix,
		Suffix:    suffix,
		Variant:   d.variant(),
	}
	dims := btn.Layout(gtx, th, ic)

	if !d.Open || menu == nil {
		return dims
	}

	menuWidth := dims.Size.X
	configuredWidth := gtx.Dp(d.width())
	if configuredWidth > menuWidth {
		menuWidth = configuredWidth
	}

	x := 0
	if d.AlignRight {
		x = dims.Size.X - menuWidth
		if x < 0 {
			x = 0
		}
	}

	macro := op.Record(gtx.Ops)
	offset := op.Offset(image.Pt(x, gtx.Dp(d.offsetY()))).Push(gtx.Ops)

	menuGTX := gtx
	menuGTX.Constraints.Min = image.Point{}
	menuGTX.Constraints.Max = image.Pt(
		menuWidth,
		gtx.Dp(d.maxHeight()),
	)

	utils.Panel(menuGTX, th.Color.Surface, unit.Dp(th.Radius.MD), func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			d.ensureListAxis()
			return d.List.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
				return menu(gtx)
			})
		})
	})

	offset.Pop()
	op.Defer(gtx.Ops, macro.Stop())

	return dims
}

func (d *Dropdown) width() unit.Dp {
	if d.Width <= 0 {
		return unit.Dp(180)
	}
	return d.Width
}

func (d *Dropdown) maxHeight() unit.Dp {
	if d.MaxHeight <= 0 {
		return unit.Dp(240)
	}
	return d.MaxHeight
}

func (d *Dropdown) offsetY() unit.Dp {
	if d.OffsetY <= 0 {
		return unit.Dp(48)
	}
	return d.OffsetY
}

func (d *Dropdown) variant() ButtonVariant {
	if d.Variant == "" {
		return ButtonSecondary
	}
	return d.Variant
}

func (d *Dropdown) ensureListAxis() {
	if d.List.Axis != layout.Vertical {
		d.List.Axis = layout.Vertical
	}
}

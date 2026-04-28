package text

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
)

type TextBlockKind string

const (
	TextParagraph TextBlockKind = "paragraph"
	TextHeading   TextBlockKind = "heading"
	TextCode      TextBlockKind = "code"
	TextMuted     TextBlockKind = "muted"
)

type TextBlock struct {
	Kind TextBlockKind
	Text string
}

type TextView struct {
	List layout.List

	Blocks []TextBlock

	AutoScrollBottom bool
	lastBlockCount   int
}

func NewTextView() *TextView {
	return &TextView{
		List: layout.List{
			Axis: layout.Vertical,
		},
		AutoScrollBottom: true,
	}
}

func (v *TextView) SetBlocks(blocks []TextBlock) {
	v.Blocks = blocks
}

func (v *TextView) Append(block TextBlock) {
	v.Blocks = append(v.Blocks, block)
}

func (v *TextView) Layout(gtx layout.Context, th themes.Theme) layout.Dimensions {
	gioTheme := th.Gio()

	if v.AutoScrollBottom && len(v.Blocks) != v.lastBlockCount {
		v.List.ScrollToEnd = true
		v.lastBlockCount = len(v.Blocks)
	}

	return v.List.Layout(gtx, len(v.Blocks), func(gtx layout.Context, i int) layout.Dimensions {
		block := v.Blocks[i]

		return layout.UniformInset(unit.Dp(6)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			switch block.Kind {
			case TextHeading:
				lbl := material.H6(gioTheme, block.Text)
				lbl.Color = th.Color.Text
				return lbl.Layout(gtx)

			case TextCode:
				lbl := material.Body2(gioTheme, block.Text)
				lbl.Color = th.Color.Text
				return codeBox(gtx, th, lbl.Layout)

			case TextMuted:
				lbl := material.Body2(gioTheme, block.Text)
				lbl.Color = th.Color.TextMuted
				return lbl.Layout(gtx)

			default:
				lbl := material.Body1(gioTheme, block.Text)
				lbl.Color = th.Color.Text
				return lbl.Layout(gtx)
			}
		})
	})
}

func codeBox(gtx layout.Context, th themes.Theme, child layout.Widget) layout.Dimensions {
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, child)
}

package text

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/DarlingGoose/bare/pkg/ui/themes"
)

type SelectableTextBlock struct {
	State widget.Selectable
	Text  string
}

func (b *SelectableTextBlock) Layout(gtx layout.Context, th themes.Theme) layout.Dimensions {
	gt := th.Gio()

	lbl := material.Body1(gt, b.Text)
	lbl.State = &b.State

	if b.State.Update(gtx) {
		fmt.Println("selection changed")
	}

	return lbl.Layout(gtx)
}

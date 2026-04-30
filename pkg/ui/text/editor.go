package text

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/DarlingGoose/bare/pkg/ui/themes"
)

type TextEditor struct {
	Editor widget.Editor

	Hint string

	AutoScrollBottom bool
	Value            string
	Highlighted      string
	MaxHeight        unit.Dp
}

func NewTextEditor(hint string) *TextEditor {
	return &TextEditor{
		Hint: hint,
		Editor: widget.Editor{
			SingleLine: false,
			Submit:     false,
		},
		AutoScrollBottom: false,
	}
}

func (e *TextEditor) Text() string {
	return e.Editor.Text()
}

func (e *TextEditor) SetText(text string) {
	e.Editor.SetText(text)
}

func (e *TextEditor) Layout(gtx layout.Context, th themes.Theme) layout.Dimensions {
	if e.MaxHeight > 0 {
		gtx.Constraints.Max.Y = gtx.Dp(e.MaxHeight)
	}
	for {
		ev, ok := e.Editor.Update(gtx)
		if !ok {
			break
		}

		switch ev.(type) {
		case widget.SelectEvent:
			e.Highlighted = e.Editor.SelectedText()
		case widget.ChangeEvent:
			e.Value = e.Editor.Text()
		}
	}

	gioTheme := th.Gio()
	editor := material.Editor(gioTheme, &e.Editor, e.Hint)
	editor.Color = th.Color.Text
	editor.HintColor = th.Color.TextMuted

	return layout.UniformInset(unit.Dp(10)).Layout(gtx, editor.Layout)
}

func (e *TextEditor) SelectedText() string {
	return e.Editor.SelectedText()
}

func (e *TextEditor) Append(text string) {
	current := e.Editor.Text()
	e.Editor.SetText(current + text)
	e.Editor.MoveCaret(len(e.Editor.Text()), len(e.Editor.Text()))
}

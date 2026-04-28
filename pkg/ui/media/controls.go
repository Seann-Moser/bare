package media

import (
	"fmt"
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
)

type MediaControls struct {
	PlayPause widget.Clickable
	Stop      widget.Clickable
	Seek      widget.Float
	Volume    widget.Float
}

func NewMediaControls() *MediaControls {
	return &MediaControls{
		Volume: widget.Float{
			Value: 0.8,
		},
	}
}

func (c *MediaControls) Layout(
	gtx layout.Context,
	th themes.Theme,
	p Player,
) layout.Dimensions {
	if p == nil {
		return layout.Dimensions{}
	}

	for c.PlayPause.Clicked(gtx) {
		if p.State() == StatePlaying {
			_ = p.Pause()
		} else {
			_ = p.Play()
		}
	}

	for c.Stop.Clicked(gtx) {
		_ = p.Stop()
	}

	duration := p.Duration()
	position := p.Position()

	if duration > 0 {
		c.Seek.Value = float32(position) / float32(duration)
	}

	if c.Seek.Dragging() && duration > 0 {
		target := time.Duration(float32(duration) * c.Seek.Value)
		_ = p.Seek(target)
	}

	if c.Volume.Dragging() {
		_ = p.SetVolume(c.Volume.Value)
	}

	label := "Play"
	if p.State() == StatePlaying {
		label = "Pause"
	}

	gt := th.Gio()

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(material.Button(gt, &c.PlayPause, label).Layout),
				layout.Rigid(spacerW(8)),
				layout.Rigid(material.Button(gt, &c.Stop, "Stop").Layout),
				layout.Rigid(spacerW(12)),
				layout.Rigid(material.Body2(gt, fmtDuration(position)+" / "+fmtDuration(duration)).Layout),
			)
		}),
		layout.Rigid(spacerH(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(gt, &c.Seek).Layout(gtx)
		}),
		layout.Rigid(spacerH(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(material.Body2(gt, "Volume").Layout),
				layout.Rigid(spacerW(8)),
				layout.Flexed(1, material.Slider(gt, &c.Volume).Layout),
			)
		}),
	)
}

func fmtDuration(d time.Duration) string {
	if d <= 0 {
		return "0:00"
	}

	total := int(d.Seconds())
	min := total / 60
	sec := total % 60

	return fmt.Sprintf("%d:%02d", min, sec)
}

func spacerW(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(gtx.Dp(dp), 0)}
	}
}

func spacerH(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(0, gtx.Dp(dp))}
	}
}

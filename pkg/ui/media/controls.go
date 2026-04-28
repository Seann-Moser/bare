package media

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
	uiutils "github.com/Seann-Moser/bare/pkg/ui/utils"
)

type MediaControls struct {
	PlayPause widget.Clickable
	Stop      widget.Clickable
	Seek      widget.Float
	Volume    widget.Float

	seekDragging   bool
	volumeDragging bool
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
	displayPosition := position
	seekValue := c.Seek.Value

	if c.Seek.Dragging() && duration > 0 {
		displayPosition = time.Duration(float32(duration) * seekValue)
		c.seekDragging = true
	} else if c.seekDragging && duration > 0 {
		target := time.Duration(float32(duration) * seekValue)
		_ = p.Seek(target)
		displayPosition = target
		c.seekDragging = false
	} else if duration > 0 && !c.Seek.Dragging() {
		c.Seek.Value = float32(position) / float32(duration)
	}

	if c.Volume.Dragging() {
		c.volumeDragging = true
	} else if c.volumeDragging {
		_ = p.SetVolume(c.Volume.Value)
		c.volumeDragging = false
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
				layout.Rigid(uiutils.SpacerW(8)),
				layout.Rigid(material.Button(gt, &c.Stop, "Stop").Layout),
				layout.Rigid(uiutils.SpacerW(12)),
				layout.Rigid(material.Body2(gt, fmtDuration(displayPosition)+" / "+fmtDuration(duration)).Layout),
			)
		}),
		layout.Rigid(uiutils.SpacerH(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Slider(gt, &c.Seek).Layout(gtx)
		}),
		layout.Rigid(uiutils.SpacerH(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(material.Body2(gt, "Volume").Layout),
				layout.Rigid(uiutils.SpacerW(8)),
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

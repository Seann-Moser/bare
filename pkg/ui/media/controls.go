package media

import (
	"fmt"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
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

var (
	playIcon  = mustNewIcon(icons.AVPlayArrow)
	pauseIcon = mustNewIcon(icons.AVPause)
	stopIcon  = mustNewIcon(icons.AVStop)
)

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

	gt := th.Gio()
	playPauseIcon := playIcon
	playPauseLabel := "Play"
	if p.State() == StatePlaying {
		playPauseIcon = pauseIcon
		playPauseLabel = "Pause"
	}
	playPauseBtn := material.IconButton(gt, &c.PlayPause, playPauseIcon, playPauseLabel)
	playPauseBtn.Size = unit.Dp(20)
	playPauseBtn.Inset = layout.UniformInset(unit.Dp(14))

	stopBtn := material.IconButton(gt, &c.Stop, stopIcon, "Stop")
	stopBtn.Size = unit.Dp(20)
	stopBtn.Inset = layout.UniformInset(unit.Dp(14))

	seekSlider := material.Slider(gt, &c.Seek)
	seekSlider.FingerSize = unit.Dp(42)

	volumeSlider := material.Slider(gt, &c.Volume)
	volumeSlider.FingerSize = unit.Dp(38)

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(playPauseBtn.Layout),
				layout.Rigid(uiutils.SpacerW(8)),
				layout.Rigid(stopBtn.Layout),
				layout.Rigid(uiutils.SpacerW(12)),
				layout.Rigid(material.Body2(gt, fmtDuration(displayPosition)+" / "+fmtDuration(duration)).Layout),
			)
		}),
		layout.Rigid(uiutils.SpacerH(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(4),
				Bottom: unit.Dp(4),
			}.Layout(gtx, seekSlider.Layout)
		}),
		layout.Rigid(uiutils.SpacerH(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(material.Body2(gt, "Volume").Layout),
				layout.Rigid(uiutils.SpacerW(8)),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    unit.Dp(2),
						Bottom: unit.Dp(2),
					}.Layout(gtx, volumeSlider.Layout)
				}),
			)
		}),
	)
}

func mustNewIcon(data []byte) *widget.Icon {
	ic, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return ic
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

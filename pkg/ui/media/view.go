package media

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
	uiutils "github.com/Seann-Moser/bare/pkg/ui/utils"
)

type MediaView struct {
	Kind Kind
	Path string

	Image       ImageView
	Poster      ImageView
	Player      Player
	VideoPlayer *InlineVideoPlayer
	Controls    *MediaControls
	LastError   error
}

func NewMediaView(player Player) *MediaView {
	return &MediaView{
		Player:      player,
		VideoPlayer: NewInlineVideoPlayer(),
		Controls:    NewMediaControls(),
	}
}

func (v *MediaView) Load(kind Kind, path string) error {
	if v.Path != "" && v.Path != path {
		v.Close()
	}

	v.Kind = kind
	v.Path = path

	switch kind {
	case KindImage:
		v.LastError = nil
		return v.Image.Load(path)

	case KindAudio:
		if v.Player == nil {
			v.LastError = fmt.Errorf("no media player backend configured")
			return v.LastError
		}
		err := v.Player.Load(path)
		v.LastError = err
		return err

	case KindVideo:
		if v.VideoPlayer == nil {
			v.VideoPlayer = NewInlineVideoPlayer()
		}
		err := v.VideoPlayer.Load(path)
		v.LastError = err
		return err

	case KindDocument:
		v.LastError = nil
		return nil

	default:
		v.LastError = fmt.Errorf("unsupported media kind %q", kind)
		return v.LastError
	}
}

func (v *MediaView) Close() error {
	var err error

	if v.VideoPlayer != nil {
		if closeErr := v.VideoPlayer.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	if v.Player != nil {
		if stopErr := v.Player.Stop(); stopErr != nil && err == nil {
			err = stopErr
		}
	}

	v.Kind = ""
	v.Path = ""
	v.LastError = err
	return err
}

func (v *MediaView) Layout(gtx layout.Context, th themes.Theme) layout.Dimensions {
	switch v.Kind {
	case KindImage:
		if v.Image.Loading() {
			return material.Body2(th.Gio(), "Loading image preview...").Layout(gtx)
		}
		return v.Image.Draw(gtx)

	case KindVideo:
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if v.VideoPlayer != nil {
					return v.VideoPlayer.Layout(gtx)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(uiutils.SpacerH(12)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return v.Controls.Layout(gtx, th, v.VideoPlayer)
			}),
		)

	case KindAudio:
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th.Gio(), "Audio Preview").Layout(gtx)
			}),
			layout.Rigid(uiutils.SpacerH(12)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return v.Controls.Layout(gtx, th, v.Player)
			}),
		)

	case KindDocument:
		return layout.Dimensions{}

	default:
		return layout.Dimensions{}
	}
}

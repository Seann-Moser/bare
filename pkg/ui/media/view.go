package media

import (
	"context"
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
)

type MediaView struct {
	Kind Kind
	Path string

	Image     ImageView
	Poster    ImageView
	Player    Player
	Controls  *MediaControls
	LastError error
}

func NewMediaView(player Player) *MediaView {
	return &MediaView{
		Player:   player,
		Controls: NewMediaControls(),
	}
}

func (v *MediaView) Load(kind Kind, path string) error {
	v.Kind = kind
	v.Path = path

	switch kind {
	case KindImage:
		v.LastError = nil
		return v.Image.Load(path)

	case KindAudio, KindVideo:
		if v.Player == nil {
			v.LastError = fmt.Errorf("no media player backend configured")
			return v.LastError
		}
		if kind == KindVideo {
			if thumb, err := extractVideoThumbnail(context.Background(), path); err == nil {
				_ = v.Poster.Load(thumb)
			}
		}
		err := v.Player.Load(path)
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
				if v.Poster.img != nil {
					return v.Poster.Draw(gtx)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(spacerH(12)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return v.Controls.Layout(gtx, th, v.Player)
			}),
		)

	case KindAudio:
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th.Gio(), "Audio Preview").Layout(gtx)
			}),
			layout.Rigid(spacerH(12)),
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

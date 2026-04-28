package media

import (
	"fmt"

	"gioui.org/layout"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
)

type MediaView struct {
	Kind Kind
	Path string

	Image    ImageView
	Player   Player
	Controls *MediaControls
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
		return v.Image.Load(path)

	case KindAudio, KindVideo:
		if v.Player == nil {
			return fmt.Errorf("no media player backend configured")
		}
		return v.Player.Load(path)

	default:
		return fmt.Errorf("unsupported media kind %q", kind)
	}
}

func (v *MediaView) Layout(gtx layout.Context, th themes.Theme) layout.Dimensions {
	switch v.Kind {
	case KindImage:
		return v.Image.Draw(gtx)

	case KindAudio, KindVideo:
		return v.Controls.Layout(gtx, th, v.Player)

	default:
		return layout.Dimensions{}
	}
}

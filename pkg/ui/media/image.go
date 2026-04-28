package media

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	_ "golang.org/x/image/webp"
)

type ImageView struct {
	Path string

	img image.Image
	op  paint.ImageOp

	loadedPath string
	err        error
}

func (v *ImageView) Load(path string) error {
	if path == v.loadedPath && v.img != nil {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		v.err = err
		return err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		v.err = err
		return err
	}

	v.Path = path
	v.img = img
	v.op = paint.NewImageOp(img)
	v.loadedPath = path
	v.err = nil

	return nil
}

func (v *ImageView) Layout(gtx layout.Context) layout.Dimensions {
	if v.Path != "" && v.Path != v.loadedPath {
		_ = v.Load(v.Path)
	}

	if v.img == nil {
		return layout.Dimensions{}
	}

	size := v.img.Bounds().Size()

	return layout.Dimensions{
		Size: image.Pt(
			min(size.X, gtx.Constraints.Max.X),
			min(size.Y, gtx.Constraints.Max.Y),
		),
	}
}

func (v *ImageView) Draw(gtx layout.Context) layout.Dimensions {
	if v.Path != "" && v.Path != v.loadedPath {
		_ = v.Load(v.Path)
	}

	if v.img == nil {
		return layout.Dimensions{}
	}

	return widget.Image{
		Src:      v.op,
		Fit:      widget.ScaleDown,
		Position: layout.Center,
		Scale:    1.0 / gtx.Metric.PxPerDp,
	}.Layout(gtx)
}

func minFloat(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

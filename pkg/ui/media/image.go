package media

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/layout"
	"gioui.org/op/paint"
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

	imgSize := v.img.Bounds().Size()
	max := gtx.Constraints.Max

	scaleX := float32(max.X) / float32(imgSize.X)
	scaleY := float32(max.Y) / float32(imgSize.Y)

	scale := minFloat(scaleX, scaleY)
	if scale > 1 {
		scale = 1
	}

	w := int(float32(imgSize.X) * scale)
	h := int(float32(imgSize.Y) * scale)

	v.op.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{
		Size: image.Pt(w, h),
	}
}

func minFloat(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

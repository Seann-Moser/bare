package icons

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

const FallbackIcon = "subway:missing"

var errIconNotCached = errors.New("icon not cached locally")

type Iconify struct {
	CacheDir string
	Client   *http.Client

	cache       map[string]*SVGIcon
	failedUntil map[string]time.Time
}

type SVGIcon struct {
	src      []byte
	op       paint.ImageOp
	imgSize  int
	imgColor color.NRGBA
}

func NewIconify() *Iconify {
	cfg, err := os.UserCacheDir()
	if err != nil {
		cfg = "."
	}

	return &Iconify{
		CacheDir: filepath.Join(cfg, "icons", "iconify"),
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache:       map[string]*SVGIcon{},
		failedUntil: map[string]time.Time{},
	}
}

func (i *Iconify) Layout(
	gtx layout.Context,
	name string, // "mdi:home"
	size unit.Dp,
	col color.NRGBA,
) layout.Dimensions {
	ic, err := i.Icon(context.Background(), name)
	if err != nil {
		return layout.Dimensions{}
	}

	return ic.Layout(gtx, size, col)
}

func (i *Iconify) EnsureFallback(ctx context.Context) {
	_, _ = i.Icon(ctx, FallbackIcon)
}

func (i *Iconify) Icon(ctx context.Context, name string) (*SVGIcon, error) {
	// already cached
	if ic, ok := i.cache[name]; ok {
		return ic, nil
	}
	if until, ok := i.failedUntil[name]; ok && time.Now().Before(until) {
		return nil, errIconNotCached
	}

	data, err := i.LoadSVG(ctx, name)
	if err != nil {
		i.failedUntil[name] = time.Now().Add(5 * time.Second)
		// fallback
		if name != FallbackIcon {
			return i.Icon(ctx, FallbackIcon)
		}
		return nil, err
	}

	ic := &SVGIcon{src: data}

	i.cache[name] = ic
	delete(i.failedUntil, name)
	return ic, nil
}

func (i *Iconify) LoadSVG(ctx context.Context, name string) ([]byte, error) {
	_ = ctx

	path, err := i.cachePath(name)
	if err != nil {
		return nil, err
	}

	// try cache first
	if data, err := os.ReadFile(path); err == nil {
		return data, nil
	}

	data, err := i.download(ctx, name)
	if err != nil {
		return nil, err
	}

	// save (best-effort)
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, data, 0o644)
	return data, nil
}

func (i *Iconify) download(ctx context.Context, name string) ([]byte, error) {
	prefix, icon, err := splitIconName(name)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"https://api.iconify.design/%s/%s.svg?height=none",
		prefix,
		icon,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := i.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download icon %q: %s", name, res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(string(data), "<svg") {
		return nil, fmt.Errorf("download icon %q: response was not SVG", name)
	}

	return data, nil
}

func (i *Iconify) cachePath(name string) (string, error) {
	prefix, icon, err := splitIconName(name)
	if err != nil {
		return "", err
	}

	return filepath.Join(i.CacheDir, prefix, icon+".svg"), nil
}

func splitIconName(name string) (prefix string, icon string, err error) {
	parts := strings.Split(name, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid icon name %q, expected prefix:name", name)
	}

	prefix = strings.TrimSpace(parts[0])
	icon = strings.TrimSpace(parts[1])

	if prefix == "" || icon == "" {
		return "", "", fmt.Errorf("invalid icon name %q", name)
	}

	return prefix, icon, nil
}

func (i *Iconify) LayoutWithSize(
	gtx layout.Context,
	name string,
	size unit.Dp,
	col color.NRGBA,
) layout.Dimensions {
	ic, err := i.Icon(context.Background(), name)
	if err != nil {
		return layout.Dimensions{}
	}

	return ic.Layout(gtx, size, col)
}

func (ic *SVGIcon) Layout(
	gtx layout.Context,
	size unit.Dp,
	col color.NRGBA,
) layout.Dimensions {
	px := gtx.Dp(size)
	if px <= 0 {
		return layout.Dimensions{}
	}

	dims := image.Pt(px, px)
	defer clip.Rect{Max: dims}.Push(gtx.Ops).Pop()

	op := ic.image(px, col)
	op.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: dims}
}

func (ic *SVGIcon) image(sz int, col color.NRGBA) paint.ImageOp {
	if sz == ic.imgSize && col == ic.imgColor {
		return ic.op
	}

	img := image.NewRGBA(image.Rect(0, 0, sz, sz))
	icon, err := oksvg.ReadIconStream(bytes.NewReader(replaceCurrentColor(ic.src, col)))
	if err != nil {
		ic.op = paint.NewImageOp(img)
		ic.imgSize = sz
		ic.imgColor = col
		return ic.op
	}

	icon.SetTarget(0, 0, float64(sz), float64(sz))
	scanner := rasterx.NewScannerGV(sz, sz, img, img.Bounds())
	dasher := rasterx.NewDasher(sz, sz, scanner)
	icon.Draw(dasher, 1)

	ic.op = paint.NewImageOp(img)
	ic.imgSize = sz
	ic.imgColor = col
	return ic.op
}

func replaceCurrentColor(src []byte, col color.NRGBA) []byte {
	hex := fmt.Sprintf("#%02x%02x%02x", col.R, col.G, col.B)
	replacer := strings.NewReplacer(
		"currentColor", hex,
		"currentcolor", hex,
	)
	return []byte(replacer.Replace(string(src)))
}

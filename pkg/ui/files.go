package ui

import (
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui/icons"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
)

type FileBrowser struct {
	Dir string

	SelectedPath string

	ShowHidden bool
	Extensions []string // example: []string{".png", ".jpg", ".mp4"}

	List layout.List

	rows map[string]*widget.Clickable

	cachedEntries    []FileEntry
	cachedDir        string
	cachedShowHidden bool
	cachedExtKey     string
}

type FileEntry struct {
	Name  string
	Path  string
	IsDir bool
}

func NewFileBrowser(dir string) *FileBrowser {
	if dir == "" {
		dir, _ = os.UserHomeDir()
	}

	return &FileBrowser{
		Dir:  dir,
		List: layout.List{Axis: layout.Vertical},
		rows: map[string]*widget.Clickable{},
	}
}

func (b *FileBrowser) Layout(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	entries, _ := b.entries()

	for _, entry := range entries {
		btn := b.row(entry.Path)

		for btn.Clicked(gtx) {
			if entry.IsDir {
				b.Dir = entry.Path
				b.SelectedPath = ""
				b.invalidateEntries()
			} else {
				b.SelectedPath = entry.Path
			}
		}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return material.Body2(th.Gio(), b.Dir).Layout(gtx)
		}),
		layout.Rigid(SpacerH(8)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return b.List.Layout(gtx, len(entries), func(gtx layout.Context, index int) layout.Dimensions {
				return b.layoutRow(gtx, th, ic, entries[index])
			})
		}),
	)
}

func (b *FileBrowser) entries() ([]FileEntry, error) {
	if b.entriesDirty() {
		entries, err := b.readEntries()
		if err != nil {
			return nil, err
		}

		b.cachedEntries = entries
		b.cachedDir = b.Dir
		b.cachedShowHidden = b.ShowHidden
		b.cachedExtKey = b.extensionsKey()
	}

	return b.cachedEntries, nil
}

func (b *FileBrowser) readEntries() ([]FileEntry, error) {
	items, err := os.ReadDir(b.Dir)
	if err != nil {
		return nil, err
	}

	entries := make([]FileEntry, 0, len(items)+1)

	if parent := filepath.Dir(b.Dir); parent != b.Dir {
		entries = append(entries, FileEntry{
			Name:  "..",
			Path:  parent,
			IsDir: true,
		})
	}

	for _, item := range items {
		name := item.Name()

		if !b.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		path := filepath.Join(b.Dir, name)
		isDir := item.IsDir()

		if !isDir && len(b.Extensions) > 0 && !b.allowedExt(path) {
			continue
		}

		entries = append(entries, FileEntry{
			Name:  name,
			Path:  path,
			IsDir: isDir,
		})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Name == ".." {
			return true
		}
		if entries[j].Name == ".." {
			return false
		}
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	return entries, nil
}

func (b *FileBrowser) entriesDirty() bool {
	return b.cachedEntries == nil ||
		b.cachedDir != b.Dir ||
		b.cachedShowHidden != b.ShowHidden ||
		b.cachedExtKey != b.extensionsKey()
}

func (b *FileBrowser) invalidateEntries() {
	b.cachedEntries = nil
}

func (b *FileBrowser) extensionsKey() string {
	if len(b.Extensions) == 0 {
		return ""
	}

	return strings.Join(b.Extensions, "\x00")
}

func (b *FileBrowser) allowedExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))

	for _, allowed := range b.Extensions {
		if ext == strings.ToLower(allowed) {
			return true
		}
	}

	return false
}

func (b *FileBrowser) row(path string) *widget.Clickable {
	if b.rows == nil {
		b.rows = map[string]*widget.Clickable{}
	}

	if b.rows[path] == nil {
		b.rows[path] = new(widget.Clickable)
	}

	return b.rows[path]
}

func (b *FileBrowser) layoutRow(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	entry FileEntry,
) layout.Dimensions {
	btn := b.row(entry.Path)

	bg := th.Color.Surface
	if entry.Path == b.SelectedPath {
		bg = th.Color.SurfaceAlt
	}

	return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return surface(gtx, bg, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(8),
				Bottom: unit.Dp(8),
				Left:   unit.Dp(10),
				Right:  unit.Dp(10),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				iconName := "mdi:file-outline"
				if entry.IsDir {
					iconName = "mdi:folder-outline"
				}
				if entry.Name == ".." {
					iconName = "mdi:arrow-up-bold"
				}

				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if ic == nil {
							return layout.Dimensions{}
						}

						return layout.Inset{
							Right: unit.Dp(8),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return ic.Layout(gtx, iconName, unit.Dp(20), th.Color.TextMuted)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body1(th.Gio(), entry.Name)
						lbl.Color = th.Color.Text
						return lbl.Layout(gtx)
					}),
				)
			})
		})
	})
}

func surface(gtx layout.Context, col color.NRGBA, child layout.Widget) layout.Dimensions {
	paint.FillShape(
		gtx.Ops,
		col,
		clip.Rect{Max: gtx.Constraints.Max}.Op(),
	)

	return child(gtx)
}

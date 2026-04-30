package filemanager

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/DarlingGoose/bare/pkg/ui"
	"github.com/DarlingGoose/bare/pkg/ui/icons"
	"github.com/DarlingGoose/bare/pkg/ui/media"
	uitext "github.com/DarlingGoose/bare/pkg/ui/text"
	"github.com/DarlingGoose/bare/pkg/ui/themes"
	uiutils "github.com/DarlingGoose/bare/pkg/ui/utils"
)

type FileBrowser struct {
	Dir string

	SelectedPath string
	CursorPath   string
	PathInput    widget.Editor
	HomeButton   widget.Clickable
	GoButton     widget.Clickable
	PathError    string
	ActionError  string
	SearchInput  widget.Editor
	SortDropdown ui.Dropdown
	OrderButton  widget.Clickable
	CreateInput  widget.Editor
	CreateButton widget.Clickable

	ShowHidden bool
	Extensions []string // example: []string{".png", ".jpg", ".mp4"}
	SortMode   FileSortMode
	SortDesc   bool

	List layout.List

	rows        map[string]*widget.Clickable
	deleteRows  map[string]*widget.Clickable
	sortOptions map[FileSortMode]*widget.Clickable

	cachedEntries    []FileEntry
	cachedDir        string
	cachedShowHidden bool
	cachedExtKey     string
	cachedSearch     string
	cachedSortMode   FileSortMode
	cachedSortDesc   bool

	Preview     *media.MediaView
	TextPreview *uitext.TextEditor
}

type FileEntry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

type FileSortMode string

const (
	FileSortName     FileSortMode = "name"
	FileSortModified FileSortMode = "modified"
	FileSortSize     FileSortMode = "size"
)

func NewFileBrowser(dir string) *FileBrowser {
	if dir == "" {
		dir, _ = os.Getwd()
	}

	return &FileBrowser{
		Dir:        dir,
		List:       layout.List{Axis: layout.Vertical},
		rows:       map[string]*widget.Clickable{},
		deleteRows: map[string]*widget.Clickable{},
		sortOptions: map[FileSortMode]*widget.Clickable{
			FileSortName:     new(widget.Clickable),
			FileSortModified: new(widget.Clickable),
			FileSortSize:     new(widget.Clickable),
		},
		Preview: media.NewMediaView(media.NewMPVPlayer()),
		TextPreview: func() *uitext.TextEditor {
			ed := uitext.NewTextEditor("")
			ed.Editor.ReadOnly = true
			return ed
		}(),
		PathInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		SearchInput: widget.Editor{
			SingleLine: true,
		},
		CreateInput: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		SortDropdown: ui.Dropdown{
			Prefix:     "mdi:sort",
			Variant:    ui.ButtonSecondary,
			Width:      unit.Dp(180),
			MaxHeight:  unit.Dp(180),
			OffsetY:    unit.Dp(48),
			AlignRight: true,
		},
		SortMode: FileSortName,
	}
}

func (b *FileBrowser) Layout(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	b.syncPathInput()

	for {
		ev, ok := b.PathInput.Update(gtx)
		if !ok {
			break
		}
		if _, ok := ev.(widget.SubmitEvent); ok {
			b.navigateToInput()
			gtx.Execute(key.FocusCmd{})
		}
	}

	//for b.GoButton.Clicked(gtx) {
	//	b.navigateToInput()
	//}
	for b.HomeButton.Clicked(gtx) {
		b.navigateHome()
		gtx.Execute(key.FocusCmd{})
	}
	for {
		if _, ok := b.SearchInput.Update(gtx); !ok {
			break
		}
	}
	for {
		ev, ok := b.CreateInput.Update(gtx)
		if !ok {
			break
		}
		if _, ok := ev.(widget.SubmitEvent); ok {
			b.createFolder()
		}
	}
	b.SortDropdown.Update(gtx)
	for mode, btn := range b.sortOptions {
		for btn.Clicked(gtx) {
			b.SortMode = mode
			b.SortDropdown.Close()
		}
	}
	for b.OrderButton.Clicked(gtx) {
		b.SortDesc = !b.SortDesc
	}
	for b.CreateButton.Clicked(gtx) {
		b.createFolder()
	}

	entries, _ := b.entries()
	b.handleListKeys(gtx, entries)

	for _, entry := range entries {
		for b.deleteRow(entry.Path).Clicked(gtx) {
			b.deleteEntry(entry.Path)
		}

		btn := b.row(entry.Path)

		for btn.Clicked(gtx) {
			b.openEntry(entry)
			gtx.Execute(key.FocusCmd{})
		}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return b.layoutPathBar(gtx, th, ic)
		}),
		layout.Rigid(uiutils.SpacerH(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return b.layoutFilterBar(gtx, th, ic, len(entries))
		}),
		layout.Rigid(uiutils.SpacerH(8)),
		//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		//	return b.layoutActionBar(gtx, th, ic)
		//}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if b.PathError == "" {
				return layout.Dimensions{}
			}

			lbl := material.Body2(th.Gio(), b.PathError)
			lbl.Color = th.Color.Error
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if b.PathError == "" {
				return layout.Dimensions{}
			}
			return uiutils.SpacerH(8)(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if b.ActionError == "" {
				return layout.Dimensions{}
			}

			lbl := material.Body2(th.Gio(), b.ActionError)
			lbl.Color = th.Color.Error
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if b.ActionError == "" {
				return layout.Dimensions{}
			}
			return uiutils.SpacerH(8)(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if b.SelectedPath == "" {
				return b.layoutList(gtx, th, ic, entries)
			}

			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Flexed(0.55, func(gtx layout.Context) layout.Dimensions {
					return b.layoutList(gtx, th, ic, entries)
				}),
				layout.Rigid(uiutils.SpacerW(16)),
				layout.Flexed(0.45, func(gtx layout.Context) layout.Dimensions {
					return b.layoutPreview(gtx, th)
				}),
			)
		}),
	)
}

func (b *FileBrowser) layoutFilterBar(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	entryCount int,
) layout.Dimensions {
	editor := material.Editor(th.Gio(), &b.SearchInput, "Search current directory")
	editor.Color = th.Color.Text
	editor.HintColor = th.Color.TextMuted

	sortLabel := sortModeLabel(b.SortMode)
	orderLabel := "Asc"
	orderIcon := "mdi:sort-ascending"
	if b.SortDesc {
		orderLabel = "Desc"
		orderIcon = "mdi:sort-descending"
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return surface(gtx, th.Color.Surface, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(6),
					Bottom: unit.Dp(6),
					Left:   unit.Dp(10),
					Right:  unit.Dp(10),
				}.Layout(gtx, editor.Layout)
			})
		}),
		layout.Rigid(uiutils.SpacerW(12)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return b.SortDropdown.Layout(gtx, th, ic, "Sort: "+sortLabel, func(gtx layout.Context) layout.Dimensions {
				return b.layoutSortMenu(gtx, th, ic)
			})
		}),
		layout.Rigid(uiutils.SpacerW(8)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := ui.Button{
				Clickable: &b.OrderButton,
				Text:      orderLabel,
				Prefix:    orderIcon,
				Variant:   ui.ButtonSecondary,
			}
			return btn.Layout(gtx, th, ic)
		}),
		layout.Rigid(uiutils.SpacerW(12)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(th.Gio(), fmt.Sprintf("%d items", entryCount))
			lbl.Color = th.Color.TextMuted
			return lbl.Layout(gtx)
		}),
	)
}

func (b *FileBrowser) layoutSortMenu(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		b.layoutSortOption(gtx, th, ic, FileSortName),
		layout.Rigid(uiutils.SpacerH(6)),
		b.layoutSortOption(gtx, th, ic, FileSortModified),
		layout.Rigid(uiutils.SpacerH(6)),
		b.layoutSortOption(gtx, th, ic, FileSortSize),
	)
}

func (b *FileBrowser) layoutSortOption(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	mode FileSortMode,
) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		variant := ui.ButtonSecondary
		if b.SortMode == mode {
			variant = ui.ButtonPrimary
		}
		btn := ui.Button{
			Clickable: b.sortOptions[mode],
			Text:      sortModeLabel(mode),
			FullWidth: true,
			Variant:   variant,
		}
		return btn.Layout(gtx, th, ic)
	})
}

func (b *FileBrowser) layoutActionBar(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	editor := material.Editor(th.Gio(), &b.CreateInput, "New folder name")
	editor.Color = th.Color.Text
	editor.HintColor = th.Color.TextMuted

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := ui.Button{
				Clickable: &b.HomeButton,
				Prefix:    "mdi:home-outline",
				Variant:   ui.ButtonSecondary,
				Icon:      true,
				Text:      "mdi:home-outline",
			}
			return btn.Layout(gtx, th, ic)
		}),
		layout.Rigid(uiutils.SpacerW(12)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return surface(gtx, th.Color.Surface, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(6),
					Bottom: unit.Dp(6),
					Left:   unit.Dp(10),
					Right:  unit.Dp(10),
				}.Layout(gtx, editor.Layout)
			})
		}),
		layout.Rigid(uiutils.SpacerW(12)),
		//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		//	btn := Button{
		//		Clickable: &b.CreateButton,
		//		Text:      "Create Folder",
		//		Prefix:    "mdi:folder-plus-outline",
		//		Variant:   ButtonSecondary,
		//	}
		//	return btn.Layout(gtx, th, ic)
		//}),
	)
}

func (b *FileBrowser) layoutPathBar(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	editor := material.Editor(th.Gio(), &b.PathInput, "Enter file or directory path")
	editor.Color = th.Color.Text
	editor.HintColor = th.Color.TextMuted

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := ui.Button{
				Clickable: &b.HomeButton,
				Prefix:    "mdi:home-outline",
				Variant:   ui.ButtonPrimary,
				Icon:      true,
				Text:      "mdi:home-outline",
			}
			return btn.Layout(gtx, th, ic)
		}),
		layout.Rigid(uiutils.SpacerW(12)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return surface(gtx, th.Color.Surface, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(6),
					Bottom: unit.Dp(6),
					Left:   unit.Dp(10),
					Right:  unit.Dp(10),
				}.Layout(gtx, editor.Layout)
			})
		}),
		//	layout.Rigid(SpacerW(12)),
		//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		//	btn := Button{
		//		Clickable: &b.GoButton,
		//		Text:      "Go",
		//		Prefix:    "mdi:arrow-right",
		//		Variant:   ButtonSecondary,
		//	}
		//	return btn.Layout(gtx, th, ic)
		//}),
	)
}

func (b *FileBrowser) layoutList(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	entries []FileEntry,
) layout.Dimensions {
	return b.List.Layout(gtx, len(entries), func(gtx layout.Context, index int) layout.Dimensions {
		return b.layoutRow(gtx, th, ic, entries[index])
	})
}

func (b *FileBrowser) entries() ([]FileEntry, error) {
	if b.entriesDirty() {
		entries, err := b.readEntries()
		if err != nil {
			return nil, err
		}

		entries = b.filterEntries(entries)
		b.sortEntries(entries)
		b.cachedEntries = entries
		b.cachedDir = b.Dir
		b.cachedShowHidden = b.ShowHidden
		b.cachedExtKey = b.extensionsKey()
		b.cachedSearch = b.searchQuery()
		b.cachedSortMode = b.SortMode
		b.cachedSortDesc = b.SortDesc
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
		info, _ := item.Info()

		if !isDir && len(b.Extensions) > 0 && !b.allowedExt(path) {
			continue
		}

		var size int64
		var modTime time.Time
		if info != nil {
			size = info.Size()
			modTime = info.ModTime()
		}

		entries = append(entries, FileEntry{
			Name:    name,
			Path:    path,
			IsDir:   isDir,
			Size:    size,
			ModTime: modTime,
		})
	}

	return entries, nil
}

func (b *FileBrowser) entriesDirty() bool {
	return b.cachedEntries == nil ||
		b.cachedDir != b.Dir ||
		b.cachedShowHidden != b.ShowHidden ||
		b.cachedExtKey != b.extensionsKey() ||
		b.cachedSearch != b.searchQuery() ||
		b.cachedSortMode != b.SortMode ||
		b.cachedSortDesc != b.SortDesc
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

func (b *FileBrowser) deleteRow(path string) *widget.Clickable {
	if b.deleteRows == nil {
		b.deleteRows = map[string]*widget.Clickable{}
	}

	if b.deleteRows[path] == nil {
		b.deleteRows[path] = new(widget.Clickable)
	}

	return b.deleteRows[path]
}

func (b *FileBrowser) layoutRow(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	entry FileEntry,
) layout.Dimensions {
	btn := b.row(entry.Path)

	bg := th.Color.Surface
	if entry.Path == b.CursorPath || (b.CursorPath == "" && entry.Path == b.SelectedPath) {
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
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if entry.Name == ".." {
							return layout.Dimensions{}
						}

						return layout.Inset{
							Left: unit.Dp(8),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							btn := ui.Button{
								Clickable: b.deleteRow(entry.Path),
								Prefix:    "mdi:trash-can-outline",
								Variant:   ui.ButtonGhost,
								Icon:      true,
								Text:      "mdi:trash-can-outline",
							}
							return btn.Layout(gtx, th, ic)
						})
					}),
				)
			})
		})
	})
}

func (b *FileBrowser) loadPreview(path string) {
	if b.Preview == nil {
		return
	}

	if isTextPreviewable(path) {
		b.clearMediaPreview()
		b.loadTextPreview(path)
		return
	}

	kind, ok := mediaKind(path)
	if !ok {
		b.clearMediaPreview()
		if b.TextPreview != nil {
			b.TextPreview.SetText("")
		}
		return
	}

	if b.TextPreview != nil {
		b.TextPreview.SetText("")
	}
	_ = b.Preview.Load(kind, path)
}

func (b *FileBrowser) clearMediaPreview() {
	if b.Preview != nil {
		_ = b.Preview.Close()
	}
}

func (b *FileBrowser) handleListKeys(gtx layout.Context, entries []FileEntry) {
	if gtx.Focused(&b.PathInput) || gtx.Focused(&b.SearchInput) || gtx.Focused(&b.CreateInput) {
		return
	}

	for {
		ev, ok := gtx.Event(
			key.Filter{Name: key.NameUpArrow},
			key.Filter{Name: key.NameDownArrow},
			key.Filter{Name: key.NameLeftArrow},
			key.Filter{Name: key.NameRightArrow},
			key.Filter{Name: key.NameReturn},
			key.Filter{Name: key.NameEnter},
			key.Filter{Name: key.NameEscape},
		)
		if !ok {
			break
		}

		ke, ok := ev.(key.Event)
		if !ok || ke.State != key.Press {
			continue
		}

		switch ke.Name {
		case key.NameUpArrow, key.NameLeftArrow:
			b.moveCursor(entries, -1)
		case key.NameDownArrow, key.NameRightArrow:
			b.moveCursor(entries, 1)
		case key.NameReturn, key.NameEnter:
			if entry, ok := b.currentEntry(entries); ok {
				b.openEntry(entry)
			}
		case key.NameEscape:
			b.navigateUp()
		}
	}
}

func (b *FileBrowser) moveCursor(entries []FileEntry, delta int) {
	if len(entries) == 0 {
		return
	}

	current := b.cursorIndex(entries)
	if current < 0 {
		if delta < 0 {
			current = len(entries)
		} else {
			current = -1
		}
	}

	next := current + delta
	if next < 0 {
		next = 0
	}
	if next >= len(entries) {
		next = len(entries) - 1
	}

	b.CursorPath = entries[next].Path
}

func (b *FileBrowser) cursorIndex(entries []FileEntry) int {
	target := b.CursorPath
	if target == "" {
		target = b.SelectedPath
	}

	for i, entry := range entries {
		if entry.Path == target {
			return i
		}
	}

	return -1
}

func (b *FileBrowser) currentEntry(entries []FileEntry) (FileEntry, bool) {
	idx := b.cursorIndex(entries)
	if idx < 0 || idx >= len(entries) {
		return FileEntry{}, false
	}
	return entries[idx], true
}

func (b *FileBrowser) loadTextPreview(path string) {
	if b.TextPreview == nil {
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		b.TextPreview.SetText("Unable to read file: " + err.Error())
		return
	}

	const maxPreviewBytes = 128 * 1024
	preview := data
	truncated := false
	if len(preview) > maxPreviewBytes {
		preview = preview[:maxPreviewBytes]
		truncated = true
	}

	text := string(preview)
	if truncated {
		text += "\n\n...[preview truncated]..."
	}
	b.TextPreview.SetText(text)
}

func (b *FileBrowser) navigateToInput() {
	path := strings.TrimSpace(b.PathInput.Text())
	if path == "" {
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		b.PathError = err.Error()
		return
	}

	b.PathError = ""
	if info.IsDir() {
		b.clearMediaPreview()
		b.Dir = path
		b.SelectedPath = ""
		b.CursorPath = ""
		b.invalidateEntries()
		return
	}

	b.Dir = filepath.Dir(path)
	b.SelectedPath = path
	b.invalidateEntries()
	b.loadPreview(path)
}

func (b *FileBrowser) openEntry(entry FileEntry) {
	b.CursorPath = entry.Path
	b.ActionError = ""
	if entry.IsDir {
		b.clearMediaPreview()
		b.Dir = entry.Path
		b.SelectedPath = ""
		b.PathError = ""
		b.PathInput.SetText(entry.Path)
		b.CursorPath = ""
		b.invalidateEntries()
		return
	}

	b.SelectedPath = entry.Path
	b.PathError = ""
	b.PathInput.SetText(entry.Path)
	b.loadPreview(entry.Path)
}

func (b *FileBrowser) navigateHome() {
	home, err := os.UserHomeDir()
	if err != nil {
		b.PathError = err.Error()
		return
	}

	b.clearMediaPreview()
	b.Dir = home
	b.SelectedPath = ""
	b.CursorPath = ""
	b.PathError = ""
	b.ActionError = ""
	b.PathInput.SetText(home)
	b.invalidateEntries()
}

func (b *FileBrowser) navigateUp() {
	parent := filepath.Dir(b.Dir)
	if parent == b.Dir {
		return
	}

	b.Dir = parent
	b.clearMediaPreview()
	b.SelectedPath = ""
	b.CursorPath = ""
	b.PathError = ""
	b.ActionError = ""
	b.PathInput.SetText(parent)
	b.invalidateEntries()
}

func (b *FileBrowser) createFolder() {
	name := strings.TrimSpace(b.CreateInput.Text())
	if name == "" {
		return
	}

	path := name
	if !filepath.IsAbs(path) {
		path = filepath.Join(b.Dir, name)
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		b.ActionError = err.Error()
		return
	}

	b.ActionError = ""
	b.CreateInput.SetText("")
	b.invalidateEntries()
}

func (b *FileBrowser) deleteEntry(path string) {
	if path == "" || path == b.Dir || filepath.Base(path) == ".." {
		return
	}

	info, err := os.Stat(path)
	if err != nil {
		b.ActionError = err.Error()
		return
	}

	if info.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		b.ActionError = err.Error()
		return
	}

	if b.SelectedPath == path {
		b.SelectedPath = ""
	}
	if b.CursorPath == path {
		b.CursorPath = ""
	}
	if b.Preview != nil && b.Preview.Path == path {
		b.clearMediaPreview()
	}
	if b.TextPreview != nil {
		b.TextPreview.SetText("")
	}

	b.ActionError = ""
	b.PathError = ""
	b.PathInput.SetText(b.Dir)
	b.invalidateEntries()
}

func (b *FileBrowser) syncPathInput() {
	target := b.Dir
	if b.SelectedPath != "" {
		target = b.SelectedPath
	}

	if b.PathInput.Text() != target && b.PathInput.Text() == "" {
		b.PathInput.SetText(target)
	}
}

func (b *FileBrowser) layoutPreview(
	gtx layout.Context,
	th themes.Theme,
) layout.Dimensions {
	info, err := os.Stat(b.SelectedPath)

	return surface(gtx, th.Color.Surface, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			children := []layout.FlexChild{
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.H6(th.Gio(), filepath.Base(b.SelectedPath))
					lbl.Color = th.Color.Text
					return lbl.Layout(gtx)
				}),
				layout.Rigid(uiutils.SpacerH(8)),
			}

			if err == nil {
				children = append(children,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(th.Gio(), b.previewSummary())
						lbl.Color = th.Color.TextMuted
						return lbl.Layout(gtx)
					}),
					layout.Rigid(uiutils.SpacerH(16)),
				)
			}

			if b.Preview != nil && b.Preview.Path == b.SelectedPath && b.Preview.Kind != "" {
				children = append(children,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							if b.Preview.Kind == media.KindDocument {
								return b.layoutDocumentPreview(gtx, th)
							}
							return b.Preview.Layout(gtx, th)
						})
					}),
					layout.Rigid(uiutils.SpacerH(16)),
				)
			} else if isTextPreviewable(b.SelectedPath) && b.TextPreview != nil {
				children = append(children,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return surface(gtx, th.Color.SurfaceAlt, func(gtx layout.Context) layout.Dimensions {
							return b.TextPreview.Layout(gtx, th)
						})
					}),
					layout.Rigid(uiutils.SpacerH(16)),
				)
			}

			children = append(children,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return b.layoutMetadata(gtx, th, info, err)
				}),
			)

			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx, children...)
		})
	})
}

func (b *FileBrowser) previewSummary() string {
	if kind, ok := mediaKind(b.SelectedPath); ok {
		switch kind {
		case media.KindImage:
			return "Image preview"
		case media.KindAudio:
			return "Audio preview"
		case media.KindVideo:
			return "Video preview"
		case media.KindDocument:
			return "Document preview"
		}
	}

	if isTextPreviewable(b.SelectedPath) {
		return "Text preview"
	}

	return fileTypeLabel(b.SelectedPath)
}

func (b *FileBrowser) layoutMetadata(
	gtx layout.Context,
	th themes.Theme,
	info os.FileInfo,
	err error,
) layout.Dimensions {
	rows := []string{
		"Path: " + b.SelectedPath,
	}

	if err == nil {
		rows = append(rows,
			"Type: "+fileTypeLabel(b.SelectedPath),
			"Size: "+formatBytes(info.Size()),
			"Modified: "+info.ModTime().Format("2006-01-02 15:04:05"),
		)
	} else {
		rows = append(rows, "Info: unavailable")
	}

	children := make([]layout.FlexChild, 0, len(rows)*2)
	for idx, row := range rows {
		if idx > 0 {
			children = append(children, layout.Rigid(uiutils.SpacerH(6)))
		}
		line := row
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(th.Gio(), line)
			lbl.Color = th.Color.Text
			return lbl.Layout(gtx)
		}))
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx, children...)
}

func (b *FileBrowser) layoutDocumentPreview(
	gtx layout.Context,
	th themes.Theme,
) layout.Dimensions {
	return surface(gtx, th.Color.SurfaceAlt, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.H6(th.Gio(), "PDF Document")
					lbl.Color = th.Color.Text
					return lbl.Layout(gtx)
				}),
				layout.Rigid(uiutils.SpacerH(8)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th.Gio(), "Inline PDF page rendering is not available yet.")
					lbl.Color = th.Color.TextMuted
					return lbl.Layout(gtx)
				}),
			)
		})
	})
}

func mediaKind(path string) (media.Kind, bool) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp":
		return media.KindImage, true
	case ".mp3", ".wav", ".ogg", ".oog", ".m4a", ".flac":
		return media.KindAudio, true
	case ".mp4", ".mov", ".mkv", ".webm", ".avi":
		return media.KindVideo, true
	case ".pdf":
		return media.KindDocument, true
	default:
		return "", false
	}
}

func isTextPreviewable(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".txt", ".md", ".json", ".yaml", ".yml", ".toml", ".ini", ".log", ".csv", ".xml":
		return true
	default:
		return false
	}
}

func fileTypeLabel(path string) string {
	if kind, ok := mediaKind(path); ok {
		switch kind {
		case media.KindImage:
			return "Image"
		case media.KindAudio:
			return "Audio"
		case media.KindVideo:
			return "Video"
		case media.KindDocument:
			return "PDF document"
		}
	}

	if isTextPreviewable(path) {
		return "Text document"
	}

	ext := strings.TrimPrefix(strings.ToUpper(filepath.Ext(path)), ".")
	if ext == "" {
		return "File"
	}
	return ext + " file"
}

func formatBytes(size int64) string {
	const base = 1024
	if size < base {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(base), 0
	for n := size / base; n >= base; n /= base {
		div *= base
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

func (b *FileBrowser) searchQuery() string {
	return strings.TrimSpace(strings.ToLower(b.SearchInput.Text()))
}

func (b *FileBrowser) filterEntries(entries []FileEntry) []FileEntry {
	query := b.searchQuery()
	if query == "" {
		return entries
	}

	filtered := make([]FileEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.Name == ".." || strings.Contains(strings.ToLower(entry.Name), query) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (b *FileBrowser) sortEntries(entries []FileEntry) {
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

		less := false
		switch b.SortMode {
		case FileSortModified:
			if entries[i].ModTime.Equal(entries[j].ModTime) {
				less = naturalLess(entries[i].Name, entries[j].Name)
			} else {
				less = entries[i].ModTime.Before(entries[j].ModTime)
			}
		case FileSortSize:
			if entries[i].Size == entries[j].Size {
				less = naturalLess(entries[i].Name, entries[j].Name)
			} else {
				less = entries[i].Size < entries[j].Size
			}
		default:
			less = naturalLess(entries[i].Name, entries[j].Name)
		}

		if b.SortDesc {
			return !less
		}
		return less
	})
}

func naturalLess(a, b string) bool {
	aParts := splitNaturalParts(strings.ToLower(a))
	bParts := splitNaturalParts(strings.ToLower(b))

	limit := len(aParts)
	if len(bParts) < limit {
		limit = len(bParts)
	}

	for i := 0; i < limit; i++ {
		ap := aParts[i]
		bp := bParts[i]

		if ap.numeric && bp.numeric {
			if ap.number != bp.number {
				return ap.number < bp.number
			}
			if ap.text != bp.text {
				return len(ap.text) < len(bp.text)
			}
			continue
		}

		if ap.text != bp.text {
			return ap.text < bp.text
		}
	}

	if len(aParts) != len(bParts) {
		return len(aParts) < len(bParts)
	}

	return strings.ToLower(a) < strings.ToLower(b)
}

type naturalPart struct {
	text    string
	number  uint64
	numeric bool
}

func splitNaturalParts(s string) []naturalPart {
	if s == "" {
		return nil
	}

	parts := make([]naturalPart, 0, 8)
	for i := 0; i < len(s); {
		start := i
		isDigit := unicode.IsDigit(rune(s[i]))
		for i < len(s) && unicode.IsDigit(rune(s[i])) == isDigit {
			i++
		}

		part := s[start:i]
		if isDigit {
			n, err := strconv.ParseUint(part, 10, 64)
			if err == nil {
				parts = append(parts, naturalPart{text: part, number: n, numeric: true})
				continue
			}
		}

		parts = append(parts, naturalPart{text: part})
	}

	return parts
}

func sortModeLabel(mode FileSortMode) string {
	switch mode {
	case FileSortModified:
		return "Modified"
	case FileSortSize:
		return "Size"
	default:
		return "Name"
	}
}

func surface(gtx layout.Context, col color.NRGBA, child layout.Widget) layout.Dimensions {
	return uiutils.Surface(gtx, col, child)
}

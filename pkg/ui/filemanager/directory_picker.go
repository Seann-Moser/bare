package filemanager

import (
	"os"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/DarlingGoose/bare/pkg/ui"
	"github.com/DarlingGoose/bare/pkg/ui/icons"
	"github.com/DarlingGoose/bare/pkg/ui/themes"
	uiutils "github.com/DarlingGoose/bare/pkg/ui/utils"
)

type DirectoryPicker struct {
	Modal   ui.Modal
	Browser *FileBrowser

	Title       string
	SelectedDir string
	OnPick      func(string)

	PickButton   widget.Clickable
	CancelButton widget.Clickable
}

func NewDirectoryPicker(dir string) *DirectoryPicker {
	browser := NewFileBrowser(dir)
	browser.DirectoriesOnly = true
	browser.ShowDelete = false
	browser.ShowPreview = false

	return &DirectoryPicker{
		Browser: browser,
		Title:   "Select Folder",
		Modal: ui.Modal{
			CloseOnScrim: true,
		},
	}
}

func (p *DirectoryPicker) Open(dir string) {
	p.ensureBrowser()

	if strings.TrimSpace(dir) != "" {
		p.setDir(dir)
	}

	p.Modal.Open = true
}

func (p *DirectoryPicker) Close() {
	p.Modal.Open = false
}

func (p *DirectoryPicker) IsOpen() bool {
	return p.Modal.Open
}

func (p *DirectoryPicker) Layout(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) (string, bool, layout.Dimensions) {
	p.ensureBrowser()

	var pickedDir string
	var picked bool

	for p.PickButton.Clicked(gtx) {
		pickedDir = p.Browser.Dir
		p.SelectedDir = pickedDir
		if p.OnPick != nil {
			p.OnPick(pickedDir)
		}
		p.Modal.Open = false
		picked = true
	}

	for p.CancelButton.Clicked(gtx) {
		p.Modal.Open = false
	}

	title := p.Title
	if title == "" {
		title = "Select Folder"
	}

	dims := p.Modal.Layout(gtx, th, title, func(gtx layout.Context) layout.Dimensions {
		return p.layoutContent(gtx, th, ic)
	})

	return pickedDir, picked, dims
}

func (p *DirectoryPicker) layoutContent(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	gtx.Constraints.Min.Y = min(gtx.Constraints.Max.Y, gtx.Dp(unit.Dp(460)))

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body2(th.Gio(), "Current folder: "+p.Browser.Dir)
			lbl.Color = th.Color.TextMuted
			return lbl.Layout(gtx)
		}),
		layout.Rigid(uiutils.SpacerH(12)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.Browser.Layout(gtx, th, ic)
		}),
		layout.Rigid(uiutils.SpacerH(12)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := ui.Button{
						Clickable: &p.CancelButton,
						Text:      "Cancel",
						Variant:   ui.ButtonSecondary,
					}
					return btn.Layout(gtx, th, ic)
				}),
				layout.Rigid(uiutils.SpacerW(8)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := ui.Button{
						Clickable: &p.PickButton,
						Text:      "Select Folder",
						Prefix:    "mdi:folder-check-outline",
						Variant:   ui.ButtonPrimary,
					}
					return btn.Layout(gtx, th, ic)
				}),
			)
		}),
	)
}

func (p *DirectoryPicker) ensureBrowser() {
	if p.Browser != nil {
		return
	}

	p.Browser = NewFileBrowser("")
	p.Browser.DirectoriesOnly = true
	p.Browser.ShowDelete = false
	p.Browser.ShowPreview = false
}

func (p *DirectoryPicker) setDir(dir string) {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return
	}

	p.Browser.clearMediaPreview()
	p.Browser.Dir = dir
	p.Browser.SelectedPath = ""
	p.Browser.CursorPath = ""
	p.Browser.PathError = ""
	p.Browser.ActionError = ""
	p.Browser.PathInput.SetText(dir)
	p.Browser.invalidateEntries()
}

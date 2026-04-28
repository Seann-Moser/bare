package dashboard

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Seann-Moser/bare/pkg/ui"
	"github.com/Seann-Moser/bare/pkg/ui/filemanager"
	"github.com/Seann-Moser/bare/pkg/ui/icons"
	"github.com/Seann-Moser/bare/pkg/ui/text"
	"github.com/Seann-Moser/bare/pkg/ui/themes"
	"github.com/Seann-Moser/bare/pkg/ui/utils"
)

type Dashboard struct {
	SidebarTabs *ui.Tabs
	ContentTabs *ui.Tabs

	ThemeSelector *themes.ThemeSelector
	FileBrowser   *filemanager.FileBrowser
	TextView      *text.TextView

	OpenSettings widget.Clickable
	Refresh      widget.Clickable
	Settings     ui.Modal
	Shell        ui.AppShell
}

func NewDashboard() *Dashboard {
	tv := text.NewTextView()
	tv.Append(text.TextBlock{Kind: text.TextHeading, Text: "System Log"})
	tv.Append(text.TextBlock{Kind: text.TextMuted, Text: "Dashboard initialized"})
	tv.Append(text.TextBlock{Kind: text.TextParagraph, Text: "Ready."})

	sidebarTabs := ui.NewTabs([]ui.TabItem{
		{ID: "overview", Label: "Overview", Icon: "mdi:view-dashboard"},
		{ID: "files", Label: "Files", Icon: "mdi:folder"},
		{ID: "media", Label: "Media", Icon: "mdi:play-box"},
		{ID: "logs", Label: "Logs", Icon: "mdi:text-box"},
		{ID: "settings", Label: "Settings", Icon: "mdi:cog"},
	}, "overview")
	sidebarTabs.Axis = layout.Vertical

	return &Dashboard{
		SidebarTabs: sidebarTabs,

		ContentTabs: ui.NewTabs([]ui.TabItem{
			{ID: "summary", Label: "Summary", Icon: "mdi:chart-line"},
			{ID: "activity", Label: "Activity", Icon: "mdi:timeline"},
			{ID: "details", Label: "Details", Icon: "mdi:information"},
		}, "summary"),

		ThemeSelector: themes.NewThemeSelector(),
		FileBrowser:   filemanager.NewFileBrowser(""),
		TextView:      tv,

		Settings: ui.Modal{
			CloseOnScrim: true,
		},
		Shell: ui.AppShell{
			SidebarWidth: unit.Dp(240),
		},
	}
}

func (d *Dashboard) Layout(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	systemDark bool,
) (themes.Theme, layout.Dimensions) {
	for d.OpenSettings.Clicked(gtx) {
		d.Settings.Open = true
	}
	//
	for d.Refresh.Clicked(gtx) {
		d.TextView.Append(text.TextBlock{
			Kind: text.TextMuted,
			Text: "Refresh clicked",
		})
	}

	dims := d.Shell.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			return d.sidebar(gtx, th, ic)
		},
		func(gtx layout.Context) layout.Dimensions {
			var mainDims layout.Dimensions
			th, mainDims = d.main(gtx, th, ic, systemDark)
			return mainDims
		},
		func(gtx layout.Context) layout.Dimensions {
			return d.Settings.Layout(gtx, th, "Settings", func(gtx layout.Context) layout.Dimensions {
				var dims layout.Dimensions
				th, dims = d.ThemeSelector.LayoutThemeSelector(gtx, th, systemDark)
				return dims
			})
		},
	)

	return th, dims
}

func (d *Dashboard) sidebar(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	return utils.Panel(gtx, th.Color.SurfaceAlt, unit.Dp(th.Radius.LG), func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return ui.LayoutSidebarTitle(gtx, th, "Gio Dashboard")
				}),
				layout.Rigid(utils.SpacerH(20)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return d.SidebarTabs.Layout(gtx, th, ic)
				}),
			)
		})
	})
}

func (d *Dashboard) main(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
	systemDark bool,
) (themes.Theme, layout.Dimensions) {
	sidebarSelection := d.SidebarTabs.Selected()

	dims := layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		children := []layout.FlexChild{
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.topbar(gtx, th, ic)
			}),
		}

		if sidebarSelection == "overview" {
			children = append(children,
				layout.Rigid(utils.SpacerH(20)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return d.ContentTabs.Layout(gtx, th, ic)
				}),
			)
		}

		children = append(children,
			layout.Rigid(utils.SpacerH(20)),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				switch sidebarSelection {
				case "files":
					return d.FileBrowser.Layout(gtx, th, ic)
				case "logs":
					return d.TextView.Layout(gtx, th)
				case "settings":
					var dims layout.Dimensions
					th, dims = d.ThemeSelector.LayoutThemeSelector(gtx, th, systemDark)
					return dims
				default:
					return d.content(gtx, th, ic)
				}
			}),
		)

		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx, children...)
	})

	return th, dims
}

func (d *Dashboard) content(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	switch d.ContentTabs.Selected() {
	case "activity":
		return d.activity(gtx, th, ic)
	case "details":
		return d.details(gtx, th, ic)
	default:
		return d.overview(gtx, th, ic)
	}
}

func (d *Dashboard) topbar(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	return ui.LayoutTopbar(gtx, th, ic, "Dashboard", []ui.TopbarAction{
		{
			Clickable: &d.Refresh,
			Text:      "Refresh",
			Prefix:    "mdi:refresh",
			Variant:   ui.ButtonSecondary,
		},
		{
			Clickable: &d.OpenSettings,
			Text:      "Settings",
			Prefix:    "mdi:cog",
			Variant:   ui.ButtonPrimary,
		},
	})
}

func (d *Dashboard) overview(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Flexed(1, statCard(th, ic, "Files", "128", "mdi:folder")),
				layout.Rigid(utils.SpacerW(12)),
				layout.Flexed(1, statCard(th, ic, "Media", "42", "mdi:play-box")),
				layout.Rigid(utils.SpacerW(12)),
				layout.Flexed(1, statCard(th, ic, "Logs", "9.2k", "mdi:text-box")),
			)
		}),
		layout.Rigid(utils.SpacerH(16)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return utils.Panel(gtx, th.Color.Surface, unit.Dp(th.Radius.LG), func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return d.TextView.Layout(gtx, th)
				})
			})
		}),
	)
}

func (d *Dashboard) activity(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	return utils.Panel(gtx, th.Color.Surface, unit.Dp(th.Radius.LG), func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.H6(th.Gio(), "Recent Activity")
					lbl.Color = th.Color.Text
					return lbl.Layout(gtx)
				}),
				layout.Rigid(utils.SpacerH(8)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th.Gio(), "Activity feed and events will appear here.")
					lbl.Color = th.Color.TextMuted
					return lbl.Layout(gtx)
				}),
				layout.Rigid(utils.SpacerH(16)),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return d.TextView.Layout(gtx, th)
				}),
			)
		})
	})
}

func (d *Dashboard) details(
	gtx layout.Context,
	th themes.Theme,
	ic *icons.Iconify,
) layout.Dimensions {
	_ = ic

	return utils.Panel(gtx, th.Color.Surface, unit.Dp(th.Radius.LG), func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.H6(th.Gio(), "Details")
					lbl.Color = th.Color.Text
					return lbl.Layout(gtx)
				}),
				layout.Rigid(utils.SpacerH(8)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th.Gio(), "Detailed metadata and diagnostics will appear here.")
					lbl.Color = th.Color.TextMuted
					return lbl.Layout(gtx)
				}),
			)
		})
	})
}

func statCard(
	th themes.Theme,
	ic *icons.Iconify,
	title string,
	value string,
	iconName string,
) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return utils.Panel(gtx, th.Color.Surface, unit.Dp(th.Radius.LG), func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if ic == nil {
							return layout.Dimensions{}
						}
						return ic.Layout(gtx, iconName, unit.Dp(32), th.Color.Primary)
					}),
					layout.Rigid(utils.SpacerW(12)),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(th.Gio(), title)
								lbl.Color = th.Color.TextMuted
								return lbl.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.H5(th.Gio(), value)
								lbl.Color = th.Color.Text
								return lbl.Layout(gtx)
							}),
						)
					}),
				)
			})
		})
	}
}

type App struct {
	Theme      themes.Theme
	Icons      *icons.Iconify
	Dashboard  *Dashboard
	SystemDark bool

	lastSavedTheme themes.Config
}

func NewApp() *App {
	cfg, err := themes.LoadConfig()
	if err != nil {
		cfg = themes.DefaultConfig()
	}

	return &App{
		Theme:          cfg.Theme(false),
		Icons:          icons.NewIconify(),
		Dashboard:      NewDashboard(),
		lastSavedTheme: cfg,
	}
}

func (a *App) Layout(gtx layout.Context) layout.Dimensions {
	var dims layout.Dimensions

	a.Theme, dims = a.Dashboard.Layout(
		gtx,
		a.Theme,
		a.Icons,
		a.SystemDark,
	)

	currentCfg := themes.ConfigFromTheme(a.Theme)
	if currentCfg != a.lastSavedTheme {
		if err := themes.SaveConfig(currentCfg); err == nil {
			a.lastSavedTheme = currentCfg
		}
	}

	return dims
}

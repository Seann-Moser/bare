# `pkg/ui`

`pkg/ui` contains shared Gio UI components used across the app.

This README covers the main reusable components currently exposed directly from this package:

- `Button`
- `Dropdown`
- `Modal`
- `Tabs`
- `Topbar`

## `Button`

### What It Does

`Button` is a themed wrapper around `widget.Clickable` with support for:

- text buttons
- icon-only buttons
- prefix and suffix icons
- `primary`, `secondary`, and `ghost` variants
- hover styling

### Main Types

- `Button`
- `ButtonVariant`

Variants:

- `ButtonPrimary`
- `ButtonSecondary`
- `ButtonGhost`

### How To Use It

Create and keep a `widget.Clickable` on your struct, then render a `Button` in layout.

### Example

```go
type Screen struct {
	Refresh widget.Clickable
}

func (s *Screen) Layout(gtx layout.Context, th themes.Theme, ic *icons.Iconify) layout.Dimensions {
	btn := ui.Button{
		Clickable: &s.Refresh,
		Text:      "Refresh",
		Prefix:    "mdi:refresh",
		Variant:   ui.ButtonSecondary,
	}
	return btn.Layout(gtx, th, ic)
}
```

For an icon-only button:

```go
btn := ui.Button{
	Clickable: &s.Refresh,
	Icon:      true,
	Text:      "mdi:refresh",
	Variant:   ui.ButtonGhost,
}
```

## `Dropdown`

### What It Does

`Dropdown` provides a reusable overlay-style dropdown trigger and menu container.

It handles:

- toggle open/close state
- overlay positioning outside normal layout flow
- configurable width and max height
- optional right alignment
- shared button styling for the toggle

### Main Type

- `Dropdown`

Important fields:

- `Open`
- `Prefix`
- `Variant`
- `Width`
- `MaxHeight`
- `OffsetY`
- `AlignRight`

### How To Use It

Call `Update(gtx)` before layout so clicks toggle the state, then call `Layout(...)` with a label and a menu widget.

### Example

```go
type Screen struct {
	Sort ui.Dropdown
	Name widget.Clickable
	Size widget.Clickable
}

func (s *Screen) Layout(gtx layout.Context, th themes.Theme, ic *icons.Iconify) layout.Dimensions {
	s.Sort.Update(gtx)

	for s.Name.Clicked(gtx) {
		s.Sort.Close()
	}
	for s.Size.Clicked(gtx) {
		s.Sort.Close()
	}

	return s.Sort.Layout(gtx, th, ic, "Sort: Name", func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := ui.Button{Clickable: &s.Name, Text: "Name", Variant: ui.ButtonSecondary}
				return btn.Layout(gtx, th, ic)
			}),
			layout.Rigid(ui.SpacerH(6)),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := ui.Button{Clickable: &s.Size, Text: "Size", Variant: ui.ButtonSecondary}
				return btn.Layout(gtx, th, ic)
			}),
		)
	})
}
```

## `Modal`

### What It Does

`Modal` renders a full-screen scrim with centered content and a close button.

It supports:

- open/closed state
- click-to-close on the scrim
- themed modal content container

### Main Type

- `Modal`

Important fields:

- `Open`
- `CloseOnScrim`

### How To Use It

Keep a `Modal` on your screen state, toggle `Open`, then call `Layout(...)` inside a top-level `layout.Stack` or shell overlay layer.

### Example

```go
type Screen struct {
	OpenSettings widget.Clickable
	Settings     ui.Modal
}

func (s *Screen) Layout(gtx layout.Context, th themes.Theme) layout.Dimensions {
	for s.OpenSettings.Clicked(gtx) {
		s.Settings.Open = true
	}

	return s.Settings.Layout(gtx, th, "Settings", func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body1(th.Gio(), "Modal body")
		lbl.Color = th.Color.Text
		return lbl.Layout(gtx)
	})
}
```

There is also:

```go
ui.OpenPopout(title, th, content)
```

for opening content in a separate Gio window.

## `Tabs`

### What It Does

`Tabs` is a simple themed tab selector built on `widget.Clickable`.

It supports:

- horizontal or vertical layout
- text labels
- optional icons
- internal active-state tracking

### Main Types

- `Tabs`
- `TabItem`

### How To Use It

Create tabs with `NewTabs(items, activeID)`, render them with `Layout(...)`, and read the selected tab with `Selected()`.

### Example

```go
tabs := ui.NewTabs([]ui.TabItem{
	{ID: "overview", Label: "Overview", Icon: "mdi:view-dashboard"},
	{ID: "logs", Label: "Logs", Icon: "mdi:text-box"},
}, "overview")

tabs.Axis = layout.Vertical
```

In layout:

```go
_ = tabs.Layout(gtx, th, ic)

switch tabs.Selected() {
case "logs":
	// render logs
default:
	// render overview
}
```

## `Topbar`

### What It Does

The topbar helpers provide a standard title row with trailing actions.

This is useful for screens that need a consistent:

- title area
- right-aligned action buttons

### Main Types and Functions

- `TopbarAction`
- `LayoutTopbar(...)`
- `LayoutSidebarTitle(...)`

### How To Use It

Pass a title and a slice of `TopbarAction` values to `LayoutTopbar(...)`.

### Example

```go
actions := []ui.TopbarAction{
	{
		Clickable: &s.Refresh,
		Text:      "Refresh",
		Prefix:    "mdi:refresh",
		Variant:   ui.ButtonSecondary,
	},
	{
		Clickable: &s.Settings,
		Text:      "Settings",
		Prefix:    "mdi:cog",
		Variant:   ui.ButtonPrimary,
	},
}

return ui.LayoutTopbar(gtx, th, ic, "Dashboard", actions)
```

Sidebar title helper:

```go
return ui.LayoutSidebarTitle(gtx, th, "My App")
```

## Notes

- All of these components assume the shared `themes.Theme` type from `pkg/ui/themes`.
- `Button`, `Tabs`, `Dropdown`, and `Topbar` optionally use `*icons.Iconify` from `pkg/ui/icons`.
- `Dropdown` does not manage option click state for you; your screen code is still responsible for handling the option buttons and closing the menu when appropriate.

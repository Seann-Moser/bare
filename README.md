Here’s a clean, developer-friendly `README.md` you can drop into your `gio-ui` package. It documents all the components you’ve built and keeps things structured for reuse.

---

# 🧩 Gio UI Components Library

A reusable UI component library built on top of **Gio**, focused on:

* Clean theming (light/dark/system)
* Composable components
* Minimal dependencies
* Desktop-first ergonomics (Linux-friendly)

---

## ✨ Features

* 🎨 Theme system (light / dark / system)
* 🧱 Reusable UI components
* 🎯 Iconify icon support (with caching)
* 🗂 File browser
* 🧭 Tabs / navigation
* 🪟 Modal + popout windows
* 📝 Text viewer + editor
* 🔘 Buttons with icons
* 🎵 Media (image/audio/video UI layer)

---

## 📦 Installation

```bash
go get your/module/ui
```

---

## 🎨 Theme System

```go
th := theme.New(theme.ModeSystem, theme.PaletteOcean, systemDark)
gioTheme := th.Gio()
```

### Modes

* `ModeLight`
* `ModeDark`
* `ModeSystem`

### Palettes

* Sunset
* Coastal
* Sky
* Blush
* Ocean
* Pastel

---

## 🎛 Theme Selector Component

```go
selector := theme.NewThemeSelector()

th, dims := selector.LayoutThemeSelector(gtx, th, systemDark)
```

---

## 🎯 Iconify Icons

Icons auto-download + cache to:

```txt
~/.config/icons/iconify/
```

### Usage

```go
icons := icons.NewIconify()

icons.Layout(gtx, "mdi:home", unit.Dp(24), th.Color.Text)
```

### Fallback

If icon fails:

```go
subway:missing
```

---

## 🔘 Button Component

Supports:

* Prefix icon
* Suffix icon
* Icon-only mode
* Variants

```go
btn := components.Button{
	Text:    "Save",
	Prefix:  "mdi:content-save",
	Suffix:  "mdi:chevron-right",
	Variant: components.ButtonPrimary,
}

btn.Layout(gtx, th, icons)
```

---

## 🧭 Tabs / Navbar

```go
tabs := components.NewTabs([]components.TabItem{
	{ID: "home", Label: "Home", Icon: "mdi:home"},
	{ID: "settings", Label: "Settings", Icon: "mdi:cog"},
}, "home")

tabs.Layout(gtx, th, icons)

selected := tabs.Selected()
```

---

## 🪟 Modal

Overlay UI rendered in same window.

```go
modal.Open = true

modal.Layout(gtx, th, "Title", func(gtx layout.Context) layout.Dimensions {
	return material.Body1(th.Gio(), "Content").Layout(gtx)
})
```

---

## 🧱 Popout Window

Separate OS window:

```go
OpenPopout("Preview", th, func(gtx layout.Context) layout.Dimensions {
	return material.Body1(th.Gio(), "Hello").Layout(gtx)
})
```

---

## 📝 Text View (Rich Display)

Supports:

* headings
* code blocks
* muted text
* auto-scroll

```go
view := components.NewTextView()

view.SetBlocks(ParseSimpleRichText(text))
view.Layout(gtx, th)
```

---

## ✏️ Text Editor

```go
editor := components.NewTextEditor("Write something...")

editor.Layout(gtx, th)

text := editor.Text()
selected := editor.SelectedText()
```

### Notes

* Scrolls automatically when height constrained
* Supports selection tracking

---

## 📂 File Browser

Supports:

* directory navigation
* filtering
* file selection
* icons

```go
browser := components.NewFileBrowser("")

browser.Extensions = []string{".png", ".jpg", ".mp4"}

browser.Layout(gtx, th, icons)

selected := browser.SelectedPath
```

---

## 🖼 Image Viewer

```go
img := components.ImageView{}
img.Load("/path/to/image.png")

img.Draw(gtx)
```

---

## 🎵 Media Player (UI Layer)

Includes:

* Play / Pause / Stop
* Seek
* Volume

```go
controls := components.NewMediaControls()

controls.Layout(gtx, th, player)
```

### Backend

Media playback is abstracted:

```go
type Player interface {
	Load(path string)
	Play()
	Pause()
	Seek()
}
```

Recommended backends:

* mpv (best for Linux)
* ffmpeg
* gstreamer

---

## 🧩 Media View

```go
view := media.NewMediaView(player)

view.Load(media.KindVideo, "/video.mp4")
view.Layout(gtx, th)
```

---

## 🧠 Architecture Notes

### Gio Philosophy

* No implicit state
* Layout is explicit
* Scrolling must be opt-in

### This library follows:

* Stateless rendering where possible
* Explicit state structs
* Separation of:

    * UI (components)
    * logic (theme/media/icons)

---

## 🚀 Suggested Next Steps

* Theme persistence (save to disk)
* Window manager for popouts
* Virtualized lists (large file dirs / chat logs)
* Keyboard navigation support
* Drag + drop (files / UI panels)

---

## ⚠️ Limitations

* Rich text editing is **not fully supported** (viewer only)
* Media playback requires external backend
* System font loading must be implemented manually

---

## 🤝 Contributing

PRs welcome. Keep components:

* small
* composable
* dependency-light

---

## 🧵 Example App Structure

```txt
ui/
  theme/
  components/
  icons/
  media/

app/
  main.go
```

---

If you want, I can next:

* turn this into a proper Go module with examples
* generate a demo app (like Storybook for Gio)
* or build a full “dashboard UI” using all components together

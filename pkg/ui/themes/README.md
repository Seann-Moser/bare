# `pkg/ui/themes`

`pkg/ui/themes` provides theme creation, palette selection, Gio `material.Theme` integration, and simple config persistence for the app UI.

## What It Does

This package is responsible for:

- building a runtime `Theme` from a mode and palette
- exposing color, spacing, and radius tokens for UI code
- converting the theme into a Gio `*material.Theme`
- loading and saving the current theme config to `~/.config/bare/theme.yaml`
- rendering a basic theme selector UI for `system`, `light`, `dark`, and palette choices

## Main Types

- `Theme`
  - runtime theme object with `Color`, `Space`, and `Radius` tokens
- `Mode`
  - `system`, `light`, `dark`
- `PaletteName`
  - `sunset`, `coastal`, `sky`, `blush`, `ocean`, `pastel`
- `Config`
  - persisted theme settings with `Mode` and `Palette`
- `ThemeSelector`
  - stateful Gio widget for editing the current theme

## Main Functions

- `New(mode, palette, systemDark) Theme`
  - creates a theme from the requested mode and palette
- `LoadConfig() (Config, error)`
  - loads the saved config or returns defaults if the file does not exist
- `SaveConfig(cfg Config) error`
  - saves the config to `~/.config/bare/theme.yaml`
- `DefaultConfig() Config`
  - returns the package default config
- `ConfigFromTheme(th Theme) Config`
  - converts a runtime theme back into a persistable config

## Basic Usage

Create a theme directly:

```go
package main

import "github.com/DarlingGoose/bare/pkg/ui/themes"

func main() {
	systemDark := false
	th := themes.New(themes.ModeDark, themes.PaletteOcean, systemDark)

	_ = th.Color.Primary
	_ = th.Space.MD
	_ = th.Radius.LG
}
```

Use it with Gio:

```go
gioTheme := th.Gio()
```

`th.Gio()` returns a configured `*material.Theme` using the package font setup and the theme color tokens.

## Config Persistence

The package stores the active theme here:

```text
~/.config/bare/theme.yaml
```

Example file:

```yaml
mode: dark
palette: ocean
```

Load and save:

```go
cfg, err := themes.LoadConfig()
if err != nil {
	cfg = themes.DefaultConfig()
}

th := cfg.Theme(false)

cfg = themes.ConfigFromTheme(th)
_ = themes.SaveConfig(cfg)
```

## Theme Selector Example

`ThemeSelector` is stateful and should be kept on your app or screen struct.

```go
type App struct {
	Theme    themes.Theme
	Selector *themes.ThemeSelector
}

func NewApp() *App {
	cfg, _ := themes.LoadConfig()
	return &App{
		Theme:    cfg.Theme(false),
		Selector: themes.NewThemeSelector(),
	}
}

func (a *App) Layout(gtx layout.Context, systemDark bool) layout.Dimensions {
	var dims layout.Dimensions
	a.Theme, dims = a.Selector.LayoutThemeSelector(gtx, a.Theme, systemDark)
	_ = themes.SaveConfig(themes.ConfigFromTheme(a.Theme))
	return dims
}
```

## Notes

- `ModeSystem` is preserved as the chosen user mode, even though token generation resolves to light or dark using `systemDark`.
- Invalid saved values are normalized back to defaults during config load.
- `ThemeSelector.LayoutThemeSelector(...)` is the more useful selector API when the caller needs the updated `Theme` returned.

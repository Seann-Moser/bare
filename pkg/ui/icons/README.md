# `pkg/ui/icons`

`pkg/ui/icons` provides a small Gio-friendly icon client backed by Iconify SVG assets.

## What It Does

This package is responsible for:

- loading icons by Iconify name such as `mdi:home`
- caching downloaded SVG files under the user cache directory
- rasterizing SVG icons into Gio paint operations
- applying a requested icon color at render time
- falling back to a missing-icon glyph when an icon cannot be loaded

This package does **not** use Gio's `widget.Icon` / IconVG path. It downloads SVG and rasterizes it to an image before painting.

## Main Types

- `Iconify`
  - icon loader, cache manager, downloader, and renderer entrypoint
- `SVGIcon`
  - cached rasterizable icon object used internally and returned by `Icon(...)`

## Main Functions and Methods

- `NewIconify() *Iconify`
  - creates a client with a local cache directory and default HTTP client
- `(*Iconify).Layout(gtx, name, size, color)`
  - convenience helper to load and draw an icon in one call
- `(*Iconify).LayoutWithSize(gtx, name, size, color)`
  - equivalent convenience wrapper
- `(*Iconify).Icon(ctx, name)`
  - loads or resolves an icon and returns an `*SVGIcon`
- `(*Iconify).LoadSVG(ctx, name)`
  - loads raw SVG from local cache or remote Iconify
- `(*Iconify).EnsureFallback(ctx)`
  - preloads the fallback missing icon

## Icon Names

Icon names use the `prefix:name` format expected by Iconify.

Examples:

- `mdi:home`
- `mdi:cog`
- `mdi:folder`
- `subway:missing`

If the name is invalid or unavailable, the package attempts to fall back to:

```text
subway:missing
```

## Cache Behavior

Downloaded icons are stored under the user cache directory:

```text
~/.cache/icons/iconify/
```

The exact base path comes from `os.UserCacheDir()`.

For example:

```text
~/.cache/icons/iconify/mdi/home.svg
```

If a lookup fails, the package keeps a short failure cooldown so it does not retry the same bad fetch every frame.

## Basic Usage

Create a shared icon client once:

```go
ic := icons.NewIconify()
```

Render an icon inside a Gio layout function:

```go
return ic.Layout(gtx, "mdi:folder", unit.Dp(20), th.Color.TextMuted)
```

## Example

```go
package example

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/DarlingGoose/bare/pkg/ui/icons"
	"github.com/DarlingGoose/bare/pkg/ui/themes"
)

type Widget struct {
	Icons *icons.Iconify
}

func NewWidget() *Widget {
	return &Widget{
		Icons: icons.NewIconify(),
	}
}

func (w *Widget) Layout(gtx layout.Context, th themes.Theme) layout.Dimensions {
	return w.Icons.Layout(gtx, "mdi:home", unit.Dp(24), th.Color.Primary)
}
```

## Network and Rendering Notes

- First-time icon loads may require network access to `https://api.iconify.design/`.
- Once downloaded, icons are read from the local cache.
- Color is applied by replacing `currentColor` in the SVG before rasterization.
- Rasterization output is cached per icon size and color on the `SVGIcon`.

## Practical Notes

- Reuse one `Iconify` instance across the app instead of creating one per widget.
- Prefer valid Iconify IDs from one consistent icon set, such as `mdi:*`.
- If an icon silently renders nothing, the name may be invalid and the fallback may also be unavailable offline.

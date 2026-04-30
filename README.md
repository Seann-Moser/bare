# bare

`bare` is a Gio-based desktop UI workspace. It currently includes a runnable dashboard app plus a reusable set of UI packages for themes, icons, text widgets, file browsing, and media preview.

The repo is split between:

- an app entrypoint in `main.go` / `cmd`
- a dashboard example in `apps/dashboard`
- reusable UI packages under `pkg/ui`

## Current State

The most complete path in the repo today is the desktop dashboard launched with the `gui` command. That dashboard exercises the shared UI packages and includes:

- a themed app shell with sidebar and topbar
- tabs, buttons, dropdowns, and modal overlays
- a file browser with filtering, sorting, and preview panes
- text viewer and editor widgets
- image, audio, and video preview support
- persistent theme selection

Some Cobra command descriptions are still scaffold text, so the code is a more reliable source of truth than the generated command help.

## Run

Build the module:

```bash
go build ./...
```

Launch the dashboard:

```bash
go run . gui
```

You can also build the binary and run:

```bash
./bare gui
```

## Runtime Requirements

This repo is pure Go at build time, but some features expect external tools at runtime:

- `mpv` for media playback control in the file browser preview
- `ffprobe` for media duration and video metadata
- `ffmpeg` for video thumbnails
- network access on first icon load, because `pkg/ui/icons` fetches SVGs from Iconify and caches them locally

Cached and persisted data is stored in standard user directories:

- theme config: `~/.config/bare/theme.yaml`
- icon cache: `~/.cache/icons/iconify/`
- media cache and sockets: `~/.cache/bare/`

## Package Guide

### `apps/dashboard`

The dashboard is the main example app in the repo. It wires together the shared UI packages into a desktop shell with overview, file, log, media, and settings surfaces.

### `pkg/ui`

Shared Gio widgets and layout helpers used by the dashboard:

- `Button`
- `Dropdown`
- `Modal`
- `Tabs`
- `Topbar`
- `AppShell`

See [pkg/ui/README.md](/home/n9s/go/src/github.com/DarlingGoose/bare/pkg/ui/README.md).

### `pkg/ui/themes`

Theme construction, palette selection, Gio `material.Theme` integration, and config persistence.

Features:

- `light`, `dark`, and `system` modes
- palette presets: `sunset`, `coastal`, `sky`, `blush`, `ocean`, `pastel`
- theme selector widget
- config load/save helpers

See [pkg/ui/themes/README.md](/home/n9s/go/src/github.com/DarlingGoose/bare/pkg/ui/themes/README.md).

### `pkg/ui/icons`

Iconify-backed icon loading for Gio. Icons are fetched as SVG, cached locally, and rasterized for painting.

See [pkg/ui/icons/README.md](/home/n9s/go/src/github.com/DarlingGoose/bare/pkg/ui/icons/README.md).

### `pkg/ui/text`

Text-oriented widgets and helpers:

- `TextView` for structured text blocks
- `TextEditor` for multiline editing
- `SelectableTextBlock`
- `ParseSimpleRichText` for lightweight formatted text parsing

See [pkg/ui/text/README.md](/home/n9s/go/src/github.com/DarlingGoose/bare/pkg/ui/text/README.md).

### `pkg/ui/filemanager`

A file browser widget with:

- path navigation
- search and sorting
- hidden file toggling
- directory and file selection
- text and media preview

The file browser currently acts as one of the main integration points for the media and text packages.

### `pkg/ui/media`

Media preview primitives used by the file browser:

- image preview
- inline video playback
- audio/video controls
- `Player` abstraction with an `mpv`-backed implementation

## Repo Layout

```text
.
├── apps/dashboard
├── cmd
├── main.go
└── pkg
    ├── app
    └── ui
        ├── filemanager
        ├── icons
        ├── media
        ├── text
        ├── themes
        └── utils
```

## Notes For Development

- The root module path is `github.com/DarlingGoose/bare`.
- The repo vendors its dependencies under `vendor/`.
- The dashboard window title is currently `Gio Dashboard`.
- The `gui --file` flag exists, but the dashboard launch path does not currently use it.

# `pkg/ui/text`

`pkg/ui/text` provides simple text-oriented UI components for Gio-based screens in this repo.

## What It Does

This package is responsible for:

- rendering structured text blocks with basic styles
- providing a multiline text editor wrapper with theme-aware colors
- exposing a selectable text block helper
- parsing a small plain-text format into `TextBlock` slices

It is intentionally lightweight. This is not a full markdown renderer or a full text editing framework.

## Main Types

- `TextBlock`
  - one display block with a `Kind` and `Text`
- `TextBlockKind`
  - `paragraph`, `heading`, `code`, `muted`
- `TextView`
  - vertical list-based viewer for `[]TextBlock`
- `TextEditor`
  - themed multiline editor wrapper around `widget.Editor`
- `SelectableTextBlock`
  - selectable one-block text display

## Main Functions

- `NewTextView() *TextView`
- `NewTextEditor(hint string) *TextEditor`
- `ParseSimpleRichText(s string) []TextBlock`

## `TextView`

`TextView` renders a scrollable list of `TextBlock` values.

Useful methods:

- `SetBlocks(blocks []TextBlock)`
- `Append(block TextBlock)`
- `Layout(gtx, th)`

If `AutoScrollBottom` is enabled, the view automatically scrolls to the end when new blocks are appended.

### Example

```go
tv := text.NewTextView()
tv.Append(text.TextBlock{Kind: text.TextHeading, Text: "System Log"})
tv.Append(text.TextBlock{Kind: text.TextMuted, Text: "Dashboard initialized"})
tv.Append(text.TextBlock{Kind: text.TextParagraph, Text: "Ready."})
```

In layout:

```go
return tv.Layout(gtx, th)
```

## `TextEditor`

`TextEditor` wraps Gio's `widget.Editor` and applies this repo's theme styling.

Useful fields and methods:

- `Hint`
- `MaxHeight`
- `Editor.ReadOnly`
- `Text() string`
- `SetText(text string)`
- `Append(text string)`
- `SelectedText() string`
- `Layout(gtx, th)`

It also tracks:

- `Value`
  - most recent text content after change events
- `Highlighted`
  - most recent selected text

### Example

```go
ed := text.NewTextEditor("Enter notes...")
ed.MaxHeight = 240
ed.SetText("hello\nworld")
```

In layout:

```go
return ed.Layout(gtx, th)
```

Read-only usage:

```go
preview := text.NewTextEditor("")
preview.Editor.ReadOnly = true
preview.SetText("Rendered file preview")
```

## `ParseSimpleRichText`

`ParseSimpleRichText` converts a very small text format into `[]TextBlock`.

Supported patterns:

- lines starting with `# ` become `TextHeading`
- lines starting with `> ` become `TextMuted`
- lines between triple backtick fences become `TextCode`
- all other lines become `TextParagraph`

### Example

```go
blocks := text.ParseSimpleRichText(`
# Notes
> generated output
plain paragraph
`)

tv := text.NewTextView()
tv.SetBlocks(blocks)
```

## `SelectableTextBlock`

`SelectableTextBlock` renders a single themed selectable label using Gio's `widget.Selectable`.

### Example

```go
var block text.SelectableTextBlock
block.Text = "Copyable text"
```

In layout:

```go
return block.Layout(gtx, th)
```

## Notes

- `TextView` is best suited for logs, status panes, and simple structured copy.
- `TextEditor` is a convenience wrapper around `widget.Editor`, not a custom text engine.
- `ParseSimpleRichText` is intentionally minimal and should be treated as a lightweight formatter, not markdown compatibility.

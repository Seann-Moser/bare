package themes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCustomThemes(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	dir, err := configDir()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	name := PaletteName("test_custom_theme")
	delete(Palettes, name)
	t.Cleanup(func() {
		delete(Palettes, name)
	})

	data := []byte(`
themes:
  test_custom_theme:
    label: Test Custom Theme
    light:
      background: "#F8F3EA"
      surface: "#FFFFFF"
      text: "#211E1B"
      text_muted: "#625B54"
      primary: "#6C4A9E"
      secondary: "#1F6D68"
      success: "#2F7D46"
      warning: "#8A5A00"
      error: "#B3261E"
    dark:
      background: "#121019"
      surface: "#1D1926"
      text: "#F4EFFA"
      text_muted: "#BDB3C9"
      primary: "#C8B6FF"
      secondary: "#8EC5B5"
      success: "#88C999"
      warning: "#F3C969"
      error: "#FF8C8C"
`)
	if err := os.WriteFile(filepath.Join(dir, "themes.yaml"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := LoadCustomThemes(); err != nil {
		t.Fatal(err)
	}

	palette, ok := Palettes[name]
	if !ok {
		t.Fatalf("expected custom palette %q to load", name)
	}
	if palette.Label != "Test Custom Theme" {
		t.Fatalf("unexpected label %q", palette.Label)
	}

	th := New(ModeDark, name, false)
	if th.Color.Primary != Hex("#C8B6FF") {
		t.Fatalf("unexpected dark primary: %#v", th.Color.Primary)
	}
}

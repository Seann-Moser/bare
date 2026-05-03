package themes

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var paletteNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

type customThemeDocument struct {
	Themes map[string]customThemeConfig `yaml:"themes"`
}

type customThemeConfig struct {
	Label string     `yaml:"label"`
	Light modeTokens `yaml:"light"`
	Dark  modeTokens `yaml:"dark"`
}

func LoadCustomThemes() error {
	paths, err := customThemePaths()
	if err != nil {
		return err
	}

	var errs []error
	for _, path := range paths {
		if err := loadCustomThemeFile(path); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func OrderedPalettes() []PaletteName {
	seen := map[PaletteName]bool{}
	names := make([]PaletteName, 0, len(Palettes))

	for _, name := range PaletteOrder {
		if _, ok := Palettes[name]; ok {
			names = append(names, name)
			seen[name] = true
		}
	}

	custom := make([]string, 0, len(Palettes))
	for name := range Palettes {
		if seen[name] {
			continue
		}
		custom = append(custom, string(name))
	}
	sort.Strings(custom)

	for _, name := range custom {
		names = append(names, PaletteName(name))
	}

	return names
}

func CustomThemesPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "themes.yaml"), nil
}

func CustomThemesDir() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "themes"), nil
}

func customThemePaths() ([]string, error) {
	filePath, err := CustomThemesPath()
	if err != nil {
		return nil, err
	}

	paths := []string{}
	if _, err := os.Stat(filePath); err == nil {
		paths = append(paths, filePath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	dir, err := CustomThemesDir()
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	paths = append(paths, matches...)

	sort.Strings(paths)
	return paths, nil
}

func loadCustomThemeFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var doc customThemeDocument
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	for rawName, cfg := range doc.Themes {
		name := PaletteName(strings.TrimSpace(rawName))
		if err := registerCustomPalette(name, cfg); err != nil {
			return fmt.Errorf("%s: %s: %w", path, rawName, err)
		}
	}

	return nil
}

func registerCustomPalette(name PaletteName, cfg customThemeConfig) error {
	if !isValidPaletteName(name) {
		return fmt.Errorf("theme names must use lowercase letters, numbers, underscores, or hyphens")
	}
	if isBuiltinPalette(name) {
		return fmt.Errorf("custom themes cannot replace built-in themes")
	}
	if !validModeTokens(cfg.Light) || !validModeTokens(cfg.Dark) {
		return fmt.Errorf("light and dark themes must define valid #RRGGBB tokens")
	}

	label := strings.TrimSpace(cfg.Label)
	if label == "" {
		label = string(name)
	}

	Palettes[name] = newStudyPalette(name, label, cfg.Light, cfg.Dark)
	return nil
}

func isValidPaletteName(name PaletteName) bool {
	return paletteNamePattern.MatchString(string(name))
}

func isBuiltinPalette(name PaletteName) bool {
	for _, builtin := range PaletteOrder {
		if name == builtin {
			return true
		}
	}

	return false
}

package themes

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configAppName = "bare"

type Config struct {
	Mode    Mode        `yaml:"mode"`
	Palette PaletteName `yaml:"palette"`
}

func DefaultConfig() Config {
	return Config{
		Mode:    ModeSystem,
		Palette: PaletteOcean,
	}
}

func (c Config) Theme(systemDark bool) Theme {
	cfg := c.normalized()
	return New(cfg.Mode, cfg.Palette, systemDark)
}

func ConfigFromTheme(th Theme) Config {
	return Config{
		Mode:    th.Mode,
		Palette: th.Palette,
	}.normalized()
}

func LoadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	return cfg.normalized(), nil
}

func SaveConfig(cfg Config) error {
	cfg = cfg.normalized()

	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, configAppName, "theme.yaml"), nil
}

func (c Config) normalized() Config {
	if !isValidMode(c.Mode) {
		c.Mode = DefaultConfig().Mode
	}
	if _, ok := Palettes[c.Palette]; !ok {
		c.Palette = DefaultConfig().Palette
	}

	return c
}

func isValidMode(mode Mode) bool {
	switch mode {
	case ModeSystem, ModeLight, ModeDark:
		return true
	default:
		return false
	}
}

package themes

import (
	"os"
	"path/filepath"

	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
)

func loadFonts() []font.FontFace {
	var faces []font.FontFace

	paths := []string{
		"/usr/share/fonts",
		"/usr/local/share/fonts",
		filepath.Join(os.Getenv("HOME"), ".fonts"),
	}

	for _, dir := range paths {
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			if filepath.Ext(path) == ".ttf" || filepath.Ext(path) == ".otf" {
				data, err := os.ReadFile(path)
				if err != nil {
					return nil
				}

				f, err := opentype.Parse(data)
				if err != nil {
					return nil
				}

				faces = append(faces, font.FontFace{
					Font: font.Font{},
					Face: f,
				})
			}

			return nil
		})
	}

	// fallback if nothing found
	if len(faces) == 0 {
		return defaultFonts()
	}

	return faces
}

func defaultFonts() []font.FontFace {
	return gofont.Collection()
}

package themes

import (
	"image/color"

	uiutils "github.com/DarlingGoose/bare/pkg/ui/utils"
)

func readableOn(bg color.NRGBA) color.NRGBA {
	return uiutils.ReadableOn(bg)
}

func parseHexByte(a, b byte) uint8 {
	return fromHex(a)<<4 | fromHex(b)
}

func fromHex(c byte) uint8 {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}

package text

import "strings"

func ParseSimpleRichText(s string) []TextBlock {
	lines := strings.Split(s, "\n")

	blocks := make([]TextBlock, 0, len(lines))
	inCode := false

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "```"):
			inCode = !inCode
			continue

		case inCode:
			blocks = append(blocks, TextBlock{
				Kind: TextCode,
				Text: line,
			})

		case strings.HasPrefix(line, "# "):
			blocks = append(blocks, TextBlock{
				Kind: TextHeading,
				Text: strings.TrimPrefix(line, "# "),
			})

		case strings.HasPrefix(line, "> "):
			blocks = append(blocks, TextBlock{
				Kind: TextMuted,
				Text: strings.TrimPrefix(line, "> "),
			})

		default:
			blocks = append(blocks, TextBlock{
				Kind: TextParagraph,
				Text: line,
			})
		}
	}

	return blocks
}

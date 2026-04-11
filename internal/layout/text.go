package layout

import (
	"strings"
)

// WrapText splits text into lines that fit within maxWidth.
// The estimation is font-size aware and also breaks long words.
func WrapText(text string, maxWidth, fontSize float64) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	maxCharsPerLine := estimatedCharsPerLine(maxWidth, fontSize)
	paragraphs := strings.Split(text, "\n")
	lines := make([]string, 0, len(paragraphs))

	for pi, paragraph := range paragraphs {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			if pi < len(paragraphs)-1 {
				lines = append(lines, "")
			}
			continue
		}

		current := ""
		for _, word := range words {
			for runeLen(word) > maxCharsPerLine {
				if current != "" {
					lines = append(lines, current)
					current = ""
				}

				part, rest := splitRunes(word, maxCharsPerLine)
				lines = append(lines, part)
				word = rest
			}

			if current == "" {
				current = word
				continue
			}

			if runeLen(current)+1+runeLen(word) <= maxCharsPerLine {
				current += " " + word
			} else {
				lines = append(lines, current)
				current = word
			}
		}

		if current != "" {
			lines = append(lines, current)
		}

		if pi < len(paragraphs)-1 {
			lines = append(lines, "")
		}
	}

	if len(lines) == 0 {
		return []string{""}
	}

	return lines
}

func wrapText(text string, maxWidth float64) []string {
	return WrapText(text, maxWidth, 12)
}

func estimatedCharsPerLine(maxWidth, fontSize float64) int {
	if fontSize <= 0 {
		fontSize = 12
	}
	avgCharWidth := fontSize * 0.55
	if avgCharWidth <= 0 {
		avgCharWidth = 6.6
	}
	maxChars := int(maxWidth / avgCharWidth)
	if maxChars < 1 {
		return 1
	}
	return maxChars
}

func runeLen(s string) int {
	return len([]rune(s))
}

func splitRunes(s string, n int) (string, string) {
	r := []rune(s)
	if n >= len(r) {
		return s, ""
	}
	return string(r[:n]), string(r[n:])
}

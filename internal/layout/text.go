package layout

import (
	"strings"
)

const CharWidth = 7.0   // rough width per char
const LineHeight = 14.0 // px/pt

func wrapText(text string, maxWidth float64) []string {
	maxCharsPerLine := int(maxWidth / CharWidth)

	words := strings.Split(text, " ")
	var lines []string

	current := ""

	for _, w := range words {
		if len(current)+len(w)+1 > maxCharsPerLine {
			lines = append(lines, current)
			current = w
		} else {
			if current == "" {
				current = w
			} else {
				current += " " + w
			}
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
}

package render

import (
	"bytes"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/layout"
)

// RenderPDF is the main entry point that converts a layout tree to PDF bytes
// It takes a root Box (from layout engine) and returns the PDF as bytes
func RenderPDF(root *layout.Box) ([]byte, error) {
	pdf := fpdf.New("P", "pt", "A4", "")
	pdf.AddPage()

	// Render the layout tree recursively
	renderBox(pdf, root)

	// Output to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// renderBox recursively renders a layout box and its children to the PDF
// It applies styles, handles text rendering, and draws box borders
func renderBox(pdf *fpdf.Fpdf, box *layout.Box) {
	if box == nil {
		return
	}

	// 🔥 RESOLVE STYLES DYNAMICALLY
	// This applies CSS, HTML attributes, and classes to get the final style
	resolvedStyle := box.Style

	// 🔥 TEXT NODE RENDERING
	if box.Node.Type == dom.TextNode {
		// Set font based on resolved style (bold, italic, size)
		fontStyle := ""
		if resolvedStyle.Bold {
			fontStyle += "B"
		}
		if resolvedStyle.Italic {
			fontStyle += "I"
		}

		pdf.SetFont("Arial", fontStyle, resolvedStyle.FontSize)

		// Split text into lines based on box width
		lines := splitLines(box.Node.Content, box.Width)

		lineHeight := resolvedStyle.FontSize + 2
		y := box.Y + resolvedStyle.Padding + resolvedStyle.FontSize

		// Render each line with proper alignment
		for _, line := range lines {
			textWidth := pdf.GetStringWidth(line)

			x := box.X + resolvedStyle.Padding

			// 🔥 REAL alignment from resolved style
			switch resolvedStyle.Align {
			case "center":
				x = box.X + (box.Width-textWidth)/2
			case "right":
				x = box.X + box.Width - textWidth - resolvedStyle.Padding
			case "justify":
				// Could implement justify here with custom spacing
				x = box.X + resolvedStyle.Padding
			}

			// Set text color from resolved style
			if resolvedStyle.Color != "" && resolvedStyle.Color != "black" {
				setColorFromString(pdf, resolvedStyle.Color)
			}

			pdf.Text(x, y, line)
			y += lineHeight
		}

		// Reset color to black
		pdf.SetTextColor(0, 0, 0)
	}

	// 🔥 DEBUG: draw box borders (optional - set to false to hide)
	drawBoxBorders := true
	if drawBoxBorders {
		pdf.SetDrawColor(200, 200, 200) // Light gray borders
		pdf.Rect(box.X, box.Y, box.Width, box.Height, "")
		pdf.SetDrawColor(0, 0, 0) // Reset to black
	}

	// 🔥 RENDER CHILDREN RECURSIVELY
	for _, child := range box.Children {
		renderBox(pdf, child)
	}
}

// splitLines breaks text into lines that fit within maxWidth
// Uses a simple character-based approach; can be improved with actual text measurement
func splitLines(text string, maxWidth float64) []string {
	// Rough estimate: 7 pixels per character at default font size
	maxChars := int(maxWidth / 7)
	if maxChars < 1 {
		maxChars = 1
	}

	words := strings.Split(text, " ")
	var lines []string

	current := ""

	for _, w := range words {
		// Check if adding this word would exceed max width
		if len(current)+len(w)+1 > maxChars {
			if current != "" {
				lines = append(lines, current)
			}
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

// setColorFromString converts color names or hex codes to RGB values
// Supports: "red", "blue", "green", "gray", "black", "#FF0000", etc.
func setColorFromString(pdf *fpdf.Fpdf, colorStr string) {
	colorMap := map[string][3]int{
		"black":     {0, 0, 0},
		"white":     {255, 255, 255},
		"red":       {255, 0, 0},
		"green":     {0, 128, 0},
		"blue":      {0, 0, 255},
		"yellow":    {255, 255, 0},
		"cyan":      {0, 255, 255},
		"magenta":   {255, 0, 255},
		"gray":      {128, 128, 128},
		"lightgray": {211, 211, 211},
		"darkgray":  {64, 64, 64},
		"orange":    {255, 165, 0},
		"purple":    {128, 0, 128},
		"brown":     {165, 42, 42},
		"pink":      {255, 192, 203},
	}

	if rgb, ok := colorMap[strings.ToLower(colorStr)]; ok {
		pdf.SetTextColor(rgb[0], rgb[1], rgb[2])
		return
	}

	// If not in map, try to parse hex color (simple implementation)
	if strings.HasPrefix(colorStr, "#") && len(colorStr) == 7 {
		// Parse #RRGGBB format
		r, g, b := 0, 0, 0
		_, _ = sscanf(colorStr[1:3], "%x", &r)
		_, _ = sscanf(colorStr[3:5], "%x", &g)
		_, _ = sscanf(colorStr[5:7], "%x", &b)
		pdf.SetTextColor(r, g, b)
		return
	}

	// Default to black if color not recognized
	pdf.SetTextColor(0, 0, 0)
}

// Simple hex parsing helper
func sscanf(input string, _ string, ptr *int) (int, error) {
	var val int
	_, err := sscanfHex(input, &val)
	*ptr = val
	return 1, err
}

func sscanfHex(input string, val *int) (int, error) {
	n := 0
	for i := 0; i < len(input); i++ {
		c := input[i]
		var digit int
		if c >= '0' && c <= '9' {
			digit = int(c - '0')
		} else if c >= 'a' && c <= 'f' {
			digit = int(c - 'a' + 10)
		} else if c >= 'A' && c <= 'F' {
			digit = int(c - 'A' + 10)
		} else {
			break
		}
		n = n*16 + digit
	}
	*val = n
	return 1, nil
}

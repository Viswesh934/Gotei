package render

import (
	"bytes"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/Viswesh934/gotei/internal/debug"
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/layout"
	"github.com/Viswesh934/gotei/internal/style"
)

// RenderPDF is the main entry point that converts a layout tree to PDF bytes
// It takes a root Box (from layout engine) and returns the PDF as bytes
func RenderPDF(root *layout.Box) ([]byte, error) {
	pdf := fpdf.New("P", "pt", "A4", "")
	pdf.AddPage()
	debug.Logf("render: start root=(x=%.2f y=%.2f w=%.2f h=%.2f)", root.X, root.Y, root.Width, root.Height)

	// Render the layout tree recursively
	renderBox(pdf, root)

	// Output to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		debug.Logf("render: pdf-output failed: %v", err)
		return nil, err
	}
	debug.Logf("render: pdf-output bytes=%d", buf.Len())

	return buf.Bytes(), nil
}

// renderBox recursively renders a layout box and its children to the PDF
// It applies styles, handles text rendering, and draws box borders
// renderBox recursively renders a layout box and its children to the PDF
func renderBox(pdf *fpdf.Fpdf, box *layout.Box) {
	if box == nil {
		return
	}

	resolvedStyle := box.Style

	// 🔥 APPLY TEXT TRANSFORMATION before rendering
	if resolvedStyle.TextTransform != "" {
		box.Node.Content = transformText(box.Node.Content, resolvedStyle.TextTransform)
	}

	// 🔥 DRAW BACKGROUND
	if resolvedStyle.BgColor != "" && resolvedStyle.BgColor != "white" {
		br := &BorderRendering{pdf: pdf}
		br.DrawBackground(box, resolvedStyle.BgColor)
	}

	// 🔥 DRAW BOX SHADOW
	if !isZeroShadow(resolvedStyle.BoxShadow) {
		br := &BorderRendering{pdf: pdf}
		br.DrawBoxShadow(box, resolvedStyle.BoxShadow)
	}

	// 🔥 TEXT NODE RENDERING
	if box.Node.Type == dom.TextNode {
		// Set font based on resolved style
		fontStyle := ""
		if resolvedStyle.Bold {
			fontStyle += "B"
		}
		if resolvedStyle.Italic {
			fontStyle += "I"
		}

		// Map font family
		fontFamily := resolvedStyle.FontFamily
		if fontFamily == "" {
			fontFamily = "Arial"
		}
		fontFamily = mapFontFamily(fontFamily)

		pdf.SetFont(fontFamily, fontStyle, resolvedStyle.FontSize)

		// Calculate line height (supports unitless, px)
		lineHeight := resolvedStyle.FontSize + 2 // default
		if resolvedStyle.LineHeight > 0 {
			if resolvedStyle.LineHeight < 5 { // unitless like 1.5
				lineHeight = resolvedStyle.FontSize * resolvedStyle.LineHeight
			} else { // px value
				lineHeight = resolvedStyle.LineHeight
			}
		}

		contentLeft := box.X + resolvedStyle.Padding.Left
		contentWidth := box.Width - resolvedStyle.Padding.Left - resolvedStyle.Padding.Right
		if contentWidth < 0 {
			contentWidth = 0
		}

		// Split text into lines using the inner content width.
		lines := splitLines(box.Node.Content, contentWidth)

		y := box.Y + resolvedStyle.Padding.Top + resolvedStyle.FontSize

		// Render each line
		for lineIdx, line := range lines {
			// Apply text shadow if present
			if !isZeroShadow(resolvedStyle.TextShadow) {
				shadowX := contentLeft + resolvedStyle.TextShadow.OffsetX
				shadowY := y + resolvedStyle.TextShadow.OffsetY
				applyTextShadow(pdf, shadowX, shadowY, line, resolvedStyle.TextShadow)
			}

			// Calculate text width with letter spacing
			textWidth := pdf.GetStringWidth(line)
			if resolvedStyle.LetterSpacing != 0 {
				textWidth += float64(len(line)) * resolvedStyle.LetterSpacing
			}
			if resolvedStyle.WordSpacing != 0 {
				textWidth += float64(strings.Count(line, " ")) * resolvedStyle.WordSpacing
			}

			x := contentLeft

			// Apply alignment
			align := resolvedStyle.TextAlign
			if align == "" {
				align = resolvedStyle.Align
			}
			if align == "start" {
				align = "left"
			}
			if align == "end" {
				align = "right"
			}

			switch align {
			case "center":
				x = contentLeft + (contentWidth-textWidth)/2
			case "right":
				x = contentLeft + contentWidth - textWidth
			case "justify":
				if lineIdx < len(lines)-1 { // Don't justify last line
					renderJustifiedLine(pdf, line, contentLeft, contentLeft+contentWidth, y, resolvedStyle)
					y += lineHeight
					continue
				}
			}

			// Set text color
			if resolvedStyle.Color != "" {
				setColorFromString(pdf, resolvedStyle.Color)
			} else {
				pdf.SetTextColor(0, 0, 0)
			}
			debug.Logf("render: text-line idx=%d x=%.2f y=%.2f align=%s color=%s text=%q", lineIdx, x, y, align, resolvedStyle.Color, line)

			// Render text with spacing when needed.
			if resolvedStyle.LetterSpacing != 0 || resolvedStyle.WordSpacing != 0 {
				renderTextWithSpacing(pdf, x, y, line, resolvedStyle.LetterSpacing, resolvedStyle.WordSpacing)
			} else {
				pdf.Text(x, y, line)
			}

			// Apply text decoration
			if resolvedStyle.TextDecoration != "none" && resolvedStyle.TextDecoration != "" {
				applyTextDecoration(pdf, x, y, textWidth, resolvedStyle.FontSize, resolvedStyle.TextDecoration)
			}

			y += lineHeight
		}

		// Reset color to black
		pdf.SetTextColor(0, 0, 0)
	}

	// 🔥 DRAW BORDERS
	if resolvedStyle.Border.Width > 0 {
		br := &BorderRendering{pdf: pdf}
		br.DrawBorder(box, resolvedStyle)
	}

	// 🔥 RENDER CHILDREN
	for _, child := range box.Children {
		renderBox(pdf, child)
	}
}

// Add this helper function for justified text
func renderJustifiedLine(pdf *fpdf.Fpdf, line string, leftX, rightX, y float64, resolvedStyle style.Style) {
	words := strings.Fields(line)
	if len(words) <= 1 {
		if resolvedStyle.LetterSpacing != 0 || resolvedStyle.WordSpacing != 0 {
			renderTextWithSpacing(pdf, leftX, y, line, resolvedStyle.LetterSpacing, resolvedStyle.WordSpacing)
		} else {
			pdf.Text(leftX, y, line)
		}
		return
	}

	// Calculate total width and space to distribute
	totalWidth := rightX - leftX
	textWidth := pdf.GetStringWidth(line)
	spaceNeeded := totalWidth - textWidth
	spacesToAdd := len(words) - 1

	if spacesToAdd <= 0 {
		pdf.Text(leftX, y, line)
		return
	}

	extraSpace := spaceNeeded / float64(spacesToAdd)

	// Render each word with distributed spacing
	x := leftX
	for i, word := range words {
		if resolvedStyle.LetterSpacing != 0 || resolvedStyle.WordSpacing != 0 {
			renderTextWithSpacing(pdf, x, y, word, resolvedStyle.LetterSpacing, resolvedStyle.WordSpacing)
		} else {
			pdf.Text(x, y, word)
		}
		x += pdf.GetStringWidth(word) + extraSpace
		if i < len(words)-1 {
			x += pdf.GetStringWidth(" ")
			if resolvedStyle.WordSpacing != 0 {
				x += resolvedStyle.WordSpacing
			}
		}
	}
}

// Add text transformation function
func transformText(content, transform string) string {
	switch strings.ToLower(transform) {
	case "uppercase":
		return strings.ToUpper(content)
	case "lowercase":
		return strings.ToLower(content)
	case "capitalize":
		return strings.Title(content)
	case "full-width":
		return toFullWidth(content)
	default:
		return content
	}
}

func toFullWidth(s string) string {
	var out []rune
	for _, r := range s {
		switch {
		case r == ' ':
			out = append(out, '\u3000')
		case r >= 33 && r <= 126:
			out = append(out, r+0xFEE0)
		default:
			out = append(out, r)
		}
	}
	return string(out)
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
// Supports:
// - Named colors (140+ CSS colors from theme config)
// - Hex codes (#FF0000, #F00)
// - RGB/RGBA (rgb(255,0,0), rgba(255,0,0,0.5))
// - HSL/HSLA (hsl(0,100%,50%), hsla(0,100%,50%,0.5))
func setColorFromString(pdf *fpdf.Fpdf, colorStr string) {
	r, g, b, ok := style.ColorToRGB(colorStr)
	if !ok {
		r, g, b = 0, 0, 0 // Default to black on parse failure
	}

	// Keep draw/fill/text colors in sync for borders, backgrounds, and text.
	pdf.SetTextColor(r, g, b)
	pdf.SetDrawColor(r, g, b)
	pdf.SetFillColor(r, g, b)
}

// isZeroShadow checks if a shadow has no visual effect
func isZeroShadow(s style.Shadow) bool {
	return s.Blur == 0 && s.OffsetX == 0 && s.OffsetY == 0 && s.Spread == 0
}

// mapFontFamily maps generic or custom font names to FPDF-supported fonts
// FPDF supports: Arial, Courier, Times, and custom embedded fonts
func mapFontFamily(fontName string) string {
	fontName = strings.TrimSpace(fontName)
	if fontName == "" {
		return "Arial"
	}
	if strings.Contains(fontName, ",") {
		fontName = strings.Split(fontName, ",")[0]
	}
	fontName = strings.TrimSpace(strings.Trim(fontName, `"'`))

	switch strings.ToLower(fontName) {
	case "monospace", "courier", "courier new", "fixed":
		return "Courier"
	case "serif", "georgia", "times", "times new roman":
		return "Times"
	case "sans-serif", "arial", "verdana", "helvetica":
		return "Arial"
	case "code":
		return "Courier"
	default:
		return "Arial" // Fallback to Arial for unknown fonts
	}
}

// renderTextWithSpacing renders text with custom letter-spacing
// Renders character-by-character to apply spacing between characters
func renderTextWithSpacing(pdf *fpdf.Fpdf, x, y float64, text string, letterSpacing, wordSpacing float64) {
	currentX := x
	for idx, char := range text {
		glyph := string(char)
		pdf.Text(currentX, y, glyph)
		glyphWidth := pdf.GetStringWidth(glyph)
		currentX += glyphWidth + letterSpacing
		if glyph == " " && wordSpacing != 0 && idx < len(text)-1 {
			currentX += wordSpacing
		}
	}
}

// applyTextDecoration applies text-decoration (underline, line-through, overline)
func applyTextDecoration(pdf *fpdf.Fpdf, x, y float64, width, fontSize float64, decoration string) {
	lineY := y
	lineWidth := 0.5

	switch strings.ToLower(decoration) {
	case "underline":
		lineY = y + fontSize/4
		pdf.Line(x, lineY, x+width, lineY)
	case "line-through":
		lineY = y - fontSize/3
		pdf.Line(x, lineY, x+width, lineY)
	case "overline":
		lineY = y - fontSize
		pdf.Line(x, lineY, x+width, lineY)
	}

	pdf.SetLineWidth(lineWidth)
	pdf.SetDrawColor(0, 0, 0) // Use current text color
}

// applyTextShadow renders text shadow effect
func applyTextShadow(pdf *fpdf.Fpdf, x, y float64, text string, shadow style.Shadow) {
	// Save current color
	r, g, b := pdf.GetTextColor()

	// Apply shadow color (usually darkened version or gray)
	pdf.SetTextColor(100, 100, 100)

	// Render shadow at offset position
	pdf.Text(x+shadow.OffsetX, y+shadow.OffsetY, text)

	// Restore original color
	pdf.SetTextColor(r, g, b)
}

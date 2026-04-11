package render

import (
	"bytes"
	"os"
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
	pageHeight := layout.DefaultHeight
	totalPages := estimatePageCount(root, pageHeight)
	for i := 0; i < totalPages; i++ {
		pdf.AddPage()
	}
	debug.Logf("render: start root=(x=%.2f y=%.2f w=%.2f h=%.2f)", root.X, root.Y, root.Width, root.Height)

	// Render the layout tree recursively
	renderBox(pdf, root, pageHeight)

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
func renderBox(pdf *fpdf.Fpdf, box *layout.Box, pageHeight float64) {
	if box == nil {
		return
	}

	resolvedStyle := box.Style

	if box.Node != nil && box.Node.Type == dom.ElementNode && strings.ToLower(box.Node.Tag) == "img" {
		src := ""
		if box.Node.Attr != nil {
			src = strings.TrimSpace(box.Node.Attr["src"])
		}

		if src != "" {
			if _, err := os.Stat(src); err == nil {
				imgPage := setRenderPage(pdf, box.Y, pageHeight)
				localY := localYForPage(box.Y, imgPage, pageHeight)
				opts := fpdf.ImageOptions{ReadDpi: true}
				pdf.ImageOptions(src, box.X, localY, box.Width, box.Height, false, opts, 0, "")
			} else {
				debug.Logf("render: image source not found src=%q err=%v", src, err)
			}
		}

		return
	}

	// 🔥 APPLY TEXT TRANSFORMATION before rendering
	if resolvedStyle.TextTransform != "" {
		box.Node.Content = transformText(box.Node.Content, resolvedStyle.TextTransform)
	}

	// 🔥 DRAW BACKGROUND
	if resolvedStyle.BgColor != "" && resolvedStyle.BgColor != "white" {
		bgPage := setRenderPage(pdf, box.Y, pageHeight)
		localBox := *box
		localBox.Y = localYForPage(box.Y, bgPage, pageHeight)
		br := &BorderRendering{pdf: pdf}
		br.DrawBackground(&localBox, resolvedStyle.BgColor)
	}

	// 🔥 DRAW BOX SHADOW
	if !isZeroShadow(resolvedStyle.BoxShadow) {
		shadowPage := setRenderPage(pdf, box.Y, pageHeight)
		localBox := *box
		localBox.Y = localYForPage(box.Y, shadowPage, pageHeight)
		br := &BorderRendering{pdf: pdf}
		br.DrawBoxShadow(&localBox, resolvedStyle.BoxShadow)
	}

	// 🔥 TEXT NODE RENDERING
	if box.Node.Type == dom.TextNode {
		linkURL := findLinkHref(box)
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
		lineHeight := resolvedStyle.FontSize // default
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

		// Split text into lines using the same wrapping logic as layout.
		lines := layout.WrapText(box.Node.Content, contentWidth, resolvedStyle.FontSize)

		yAbs := box.Y + resolvedStyle.Padding.Top + resolvedStyle.FontSize

		// Render each line
		for lineIdx, line := range lines {
			linePage := setRenderPage(pdf, yAbs, pageHeight)
			lineY := localYForPage(yAbs, linePage, pageHeight)

			// Apply text shadow if present
			if !isZeroShadow(resolvedStyle.TextShadow) {
				shadowX := contentLeft + resolvedStyle.TextShadow.OffsetX
				shadowYAbs := yAbs + resolvedStyle.TextShadow.OffsetY
				shadowPage := setRenderPage(pdf, shadowYAbs, pageHeight)
				shadowY := localYForPage(shadowYAbs, shadowPage, pageHeight)
				applyTextShadow(pdf, shadowX, shadowY, line, resolvedStyle.TextShadow)
				setRenderPage(pdf, yAbs, pageHeight)
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
					renderJustifiedLine(pdf, line, contentLeft, contentLeft+contentWidth, lineY, resolvedStyle)
					yAbs += lineHeight
					continue
				}
			}

			// Set text color
			if resolvedStyle.Color != "" {
				setColorFromString(pdf, resolvedStyle.Color)
			} else {
				pdf.SetTextColor(0, 0, 0)
			}
			debug.Logf("render: text-line idx=%d x=%.2f y=%.2f align=%s color=%s text=%q", lineIdx, x, yAbs, align, resolvedStyle.Color, line)

			// Render text with spacing when needed.
			if resolvedStyle.LetterSpacing != 0 || resolvedStyle.WordSpacing != 0 {
				renderTextWithSpacing(pdf, x, lineY, line, resolvedStyle.LetterSpacing, resolvedStyle.WordSpacing)
			} else {
				pdf.Text(x, lineY, line)
			}

			if linkURL != "" {
				linkTopY := lineY - resolvedStyle.FontSize
				if linkTopY < 0 {
					linkTopY = 0
				}
				pdf.LinkString(x, linkTopY, textWidth, lineHeight, linkURL)
			}

			// Apply text decoration
			if resolvedStyle.TextDecoration != "none" && resolvedStyle.TextDecoration != "" {
				applyTextDecoration(pdf, x, lineY, textWidth, resolvedStyle.FontSize, resolvedStyle.TextDecoration)
			}

			yAbs += lineHeight
		}

		// Reset color to black
		pdf.SetTextColor(0, 0, 0)
	}

	// 🔥 DRAW BORDERS
	if resolvedStyle.Border.Width > 0 {
		borderPage := setRenderPage(pdf, box.Y, pageHeight)
		localBox := *box
		localBox.Y = localYForPage(box.Y, borderPage, pageHeight)

		if isTableCellBox(box) {
			drawCollapsedTableCellBorder(pdf, &localBox, resolvedStyle, box)
		} else {
			br := &BorderRendering{pdf: pdf}
			br.DrawBorder(&localBox, resolvedStyle)
		}
	}

	// 🔥 RENDER CHILDREN
	for _, child := range box.Children {
		renderBox(pdf, child, pageHeight)
	}
}

func findLinkHref(box *layout.Box) string {
	for current := box; current != nil; current = current.Parent {
		if current.Node == nil || current.Node.Type != dom.ElementNode {
			continue
		}

		if strings.ToLower(current.Node.Tag) != "a" || current.Node.Attr == nil {
			continue
		}

		href := strings.TrimSpace(current.Node.Attr["href"])
		if href != "" {
			return href
		}
	}

	return ""
}

func estimatePageCount(root *layout.Box, pageHeight float64) int {
	if pageHeight <= 0 {
		return 1
	}

	maxBottom := maxBoxBottom(root)
	if maxBottom <= 0 {
		return 1
	}

	pages := int(maxBottom/pageHeight) + 1
	if pages < 1 {
		return 1
	}
	return pages
}

func maxBoxBottom(box *layout.Box) float64 {
	if box == nil {
		return 0
	}

	maxBottom := box.Y + box.Height
	for _, child := range box.Children {
		childBottom := maxBoxBottom(child)
		if childBottom > maxBottom {
			maxBottom = childBottom
		}
	}

	return maxBottom
}

func setRenderPage(pdf *fpdf.Fpdf, y, pageHeight float64) int {
	page := pageForY(y, pageHeight)
	pdf.SetPage(page)
	return page
}

func pageForY(y, pageHeight float64) int {
	if pageHeight <= 0 {
		return 1
	}
	if y < 0 {
		y = 0
	}
	return int(y/pageHeight) + 1
}

func localYForPage(y float64, page int, pageHeight float64) float64 {
	if pageHeight <= 0 {
		return y
	}
	if page < 1 {
		page = 1
	}
	return y - float64(page-1)*pageHeight
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

	// Distribute extra space evenly, but cap it at reasonable limits
	extraSpace := spaceNeeded / float64(spacesToAdd)
	
	// Cap max extra space to avoid huge gaps (max 4pt additional spacing)
	maxExtraPerGap := 4.0
	if extraSpace > maxExtraPerGap {
		extraSpace = maxExtraPerGap
	}

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

func isTableCellBox(box *layout.Box) bool {
	if box == nil || box.Node == nil || box.Node.Type != dom.ElementNode {
		return false
	}
	tag := strings.ToLower(box.Node.Tag)
	return tag == "td" || tag == "th"
}

func drawCollapsedTableCellBorder(pdf *fpdf.Fpdf, localBox *layout.Box, s style.Style, original *layout.Box) {
	if localBox == nil || original == nil || s.Border.Width <= 0 {
		return
	}

	row := original.Parent
	if row == nil || row.Node == nil {
		br := &BorderRendering{pdf: pdf}
		br.DrawBorder(localBox, s)
		return
	}

	pdf.SetLineWidth(s.Border.Width)
	setColorFromString(pdf, s.Border.Color)

	x0 := localBox.X
	y0 := localBox.Y
	x1 := localBox.X + localBox.Width
	y1 := localBox.Y + localBox.Height

	firstCol := isFirstCellInRow(original, row)
	firstRow := isFirstRowInTable(row)

	// Draw a single shared grid by drawing only right+bottom edges for each cell,
	// plus top/left edges only on the outer boundary.
	if firstRow {
		pdf.Line(x0, y0, x1, y0)
	}
	if firstCol {
		pdf.Line(x0, y0, x0, y1)
	}
	pdf.Line(x1, y0, x1, y1)
	pdf.Line(x0, y1, x1, y1)

	pdf.SetLineWidth(0.5)
	pdf.SetDrawColor(0, 0, 0)
}

func isFirstCellInRow(cell, row *layout.Box) bool {
	if cell == nil || row == nil {
		return false
	}
	for _, child := range row.Children {
		if child == nil || child.Node == nil || child.Node.Type != dom.ElementNode {
			continue
		}
		tag := strings.ToLower(child.Node.Tag)
		if tag != "td" && tag != "th" {
			continue
		}
		return child == cell
	}
	return false
}

func isFirstRowInTable(row *layout.Box) bool {
	if row == nil {
		return false
	}
	table := findAncestorTable(row)
	if table == nil {
		return false
	}

	for _, child := range table.Children {
		if child == nil || child.Node == nil || child.Node.Type != dom.ElementNode {
			continue
		}
		tag := strings.ToLower(child.Node.Tag)
		if tag == "tr" {
			return child == row
		}
		if tag == "thead" || tag == "tbody" || tag == "tfoot" {
			for _, grandChild := range child.Children {
				if grandChild == nil || grandChild.Node == nil || grandChild.Node.Type != dom.ElementNode {
					continue
				}
				if strings.ToLower(grandChild.Node.Tag) != "tr" {
					continue
				}
				return grandChild == row
			}
		}
	}

	return false
}

func findAncestorTable(box *layout.Box) *layout.Box {
	for current := box; current != nil; current = current.Parent {
		if current.Node == nil || current.Node.Type != dom.ElementNode {
			continue
		}
		if strings.ToLower(current.Node.Tag) == "table" {
			return current
		}
	}
	return nil
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

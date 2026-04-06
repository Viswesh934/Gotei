package render

import (
	"bytes"
	"strings"

	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/layout"
	"github.com/jung-kurt/gofpdf"
)

func RenderPDF(root *layout.Box) ([]byte, error) {
	pdf := gofpdf.New("P", "pt", "A4", "")
	pdf.AddPage()

	// 🔥 Set font (important)
	pdf.SetFont("Arial", "", 12)

	renderBox(pdf, root)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// helper to write PDF into []byte
type bufferWriter struct {
	buf *[]byte
}

func (w *bufferWriter) Write(p []byte) (int, error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}

func renderBox(pdf *gofpdf.Fpdf, box *layout.Box) {
	if box == nil {
		return
	}

	// 🔥 TEXT NODE
	if box.Node.Type == dom.TextNode {
		lines := splitLines(box.Node.Content, box.Width)

		y := box.Y + 12

		for _, line := range lines {
			pdf.Text(box.X+5, y, line)
			y += 14
		}
	}

	// DEBUG: draw box borders (optional but SUPER useful)
	pdf.Rect(box.X, box.Y, box.Width, box.Height, "")

	for _, child := range box.Children {
		renderBox(pdf, child)
	}
}

func splitLines(text string, maxWidth float64) []string {
	maxChars := int(maxWidth / 7)

	words := strings.Split(text, " ")
	var lines []string

	current := ""

	for _, w := range words {
		if len(current)+len(w)+1 > maxChars {
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

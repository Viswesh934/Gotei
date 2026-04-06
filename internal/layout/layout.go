package layout

import "github.com/Viswesh934/gotei/internal/dom"

const DefaultWidth = 595.0  // A4 width
const DefaultHeight = 842.0 // A4 height

func Layout(box *Box, x, y, maxWidth float64) float64 {
	box.X = x
	box.Y = y
	box.Width = maxWidth

	// TEXT NODE handling
	if box.Node.Type == dom.TextNode {
		lines := wrapText(box.Node.Content, maxWidth)
		box.Height = float64(len(lines)) * LineHeight
		return box.Height
	}

	currentY := y

	for _, child := range box.Children {
		childHeight := Layout(child, x, currentY, maxWidth)
		currentY += childHeight
	}

	box.Height = currentY - y
	return box.Height
}

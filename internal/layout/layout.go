package layout

import "github.com/Viswesh934/gotei/internal/dom"

const DefaultWidth = 595.0  // A4 width
const DefaultHeight = 842.0 // A4 height

func Layout(box *Box, x, y, maxWidth float64) float64 {
	// Apply margin
	x += box.Style.Margin
	y += box.Style.Margin

	box.X = x
	box.Y = y
	box.Width = maxWidth - (2 * box.Style.Margin)

	// TEXT NODE
	if box.Node.Type == dom.TextNode {
		lines := wrapText(box.Node.Content, box.Width-(2*box.Style.Padding))

		box.Height = float64(len(lines))*LineHeight + (2 * box.Style.Padding)

		return box.Height + (2 * box.Style.Margin)
	}

	// Children layout with padding
	currentY := y + box.Style.Padding

	for _, child := range box.Children {
		childHeight := Layout(
			child,
			x+box.Style.Padding,
			currentY,
			box.Width-(2*box.Style.Padding),
		)
		currentY += childHeight
	}

	contentHeight := currentY - (y + box.Style.Padding)

	box.Height = contentHeight + (2 * box.Style.Padding)

	return box.Height + (2 * box.Style.Margin)
}

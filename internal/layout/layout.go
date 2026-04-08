package layout

import "github.com/Viswesh934/gotei/internal/dom"

const DefaultWidth = 595.0  // A4 width
const DefaultHeight = 842.0 // A4 height
const LineHeight = 16.0     // Default line height in pixels

func Layout(box *Box, x, y, maxWidth float64) float64 {
	// Apply margin
	x += box.Style.Margin.Left
	y += box.Style.Margin.Top

	box.X = x
	box.Y = y
	box.Width = maxWidth - (box.Style.Margin.Left + box.Style.Margin.Right)

	// TEXT NODE
	if box.Node.Type == dom.TextNode {
		paddingHorizontal := box.Style.Padding.Left + box.Style.Padding.Right
		paddingVertical := box.Style.Padding.Top + box.Style.Padding.Bottom
		lines := wrapText(box.Node.Content, box.Width-paddingHorizontal)

		lineHeight := CalculateLineHeight(box.Style)
		box.Height = float64(len(lines))*lineHeight + paddingVertical

		return box.Height + box.Style.Margin.Top + box.Style.Margin.Bottom
	}

	// Route to appropriate layout algorithm based on display property
	switch box.Style.Display {
	case "flex":
		return LayoutFlexbox(box, x, y, box.Width)
	case "inline", "inline-block":
		// For now, treat inline-block similar to block
		fallthrough
	default: // "block" or empty
		return layoutBlock(box, x, y)
	}
}

// layoutBlock implements standard block layout (default)
func layoutBlock(box *Box, x, y float64) float64 {
	// Children layout with padding
	currentY := y + box.Style.Padding.Top

	for _, child := range box.Children {
		paddingHorizontal := box.Style.Padding.Left + box.Style.Padding.Right
		childHeight := Layout(
			child,
			x+box.Style.Padding.Left,
			currentY,
			box.Width-paddingHorizontal,
		)
		currentY += childHeight
	}

	contentHeight := currentY - (y + box.Style.Padding.Top)

	box.Height = contentHeight + box.Style.Padding.Top + box.Style.Padding.Bottom

	return box.Height + box.Style.Margin.Top + box.Style.Margin.Bottom
}

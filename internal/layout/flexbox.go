package layout

import (
	"math"

	"github.com/Viswesh934/gotei/internal/style"
)

// FlexItem represents a child in a flex container
type FlexItem struct {
	Box        *Box
	FlexGrow   float64
	FlexShrink float64
	FlexBasis  float64
	MinWidth   float64
	MaxWidth   float64
	Margin     style.BoxDimensions
}

// LayoutFlexbox arranges flex children in a container
// Implements W3C Flexbox algorithm (simplified)
func LayoutFlexbox(box *Box, x, y, maxWidth float64) float64 {
	if len(box.Children) == 0 {
		box.X = x
		box.Y = y
		box.Width = maxWidth
		box.Height = 0
		return 0
	}

	isRow := box.Style.FlexDirection == "row" || box.Style.FlexDirection == ""
	isWrap := box.Style.FlexWrap == "wrap"

	if isRow {
		return layoutFlexRow(box, x, y, maxWidth)
	}
	var _ bool = isWrap
	return layoutFlexColumn(box, x, y, maxWidth)
}

func layoutFlexRow(box *Box, x, y, maxWidth float64) float64 {
	paddingH := box.Style.Padding.Left + box.Style.Padding.Right
	marginH := box.Style.Margin.Left + box.Style.Margin.Right
	availableWidth := maxWidth - marginH - paddingH

	box.X = x + box.Style.Margin.Left
	box.Y = y + box.Style.Margin.Top
	box.Width = maxWidth - marginH

	// Collect flex items for main axis distribution
	items := make([]*FlexItem, len(box.Children))
	totalFlexGrow := 0.0
	totalFlexShrink := 0.0
	usedWidth := 0.0

	for i, child := range box.Children {
		items[i] = &FlexItem{
			Box:        child,
			FlexGrow:   child.Style.FlexGrow,
			FlexShrink: child.Style.FlexShrink,
			FlexBasis:  child.Style.FlexBasis,
			MinWidth:   child.Style.MinWidth,
			MaxWidth:   child.Style.MaxWidth,
			Margin:     child.Style.Margin,
		}

		totalFlexGrow += items[i].FlexGrow
		totalFlexShrink += items[i].FlexShrink

		// Calculate base width
		if items[i].FlexBasis > 0 {
			usedWidth += items[i].FlexBasis
		} else if items[i].Box.Width > 0 {
			usedWidth += items[i].Box.Width
		}
	}

	// Distribute remaining space
	remainingSpace := availableWidth - usedWidth
	if remainingSpace > 0 && totalFlexGrow > 0 {
		// Grow items
		for _, item := range items {
			if item.FlexGrow > 0 {
				extra := (remainingSpace * item.FlexGrow) / totalFlexGrow
				item.Box.Width += extra
			}
		}
	} else if remainingSpace < 0 && totalFlexShrink > 0 {
		// Shrink items
		for _, item := range items {
			if item.FlexShrink > 0 {
				reduction := (-remainingSpace * item.FlexShrink) / totalFlexShrink
				item.Box.Width = math.Max(item.MinWidth, item.Box.Width-reduction)
			}
		}
	}

	// Position items along main axis
	currentX := box.X + box.Style.Padding.Left
	maxHeight := 0.0

	for _, child := range box.Children {
		child.X = currentX + child.Style.Margin.Left
		child.Y = box.Y + box.Style.Padding.Top + child.Style.Margin.Top

		// Align items on cross axis (vertical for row)
		switch box.Style.AlignItems {
		case "center":
			// Will be calculated after we know max height
		case "flex-end":
			// Will be calculated after we know max height
		default: // flex-start, stretch
			// Default behavior
		}

		currentX += child.Width + child.Style.Margin.Left + child.Style.Margin.Right + box.Style.Gap

		if child.Height > maxHeight {
			maxHeight = child.Height
		}
	}

	// Apply cross-axis alignment
	for _, child := range box.Children {
		switch box.Style.AlignItems {
		case "center":
			child.Y = box.Y + box.Style.Padding.Top + (maxHeight-child.Height)/2
		case "flex-end":
			child.Y = box.Y + box.Style.Padding.Top + maxHeight - child.Height
		}
	}

	box.Height = maxHeight + box.Style.Padding.Top + box.Style.Padding.Bottom

	// Distribute items along main axis if justify-content is set
	totalItemWidth := 0.0
	for _, child := range box.Children {
		totalItemWidth += child.Width + child.Style.Margin.Left + child.Style.Margin.Right
	}
	freeSpace := availableWidth - totalItemWidth

	if freeSpace > 0 {
		switch box.Style.JustifyContent {
		case "center":
			offset := freeSpace / 2
			for _, child := range box.Children {
				child.X += offset
			}
		case "flex-end":
			for _, child := range box.Children {
				child.X += freeSpace
			}
		case "space-between":
			if len(box.Children) > 1 {
				gap := freeSpace / float64(len(box.Children)-1)
				for i := 1; i < len(box.Children); i++ {
					box.Children[i].X += gap * float64(i)
				}
			}
		case "space-around":
			gap := freeSpace / float64(len(box.Children))
			for i, child := range box.Children {
				child.X += gap * (float64(i) + 0.5)
			}
		}
	}

	return box.Height + box.Style.Margin.Top + box.Style.Margin.Bottom
}

func layoutFlexColumn(box *Box, x, y, maxWidth float64) float64 {
	paddingV := box.Style.Padding.Top + box.Style.Padding.Bottom
	marginV := box.Style.Margin.Top + box.Style.Margin.Bottom
	availableHeight := 842.0 - y - marginV - paddingV

	box.X = x + box.Style.Margin.Left
	box.Y = y + box.Style.Margin.Top
	box.Width = maxWidth

	totalFlexGrow := 0.0
	totalFlexShrink := 0.0
	usedHeight := 0.0

	for _, child := range box.Children {
		totalFlexGrow += child.Style.FlexGrow
		totalFlexShrink += child.Style.FlexShrink
		usedHeight += child.Height + child.Style.Margin.Top + child.Style.Margin.Bottom
	}

	// Distribute remaining space
	remainingSpace := availableHeight - usedHeight
	if remainingSpace > 0 && totalFlexGrow > 0 {
		for _, child := range box.Children {
			if child.Style.FlexGrow > 0 {
				extra := (remainingSpace * child.Style.FlexGrow) / totalFlexGrow
				child.Height += extra
			}
		}
	}

	// Position items along main axis (vertical)
	currentY := box.Y + box.Style.Padding.Top
	for _, child := range box.Children {
		child.X = box.X + box.Style.Padding.Left + child.Style.Margin.Left
		child.Y = currentY + child.Style.Margin.Top

		currentY += child.Height + child.Style.Margin.Top + child.Style.Margin.Bottom + box.Style.Gap
	}

	box.Height = usedHeight + paddingV

	return box.Height + marginV
}

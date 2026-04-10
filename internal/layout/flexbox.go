package layout

import (
	"math"

	"github.com/Viswesh934/gotei/internal/dom"
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

	// Position container
	box.X = x
	box.Y = y
	box.Width = maxWidth
	if box.Width < 0 {
		box.Width = 0
	}

	contentWidth := box.Width - box.Style.Padding.Left - box.Style.Padding.Right
	if contentWidth < 0 {
		contentWidth = 0
	}

	// Recursive pre-measure pass so text/content gets real dimensions before flex math.
	preMeasureWidth := contentWidth
	if len(box.Children) > 0 {
		preMeasureWidth = contentWidth / float64(len(box.Children))
	}
	if preMeasureWidth < 0 {
		preMeasureWidth = 0
	}
	for _, child := range box.Children {
		Layout(child, box.X+box.Style.Padding.Left, box.Y+box.Style.Padding.Top, preMeasureWidth)
	}

	isRow := box.Style.FlexDirection == "row" || box.Style.FlexDirection == ""
	if isRow {
		return layoutFlexRowFixed(box, contentWidth)
	}
	return layoutFlexColumnFixed(box, contentWidth)
}

// layoutFlexRowFixed arranges items horizontally
func layoutFlexRowFixed(box *Box, contentWidth float64) float64 {
	availableWidth := contentWidth

	// Collect flex items
	totalFlexGrow := 0.0
	totalFlexShrink := 0.0
	totalUsedWidth := 0.0 // Includes item outer widths (not gap)

	for _, child := range box.Children {
		totalFlexGrow += child.Style.FlexGrow
		totalFlexShrink += child.Style.FlexShrink

		baseWidth := child.Style.Width
		if child.Style.FlexBasis > 0 {
			baseWidth = child.Style.FlexBasis
		} else if baseWidth <= 0 {
			if child.Style.FlexGrow > 0 {
				baseWidth = 0
			} else {
				// Use intrinsic content width for non-growing items so justify-content
				// can center/end-align correctly instead of stretching every item.
				baseWidth = intrinsicContentWidth(child)
			}
		}
		if baseWidth <= 0 {
			baseWidth = child.Width
		}
		if baseWidth < 0 {
			baseWidth = 0
		}

		child.Width = baseWidth
		totalUsedWidth += child.Width + child.Style.Margin.Left + child.Style.Margin.Right
	}

	gapTotal := 0.0
	if len(box.Children) > 1 {
		gapTotal = box.Style.Gap * float64(len(box.Children)-1)
	}

	// Calculate remaining space for flex distribution
	remainingSpace := availableWidth - totalUsedWidth - gapTotal

	// Distribute space based on flex-grow or flex-shrink
	if remainingSpace > 0 && totalFlexGrow > 0 {
		// Grow items proportionally.
		for _, child := range box.Children {
			if child.Style.FlexGrow > 0 {
				extra := (remainingSpace * child.Style.FlexGrow) / totalFlexGrow
				child.Width += extra
			}
		}
	} else if remainingSpace < 0 && totalFlexShrink > 0 {
		// Shrink items proportionally.
		for _, child := range box.Children {
			if child.Style.FlexShrink > 0 {
				reduction := (-remainingSpace * child.Style.FlexShrink) / totalFlexShrink
				child.Width = math.Max(child.Style.MinWidth, child.Width-reduction)
			}
		}
	}

	// Final child layout pass with resolved outer widths.
	currentX := box.X + box.Style.Padding.Left
	maxOuterHeight := 0.0

	for i, child := range box.Children {
		outerWidth := child.Width + child.Style.Margin.Left + child.Style.Margin.Right
		if outerWidth < 0 {
			outerWidth = 0
		}

		outerHeight := Layout(child, currentX, box.Y+box.Style.Padding.Top, outerWidth)
		if outerHeight > maxOuterHeight {
			maxOuterHeight = outerHeight
		}

		currentX += outerWidth
		if i < len(box.Children)-1 {
			currentX += box.Style.Gap
		}
	}

	// Apply cross-axis alignment (vertical alignment for flex-row)
	for _, child := range box.Children {
		childOuterHeight := child.Height + child.Style.Margin.Top + child.Style.Margin.Bottom
		dy := 0.0
		switch box.Style.AlignItems {
		case "center":
			dy = (maxOuterHeight - childOuterHeight) / 2
		case "flex-end":
			dy = maxOuterHeight - childOuterHeight
		case "stretch":
			target := maxOuterHeight - child.Style.Margin.Top - child.Style.Margin.Bottom
			if target < 0 {
				target = 0
			}
			child.Height = target
		default: // "flex-start"
			// No shift.
		}

		if dy != 0 {
			shiftSubtree(child, 0, dy)
		}
	}

	// Calculate container height
	box.Height = maxOuterHeight + box.Style.Padding.Top + box.Style.Padding.Bottom

	// Apply extra distribution from justify-content in addition to base gap spacing.
	if box.Style.JustifyContent != "" && box.Style.JustifyContent != "flex-start" {
		totalItemWidth := 0.0
		for _, child := range box.Children {
			totalItemWidth += child.Width + child.Style.Margin.Left + child.Style.Margin.Right
		}

		freeSpace := availableWidth - totalItemWidth - gapTotal

		if freeSpace > 0 {
			switch box.Style.JustifyContent {
			case "center":
				offset := freeSpace / 2
				for _, child := range box.Children {
					shiftSubtree(child, offset, 0)
				}

			case "flex-end":
				for _, child := range box.Children {
					shiftSubtree(child, freeSpace, 0)
				}

			case "space-between":
				if len(box.Children) > 1 {
					spaceBetweenGap := freeSpace / float64(len(box.Children)-1)
					for i := 1; i < len(box.Children); i++ {
						shiftSubtree(box.Children[i], spaceBetweenGap*float64(i), 0)
					}
				}

			case "space-around":
				spacePer := freeSpace / float64(len(box.Children))
				for i, child := range box.Children {
					shiftSubtree(child, spacePer*(float64(i)+0.5), 0)
				}

			case "space-evenly":
				spacePer := freeSpace / float64(len(box.Children)+1)
				for i, child := range box.Children {
					shiftSubtree(child, spacePer*float64(i+1), 0)
				}
			}
		}
	}

	return box.Height + box.Style.Margin.Top + box.Style.Margin.Bottom
}

func intrinsicContentWidth(box *Box) float64 {
	if box == nil {
		return 0
	}

	if box.Node != nil && box.Node.Type == dom.TextNode {
		fontSize := box.Style.FontSize
		if fontSize <= 0 {
			fontSize = 12
		}
		textWidth := float64(len([]rune(box.Node.Content))) * fontSize * 0.55
		return textWidth + box.Style.Padding.Left + box.Style.Padding.Right
	}

	if len(box.Children) == 0 {
		if box.Style.Width > 0 {
			return box.Style.Width
		}
		return box.Style.Padding.Left + box.Style.Padding.Right
	}

	maxChild := 0.0
	for _, child := range box.Children {
		w := intrinsicContentWidth(child) + child.Style.Margin.Left + child.Style.Margin.Right
		if w > maxChild {
			maxChild = w
		}
	}

	result := maxChild + box.Style.Padding.Left + box.Style.Padding.Right
	if box.Style.Width > result {
		result = box.Style.Width
	}
	return result
}

// layoutFlexColumnFixed arranges items vertically
func layoutFlexColumnFixed(box *Box, contentWidth float64) float64 {
	paddingV := box.Style.Padding.Top + box.Style.Padding.Bottom
	availableHeight := 842.0 - box.Y - box.Style.Margin.Bottom - paddingV // A4 page height

	// Collect flex items
	totalFlexGrow := 0.0
	totalFlexShrink := 0.0
	totalUsedHeight := 0.0

	for _, child := range box.Children {
		totalFlexGrow += child.Style.FlexGrow
		totalFlexShrink += child.Style.FlexShrink
		totalUsedHeight += child.Height + child.Style.Margin.Top + child.Style.Margin.Bottom
	}

	// Account for gaps
	numGaps := math.Max(0, float64(len(box.Children)-1))
	gapTotal := box.Style.Gap * numGaps
	totalUsedHeight += gapTotal

	// Distribute remaining space
	remainingSpace := availableHeight - totalUsedHeight

	if remainingSpace > 0 && totalFlexGrow > 0 {
		// Grow items
		for _, child := range box.Children {
			if child.Style.FlexGrow > 0 {
				extra := (remainingSpace * child.Style.FlexGrow) / totalFlexGrow
				child.Height += extra
			}
		}
	} else if remainingSpace < 0 && totalFlexShrink > 0 {
		// Shrink items
		for _, child := range box.Children {
			if child.Style.FlexShrink > 0 {
				reduction := (-remainingSpace * child.Style.FlexShrink) / totalFlexShrink
				child.Height = math.Max(child.Style.MinHeight, child.Height-reduction)
			}
		}
	}

	// Position items along main axis (top to bottom)
	currentY := box.Y + box.Style.Padding.Top
	usedHeight := 0.0
	alignItems := box.Style.AlignItems
	if alignItems == "" {
		alignItems = "stretch"
	}

	for i, child := range box.Children {
		desiredContentWidth := child.Style.Width
		if child.Style.FlexBasis > 0 {
			desiredContentWidth = child.Style.FlexBasis
		}

		if desiredContentWidth <= 0 {
			if alignItems == "stretch" {
				desiredContentWidth = contentWidth - child.Style.Margin.Left - child.Style.Margin.Right
			} else {
				desiredContentWidth = intrinsicContentWidth(child)
			}
		}

		if desiredContentWidth < 0 {
			desiredContentWidth = 0
		}

		outerWidth := desiredContentWidth + child.Style.Margin.Left + child.Style.Margin.Right
		outerHeight := Layout(child, box.X+box.Style.Padding.Left, currentY, outerWidth)

		currentY += outerHeight
		usedHeight += outerHeight
		if i < len(box.Children)-1 {
			currentY += box.Style.Gap
			usedHeight += box.Style.Gap
		}
	}

	// Apply cross-axis alignment (horizontal alignment for flex-column)
	for _, child := range box.Children {
		childOuterWidth := child.Width + child.Style.Margin.Left + child.Style.Margin.Right
		freeCross := contentWidth - childOuterWidth
		if freeCross < 0 {
			freeCross = 0
		}

		switch box.Style.AlignItems {
		case "center":
			shiftSubtree(child, freeCross/2, 0)
		case "flex-end":
			shiftSubtree(child, freeCross, 0)
		case "stretch":
			target := contentWidth - child.Style.Margin.Left - child.Style.Margin.Right
			if target < 0 {
				target = 0
			}
			child.Width = target
		default: // "flex-start"
			// No shift.
		}
	}

	// Calculate container height
	box.Height = usedHeight + paddingV

	// Apply main-axis alignment (justify-content for column)
	if box.Style.JustifyContent != "" && box.Style.JustifyContent != "flex-start" {
		totalItemHeight := 0.0
		for _, child := range box.Children {
			totalItemHeight += child.Height + child.Style.Margin.Top + child.Style.Margin.Bottom
		}

		freeSpace := availableHeight - totalItemHeight - gapTotal

		if freeSpace > 0 {
			switch box.Style.JustifyContent {
			case "center":
				offset := freeSpace / 2
				for _, child := range box.Children {
					shiftSubtree(child, 0, offset)
				}

			case "flex-end":
				for _, child := range box.Children {
					shiftSubtree(child, 0, freeSpace)
				}

			case "space-between":
				if len(box.Children) > 1 {
					spaceBetweenGap := freeSpace / float64(len(box.Children)-1)
					for i := 1; i < len(box.Children); i++ {
						shiftSubtree(box.Children[i], 0, spaceBetweenGap*float64(i))
					}
				}

			case "space-around":
				spacePer := freeSpace / float64(len(box.Children))
				for i, child := range box.Children {
					shiftSubtree(child, 0, spacePer*(float64(i)+0.5))
				}
			}
		}
	}

	return box.Height + box.Style.Margin.Top + box.Style.Margin.Bottom
}

func shiftSubtree(box *Box, dx, dy float64) {
	if box == nil || (dx == 0 && dy == 0) {
		return
	}

	box.X += dx
	box.Y += dy
	for _, child := range box.Children {
		shiftSubtree(child, dx, dy)
	}
}

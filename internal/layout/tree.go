package layout

import (
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/style"
)

func BuildLayoutTree(n *dom.Node, parentStyle style.Style) *Box {
	if n == nil {
		return nil
	}

	currentStyle := style.Resolve(n)

	// Text nodes inherit ALL inheritable properties from parent (W3C spec)
	if n.Type == dom.TextNode {
		currentStyle = inheritAllProperties(currentStyle, parentStyle)
	} else {
		// Element nodes inherit only if not explicitly set
		currentStyle = shouldInherit(currentStyle, parentStyle)
	}

	box := &Box{
		Node:  n,
		Style: currentStyle,
	}

	for _, child := range n.Children {
		childBox := BuildLayoutTree(child, box.Style)
		if childBox != nil {
			childBox.Parent = box
			box.Children = append(box.Children, childBox)
		}
	}

	return box
}

// BuildLayoutTreeWithCascade builds the layout tree using cascade-computed styles
// This version uses pre-computed styles from the CSS cascade engine
func BuildLayoutTreeWithCascade(n *dom.Node, styleMap map[*dom.Node]style.Style, parentStyle style.Style) *Box {
	if n == nil {
		return nil
	}

	// Get pre-computed style from cascade, or fall back to Resolve
	currentStyle, exists := styleMap[n]
	if !exists {
		currentStyle = style.Resolve(n)
	}

	// Text nodes inherit ALL inheritable properties from parent (W3C spec)
	if n.Type == dom.TextNode {
		currentStyle = inheritAllProperties(currentStyle, parentStyle)
	} else {
		// Element nodes inherit only if not explicitly set
		currentStyle = shouldInherit(currentStyle, parentStyle)
	}

	box := &Box{
		Node:  n,
		Style: currentStyle,
	}

	for _, child := range n.Children {
		childBox := BuildLayoutTreeWithCascade(child, styleMap, box.Style)
		if childBox != nil {
			childBox.Parent = box
			box.Children = append(box.Children, childBox)
		}
	}

	return box
}

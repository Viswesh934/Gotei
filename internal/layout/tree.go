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

	// Only inherit parent align if node doesn't have explicit align attribute
	switch n.Type {
	case dom.ElementNode:
		if _, hasAlignAttr := n.Attr["align"]; !hasAlignAttr {
			// No explicit align, so inherit from parent
			currentStyle.Align = parentStyle.Align
		}
	case dom.TextNode:
		// Text nodes always inherit from parent
		currentStyle.Align = parentStyle.Align
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

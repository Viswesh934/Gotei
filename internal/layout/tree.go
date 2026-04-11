package layout

import (
	"strings"

	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/style"
)

// BuildLayoutTreeWithCascade builds the layout tree using cascade-computed styles
// This version uses pre-computed styles from the CSS cascade engine
func BuildLayoutTreeWithCascade(n *dom.Node, styleMap map[*dom.Node]style.Style) *Box {
	if n == nil {
		return nil
	}

	if n.Type == dom.ElementNode {
		tag := strings.ToLower(strings.TrimSpace(n.Tag))
		switch tag {
		case "head", "style", "script", "meta", "title", "link", "noscript":
			return nil
		}
	}

	// The cascade step already resolved inheritance and defaults.
	// This path should only consume the computed style map.
	currentStyle, exists := styleMap[n]
	if !exists {
		currentStyle = style.DefaultStyle
	}

	box := &Box{
		Node:  n,
		Style: currentStyle,
	}

	for _, child := range n.Children {
		childBox := BuildLayoutTreeWithCascade(child, styleMap)
		if childBox != nil {
			childBox.Parent = box
			box.Children = append(box.Children, childBox)
		}
	}

	return box
}

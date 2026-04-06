package layout

import "github.com/Viswesh934/gotei/internal/dom"

func BuildLayoutTree(n *dom.Node) *Box {
	if n == nil {
		return nil
	}

	box := &Box{
		Node: n,
	}

	for _, child := range n.Children {
		childBox := BuildLayoutTree(child)
		if childBox != nil {
			box.Children = append(box.Children, childBox)
		}
	}

	return box
}

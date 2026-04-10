package css

import (
	"golang.org/x/net/html"

	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/style"
)

// ApplyCascadeToTree applies CSS cascade to all nodes in the DOM tree
// This computes the final styles for each node considering selectors, specificity, inheritance
func ApplyCascadeToTree(root *dom.Node, sheet *StyleSheet) map[*dom.Node]style.Style {
	styleMap := make(map[*dom.Node]style.Style)
	_, domToHTML := mirrorToHTML(root)
	applyCascadeRecursive(root, domToHTML[root], sheet, style.DefaultStyle, styleMap, domToHTML)
	return styleMap
}

// applyCascadeRecursive walks the tree and computes styles
func applyCascadeRecursive(node *dom.Node, htmlNode *html.Node, sheet *StyleSheet, parentStyle style.Style, styleMap map[*dom.Node]style.Style, domToHTML map[*dom.Node]*html.Node) {
	if node == nil {
		return
	}
	// Compute style for this node considering cascade and inheritance
	computed := ComputeStyleWithHTMLNode(node, htmlNode, sheet, parentStyle)
	styleMap[node] = computed

	// Recursively apply cascade to children
	for _, child := range node.Children {
		applyCascadeRecursive(child, domToHTML[child], sheet, computed, styleMap, domToHTML)
	}
}

func mirrorToHTML(root *dom.Node) (*html.Node, map[*dom.Node]*html.Node) {
	mapping := make(map[*dom.Node]*html.Node)
	return mirrorToHTMLRecursive(root, mapping), mapping
}

func mirrorToHTMLRecursive(n *dom.Node, mapping map[*dom.Node]*html.Node) *html.Node {
	if n == nil {
		return nil
	}

	h := &html.Node{}
	switch n.Type {
	case dom.ElementNode:
		h.Type = html.ElementNode
		h.Data = n.Tag
		if len(n.Attr) > 0 {
			h.Attr = make([]html.Attribute, 0, len(n.Attr))
			for k, v := range n.Attr {
				h.Attr = append(h.Attr, html.Attribute{Key: k, Val: v})
			}
		}
	case dom.TextNode:
		h.Type = html.TextNode
		h.Data = n.Content
	default:
		h.Type = html.ElementNode
		h.Data = "div"
	}

	mapping[n] = h

	for _, child := range n.Children {
		hc := mirrorToHTMLRecursive(child, mapping)
		if hc == nil {
			continue
		}
		h.AppendChild(hc)
	}

	return h
}

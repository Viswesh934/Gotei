package css

import (
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/style"
)

// ApplyCascadeToTree applies CSS cascade to all nodes in the DOM tree
// This computes the final styles for each node considering selectors, specificity, inheritance
func ApplyCascadeToTree(root *dom.Node, sheet *StyleSheet) map[*dom.Node]style.Style {
	styleMap := make(map[*dom.Node]style.Style)
	applyCascadeRecursive(root, sheet, style.DefaultStyle, styleMap)
	return styleMap
}

// applyCascadeRecursive walks the tree and computes styles
func applyCascadeRecursive(node *dom.Node, sheet *StyleSheet, parentStyle style.Style, styleMap map[*dom.Node]style.Style) {
	if node == nil {
		return
	}

	// Compute style for this node considering cascade and inheritance
	computed := ComputeStyle(node, sheet, parentStyle)
	styleMap[node] = computed

	// Recursively apply cascade to children
	for _, child := range node.Children {
		applyCascadeRecursive(child, sheet, computed, styleMap)
	}
}

package engine

import (
	"github.com/Viswesh934/gotei/internal/css"
	"github.com/Viswesh934/gotei/internal/debug"
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/layout"
	"github.com/Viswesh934/gotei/internal/render"
)

func Render(htmlStr string) ([]byte, error) {
	debug.Logf("render start: html-bytes=%d", len(htmlStr))

	// Phase 1: Parse HTML
	tree, err := dom.ParseHTML(htmlStr)
	if err != nil {
		debug.Logf("phase=parse-html failed: %v", err)
		return nil, err
	}
	debug.Logf("phase=parse-html ok: dom-nodes=%d", countDOMNodes(tree))

	// Phase 2: Extract CSS from <style> tags and build stylesheet
	sheet := css.ParseStyleSheet(tree)
	if sheet != nil {
		debug.Logf("phase=parse-css ok: css-rules=%d", len(sheet.Rules))
	} else {
		debug.Logf("phase=parse-css ok: css-rules=0 (nil sheet)")
	}

	// Phase 3: Apply CSS cascade to all nodes
	styleMap := css.ApplyCascadeToTree(tree, sheet)
	debug.Logf("phase=cascade ok: styled-nodes=%d", len(styleMap))

	// Phase 4: Build layout tree using cascade-computed styles
	layoutTree := layout.BuildLayoutTreeWithCascade(tree, styleMap)
	debug.Logf("phase=build-layout-tree ok")

	// Phase 5: Calculate layout
	layout.Layout(layoutTree, 0, 0, layout.DefaultWidth)
	debug.Logf("phase=layout ok: root=(x=%.2f y=%.2f w=%.2f h=%.2f)", layoutTree.X, layoutTree.Y, layoutTree.Width, layoutTree.Height)

	// Phase 6: Render to PDF
	pdfBytes, err := render.RenderPDF(layoutTree)
	if err != nil {
		debug.Logf("phase=render-pdf failed: %v", err)
		return nil, err
	}
	debug.Logf("phase=render-pdf ok: bytes=%d", len(pdfBytes))
	return pdfBytes, nil
}

func countDOMNodes(n *dom.Node) int {
	if n == nil {
		return 0
	}
	total := 1
	for _, child := range n.Children {
		total += countDOMNodes(child)
	}
	return total
}

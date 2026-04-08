package engine

import (
	"github.com/Viswesh934/gotei/internal/css"
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/layout"
	"github.com/Viswesh934/gotei/internal/render"
	"github.com/Viswesh934/gotei/internal/style"
)

func Render(htmlStr string) ([]byte, error) {
	// Phase 1: Parse HTML
	tree, err := dom.ParseHTML(htmlStr)
	if err != nil {
		return nil, err
	}

	// Phase 2: Extract CSS from <style> tags and build stylesheet
	sheet := css.ParseStyleSheet(tree)

	// Phase 3: Apply CSS cascade to all nodes
	styleMap := css.ApplyCascadeToTree(tree, sheet)

	// Phase 4: Build layout tree using cascade-computed styles
	layoutTree := layout.BuildLayoutTreeWithCascade(tree, styleMap, style.Style{})

	// Phase 5: Calculate layout
	layout.Layout(layoutTree, 0, 0, layout.DefaultWidth)

	// Phase 6: Render to PDF
	return render.RenderPDF(layoutTree)
}

package engine

import (
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/layout"
	"github.com/Viswesh934/gotei/internal/render"
)

func Render(htmlStr string) ([]byte, error) {
	tree, err := dom.ParseHTML(htmlStr)
	if err != nil {
		return nil, err
	}

	layoutTree := layout.BuildLayoutTree(tree)
	layout.Layout(layoutTree, 0, 0, layout.DefaultWidth)

	return render.RenderPDF(layoutTree)
}

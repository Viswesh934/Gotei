package layout

import (
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/style"
)

type Box struct {
	X, Y          float64
	Width, Height float64

	Style style.Style

	Node     *dom.Node
	Children []*Box
	Parent   *Box
}

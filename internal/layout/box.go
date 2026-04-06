package layout

import "github.com/Viswesh934/gotei/internal/dom"

type Box struct {
	X, Y          float64
	Width, Height float64

	Node     *dom.Node
	Children []*Box
}

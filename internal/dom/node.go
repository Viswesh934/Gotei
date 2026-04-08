package dom

type NodeType int

const (
	ElementNode NodeType = iota
	TextNode
)

type Node struct {
	Type     NodeType
	Tag      string
	Content  string
	Attr     map[string]string
	Children []*Node
}

package dom

import (
	"strings"

	"golang.org/x/net/html"
)

func ParseHTML(input string) (*Node, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return nil, err
	}

	return convertNode(doc), nil
}

func convertNode(n *html.Node) *Node {
	if n.Type == html.TextNode {
		content := strings.TrimSpace(n.Data)
		if content == "" {
			return nil
		}
		return &Node{
			Type:    TextNode,
			Content: content,
		}
	}

	if n.Type == html.ElementNode {
		node := &Node{
			Type: ElementNode,
			Tag:  n.Data,
			Attr: map[string]string{},
		}

		for _, attr := range n.Attr {
			node.Attr[attr.Key] = attr.Val
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			child := convertNode(c)
			if child != nil {
				node.Children = append(node.Children, child)
			}
		}

		return node
	}

	// Skip comments, doctypes etc
	var children []*Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		child := convertNode(c)
		if child != nil {
			children = append(children, child)
		}
	}

	if len(children) == 1 {
		return children[0]
	}

	return &Node{
		Type:     ElementNode,
		Tag:      "root",
		Children: children,
	}
}

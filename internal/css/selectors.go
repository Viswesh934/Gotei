package css

import (
	"strings"

	"github.com/Viswesh934/gotei/internal/dom"
)

// Selector represents a CSS selector with specificity calculation
type Selector struct {
	Element     string // element name or "*"
	Classes     []string
	IDs         []string
	Attributes  []AttrSelector
	PseudoClass string
	Specificity Specificity
}

// AttrSelector represents an attribute selector
type AttrSelector struct {
	Name  string
	Value string
	Match string // "=", "~=", "|=", "^=", "$=", "*="
}

// Specificity calculates CSS specificity (a,b,c)
// a = IDs, b = classes+attributes+pseudo-classes, c = elements
type Specificity struct {
	IDs      int // a
	Classes  int // b (includes classes, attributes, pseudo-classes)
	Elements int // c
}

// Compare returns:
// -1 if s < other
//
//	0 if s == other
//	1 if s > other
func (s Specificity) Compare(other Specificity) int {
	if s.IDs != other.IDs {
		if s.IDs > other.IDs {
			return 1
		}
		return -1
	}
	if s.Classes != other.Classes {
		if s.Classes > other.Classes {
			return 1
		}
		return -1
	}
	if s.Elements != other.Elements {
		if s.Elements > other.Elements {
			return 1
		}
		return -1
	}
	return 0
}

// ParseSelector parses a CSS selector string
func ParseSelector(sel string) *Selector {
	sel = strings.TrimSpace(sel)
	selector := &Selector{
		Classes:    []string{},
		IDs:        []string{},
		Attributes: []AttrSelector{},
	}

	parts := strings.FieldsFunc(sel, func(r rune) bool {
		return r == '.' || r == '#' || r == '[' || r == ']'
	})

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check what this part is
		if strings.HasPrefix(sel, "#") {
			// ID selector
			idPart := strings.SplitN(part, "#", 2)
			if len(idPart) > 1 {
				selector.IDs = append(selector.IDs, idPart[1])
			}
		} else if strings.HasPrefix(sel, ".") {
			// Class selector
			classPart := strings.SplitN(part, ".", 2)
			if len(classPart) > 1 {
				selector.Classes = append(selector.Classes, classPart[1])
			}
		} else if selector.Element == "" {
			// Element selector
			if part != "*" && !strings.Contains(part, ":") {
				selector.Element = part
			} else if part == "*" {
				selector.Element = "*"
			}
		}
	}

	// Parse ID selectors from original
	for _, idMatch := range strings.Split(sel, "#") {
		if idMatch == "" || idMatch == sel {
			continue
		}
		id := strings.FieldsFunc(idMatch, func(r rune) bool {
			return r == '.' || r == '[' || r == ':' || r == ' '
		})
		if len(id) > 0 {
			selector.IDs = append(selector.IDs, id[0])
		}
	}

	// Parse class selectors from original
	for _, classMatch := range strings.Split(sel, ".") {
		if classMatch == "" || classMatch == sel {
			continue
		}
		class := strings.FieldsFunc(classMatch, func(r rune) bool {
			return r == '#' || r == '[' || r == ':' || r == ' '
		})
		if len(class) > 0 {
			selector.Classes = append(selector.Classes, class[0])
		}
	}

	// Set element if not already set and it's not a class/id-only selector
	if selector.Element == "" && len(selector.IDs) == 0 && len(selector.Classes) == 0 {
		// Extract element name from beginning
		first := strings.FieldsFunc(sel, func(r rune) bool {
			return r == '.' || r == '#' || r == '[' || r == ':' || r == ' '
		})
		if len(first) > 0 && first[0] != "*" {
			selector.Element = first[0]
		}
	}

	// Calculate specificity
	selector.Specificity = Specificity{
		IDs:      len(selector.IDs),
		Classes:  len(selector.Classes) + len(selector.Attributes),
		Elements: 1,
	}
	if selector.Element == "*" || selector.Element == "" {
		selector.Specificity.Elements = 0
	}

	return selector
}

// Matches checks if a selector matches a DOM node
func (sel *Selector) Matches(node *dom.Node) bool {
	if node.Type != dom.ElementNode {
		return false
	}

	// Check element
	if sel.Element != "" && sel.Element != "*" && node.Tag != sel.Element {
		return false
	}

	// Check ID selectors
	for _, id := range sel.IDs {
		if nodeID, ok := node.Attr["id"]; !ok || nodeID != id {
			return false
		}
	}

	// Check class selectors
	nodeClasses := strings.Fields(strings.TrimSpace(node.Attr["class"]))
	for _, class := range sel.Classes {
		found := false
		for _, nodeClass := range nodeClasses {
			if nodeClass == class {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check attribute selectors
	for _, attr := range sel.Attributes {
		if !sel.matchesAttribute(node, attr) {
			return false
		}
	}

	return true
}

func (sel *Selector) matchesAttribute(node *dom.Node, attr AttrSelector) bool {
	val, exists := node.Attr[attr.Name]
	if !exists && attr.Match != "" {
		return false
	}

	switch attr.Match {
	case "=":
		return val == attr.Value
	case "~=":
		// Space-separated list
		return containsWord(val, attr.Value)
	case "|=":
		// Exact or followed by hyphen
		return val == attr.Value || strings.HasPrefix(val, attr.Value+"-")
	case "^=":
		// Starts with
		return strings.HasPrefix(val, attr.Value)
	case "$=":
		// Ends with
		return strings.HasSuffix(val, attr.Value)
	case "*=":
		// Contains
		return strings.Contains(val, attr.Value)
	default:
		// Just check existence
		return exists
	}
}

func containsWord(haystack, word string) bool {
	words := strings.Fields(haystack)
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}

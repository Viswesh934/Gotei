package css

import (
	"strings"

	"github.com/Viswesh934/gotei/internal/dom"
)

// ParseStyleSheet extracts CSS from <style> tags in the DOM and builds a StyleSheet
func ParseStyleSheet(htmlRoot *dom.Node) *StyleSheet {
	sheet := &StyleSheet{
		Rules: []*Rule{},
	}

	// Find all <style> tags in the DOM
	walkDOM(htmlRoot, func(n *dom.Node) {
		if n.Type == dom.ElementNode && n.Tag == "style" && len(n.Children) > 0 {
			// Extract CSS text from style tag
			cssText := ""
			for _, child := range n.Children {
				if child.Type == dom.TextNode {
					cssText += child.Content
				}
			}

			// Parse CSS and add rules to sheet
			rules := parseCSS(cssText)
			sheet.Rules = append(sheet.Rules, rules...)
		}
	})

	return sheet
}

// walkDOM recursively walks through DOM nodes and calls fn on each
func walkDOM(n *dom.Node, fn func(*dom.Node)) {
	if n == nil {
		return
	}
	fn(n)
	for _, child := range n.Children {
		walkDOM(child, fn)
	}
}

// parseCSS parses CSS text and returns a slice of Rules
// Handles simple selectors and declarations
func parseCSS(cssText string) []*Rule {
	var rules []*Rule

	// Split by closing brace to find rule blocks
	blocks := strings.Split(cssText, "}")
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		// Split selector from declarations
		parts := strings.SplitN(block, "{", 2)
		if len(parts) != 2 {
			continue
		}

		selector := strings.TrimSpace(parts[0])
		declarations := strings.TrimSpace(parts[1])

		// Parse selector
		sel := ParseSelector(selector)
		if sel == nil {
			continue
		}

		// Parse declarations into properties map
		properties := parseDeclarations(declarations)

		rule := &Rule{
			Selector:    sel,
			Properties:  properties,
			Specificity: sel.Specificity,
		}

		rules = append(rules, rule)
	}

	return rules
}

// parseDeclarations parses CSS declarations (property: value; property: value;)
func parseDeclarations(decl string) map[string]string {
	properties := make(map[string]string)

	// Split by semicolon
	parts := strings.Split(decl, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split by colon
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			properties[key] = val
		}
	}

	return properties
}

package css

import (
	"strings"

	wrparser "github.com/benoitkugler/webrender/css/parser"
	wrselector "github.com/benoitkugler/webrender/css/selector"
	"github.com/Viswesh934/gotei/internal/debug"
	"github.com/Viswesh934/gotei/internal/dom"
)

// ParseStyleSheet extracts CSS from <style> tags in the DOM and builds a StyleSheet
func ParseStyleSheet(htmlRoot *dom.Node) *StyleSheet {
	sheet := &StyleSheet{
		Rules: []*Rule{},
	}
	sourceOrder := 0
	styleTagCount := 0

	// Find all <style> tags in the DOM
	walkDOM(htmlRoot, func(n *dom.Node) {
		if n.Type == dom.ElementNode && n.Tag == "style" && len(n.Children) > 0 {
			styleTagCount++
			// Extract CSS text from style tag
			cssText := ""
			for _, child := range n.Children {
				if child.Type == dom.TextNode {
					cssText += child.Content
				}
			}

			// Parse CSS using WebRender's parser + selector engine.
			rules, nextOrder := parseCSS(cssText, sourceOrder)
			sourceOrder = nextOrder
			sheet.Rules = append(sheet.Rules, rules...)
			debug.Logf("css.parse: style-tag=%d css-bytes=%d parsed-rules=%d", styleTagCount, len(cssText), len(rules))
		}
	})

	debug.Logf("css.parse: done style-tags=%d total-rules=%d", styleTagCount, len(sheet.Rules))

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

// parseCSS parses CSS text into rule+declaration entries using WebRender's CSS parser.
func parseCSS(cssText string, sourceOrder int) ([]*Rule, int) {
	var rules []*Rule

	compounds := wrparser.ParseStylesheetBytes([]byte(cssText), true, true)
	for _, compound := range compounds {
		qualified, ok := compound.(wrparser.QualifiedRule)
		if !ok {
			continue
		}

		selectorText := strings.TrimSpace(wrparser.Serialize(qualified.Prelude))
		if selectorText == "" {
			continue
		}

		selectorGroup, err := wrselector.ParseGroup(selectorText)
		if err != nil {
			debug.Logf("css.parse: selector-parse-failed selector=%q err=%v", selectorText, err)
			continue
		}

		properties := parseDeclarations(qualified.Content)
		if len(properties) == 0 {
			continue
		}

		for _, sel := range selectorGroup {
			// Pseudo-elements are not represented in our DOM layout tree.
			if sel.PseudoElement() != "" {
				continue
			}

			sp := sel.Specificity()
			rule := &Rule{
				Matcher: sel,
				Properties: properties,
				Specificity: Specificity{
					IDs:      sp[0],
					Classes:  sp[1],
					Elements: sp[2],
				},
				SourceOrder: sourceOrder,
			}
			sourceOrder++
			rules = append(rules, rule)
		}
		debug.Logf("css.parse: selector=%q matcher-count=%d properties=%d", selectorText, len(selectorGroup), len(properties))
	}

	return rules, sourceOrder
}

// parseDeclarations parses qualified rule content into property/value pairs.
func parseDeclarations(tokens []wrparser.Token) map[string]string {
	properties := make(map[string]string)

	declarations := wrparser.ParseDeclarationList(tokens, true, true)
	for _, compound := range declarations {
		decl, ok := compound.(wrparser.Declaration)
		if !ok {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(decl.Name))
		if key == "" {
			continue
		}

		val := strings.TrimSpace(wrparser.Serialize(decl.Value))
		if val == "" {
			continue
		}

		properties[key] = val
	}

	return properties
}

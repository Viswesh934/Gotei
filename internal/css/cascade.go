package css

import (
	"strconv"
	"strings"

	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/style"
)

// Rule represents a CSS rule with selector and properties
type Rule struct {
	Selector    *Selector
	Properties  map[string]string
	Specificity Specificity
}

// StyleSheet contains CSS rules
type StyleSheet struct {
	Rules []*Rule
}

// ComputedStyle calculates the final computed style for a node
// considering cascade rules, specificity, and inheritance
func ComputeStyle(node *dom.Node, sheet *StyleSheet, parentStyle style.Style) style.Style {

	// Start with default
	computed := style.DefaultStyle

	// Apply tag defaults if element node
	if node.Type == dom.ElementNode {
		if tagStyle, exists := style.TagDefaults[node.Tag]; exists {
			computed = mergeStylesForCompute(computed, tagStyle)
		}
	}

	// Collect matching rules sorted by specificity
	matchingRules := collectMatchingRules(node, sheet)
	sortRulesBySpecificity(matchingRules)

	// Apply rules in cascading order (lowest to highest specificity)
	for _, rule := range matchingRules {
		computed = applyRuleProperties(computed, rule.Properties)
	}

	// Apply inline styles (highest specificity)
	if node.Type == dom.ElementNode {
		computed = applyElementAttributes(computed, node)
		computed = applyInlineStyleAttribute(computed, node)
	}

	// Apply inheritance for inheritable properties
	computed = applyInheritance(computed, parentStyle)

	return computed
}

func collectMatchingRules(node *dom.Node, sheet *StyleSheet) []*Rule {
	if sheet == nil {
		return []*Rule{}
	}

	matching := []*Rule{}
	for _, rule := range sheet.Rules {
		if rule.Selector.Matches(node) {
			matching = append(matching, rule)
		}
	}
	return matching
}

func sortRulesBySpecificity(rules []*Rule) {
	// Simple bubble sort by specificity
	for i := 0; i < len(rules)-1; i++ {
		for j := 0; j < len(rules)-i-1; j++ {
			if rules[j].Specificity.Compare(rules[j+1].Specificity) > 0 {
				rules[j], rules[j+1] = rules[j+1], rules[j]
			}
		}
	}
}

func applyRuleProperties(s style.Style, props map[string]string) style.Style {
	for key, val := range props {
		s = applyProperty(s, key, val)
	}
	return s
}

func applyProperty(s style.Style, key, val string) style.Style {
	switch key {
	case "font-size":
		if size, err := parseSize(val); err == nil {
			s.FontSize = size
		}
	case "font-family":
		s.FontFamily = val
	case "font-weight":
		s.FontWeight = val
		if val == "bold" || val == "700" || val == "800" || val == "900" {
			s.Bold = true
		}
	case "font-style":
		s.FontStyle = val
		if val == "italic" || val == "oblique" {
			s.Italic = true
		}
	case "color":
		s.Color = val
	case "background-color", "background":
		s.BgColor = val
	case "text-align":
		s.TextAlign = val
		s.Align = val // legacy
	case "line-height":
		if lh, err := parseSize(val); err == nil {
			s.LineHeight = lh
		}
	case "letter-spacing":
		if ls, err := parseSize(val); err == nil {
			s.LetterSpacing = ls
		}
	case "word-spacing":
		if ws, err := parseSize(val); err == nil {
			s.WordSpacing = ws
		}
	case "text-transform":
		s.TextTransform = val
	case "text-decoration":
		s.TextDecoration = val
	case "display":
		s.Display = val
	case "position":
		s.Position = val
	case "margin":
		s.Margin = parseBoxDimensions(val, s.Margin)
	case "margin-top":
		if m, err := parseSize(val); err == nil {
			s.Margin.Top = m
		}
	case "margin-right":
		if m, err := parseSize(val); err == nil {
			s.Margin.Right = m
		}
	case "margin-bottom":
		if m, err := parseSize(val); err == nil {
			s.Margin.Bottom = m
		}
	case "margin-left":
		if m, err := parseSize(val); err == nil {
			s.Margin.Left = m
		}
	case "padding":
		s.Padding = parseBoxDimensions(val, s.Padding)
	case "padding-top":
		if p, err := parseSize(val); err == nil {
			s.Padding.Top = p
		}
	case "padding-right":
		if p, err := parseSize(val); err == nil {
			s.Padding.Right = p
		}
	case "padding-bottom":
		if p, err := parseSize(val); err == nil {
			s.Padding.Bottom = p
		}
	case "padding-left":
		if p, err := parseSize(val); err == nil {
			s.Padding.Left = p
		}
	case "border":
		s.Border = parseBorder(val, s.Border)
	case "border-width":
		if bw, err := parseSize(val); err == nil {
			s.Border.Width = bw
		}
	case "border-color":
		s.Border.Color = val
	case "border-style":
		s.Border.Style = val
	case "border-radius":
		if br, err := parseSize(val); err == nil {
			s.Border.Radius = br
		}
	case "width":
		if w, err := parseSize(val); err == nil {
			s.Width = w
		}
	case "height":
		if h, err := parseSize(val); err == nil {
			s.Height = h
		}
	case "max-width":
		if mw, err := parseSize(val); err == nil {
			s.MaxWidth = mw
		}
	case "min-width":
		if mw, err := parseSize(val); err == nil {
			s.MinWidth = mw
		}
	case "max-height":
		if mh, err := parseSize(val); err == nil {
			s.MaxHeight = mh
		}
	case "min-height":
		if mh, err := parseSize(val); err == nil {
			s.MinHeight = mh
		}
	case "opacity":
		if op, err := parseSize(val); err == nil {
			s.Opacity = op
		}
	case "flex-direction":
		s.FlexDirection = val
	case "flex-wrap":
		s.FlexWrap = val
	case "justify-content":
		s.JustifyContent = val
	case "align-items":
		s.AlignItems = val
	case "align-content":
		s.AlignContent = val
	case "gap":
		if g, err := parseSize(val); err == nil {
			s.Gap = g
		}
	case "flex-grow":
		if fg, err := parseSize(val); err == nil {
			s.FlexGrow = fg
		}
	case "flex-shrink":
		if fs, err := parseSize(val); err == nil {
			s.FlexShrink = fs
		}
	case "flex-basis":
		if fb, err := parseSize(val); err == nil {
			s.FlexBasis = fb
		}
	case "flex":
		s.Flex = val
	case "overflow":
		s.Overflow = val
	case "overflow-x":
		s.OverflowX = val
	case "overflow-y":
		s.OverflowY = val
	case "visibility":
		s.Visibility = val
	case "box-shadow":
		s.BoxShadow = parseBoxShadow(val, s.BoxShadow)
	case "text-shadow":
		s.TextShadow = parseBoxShadow(val, s.TextShadow)
	}
	return s
}

func applyElementAttributes(s style.Style, node *dom.Node) style.Style {
	if node.Attr == nil {
		return s
	}

	if val, ok := node.Attr["align"]; ok {
		s.TextAlign = val
		s.Align = val
	}

	if val, ok := node.Attr["width"]; ok {
		if width, err := parseSize(val); err == nil {
			s.Width = width
		}
	}

	if val, ok := node.Attr["height"]; ok {
		if height, err := parseSize(val); err == nil {
			s.Height = height
		}
	}

	return s
}

func applyInlineStyleAttribute(s style.Style, node *dom.Node) style.Style {
	if node.Attr == nil {
		return s
	}

	styleStr, ok := node.Attr["style"]
	if !ok {
		return s
	}

	parts := strings.Split(styleStr, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])
		s = applyProperty(s, key, val)
	}

	return s
}

// parseSize converts CSS size strings to float64
func parseSize(sizeStr string) (float64, error) {
	sizeStr = strings.TrimSpace(sizeStr)

	// Remove common CSS units
	units := []string{"px", "pt", "em", "rem", "%"}
	for _, unit := range units {
		sizeStr = strings.TrimSuffix(sizeStr, unit)
	}

	return strconv.ParseFloat(sizeStr, 64)
}

// parseBoxDimensions parses CSS box model shorthand (margin/padding)
func parseBoxDimensions(sizeStr string, current style.BoxDimensions) style.BoxDimensions {
	values := strings.Fields(sizeStr)
	dims := current

	if len(values) == 0 {
		return dims
	}

	sizes := make([]float64, len(values))
	for i, v := range values {
		if size, err := parseSize(v); err == nil {
			sizes[i] = size
		}
	}

	if len(sizes) == 1 {
		// All sides equal
		dims.Top = sizes[0]
		dims.Right = sizes[0]
		dims.Bottom = sizes[0]
		dims.Left = sizes[0]
	} else if len(sizes) == 2 {
		// Top/Bottom, Left/Right
		dims.Top = sizes[0]
		dims.Bottom = sizes[0]
		dims.Left = sizes[1]
		dims.Right = sizes[1]
	} else if len(sizes) == 3 {
		// Top, Left/Right, Bottom
		dims.Top = sizes[0]
		dims.Left = sizes[1]
		dims.Right = sizes[1]
		dims.Bottom = sizes[2]
	} else if len(sizes) >= 4 {
		// Top, Right, Bottom, Left
		dims.Top = sizes[0]
		dims.Right = sizes[1]
		dims.Bottom = sizes[2]
		dims.Left = sizes[3]
	}

	return dims
}

func applyInheritance(s style.Style, parent style.Style) style.Style {
	// Apply inherited properties from parent if node has defaults
	for propName := range style.InheritableFields {
		if !style.InheritableFields[propName] {
			continue // Skip non-inheritable
		}

		// Simple inheritance: if property is at default and parent has custom value, inherit
		// This would need reflection for full implementation; for now we handle key properties
		switch propName {
		case "FontSize":
			if s.FontSize == style.DefaultStyle.FontSize && parent.FontSize != style.DefaultStyle.FontSize {
				s.FontSize = parent.FontSize
			}
		case "Color":
			if s.Color == "black" && parent.Color != "black" {
				s.Color = parent.Color
			}
		case "FontFamily":
			if s.FontFamily == "" && parent.FontFamily != "" {
				s.FontFamily = parent.FontFamily
			}
		case "LineHeight":
			if s.LineHeight == 0 && parent.LineHeight != 0 {
				s.LineHeight = parent.LineHeight
			}
		}
	}
	return s
}

func mergeStylesForCompute(base, override style.Style) style.Style {
	if override.FontSize != 0 {
		base.FontSize = override.FontSize
	}
	if override.FontWeight != "" {
		base.FontWeight = override.FontWeight
	}
	if override.FontFamily != "" {
		base.FontFamily = override.FontFamily
	}
	if override.Color != "" {
		base.Color = override.Color
	}
	if override.Bold {
		base.Bold = override.Bold
	}
	if override.Italic {
		base.Italic = override.Italic
	}
	if (override.Margin != style.BoxDimensions{}) {
		base.Margin = override.Margin
	}
	if (override.Padding != style.BoxDimensions{}) {
		base.Padding = override.Padding
	}
	return base
}

// parseBorder parses CSS border shorthand: "2px solid #333333"
func parseBorder(val string, current style.Border) style.Border {
	parts := strings.Fields(val)
	border := current

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Try parsing as size (border-width)
		if size, err := parseSize(part); err == nil {
			border.Width = size
			continue
		}

		// Try parsing as style
		if part == "solid" || part == "dashed" || part == "dotted" || part == "double" {
			border.Style = part
			continue
		}

		// Otherwise treat as color
		if part != "" {
			border.Color = part
		}
	}

	return border
}

// parseBoxShadow parses CSS box-shadow or text-shadow: "2px 2px 4px #000000"
func parseBoxShadow(val string, current style.Shadow) style.Shadow {
	parts := strings.Fields(val)
	shadow := current
	values := []float64{}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Try parsing as size
		if size, err := parseSize(part); err == nil {
			values = append(values, size)
		} else {
			// Treat as color
			if part != "" {
				shadow.Color = part
			}
		}
	}

	// Apply values: offsetX, offsetY, blur, spread
	if len(values) >= 1 {
		shadow.OffsetX = values[0]
	}
	if len(values) >= 2 {
		shadow.OffsetY = values[1]
	}
	if len(values) >= 3 {
		shadow.Blur = values[2]
	}
	if len(values) >= 4 {
		shadow.Spread = values[3]
	}

	return shadow
}

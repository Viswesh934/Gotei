package style

import (
	"strconv"
	"strings"

	"github.com/Viswesh934/gotei/internal/dom"
)

// Resolve returns the computed style for a given DOM node
// It applies styles in cascade order: defaults → tag defaults → attributes → inline styles → classes
func Resolve(n *dom.Node) Style {
	// Start with default style
	s := DefaultStyle

	// Apply tag-specific defaults if it's an element node
	if n.Type == dom.ElementNode {
		if tagStyle, exists := TagDefaults[n.Tag]; exists {
			s = mergeStyles(s, tagStyle)
		}
	}

	// Apply inline attributes (highest priority)
	if n.Type == dom.ElementNode {
		s = applyAttributes(s, n.Attr)
		s = applyInlineStyles(s, n.Attr)
		s = applyClasses(s, n.Attr)
	}

	if n.Attr != nil {
	if val, ok := n.Attr["align"]; ok {
		s.Align = val
	}
}

	return s
}

// applyAttributes handles common HTML attributes
func applyAttributes(s Style, attr map[string]string) Style {
	// Alignment
	if val, ok := attr["align"]; ok {
		if isValidAlign(val) {
			s.Align = val
		}
	}

	// Width and Height
	if val, ok := attr["width"]; ok {
		if width, err := parseSize(val); err == nil {
			s.Width = width
		}
	}

	if val, ok := attr["height"]; ok {
		if height, err := parseSize(val); err == nil {
			s.Height = height
		}
	}

	return s
}

// applyInlineStyles handles style attribute
// Example: style="font-size: 16px; color: red; text-align: center"
func applyInlineStyles(s Style, attr map[string]string) Style {
	styleStr, ok := attr["style"]
	if !ok {
		return s
	}

	// Split by semicolon to get individual properties
	properties := strings.Split(styleStr, ";")

	for _, prop := range properties {
		prop = strings.TrimSpace(prop)
		if prop == "" {
			continue
		}

		// Split by colon to get property and value
		parts := strings.Split(prop, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch strings.ToLower(key) {
		case "font-size":
			if size, err := parseSize(val); err == nil {
				s.FontSize = size
			}
		case "margin":
			s.Margin = parseBoxDimensions(val, s.Margin)
		case "padding":
			s.Padding = parseBoxDimensions(val, s.Padding)
		case "color":
			s.Color = val
		case "background-color", "background":
			s.BgColor = val
		case "text-align":
			if isValidAlign(val) {
				s.TextAlign = val
				s.Align = val // legacy support
			}
		case "font-weight":
			if val == "bold" || val == "700" || val == "800" || val == "900" {
				s.Bold = true
			}
		case "font-style":
			if val == "italic" || val == "oblique" {
				s.Italic = true
			}
		case "width":
			if width, err := parseSize(val); err == nil {
				s.Width = width
			}
		case "height":
			if height, err := parseSize(val); err == nil {
				s.Height = height
			}
		}
	}

	return s
}

// applyClasses handles CSS classes
// Classes are defined in the ClassStyles map in style.go
// Example: class="title highlight center"
func applyClasses(s Style, attr map[string]string) Style {
	classStr, ok := attr["class"]
	if !ok {
		return s
	}

	classes := strings.Split(classStr, " ")

	for _, class := range classes {
		class = strings.TrimSpace(class)
		if class == "" {
			continue
		}

		if classStyle, exists := ClassStyles[class]; exists {
			s = mergeStyles(s, classStyle)
		}
	}

	return s
}

// mergeStyles combines two styles, with the second one taking precedence
// Used for cascading style application
func mergeStyles(base, override Style) Style {
	if override.FontSize != 0 {
		base.FontSize = override.FontSize
	}
	if (override.Margin != BoxDimensions{}) {
		base.Margin = override.Margin
	}
	if (override.Padding != BoxDimensions{}) {
		base.Padding = override.Padding
	}
	if override.TextAlign != "" {
		base.TextAlign = override.TextAlign
	}
	if override.Align != "" {
		base.Align = override.Align
	}
	if override.Color != "" {
		base.Color = override.Color
	}
	if override.BgColor != "" {
		base.BgColor = override.BgColor
	}
	if override.Width != 0 {
		base.Width = override.Width
	}
	if override.Height != 0 {
		base.Height = override.Height
	}
	if override.Bold {
		base.Bold = override.Bold
	}
	if override.Italic {
		base.Italic = override.Italic
	}

	return base
}

// parseSize converts CSS size strings to float64
// Supports: "16", "16px", "2em", "1.5rem", "50%"
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
// Supports: "10px", "10px 20px", "10px 20px 30px", "10px 20px 30px 40px"
// Returns top, right, bottom, left
func parseBoxDimensions(sizeStr string, current BoxDimensions) BoxDimensions {
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

// isValidAlign checks if alignment value is valid
func isValidAlign(align string) bool {
	validAligns := map[string]bool{
		"left":    true,
		"center":  true,
		"right":   true,
		"justify": true,
	}
	return validAligns[strings.ToLower(align)]
}

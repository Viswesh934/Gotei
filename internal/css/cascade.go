package css

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Viswesh934/gotei/internal/debug"
	"github.com/Viswesh934/gotei/internal/dom"
	"github.com/Viswesh934/gotei/internal/style"
	wrselector "github.com/benoitkugler/webrender/css/selector"
	"golang.org/x/net/html"
)

// ─────────────────────────────────────────────────────────────────────────────
// Specificity
// ─────────────────────────────────────────────────────────────────────────────

// Specificity holds a CSS specificity triple (IDs, Classes, Elements).
type Specificity struct {
	IDs      int
	Classes  int
	Elements int
}

// Compare returns -1, 0, or 1.
func (s Specificity) Compare(other Specificity) int {
	switch {
	case s.IDs != other.IDs:
		if s.IDs > other.IDs {
			return 1
		}
		return -1
	case s.Classes != other.Classes:
		if s.Classes > other.Classes {
			return 1
		}
		return -1
	case s.Elements != other.Elements:
		if s.Elements > other.Elements {
			return 1
		}
		return -1
	}
	return 0
}

// ─────────────────────────────────────────────────────────────────────────────
// Rule / StyleSheet
// ─────────────────────────────────────────────────────────────────────────────

// Rule represents a CSS rule with its selector, parsed properties, and position.
type Rule struct {
	Matcher     wrselector.Matcher
	Properties  map[string]string
	Specificity Specificity
	SourceOrder int
}

// StyleSheet holds all parsed CSS rules for a document.
type StyleSheet struct {
	Rules []*Rule
}

// ─────────────────────────────────────────────────────────────────────────────
// Cascade entry point
// ─────────────────────────────────────────────────────────────────────────────

// ComputeStyleWithHTMLNode calculates the final computed style for a node by
// following the standard CSS cascade order:
//
//  1. DefaultStyle (browser-like UA defaults)
//  2. TagDefaults  (element-level defaults, e.g. h1, p, td)
//  3. Stylesheet rules, sorted by specificity then source order
//  4. HTML presentation attributes (align=, width=, border=, …)
//  5. Inline style="" attribute  (highest author specificity)
//  6. Inherited properties from parent
func ComputeStyleWithHTMLNode(node *dom.Node, htmlNode *html.Node, sheet *StyleSheet, parentStyle style.Style) style.Style {

	// 1. Start from UA defaults
	computed := style.DefaultStyle

	// 2. Merge tag-level defaults
	if node.Type == dom.ElementNode {
		if tagStyle, ok := style.TagDefaults[node.Tag]; ok {
			computed = mergeTagStyle(computed, tagStyle)
		}
	}

	// 3. Apply matching stylesheet rules, lowest → highest specificity
	matched := collectMatchingRules(htmlNode, sheet)
	sortRulesBySpecificity(matched)
	debug.Logf("css.cascade: node=%s matched-rules=%d", describeNode(node), len(matched))
	for _, rule := range matched {
		computed = applyProperties(computed, rule.Properties)
	}

	// 4. HTML presentation attributes
	if node.Type == dom.ElementNode {
		computed = applyHTMLAttributes(computed, node)
	}

	// 5. Inline style attribute
	if node.Type == dom.ElementNode {
		computed = applyInlineStyle(computed, node)
	}

	// 6. CSS inheritance from parent
	computed = applyInheritance(computed, parentStyle)

	// Text nodes carry no box model — strip it after everything else so that
	// neither tag defaults nor inheritance can accidentally leave values on them.
	if node.Type == dom.TextNode {
		computed.Margin  = style.BoxDimensions{}
		computed.Padding = style.BoxDimensions{}
		computed.Border  = style.Border{}
	}

	debug.Logf("css.cascade: node=%s computed={display=%s font=%s %.1fpx color=%s bg=%q align=%s width=%.0f%%}",
		describeNode(node),
		computed.Display,
		computed.FontFamily,
		computed.FontSize,
		computed.Color,
		computed.BgColor,
		computed.TextAlign,
		computed.WidthPercent,
	)

	return computed
}

// ─────────────────────────────────────────────────────────────────────────────
// Rule collection + sorting
// ─────────────────────────────────────────────────────────────────────────────

func collectMatchingRules(htmlNode *html.Node, sheet *StyleSheet) []*Rule {
	if sheet == nil || htmlNode == nil {
		return nil
	}
	var matched []*Rule
	for _, rule := range sheet.Rules {
		if rule.Matcher != nil && rule.Matcher.Match(htmlNode) {
			matched = append(matched, rule)
		}
	}
	return matched
}

func sortRulesBySpecificity(rules []*Rule) {
	sort.SliceStable(rules, func(i, j int) bool {
		cmp := rules[i].Specificity.Compare(rules[j].Specificity)
		if cmp != 0 {
			return cmp < 0
		}
		return rules[i].SourceOrder < rules[j].SourceOrder
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// mergeTagStyle — applies TagDefaults on top of DefaultStyle
//
// Rules:
//   - Scalar fields (FontSize, FontWeight, …) are only overridden when the
//     tag explicitly sets them (non-zero / non-empty).
//   - Box model fields (Margin, Padding, Border) are replaced only when the
//     tag entry carries a non-zero value in at least one dimension. This
//     prevents an empty BoxDimensions{} in a tag entry from wiping out the
//     DefaultStyle values.
//   - WidthPercent is propagated even when 0 ONLY for inline/inline-block
//     elements, where 0 is the correct value. Block elements that don't set it
//     keep the DefaultStyle value of 100.
// ─────────────────────────────────────────────────────────────────────────────

func mergeTagStyle(base, tag style.Style) style.Style {
	// Typography
	if tag.FontSize != 0 {
		base.FontSize = tag.FontSize
	}
	if tag.FontFamily != "" {
		base.FontFamily = tag.FontFamily
	}
	if tag.FontWeight != "" {
		base.FontWeight = tag.FontWeight
	}
	if tag.FontStyle != "" {
		base.FontStyle = tag.FontStyle
	}
	if tag.LineHeight != 0 {
		base.LineHeight = tag.LineHeight
	}
	if tag.LetterSpacing != 0 {
		base.LetterSpacing = tag.LetterSpacing
	}
	if tag.WordSpacing != 0 {
		base.WordSpacing = tag.WordSpacing
	}
	if tag.TextTransform != "" {
		base.TextTransform = tag.TextTransform
	}
	if tag.TextDecoration != "" {
		base.TextDecoration = tag.TextDecoration
	}

	// Alignment & flow
	if tag.TextAlign != "" {
		base.TextAlign = tag.TextAlign
		base.Align = tag.TextAlign // keep legacy in sync
	}
	if tag.VerticalAlign != "" {
		base.VerticalAlign = tag.VerticalAlign
	}
	if tag.WhiteSpace != "" {
		base.WhiteSpace = tag.WhiteSpace
	}
	if tag.WordBreak != "" {
		base.WordBreak = tag.WordBreak
	}
	if tag.WordWrap != "" {
		base.WordWrap = tag.WordWrap
	}

	// Color & background
	if tag.Color != "" {
		base.Color = tag.Color
	}
	if tag.BgColor != "" {
		base.BgColor = tag.BgColor
	}
	if tag.Background.Color != "" {
		base.Background.Color = tag.Background.Color
	}
	if tag.Opacity != 0 {
		base.Opacity = tag.Opacity
	}

	// Display & layout
	if tag.Display != "" {
		base.Display = tag.Display
	}
	if tag.Position != "" {
		base.Position = tag.Position
	}
	if tag.Overflow != "" {
		base.Overflow = tag.Overflow
	}
	if tag.Float != "" {
		base.Float = tag.Float
	}
	if tag.Clear != "" {
		base.Clear = tag.Clear
	}
	if tag.ZIndex != 0 {
		base.ZIndex = tag.ZIndex
	}

	// Width — inline elements intentionally have WidthPercent=0; block elements
	// that don't set it keep the DefaultStyle value (100).
	isInline := tag.Display == "inline" || tag.Display == "inline-block"
	if tag.WidthPercent > 0 {
		base.WidthPercent = tag.WidthPercent
		base.Width = 0
	} else if isInline {
		// Explicitly clear block width for inline elements
		base.WidthPercent = 0
		base.Width = 0
	}
	// else: block element that didn't set WidthPercent → keep base (100%)

	if tag.Width > 0 {
		base.Width = tag.Width
		base.WidthPercent = 0
	}
	if tag.Height > 0 {
		base.Height = tag.Height
		base.HeightPercent = 0
	}
	if tag.HeightPercent > 0 {
		base.HeightPercent = tag.HeightPercent
		base.Height = 0
	}
	if tag.MaxWidth > 0 {
		base.MaxWidth = tag.MaxWidth
	}
	if tag.MinWidth > 0 {
		base.MinWidth = tag.MinWidth
	}
	if tag.MaxHeight > 0 {
		base.MaxHeight = tag.MaxHeight
	}
	if tag.MinHeight > 0 {
		base.MinHeight = tag.MinHeight
	}

	// Box model — only replace if the tag actually specifies something.
	// A completely zero BoxDimensions{} in a tag entry means "not set",
	// so we keep whatever DefaultStyle had.
	if !isZeroBox(tag.Margin) {
		base.Margin = tag.Margin
	}
	if !isZeroBox(tag.Padding) {
		base.Padding = tag.Padding
	}
	if !isZeroBorder(tag.Border) {
		base.Border = tag.Border
	}

	// Flexbox (only relevant when Display == "flex")
	if tag.FlexDirection != "" {
		base.FlexDirection = tag.FlexDirection
	}
	if tag.FlexWrap != "" {
		base.FlexWrap = tag.FlexWrap
	}
	if tag.JustifyContent != "" {
		base.JustifyContent = tag.JustifyContent
	}
	if tag.AlignItems != "" {
		base.AlignItems = tag.AlignItems
	}
	if tag.AlignContent != "" {
		base.AlignContent = tag.AlignContent
	}
	if tag.Gap != 0 {
		base.Gap = tag.Gap
	}

	// Legacy
	if tag.Bold {
		base.Bold = true
	}
	if tag.Italic {
		base.Italic = true
	}

	return base
}

// ─────────────────────────────────────────────────────────────────────────────
// applyProperties — applies a map of parsed CSS declarations to a style
// ─────────────────────────────────────────────────────────────────────────────

func applyProperties(s style.Style, props map[string]string) style.Style {
	for key, val := range props {
		s = applyProperty(s, key, val)
	}
	return s
}

func applyProperty(s style.Style, key, val string) style.Style {
	key = strings.TrimSpace(key)
	val = strings.TrimSpace(val)

	switch key {
	// ── Typography ────────────────────────────────────────────────────────────
	case "font-size":
		if v, err := parseSize(val); err == nil {
			s.FontSize = v
		}
	case "font-family":
		s.FontFamily = normalizeFontFamily(val)
	case "font-weight":
		s.FontWeight = val
		s.Bold = val == "bold" || val == "700" || val == "800" || val == "900"
	case "font-style":
		s.FontStyle = val
		s.Italic = val == "italic" || val == "oblique"
	case "line-height":
		if v, err := parseSize(val); err == nil {
			s.LineHeight = v
		}
	case "letter-spacing":
		if v, err := parseSize(val); err == nil {
			s.LetterSpacing = v
		}
	case "word-spacing":
		if v, err := parseSize(val); err == nil {
			s.WordSpacing = v
		}
	case "text-transform":
		s.TextTransform = val
	case "text-decoration":
		s.TextDecoration = val
	case "text-shadow":
		s.TextShadow = parseShadow(val)

	// ── Alignment & flow ──────────────────────────────────────────────────────
	case "text-align":
		s.TextAlign = val
		s.Align = val
	case "vertical-align":
		s.VerticalAlign = val
	case "direction":
		s.Direction = val
	case "white-space":
		s.WhiteSpace = val
	case "word-break":
		s.WordBreak = val
	case "word-wrap", "overflow-wrap":
		s.WordWrap = val

	// ── Color & background ────────────────────────────────────────────────────
	case "color":
		s.Color = val
	case "background-color":
		s.BgColor = val
		s.Background.Color = val
	case "background":
		// Treat a plain color value as background-color; ignore image/gradient.
		if !strings.Contains(val, "(") {
			s.BgColor = val
			s.Background.Color = val
		}
	case "opacity":
		if v, err := parseSize(val); err == nil {
			s.Opacity = v
		}

	// ── Display & layout ──────────────────────────────────────────────────────
	case "display":
		s.Display = val
	case "position":
		s.Position = val
	case "visibility":
		s.Visibility = val
	case "overflow":
		s.Overflow = val
	case "overflow-x":
		s.OverflowX = val
	case "overflow-y":
		s.OverflowY = val
	case "float":
		s.Float = val
	case "clear":
		s.Clear = val
	case "z-index":
		if v, err := strconv.Atoi(val); err == nil {
			s.ZIndex = v
		}
	case "top":
		if v, err := parseSize(val); err == nil {
			s.Top = v
		}
	case "right":
		if v, err := parseSize(val); err == nil {
			s.Right = v
		}
	case "bottom":
		if v, err := parseSize(val); err == nil {
			s.Bottom = v
		}
	case "left":
		if v, err := parseSize(val); err == nil {
			s.Left = v
		}

	// ── Box model ─────────────────────────────────────────────────────────────
	case "width":
		if strings.HasSuffix(val, "%") {
			if v, err := parsePercent(val); err == nil {
				s.WidthPercent = v
				s.Width = 0
			}
		} else if v, err := parseSize(val); err == nil {
			s.Width = v
			s.WidthPercent = 0
		}
	case "height":
		if strings.HasSuffix(val, "%") {
			if v, err := parsePercent(val); err == nil {
				s.HeightPercent = v
				s.Height = 0
			}
		} else if v, err := parseSize(val); err == nil {
			s.Height = v
			s.HeightPercent = 0
		}
	case "max-width":
		if v, err := parseSize(val); err == nil {
			s.MaxWidth = v
		}
	case "min-width":
		if v, err := parseSize(val); err == nil {
			s.MinWidth = v
		}
	case "max-height":
		if v, err := parseSize(val); err == nil {
			s.MaxHeight = v
		}
	case "min-height":
		if v, err := parseSize(val); err == nil {
			s.MinHeight = v
		}
	case "box-sizing":
		s.BoxSizing = val

	case "margin":
		s.Margin = parseBox(val, s.Margin)
	case "margin-top":
		if v, err := parseSize(val); err == nil {
			s.Margin.Top = v
		}
	case "margin-right":
		if v, err := parseSize(val); err == nil {
			s.Margin.Right = v
		}
	case "margin-bottom":
		if v, err := parseSize(val); err == nil {
			s.Margin.Bottom = v
		}
	case "margin-left":
		if v, err := parseSize(val); err == nil {
			s.Margin.Left = v
		}

	case "padding":
		s.Padding = parseBox(val, s.Padding)
	case "padding-top":
		if v, err := parseSize(val); err == nil {
			s.Padding.Top = v
		}
	case "padding-right":
		if v, err := parseSize(val); err == nil {
			s.Padding.Right = v
		}
	case "padding-bottom":
		if v, err := parseSize(val); err == nil {
			s.Padding.Bottom = v
		}
	case "padding-left":
		if v, err := parseSize(val); err == nil {
			s.Padding.Left = v
		}

	case "border":
		s.Border = parseBorder(val, s.Border)
	case "border-width":
		if v, err := parseSize(val); err == nil {
			s.Border.Width = v
		}
	case "border-style":
		s.Border.Style = val
	case "border-color":
		s.Border.Color = val
	case "border-radius":
		if v, err := parseSize(val); err == nil {
			s.Border.Radius = v
		}
	case "border-top-left-radius":
		if v, err := parseSize(val); err == nil {
			s.Border.TopLeft = v
		}
	case "border-top-right-radius":
		if v, err := parseSize(val); err == nil {
			s.Border.TopRight = v
		}
	case "border-bottom-left-radius":
		if v, err := parseSize(val); err == nil {
			s.Border.BottomLeft = v
		}
	case "border-bottom-right-radius":
		if v, err := parseSize(val); err == nil {
			s.Border.BottomRight = v
		}

	// ── Flexbox ───────────────────────────────────────────────────────────────
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
	case "flex-grow":
		if v, err := parseSize(val); err == nil {
			s.FlexGrow = v
		}
	case "flex-shrink":
		if v, err := parseSize(val); err == nil {
			s.FlexShrink = v
		}
	case "flex-basis":
		if v, err := parseSize(val); err == nil {
			s.FlexBasis = v
		}
	case "flex":
		s.Flex = val
	case "gap":
		if v, err := parseSize(val); err == nil {
			s.Gap = v
		}

	// ── Grid ──────────────────────────────────────────────────────────────────
	case "grid-template-columns":
		s.GridTemplateColumns = val
	case "grid-template-rows":
		s.GridTemplateRows = val
	case "grid-gap", "gap-grid":
		if v, err := parseSize(val); err == nil {
			s.GridGap = v
		}

	// ── Visual effects ────────────────────────────────────────────────────────
	case "box-shadow":
		s.BoxShadow = parseShadow(val)
	case "filter":
		s.Filter = val
	case "transform":
		s.Transform = val
	case "transition":
		s.Transition = val
	case "animation":
		s.Animation = val
	}

	return s
}

// ─────────────────────────────────────────────────────────────────────────────
// applyHTMLAttributes — maps legacy HTML presentation attributes to style
// ─────────────────────────────────────────────────────────────────────────────

func applyHTMLAttributes(s style.Style, node *dom.Node) style.Style {
	if node.Attr == nil {
		return s
	}

	if v, ok := node.Attr["align"]; ok {
		s.TextAlign = v
		s.Align = v
	}
	if v, ok := node.Attr["valign"]; ok {
		s.VerticalAlign = v
	}
	if v, ok := node.Attr["width"]; ok {
		if strings.HasSuffix(strings.TrimSpace(v), "%") {
			if p, err := parsePercent(v); err == nil {
				s.WidthPercent = p
				s.Width = 0
			}
		} else if w, err := parseSize(v); err == nil {
			s.Width = w
			s.WidthPercent = 0
		}
	}
	if v, ok := node.Attr["height"]; ok {
		if strings.HasSuffix(strings.TrimSpace(v), "%") {
			if p, err := parsePercent(v); err == nil {
				s.HeightPercent = p
				s.Height = 0
			}
		} else if h, err := parseSize(v); err == nil {
			s.Height = h
			s.HeightPercent = 0
		}
	}
	if v, ok := node.Attr["border"]; ok {
		if w, err := parseSize(v); err == nil && w > 0 {
			s.Border.Width = w
			s.Border.Style = "solid"
			if s.Border.Color == "" {
				s.Border.Color = "#000000"
			}
		}
	}
	if v, ok := node.Attr["cellpadding"]; ok {
		if p, err := parseSize(v); err == nil {
			s.Padding = style.BoxDimensions{Top: p, Right: p, Bottom: p, Left: p}
		}
	}
	if v, ok := node.Attr["cellspacing"]; ok {
		if p, err := parseSize(v); err == nil {
			s.Gap = p
		}
	}
	if v, ok := node.Attr["bgcolor"]; ok && v != "" {
		s.BgColor = v
		s.Background.Color = v
	}
	if v, ok := node.Attr["color"]; ok && v != "" {
		s.Color = v
	}

	return s
}

// ─────────────────────────────────────────────────────────────────────────────
// applyInlineStyle — parses the style="" attribute and applies each declaration
// ─────────────────────────────────────────────────────────────────────────────

func applyInlineStyle(s style.Style, node *dom.Node) style.Style {
	if node.Attr == nil {
		return s
	}
	styleStr, ok := node.Attr["style"]
	if !ok || styleStr == "" {
		return s
	}

	for _, decl := range strings.Split(styleStr, ";") {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}
		parts := strings.SplitN(decl, ":", 2)
		if len(parts) != 2 {
			continue
		}
		s = applyProperty(s, strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	return s
}

// ─────────────────────────────────────────────────────────────────────────────
// applyInheritance — propagates inheritable CSS properties from parent
//
// Only propagates when the child is still carrying the DefaultStyle value,
// meaning it hasn't been explicitly set by any rule or tag default.
// Uses DefaultStyle as the "unset" sentinel to avoid incorrectly overwriting
// deliberately set values.
// ─────────────────────────────────────────────────────────────────────────────

func applyInheritance(s, parent style.Style) style.Style {
	d := style.DefaultStyle

	// Font
	if s.FontSize == d.FontSize && parent.FontSize != d.FontSize {
		s.FontSize = parent.FontSize
	}
	if (s.FontFamily == "" || s.FontFamily == d.FontFamily) && parent.FontFamily != "" && parent.FontFamily != d.FontFamily {
		s.FontFamily = parent.FontFamily
	}
	if (s.FontWeight == "" || s.FontWeight == d.FontWeight) && parent.FontWeight != "" && parent.FontWeight != d.FontWeight {
		s.FontWeight = parent.FontWeight
		if parent.Bold {
			s.Bold = true
		}
	}
	if (s.FontStyle == "" || s.FontStyle == d.FontStyle) && parent.FontStyle != "" && parent.FontStyle != d.FontStyle {
		s.FontStyle = parent.FontStyle
		if parent.Italic {
			s.Italic = true
		}
	}
	if s.LineHeight == d.LineHeight && parent.LineHeight != d.LineHeight {
		s.LineHeight = parent.LineHeight
	}
	if s.LetterSpacing == d.LetterSpacing && parent.LetterSpacing != d.LetterSpacing {
		s.LetterSpacing = parent.LetterSpacing
	}
	if s.WordSpacing == d.WordSpacing && parent.WordSpacing != d.WordSpacing {
		s.WordSpacing = parent.WordSpacing
	}

	// Text
	if (s.TextTransform == "" || s.TextTransform == d.TextTransform) && parent.TextTransform != "" && parent.TextTransform != d.TextTransform {
		s.TextTransform = parent.TextTransform
	}
	if (s.TextDecoration == "" || s.TextDecoration == d.TextDecoration) && parent.TextDecoration != "" && parent.TextDecoration != d.TextDecoration {
		s.TextDecoration = parent.TextDecoration
	}
	if (s.TextAlign == "" || s.TextAlign == d.TextAlign) && parent.TextAlign != "" && parent.TextAlign != d.TextAlign {
		s.TextAlign = parent.TextAlign
		s.Align = parent.TextAlign
	}

	// Color — compare against DefaultStyle's value, not the string "black"
	if s.Color == d.Color && parent.Color != d.Color {
		s.Color = parent.Color
	}

	// Direction
	if (s.Direction == "" || s.Direction == d.Direction) && parent.Direction != "" && parent.Direction != d.Direction {
		s.Direction = parent.Direction
	}

	// Note: box model (margin, padding, border), display, width, background
	// are NOT inherited in CSS — we intentionally do not propagate them.

	return s
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func describeNode(node *dom.Node) string {
	if node == nil {
		return "<nil>"
	}
	if node.Type == dom.TextNode {
		t := strings.TrimSpace(node.Content)
		if t == "" {
			return "text[(whitespace)]"
		}
		if len(t) > 30 {
			t = t[:30] + "…"
		}
		return fmt.Sprintf("text[%q]", t)
	}
	if node.Tag == "" {
		return "element[unknown]"
	}
	if id := node.Attr["id"]; id != "" {
		return fmt.Sprintf("%s#%s", node.Tag, id)
	}
	if cls := node.Attr["class"]; cls != "" {
		return fmt.Sprintf("%s.%s", node.Tag, strings.ReplaceAll(strings.TrimSpace(cls), " ", "."))
	}
	return node.Tag
}

func normalizeFontFamily(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	// Take only the first family in a comma-separated list
	if idx := strings.IndexByte(raw, ','); idx >= 0 {
		raw = raw[:idx]
	}
	return strings.Trim(strings.TrimSpace(raw), `"'`)
}

// parseSize converts a CSS value string to float64, stripping any unit suffix.
func parseSize(raw string) (float64, error) {
	raw = strings.TrimSpace(raw)
	for _, unit := range []string{"px", "pt", "em", "rem", "vh", "vw", "%"} {
		raw = strings.TrimSuffix(raw, unit)
	}
	return strconv.ParseFloat(strings.TrimSpace(raw), 64)
}

func parsePercent(raw string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(raw), "%")), 64)
}

// parseBox parses CSS shorthand for margin/padding (1–4 values).
func parseBox(raw string, current style.BoxDimensions) style.BoxDimensions {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return current
	}

	vals := make([]float64, 0, len(fields))
	for _, f := range fields {
		if v, err := parseSize(f); err == nil {
			vals = append(vals, v)
		}
	}

	switch len(vals) {
	case 1:
		return style.BoxDimensions{Top: vals[0], Right: vals[0], Bottom: vals[0], Left: vals[0]}
	case 2:
		return style.BoxDimensions{Top: vals[0], Right: vals[1], Bottom: vals[0], Left: vals[1]}
	case 3:
		return style.BoxDimensions{Top: vals[0], Right: vals[1], Bottom: vals[2], Left: vals[1]}
	case 4:
		return style.BoxDimensions{Top: vals[0], Right: vals[1], Bottom: vals[2], Left: vals[3]}
	}
	return current
}

// parseBorder parses CSS border shorthand: "1px solid #cccccc"
func parseBorder(raw string, current style.Border) style.Border {
	b := current
	for _, part := range strings.Fields(raw) {
		switch part {
		case "solid", "dashed", "dotted", "double", "none":
			b.Style = part
		default:
			if v, err := parseSize(part); err == nil {
				b.Width = v
			} else {
				b.Color = part
			}
		}
	}
	return b
}

// parseShadow parses CSS box-shadow / text-shadow: "2px 2px 4px 0px #000"
func parseShadow(raw string) style.Shadow {
	var s style.Shadow
	var nums []float64
	for _, part := range strings.Fields(raw) {
		if part == "inset" {
			s.Inset = true
			continue
		}
		if v, err := parseSize(part); err == nil {
			nums = append(nums, v)
		} else {
			s.Color = part
		}
	}
	if len(nums) >= 1 { s.OffsetX = nums[0] }
	if len(nums) >= 2 { s.OffsetY = nums[1] }
	if len(nums) >= 3 { s.Blur    = nums[2] }
	if len(nums) >= 4 { s.Spread  = nums[3] }
	return s
}

// isZeroBox reports whether all dimensions are zero (i.e. "not set" in a tag entry).
func isZeroBox(b style.BoxDimensions) bool {
	return b.Top == 0 && b.Right == 0 && b.Bottom == 0 && b.Left == 0
}

// isZeroBorder reports whether a Border is completely unset.
func isZeroBorder(b style.Border) bool {
	return b.Width == 0 && b.Style == "" && b.Color == ""
}
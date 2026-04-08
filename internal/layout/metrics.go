package layout

import (
	"github.com/Viswesh934/gotei/internal/style"
)

// FontMetrics represents measurements for a font
type FontMetrics struct {
	Height    float64
	Ascender  float64
	Descender float64
	LineGap   float64
	CapHeight float64
}

// EstimateFontMetrics returns approximate metrics for a font
// Based on common typography ratios
func EstimateFontMetrics(fontSize float64, fontFamily string) FontMetrics {
	// Typical ratios relative to font size
	// These are approximations; real values depend on the font
	var ascenderRatio, descenderRatio, capHeightRatio float64

	switch fontFamily {
	case "Arial", "Helvetica", "sans-serif":
		// Sans serif typical ratios
		ascenderRatio = 0.75
		descenderRatio = -0.25
		capHeightRatio = 0.68
	case "Times", "Georgia", "serif":
		// Serif typical ratios
		ascenderRatio = 0.72
		descenderRatio = -0.28
		capHeightRatio = 0.65
	case "Courier", "monospace":
		// Monospace typical ratios
		ascenderRatio = 0.75
		descenderRatio = -0.25
		capHeightRatio = 0.70
	default:
		// Fallback
		ascenderRatio = 0.75
		descenderRatio = -0.25
		capHeightRatio = 0.68
	}

	return FontMetrics{
		Height:    fontSize,
		Ascender:  fontSize * ascenderRatio,
		Descender: fontSize * descenderRatio,
		CapHeight: fontSize * capHeightRatio,
		LineGap:   fontSize * 0.15,
	}
}

// CalculateLineHeight returns the total line height based on the style
func CalculateLineHeight(style style.Style) float64 {
	if style.LineHeight == 0 {
		// Default 1.5x line height
		return style.FontSize * 1.5
	}

	// If line-height is 1 or greater, treat as multiplier
	if style.LineHeight >= 1 {
		return style.FontSize * style.LineHeight
	}

	// Otherwise treat as absolute size
	return style.LineHeight
}

// MeasureText estimates the width of text in pixels
// This is a rough approximation; real measurement requires font metrics
func MeasureText(text string, fontSize float64, fontFamily string) float64 {
	// Rough character width estimation
	charWidth := fontSize * 0.5 // Approximate for most fonts

	// Adjust for monospace
	if fontFamily == "Courier" || fontFamily == "monospace" {
		charWidth = fontSize * 0.6
	}

	return float64(len(text)) * charWidth
}

// EstimateTextHeight estimates the height needed to render text
func EstimateTextHeight(text string, maxWidth float64, fontSize float64, lineHeight float64) float64 {
	if maxWidth <= 0 {
		return lineHeight
	}

	charWidth := fontSize * 0.5
	charsPerLine := int(maxWidth / charWidth)
	if charsPerLine < 1 {
		charsPerLine = 1
	}

	lines := (len(text) + charsPerLine - 1) / charsPerLine
	if lines < 1 {
		lines = 1
	}

	return float64(lines) * lineHeight
}

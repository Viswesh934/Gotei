package style

// CSSUnit represents a CSS size value with unit
type CSSUnit struct {
	Value float64
	Unit  string // "px", "pt", "em", "rem", "%", etc.
}

// BoxDimensions represents the CSS box model
type BoxDimensions struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// Background represents background styling
type Background struct {
	Color string
	Image string
	Size  string
}

// Border represents border styling
type Border struct {
	Width       float64
	Style       string // "solid", "dashed", "dotted"
	Color       string
	Radius      float64
	TopLeft     float64
	TopRight    float64
	BottomLeft  float64
	BottomRight float64
}

// Shadow represents text or box shadow
type Shadow struct {
	OffsetX float64
	OffsetY float64
	Blur    float64
	Spread  float64
	Color   string
	Inset   bool
}

// Style represents comprehensive CSS styling properties (W3C compliant)
type Style struct {
	// TEXT & FONT
	FontSize       float64 // in pixels
	FontFamily     string  // "Arial", "Georgia", "Courier", etc.
	FontWeight     string  // "normal", "bold", "700", "800", etc.
	FontStyle      string  // "normal", "italic", "oblique"
	LineHeight     float64 // relative to font size
	LetterSpacing  float64
	WordSpacing    float64
	TextTransform  string // "none", "uppercase", "lowercase", "capitalize"
	TextDecoration string // "none", "underline", "overline", "line-through"
	TextShadow     Shadow

	// TEXT ALIGNMENT & FLOW
	TextAlign     string // "left", "center", "right", "justify"
	VerticalAlign string // "top", "middle", "bottom", "baseline"
	Direction     string // "ltr", "rtl"
	WhiteSpace    string // "normal", "nowrap", "pre", "pre-wrap"
	WordBreak     string // "normal", "break-all", "keep-all"
	WordWrap      string // "normal", "break-word"

	// COLOR & BACKGROUND
	Color      string
	Opacity    float64 // 0.0 to 1.0
	Background Background
	BgColor    string // shorthand/legacy

	// BOX MODEL
	Margin    BoxDimensions
	Padding   BoxDimensions
	Border    Border
	Width     float64
	Height    float64
	MaxWidth  float64
	MinWidth  float64
	MaxHeight float64
	MinHeight float64
	BoxSizing string // "content-box", "border-box"

	// DISPLAY & LAYOUT
	Display    string // "block", "inline", "inline-block", "flex", "grid", "none"
	Visibility string // "visible", "hidden", "collapse"
	Overflow   string // "visible", "hidden", "scroll", "auto"
	OverflowX  string
	OverflowY  string
	Position   string // "static", "relative", "absolute", "fixed", "sticky"
	Top        float64
	Right      float64
	Bottom     float64
	Left       float64
	ZIndex     int
	Float      string // "none", "left", "right"
	Clear      string // "none", "left", "right", "both"

	// FLEXBOX
	FlexDirection  string // "row", "column", "row-reverse", "column-reverse"
	FlexWrap       string // "nowrap", "wrap", "wrap-reverse"
	JustifyContent string // "flex-start", "center", "flex-end", "space-between", "space-around"
	AlignItems     string // "flex-start", "center", "flex-end", "stretch", "baseline"
	AlignContent   string
	FlexGrow       float64
	FlexShrink     float64
	FlexBasis      float64
	Flex           string
	Gap            float64

	// GRID
	GridTemplateColumns string
	GridTemplateRows    string
	GridColumn          int
	GridRow             int
	GridGap             float64

	// VISUAL EFFECTS
	BoxShadow  Shadow
	Filter     string
	Transform  string
	Transition string
	Animation  string

	// DEPRECATED/LEGACY (for backwards compatibility)
	Bold   bool
	Italic bool
	Align  string // legacy, use TextAlign instead
}

// Default style for all elements (W3C compliant)
var DefaultStyle = Style{
	// Text & Font
	FontSize:       12,
	FontFamily:     "Arial",
	FontWeight:     "normal",
	FontStyle:      "normal",
	LineHeight:     1.5,
	LetterSpacing:  0,
	WordSpacing:    0,
	TextTransform:  "none",
	TextDecoration: "none",

	// Text Alignment
	TextAlign:     "left",
	VerticalAlign: "baseline",
	Direction:     "ltr",
	WhiteSpace:    "normal",
	WordBreak:     "normal",
	WordWrap:      "normal",

	// Color & Background
	Color:   "black",
	Opacity: 1.0,
	BgColor: "white",
	Background: Background{
		Color: "white",
		Size:  "auto",
	},

	// Box Model
	Margin: BoxDimensions{
		Top:    4,
		Right:  4,
		Bottom: 4,
		Left:   4,
	},
	Padding: BoxDimensions{
		Top:    2,
		Right:  2,
		Bottom: 2,
		Left:   2,
	},
	Border: Border{
		Width: 0,
		Style: "solid",
		Color: "black",
	},
	Width:     0,
	Height:    0,
	MaxWidth:  0,
	MinWidth:  0,
	MaxHeight: 0,
	MinHeight: 0,
	BoxSizing: "content-box",

	// Display & Layout
	Display:    "block",
	Visibility: "visible",
	Overflow:   "visible",
	OverflowX:  "visible",
	OverflowY:  "visible",
	Position:   "static",
	ZIndex:     0,
	Float:      "none",
	Clear:      "none",

	// Flexbox
	FlexDirection:  "row",
	FlexWrap:       "nowrap",
	JustifyContent: "flex-start",
	AlignItems:     "stretch",
	AlignContent:   "stretch",
	FlexGrow:       0,
	FlexShrink:     1,
	Gap:            0,

	// Grid
	GridGap: 0,

	// Legacy
	Bold:   false,
	Italic: false,
	Align:  "left",
}

// InheritableFields defines which CSS properties cascade from parent to child
// per W3C spec (most text properties inherit, box model generally doesn't)
var InheritableFields = map[string]bool{
	// Inherited by default
	"FontSize":       true,
	"FontFamily":     true,
	"FontWeight":     true,
	"FontStyle":      true,
	"LineHeight":     true,
	"LetterSpacing":  true,
	"WordSpacing":    true,
	"TextTransform":  true,
	"TextDecoration": true,
	"TextAlign":      true,
	"VerticalAlign":  true,
	"Direction":      true,
	"WhiteSpace":     true,
	"WordBreak":      true,
	"WordWrap":       true,
	"Color":          true,
	"Opacity":        true,

	// Not inherited
	"Margin":   false,
	"Padding":  false,
	"Border":   false,
	"Width":    false,
	"Height":   false,
	"Display":  false,
	"Position": false,
	"Float":    false,

	// Legacy
	"Bold":   true,
	"Italic": true,
	"Align":  true,
}

// TagDefaults maps HTML tags to their default styles (W3C compliant)
var TagDefaults = map[string]Style{
	"h1": {
		FontSize:   28,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    16,
			Bottom: 16,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    8,
			Bottom: 8,
			Left:   8,
			Right:  8,
		},
	},
	"h2": {
		FontSize:   24,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    14,
			Bottom: 14,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    6,
			Bottom: 6,
			Left:   6,
			Right:  6,
		},
	},
	"h3": {
		FontSize:   20,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    12,
			Bottom: 12,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    4,
			Bottom: 4,
			Left:   4,
			Right:  4,
		},
	},
	"h4": {
		FontSize:   18,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    10,
			Bottom: 10,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    4,
			Bottom: 4,
			Left:   4,
			Right:  4,
		},
	},
	"h5": {
		FontSize:   16,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    8,
			Bottom: 8,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    3,
			Bottom: 3,
			Left:   3,
			Right:  3,
		},
	},
	"h6": {
		FontSize:   14,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    6,
			Bottom: 6,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    2,
			Bottom: 2,
			Left:   2,
			Right:  2,
		},
	},
	"p": {
		FontSize: 12,
		Margin: BoxDimensions{
			Top:    8,
			Bottom: 8,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    2,
			Bottom: 2,
			Left:   2,
			Right:  2,
		},
	},
	"div": {
		FontSize: 12,
		Display:  "block",
		Margin: BoxDimensions{
			Top:    6,
			Bottom: 6,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    2,
			Bottom: 2,
			Left:   2,
			Right:  2,
		},
	},
	"span": {
		FontSize: 12,
		Display:  "inline",
		Margin:   BoxDimensions{},
		Padding:  BoxDimensions{},
	},
	"strong": {
		FontSize:   12,
		FontWeight: "bold",
		Bold:       true,
		Display:    "inline",
		Margin:     BoxDimensions{},
		Padding:    BoxDimensions{},
	},
	"em": {
		FontSize:  12,
		FontStyle: "italic",
		Italic:    true,
		Display:   "inline",
		Margin:    BoxDimensions{},
		Padding:   BoxDimensions{},
	},
	"b": {
		FontSize:   12,
		FontWeight: "bold",
		Bold:       true,
		Display:    "inline",
		Margin:     BoxDimensions{},
		Padding:    BoxDimensions{},
	},
	"i": {
		FontSize:  12,
		FontStyle: "italic",
		Italic:    true,
		Display:   "inline",
		Margin:    BoxDimensions{},
		Padding:   BoxDimensions{},
	},
	"code": {
		FontSize:   10,
		FontFamily: "Courier",
		Display:    "inline",
		BgColor:    "lightgray",
		Padding: BoxDimensions{
			Top:    2,
			Bottom: 2,
			Left:   2,
			Right:  2,
		},
	},
	"pre": {
		FontSize:   10,
		FontFamily: "Courier",
		Display:    "block",
		WhiteSpace: "pre",
		Margin: BoxDimensions{
			Top:    8,
			Bottom: 8,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{
			Top:    8,
			Bottom: 8,
			Left:   8,
			Right:  8,
		},
		BgColor: "lightgray",
	},
}

// ClassStyles maps CSS class names to their styles
var ClassStyles = map[string]Style{
	"title": {
		FontSize:   24,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    12,
			Bottom: 12,
			Left:   0,
			Right:  0,
		},
	},
	"subtitle": {
		FontSize:   18,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    8,
			Bottom: 8,
			Left:   0,
			Right:  0,
		},
	},
	"highlight": {
		BgColor: "yellow",
		Bold:    true,
	},
	"muted": {
		Color:   "gray",
		Opacity: 0.7,
	},
	"center": {
		TextAlign: "center",
	},
	"right": {
		TextAlign: "right",
	},
	"code": {
		FontSize:   10,
		FontFamily: "Courier",
		BgColor:    "lightgray",
		Padding: BoxDimensions{
			Top:    2,
			Bottom: 2,
			Left:   2,
			Right:  2,
		},
	},
	"text-success": {
		Color: "green",
	},
	"text-danger": {
		Color: "red",
	},
	"text-warning": {
		Color: "orange",
	},
	"text-info": {
		Color: "blue",
	},
}

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
	Margin        BoxDimensions
	Padding       BoxDimensions
	Border        Border
	Width         float64
	WidthPercent  float64
	Height        float64
	HeightPercent float64
	MaxWidth      float64
	MinWidth      float64
	MaxHeight     float64
	MinHeight     float64
	BoxSizing     string // "content-box", "border-box"

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

// Default style optimized for standard PDF documents
var DefaultStyle = Style{
	// Text & Font - standard readable size
	FontSize:       11,
	FontFamily:     "Arial",
	FontWeight:     "normal",
	FontStyle:      "normal",
	LineHeight:     1.4,  // Standard comfortable reading
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
	Color:   "#000000",
	Opacity: 1.0,
	BgColor: "#ffffff",
	Background: Background{
		Color: "#ffffff",
		Size:  "auto",
	},

	// Box Model - balanced margins
	Margin: BoxDimensions{
		Top:    0,
		Right:  0,
		Bottom: 0,
		Left:   0,
	},
	Padding: BoxDimensions{
		Top:    0,
		Right:  0,
		Bottom: 0,
		Left:   0,
	},
	Border: Border{
		Width: 0,
		Style: "solid",
		Color: "#000000",
	},
	Width:         0,
	WidthPercent:  0,
	Height:        0,
	HeightPercent: 0,
	MaxWidth:      0,
	MinWidth:      0,
	MaxHeight:     0,
	MinHeight:     0,
	BoxSizing:     "border-box",

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

// TagDefaults for standard PDF documents
var TagDefaults = map[string]Style{
	// Headings - comfortable spacing
	"h1": {
		FontSize:   24,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    12,
			Bottom: 8,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	"h2": {
		FontSize:   20,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    10,
			Bottom: 6,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	"h3": {
		FontSize:   18,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    8,
			Bottom: 5,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	"h4": {
		FontSize:   16,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    6,
			Bottom: 4,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	"h5": {
		FontSize:   14,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    5,
			Bottom: 3,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	"h6": {
		FontSize:   12,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    4,
			Bottom: 2,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	
	// Paragraph - standard spacing with justified alignment for full width
	"p": {
		FontSize:   11,
		LineHeight: 1.4,
		TextAlign:  "justify",
		Margin: BoxDimensions{
			Top:    0,
			Bottom: 8,  // Standard paragraph spacing
			Left:   0,
			Right:  0,
		},
		WidthPercent: 100,
		Padding: BoxDimensions{},
	},
	
	// Div - minimal margins
	"div": {
		FontSize:   11,
		Display:    "block",
		LineHeight: 1.4,
		Margin: BoxDimensions{
			Top:    0,
			Bottom: 4,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	
	// Span - inline
	"span": {
		FontSize:   11,
		Display:    "inline",
		LineHeight: 1.4,
		Margin:     BoxDimensions{},
		Padding:    BoxDimensions{},
	},
	
	// Links
	"a": {
		FontSize:       11,
		Display:        "inline",
		Color:          "#0066cc",
		TextDecoration: "underline",
		Margin:         BoxDimensions{},
		Padding:        BoxDimensions{},
	},
	
	// Images
	"img": {
		Display: "inline-block",
		Margin: BoxDimensions{
			Top:    4,
			Bottom: 4,
			Left:   0,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	
	// Bold/strong
	"strong": {
		FontSize:   11,
		FontWeight: "bold",
		Bold:       true,
		Display:    "inline",
		Margin:     BoxDimensions{},
		Padding:    BoxDimensions{},
	},
	"b": {
		FontSize:   11,
		FontWeight: "bold",
		Bold:       true,
		Display:    "inline",
		Margin:     BoxDimensions{},
		Padding:    BoxDimensions{},
	},
	
	// Italic/em
	"em": {
		FontSize:  11,
		FontStyle: "italic",
		Italic:    true,
		Display:   "inline",
		Margin:    BoxDimensions{},
		Padding:   BoxDimensions{},
	},
	"i": {
		FontSize:  11,
		FontStyle: "italic",
		Italic:    true,
		Display:   "inline",
		Margin:    BoxDimensions{},
		Padding:   BoxDimensions{},
	},
	
	// Code/pre
	"code": {
		FontSize:   10,
		FontFamily: "Courier",
		Display:    "inline",
		BgColor:    "#f4f4f4",
		Padding: BoxDimensions{
			Top:    1,
			Bottom: 1,
			Left:   3,
			Right:  3,
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
		BgColor: "#f4f4f4",
	},
	
	// Tables - standard spacing
	"table": {
		Display: "table",
		Margin: BoxDimensions{
			Top:    8,
			Bottom: 8,
			Left:   0,
			Right:  0,
		},
		Border: Border{
			Width: 0,
		},
	},
	"thead": {
		Display: "table-header-group",
	},
	"tbody": {
		Display: "table-row-group",
	},
	"tfoot": {
		Display: "table-footer-group",
	},
	"tr": {
		Display: "table-row",
	},
	
	// Table cells
	"td": {
		Display: "table-cell",
		FontSize: 10,
		Padding: BoxDimensions{
			Top:    4,
			Bottom: 4,
			Left:   6,
			Right:  6,
		},
		Border: Border{
			Width: 0,
		},
	},
	"th": {
		Display:    "table-cell",
		FontSize:   10,
		FontWeight: "bold",
		Bold:       true,
		Padding: BoxDimensions{
			Top:    4,
			Bottom: 4,
			Left:   6,
			Right:  6,
		},
		Border: Border{
			Width: 0,
		},
	},
	
	// Lists
	"ul": {
		Display: "block",
		Margin: BoxDimensions{
			Top:    4,
			Bottom: 4,
			Left:   24,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	"ol": {
		Display: "block",
		Margin: BoxDimensions{
			Top:    4,
			Bottom: 4,
			Left:   24,
			Right:  0,
		},
		Padding: BoxDimensions{},
	},
	"li": {
		Display: "list-item",
		Margin: BoxDimensions{
			Top:    0,
			Bottom: 2,
			Left:   0,
			Right:  0,
		},
	},
	
	// Form elements
	"input": {
		Display: "inline-block",
		Margin:  BoxDimensions{},
		Padding: BoxDimensions{
			Top:    2,
			Bottom: 2,
			Left:   4,
			Right:  4,
		},
	},
	"button": {
		Display: "inline-block",
		Margin:  BoxDimensions{},
		Padding: BoxDimensions{
			Top:    4,
			Bottom: 4,
			Left:   8,
			Right:  8,
		},
	},
}

// ClassStyles for standard documents
var ClassStyles = map[string]Style{
	"title": {
		FontSize:   24,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    8,
			Bottom: 12,
			Left:   0,
			Right:  0,
		},
	},
	"subtitle": {
		FontSize:   16,
		FontWeight: "bold",
		Bold:       true,
		Margin: BoxDimensions{
			Top:    4,
			Bottom: 8,
			Left:   0,
			Right:  0,
		},
	},
	"body-text": {
		FontSize:   11,
		LineHeight: 1.4,
		Margin: BoxDimensions{
			Top:    0,
			Bottom: 8,
			Left:   0,
			Right:  0,
		},
	},
	"compact": {
		LineHeight: 1.2,
		Margin: BoxDimensions{
			Top:    0,
			Bottom: 4,
			Left:   0,
			Right:  0,
		},
	},
	"tight": {
		Margin: BoxDimensions{
			Top:    0,
			Bottom: 2,
			Left:   0,
			Right:  0,
		},
	},
	"no-margin": {
		Margin: BoxDimensions{},
	},
	"no-padding": {
		Padding: BoxDimensions{},
	},
	"highlight": {
		BgColor: "#ffffcc",
	},
	"muted": {
		Color:   "#666666",
		Opacity: 0.8,
	},
	"center": {
		TextAlign: "center",
	},
	"right": {
		TextAlign: "right",
	},
	"justify": {
		TextAlign: "justify",
	},
	"code": {
		FontSize:   10,
		FontFamily: "Courier",
		BgColor:    "#f4f4f4",
		Padding: BoxDimensions{
			Top:    2,
			Bottom: 2,
			Left:   4,
			Right:  4,
		},
	},
	"text-success": {
		Color: "#28a745",
	},
	"text-danger": {
		Color: "#dc3545",
	},
	"text-warning": {
		Color: "#ffc107",
	},
	"text-info": {
		Color: "#17a2b8",
	},
	"border": {
		Border: Border{
			Width: 1,
			Style: "solid",
			Color: "#cccccc",
		},
		Padding: BoxDimensions{
			Top:    8,
			Bottom: 8,
			Left:   8,
			Right:  8,
		},
	},
	"border-bottom": {
		Border: Border{
			Width: 1,
			Style: "solid",
			Color: "#cccccc",
		},
	},
	"border-top": {
		Border: Border{
			Width: 1,
			Style: "solid",
			Color: "#cccccc",
		},
	},
}
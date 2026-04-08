package style

// Style represents the styling properties of an HTML element
type Style struct {
	FontSize float64
	Margin   float64
	Padding  float64
	Align    string
	Color    string
	Bold     bool
	Italic   bool
	BgColor  string
	Width    float64
	Height   float64
}

// Default style for all elements
var DefaultStyle = Style{
	FontSize: 12,
	Margin:   4,
	Padding:  2,
	Align:    "left",
	Color:    "black",
	BgColor:  "white",
	Bold:     false,
	Italic:   false,
}

// TagDefaults maps HTML tags to their default styles
var TagDefaults = map[string]Style{
	"h1": {
		FontSize: 28,
		Margin:   16,
		Padding:  8,
		Bold:     true,
		Align:    "left",
	},
	"h2": {
		FontSize: 24,
		Margin:   14,
		Padding:  6,
		Bold:     true,
		Align:    "left",
	},
	"h3": {
		FontSize: 20,
		Margin:   12,
		Padding:  4,
		Bold:     true,
		Align:    "left",
	},
	"h4": {
		FontSize: 18,
		Margin:   10,
		Padding:  4,
		Bold:     true,
		Align:    "left",
	},
	"h5": {
		FontSize: 16,
		Margin:   8,
		Padding:  3,
		Bold:     true,
		Align:    "left",
	},
	"h6": {
		FontSize: 14,
		Margin:   6,
		Padding:  2,
		Bold:     true,
		Align:    "left",
	},
	"p": {
		FontSize: 12,
		Margin:   8,
		Padding:  2,
		Align:    "left",
	},
	"div": {
		FontSize: 12,
		Margin:   6,
		Padding:  2,
		Align:    "left",
	},
	"span": {
		FontSize: 12,
		Margin:   0,
		Padding:  0,
		Align:    "left",
	},
	"strong": {
		FontSize: 12,
		Margin:   0,
		Padding:  0,
		Bold:     true,
		Align:    "left",
	},
	"em": {
		FontSize: 12,
		Margin:   0,
		Padding:  0,
		Italic:   true,
		Align:    "left",
	},
	"b": {
		FontSize: 12,
		Margin:   0,
		Padding:  0,
		Bold:     true,
		Align:    "left",
	},
	"i": {
		FontSize: 12,
		Margin:   0,
		Padding:  0,
		Italic:   true,
		Align:    "left",
	},
}

// ClassStyles maps CSS class names to their styles
var ClassStyles = map[string]Style{
	"title": {
		FontSize: 24,
		Bold:     true,
		Margin:   12,
	},
	"subtitle": {
		FontSize: 18,
		Bold:     true,
		Margin:   8,
	},
	"highlight": {
		BgColor: "yellow",
		Bold:    true,
	},
	"muted": {
		Color: "gray",
	},
	"center": {
		Align: "center",
	},
	"right": {
		Align: "right",
	},
	"code": {
		FontSize: 10,
		BgColor:  "lightgray",
		Padding:  2,
	},
}

package style

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"strings"
)

//go:embed css_theme.json
var embeddedThemeConfig []byte

type themeConfig struct {
	Fonts struct {
		DefaultFamily   string                `json:"defaultFamily"`
		MonospaceFamily string                `json:"monospaceFamily"`
		BaseSize        float64               `json:"baseSize"`
		AllProperties   *FontPropertiesConfig `json:"allProperties,omitempty"`
	} `json:"fonts"`
	Colors struct {
		DefaultText       string              `json:"defaultText"`
		DefaultBackground string              `json:"defaultBackground"`
		Named             map[string]string   `json:"named"`
		AllFormats        *AllColorFormats    `json:"allFormats,omitempty"`
		CSSFunctions      map[string][]string `json:"cssFunctions,omitempty"`
	} `json:"colors"`
}

// FontPropertiesConfig represents all CSS font properties
type FontPropertiesConfig struct {
	FontFamily            []string              `json:"fontFamily,omitempty"`
	FontSize              *FontSizeConfig       `json:"fontSize,omitempty"`
	FontWeight            *FontWeightConfig     `json:"fontWeight,omitempty"`
	FontStyle             []string              `json:"fontStyle,omitempty"`
	FontVariant           []string              `json:"fontVariant,omitempty"`
	FontStretch           []string              `json:"fontStretch,omitempty"`
	LineHeight            *LineHeightConfig     `json:"lineHeight,omitempty"`
	LetterSpacing         *SpacingConfig        `json:"letterSpacing,omitempty"`
	WordSpacing           *SpacingConfig        `json:"wordSpacing,omitempty"`
	TextAlign             []string              `json:"textAlign,omitempty"`
	TextDecoration        *TextDecorationConfig `json:"textDecoration,omitempty"`
	TextTransform         []string              `json:"textTransform,omitempty"`
	TextIndent            *TextIndentConfig     `json:"textIndent,omitempty"`
	FontSynthesis         []string              `json:"fontSynthesis,omitempty"`
	FontKerning           []string              `json:"fontKerning,omitempty"`
	FontFeatureSettings   []string              `json:"fontFeatureSettings,omitempty"`
	FontVariationSettings []string              `json:"fontVariationSettings,omitempty"`
	SystemFonts           []string              `json:"systemFonts,omitempty"`
}

type FontSizeConfig struct {
	Px       []int     `json:"px,omitempty"`
	Em       []float64 `json:"em,omitempty"`
	Rem      []float64 `json:"rem,omitempty"`
	Percent  []int     `json:"percent,omitempty"`
	Keywords []string  `json:"keywords,omitempty"`
}

type FontWeightConfig struct {
	Keywords []string `json:"keywords,omitempty"`
	Numeric  []int    `json:"numeric,omitempty"`
}

type LineHeightConfig struct {
	Unitless []float64 `json:"unitless,omitempty"`
	Px       []int     `json:"px,omitempty"`
	Percent  []int     `json:"percent,omitempty"`
	Keywords []string  `json:"keywords,omitempty"`
}

type SpacingConfig struct {
	Px     []float64 `json:"px,omitempty"`
	Em     []float64 `json:"em,omitempty"`
	Normal string    `json:"normal,omitempty"`
}

type TextDecorationConfig struct {
	Line  []string `json:"line,omitempty"`
	Style []string `json:"style,omitempty"`
	Color string   `json:"color,omitempty"`
}

type TextIndentConfig struct {
	Px      []int     `json:"px,omitempty"`
	Em      []float64 `json:"em,omitempty"`
	Percent []int     `json:"percent,omitempty"`
}

// AllColorFormats represents all CSS color formats
type AllColorFormats struct {
	Hex                 *HexConfig                 `json:"hex,omitempty"`
	RGB                 *RGBConfig                 `json:"rgb,omitempty"`
	HSL                 *HSLConfig                 `json:"hsl,omitempty"`
	HWB                 []string                   `json:"hwb,omitempty"`
	Lab                 []string                   `json:"lab,omitempty"`
	Lch                 []string                   `json:"lch,omitempty"`
	Color               *ColorFunctionConfig       `json:"color,omitempty"`
	Oklab               []string                   `json:"oklab,omitempty"`
	Oklch               []string                   `json:"oklch,omitempty"`
	Transparent         string                     `json:"transparent,omitempty"`
	CurrentColor        string                     `json:"currentColor,omitempty"`
	CSSVariables        []string                   `json:"cssVariables,omitempty"`
	NamedColorsExtended *NamedColorsExtendedConfig `json:"namedColorsExtended,omitempty"`
}

type HexConfig struct {
	Digit3 []string `json:"3digit,omitempty"`
	Digit4 []string `json:"4digit,omitempty"`
	Digit6 []string `json:"6digit,omitempty"`
	Digit8 []string `json:"8digit,omitempty"`
}

type RGBConfig struct {
	RGB     []string `json:"rgb,omitempty"`
	RGBA    []string `json:"rgba,omitempty"`
	Percent []string `json:"percent,omitempty"`
}

type HSLConfig struct {
	HSL  []string `json:"hsl,omitempty"`
	HSLA []string `json:"hsla,omitempty"`
}

type ColorFunctionConfig struct {
	Srgb       []string `json:"srgb,omitempty"`
	DisplayP3  []string `json:"display-p3,omitempty"`
	SrgbLinear []string `json:"srgb-linear,omitempty"`
	XYZ        []string `json:"xyz,omitempty"`
}

type NamedColorsExtendedConfig struct {
	Standard []string `json:"standard,omitempty"`
}

// Enhanced color structures for runtime
type RGBColor struct {
	R, G, B int
}

type RGBAColor struct {
	R, G, B int
	A       float64
}

type HSLColor struct {
	H, S, L float64
}

type HSLAColor struct {
	H, S, L float64
	A       float64
}

// NamedColors stores color-name to hex mappings loaded from css_theme.json.
var NamedColors = map[string]string{}

// ExtendedNamedColors stores all 140+ CSS named colors
var ExtendedNamedColors = map[string]string{
	"aliceblue": "#F0F8FF", "antiquewhite": "#FAEBD7", "aqua": "#00FFFF",
	"aquamarine": "#7FFFD4", "azure": "#F0FFFF", "beige": "#F5F5DC",
	"bisque": "#FFE4C4", "blanchedalmond": "#FFEBCD", "blueviolet": "#8A2BE2",
	"brown": "#A52A2A", "burlywood": "#DEB887", "cadetblue": "#5F9EA0",
	"chartreuse": "#7FFF00", "chocolate": "#D2691E", "coral": "#FF7F50",
	"cornflowerblue": "#6495ED", "cornsilk": "#FFF8DC", "crimson": "#DC143C",
	"cyan": "#00FFFF", "darkblue": "#00008B", "darkcyan": "#008B8B",
	"darkgoldenrod": "#B8860B", "darkgray": "#A9A9A9", "darkgreen": "#006400",
	"darkkhaki": "#BDB76B", "darkmagenta": "#8B008B", "darkolivegreen": "#556B2F",
	"darkorange": "#FF8C00", "darkorchid": "#9932CC", "darkred": "#8B0000",
	"darksalmon": "#E9967A", "darkseagreen": "#8FBC8F", "darkslateblue": "#483D8B",
	"darkslategray": "#2F4F4F", "darkturquoise": "#00CED1", "darkviolet": "#9400D3",
	"deeppink": "#FF1493", "deepskyblue": "#00BFFF", "dimgray": "#696969",
	"dodgerblue": "#1E90FF", "firebrick": "#B22222", "floralwhite": "#FFFAF0",
	"forestgreen": "#228B22", "fuchsia": "#FF00FF", "gainsboro": "#DCDCDC",
	"ghostwhite": "#F8F8FF", "gold": "#FFD700", "goldenrod": "#DAA520",
	"greenyellow": "#ADFF2F", "honeydew": "#F0FFF0", "hotpink": "#FF69B4",
	"indianred": "#CD5C5C", "indigo": "#4B0082", "ivory": "#FFFFF0",
	"khaki": "#F0E68C", "lavender": "#E6E6FA", "lavenderblush": "#FFF0F5",
	"lawngreen": "#7CFC00", "lemonchiffon": "#FFFACD", "lightblue": "#ADD8E6",
	"lightcoral": "#F08080", "lightcyan": "#E0FFFF", "lightgoldenrodyellow": "#FAFAD2",
	"lightgray": "#D3D3D3", "lightgreen": "#90EE90", "lightpink": "#FFB6C1",
	"lightsalmon": "#FFA07A", "lightseagreen": "#20B2AA", "lightskyblue": "#87CEFA",
	"lightslategray": "#778899", "lightsteelblue": "#B0C4DE", "lightyellow": "#FFFFE0",
	"lime": "#00FF00", "limegreen": "#32CD32", "linen": "#FAF0E6",
	"magenta": "#FF00FF", "maroon": "#800000", "mediumaquamarine": "#66CDAA",
	"mediumblue": "#0000CD", "mediumorchid": "#BA55D3", "mediumpurple": "#9370DB",
	"mediumseagreen": "#3CB371", "mediumslateblue": "#7B68EE", "mediumspringgreen": "#00FA9A",
	"mediumturquoise": "#48D1CC", "mediumvioletred": "#C71585", "midnightblue": "#191970",
	"mintcream": "#F5FFFA", "mistyrose": "#FFE4E1", "moccasin": "#FFE4B5",
	"navajowhite": "#FFDEAD", "navy": "#000080", "oldlace": "#FDF5E6",
	"olive": "#808000", "olivedrab": "#6B8E23", "orange": "#FFA500",
	"orangered": "#FF4500", "orchid": "#DA70D6", "palegoldenrod": "#EEE8AA",
	"palegreen": "#98FB98", "paleturquoise": "#AFEEEE", "palevioletred": "#DB7093",
	"papayawhip": "#FFEFD5", "peachpuff": "#FFDAB9", "peru": "#CD853F",
	"pink": "#FFC0CB", "plum": "#DDA0DD", "powderblue": "#B0E0E6",
	"purple": "#800080", "rebeccapurple": "#663399", "red": "#FF0000",
	"rosybrown": "#BC8F8F", "royalblue": "#4169E1", "saddlebrown": "#8B4513",
	"salmon": "#FA8072", "sandybrown": "#F4A460", "seagreen": "#2E8B57",
	"seashell": "#FFF5EE", "sienna": "#A0522D", "silver": "#C0C0C0",
	"skyblue": "#87CEEB", "slateblue": "#6A5ACD", "slategray": "#708090",
	"snow": "#FFFAFA", "springgreen": "#00FF7F", "steelblue": "#4682B4",
	"tan": "#D2B48C", "teal": "#008080", "thistle": "#D8BFD8",
	"tomato": "#FF6347", "turquoise": "#40E0D0", "violet": "#EE82EE",
	"wheat": "#F5DEB3", "white": "#FFFFFF", "whitesmoke": "#F5F5F5",
	"yellow": "#FFFF00", "yellowgreen": "#9ACD32",
}

// FontStyleStore holds parsed font configuration
var FontStyleStore = &FontPropertiesConfig{}

func init() {
	loadThemeConfig()
}

func loadThemeConfig() {
	cfg := themeConfig{}
	if err := json.Unmarshal(embeddedThemeConfig, &cfg); err != nil {
		return
	}

	// Load basic font config
	if cfg.Fonts.DefaultFamily != "" {
		DefaultStyle.FontFamily = cfg.Fonts.DefaultFamily
	}
	if cfg.Fonts.BaseSize > 0 {
		DefaultStyle.FontSize = cfg.Fonts.BaseSize
	}

	// Load extended font properties
	if cfg.Fonts.AllProperties != nil {
		FontStyleStore = cfg.Fonts.AllProperties
	}

	// Load color config
	if cfg.Colors.DefaultText != "" {
		DefaultStyle.Color = cfg.Colors.DefaultText
	}
	if cfg.Colors.DefaultBackground != "" {
		DefaultStyle.BgColor = cfg.Colors.DefaultBackground
		DefaultStyle.Background.Color = cfg.Colors.DefaultBackground
	}

	// Load monospace for code/pre
	if cfg.Fonts.MonospaceFamily != "" {
		if tag, ok := TagDefaults["code"]; ok {
			tag.FontFamily = cfg.Fonts.MonospaceFamily
			TagDefaults["code"] = tag
		}
		if tag, ok := TagDefaults["pre"]; ok {
			tag.FontFamily = cfg.Fonts.MonospaceFamily
			TagDefaults["pre"] = tag
		}
		if classStyle, ok := ClassStyles["code"]; ok {
			classStyle.FontFamily = cfg.Fonts.MonospaceFamily
			ClassStyles["code"] = classStyle
		}
	}

	// Load named colors from JSON
	for k, v := range cfg.Colors.Named {
		name := strings.ToLower(strings.TrimSpace(k))
		hex := strings.TrimSpace(v)
		if name == "" || !isHexColor(hex) {
			continue
		}
		NamedColors[name] = strings.ToUpper(hex)
	}

	// Load extended named colors if present
	if cfg.Colors.AllFormats != nil && cfg.Colors.AllFormats.NamedColorsExtended != nil {
		for _, name := range cfg.Colors.AllFormats.NamedColorsExtended.Standard {
			if hex, ok := ExtendedNamedColors[name]; ok {
				NamedColors[name] = hex
			}
		}
	}
}

// ColorToRGB resolves named or hex colors to RGB.
func ColorToRGB(color string) (int, int, int, bool) {
	c := strings.TrimSpace(strings.ToLower(color))
	if c == "" {
		return 0, 0, 0, false
	}

	if hex, ok := NamedColors[c]; ok {
		return hexToRGB(hex)
	}

	if strings.HasPrefix(c, "#") {
		return hexToRGB(c)
	}

	// Parse rgb/rgba format
	if strings.HasPrefix(c, "rgb") {
		return parseRGBColor(c)
	}

	// Parse hsl/hsla format
	if strings.HasPrefix(c, "hsl") {
		return parseHSLColor(c)
	}

	return 0, 0, 0, false
}

// parseRGBColor parses rgb(r,g,b) or rgba(r,g,b,a) format
func parseRGBColor(color string) (int, int, int, bool) {
	start := strings.Index(color, "(")
	end := strings.Index(color, ")")
	if start == -1 || end == -1 {
		return 0, 0, 0, false
	}

	values := strings.Split(color[start+1:end], ",")
	if len(values) < 3 {
		return 0, 0, 0, false
	}

	r, ok := parseRGBComponent(strings.TrimSpace(values[0]))
	if !ok {
		return 0, 0, 0, false
	}
	g, ok := parseRGBComponent(strings.TrimSpace(values[1]))
	if !ok {
		return 0, 0, 0, false
	}
	b, ok := parseRGBComponent(strings.TrimSpace(values[2]))
	if !ok {
		return 0, 0, 0, false
	}

	return r, g, b, true
}

func parseRGBComponent(v string) (int, bool) {
	v = strings.TrimSpace(v)
	if strings.HasSuffix(v, "%") {
		pct, err := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
		if err != nil {
			return 0, false
		}
		if pct < 0 {
			pct = 0
		}
		if pct > 100 {
			pct = 100
		}
		return int((pct / 100.0) * 255.0), true
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	if i < 0 {
		i = 0
	}
	if i > 255 {
		i = 255
	}
	return i, true
}

// parseHSLColor converts HSL to RGB
func parseHSLColor(color string) (int, int, int, bool) {
	start := strings.Index(color, "(")
	end := strings.Index(color, ")")
	if start == -1 || end == -1 {
		return 0, 0, 0, false
	}

	values := strings.Split(color[start+1:end], ",")
	if len(values) < 3 {
		return 0, 0, 0, false
	}

	h, _ := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(values[0], "deg")), 64)
	s, _ := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(values[1], "%")), 64)
	l, _ := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(values[2], "%")), 64)

	r, g, b := hslToRGB(h, s/100, l/100)
	return r, g, b, true
}

// hslToRGB converts HSL to RGB values
func hslToRGB(h, s, l float64) (int, int, int) {
	if s == 0 {
		gray := int(l * 255)
		return gray, gray, gray
	}

	var r, g, b float64

	hueToRGB := func(p, q, t float64) float64 {
		if t < 0 {
			t += 1
		}
		if t > 1 {
			t -= 1
		}
		if t < 1.0/6.0 {
			return p + (q-p)*6*t
		}
		if t < 1.0/2.0 {
			return q
		}
		if t < 2.0/3.0 {
			return p + (q-p)*(2.0/3.0-t)*6
		}
		return p
	}

	h = h / 360
	q := l + s - l*s
	p := 2*l - q
	r = hueToRGB(p, q, h+1.0/3.0)
	g = hueToRGB(p, q, h)
	b = hueToRGB(p, q, h-1.0/3.0)

	return int(r * 255), int(g * 255), int(b * 255)
}

// GetFontSizeValue returns parsed font size from config
func GetFontSizeValue(sizeType string, value interface{}) (float64, bool) {
	switch sizeType {
	case "px":
		if v, ok := value.(float64); ok {
			return v, true
		}
		if v, ok := value.(int); ok {
			return float64(v), true
		}
	case "em", "rem":
		if v, ok := value.(float64); ok {
			return v, true
		}
	case "percent":
		if v, ok := value.(float64); ok {
			return v / 100, true
		}
	}
	return 0, false
}

func isHexColor(s string) bool {
	if s == "" || s[0] != '#' {
		return false
	}
	if len(s) != 4 && len(s) != 5 && len(s) != 7 && len(s) != 9 {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func hexToRGB(hex string) (int, int, int, bool) {
	hex = strings.TrimSpace(hex)
	if !isHexColor(hex) {
		return 0, 0, 0, false
	}

	if len(hex) == 4 || len(hex) == 5 {
		h := hex[1:]
		expanded := make([]byte, 0, 7)
		expanded = append(expanded, '#')
		expanded = append(expanded, h[0], h[0], h[1], h[1], h[2], h[2])
		hex = string(expanded)
	}

	if len(hex) == 9 {
		hex = hex[:7] // ignore alpha for RGB rendering path
	}

	r, err := strconv.ParseInt(hex[1:3], 16, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	g, err := strconv.ParseInt(hex[3:5], 16, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	b, err := strconv.ParseInt(hex[5:7], 16, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	return int(r), int(g), int(b), true
}

// New helper functions for extended color support

// ColorToHex converts any color format to hex
func ColorToHex(color string) (string, bool) {
	r, g, b, ok := ColorToRGB(color)
	if !ok {
		return "", false
	}
	return rgbToHex(r, g, b), true
}

func rgbToHex(r, g, b int) string {
	return "#" + strconv.FormatInt(int64(r), 16) +
		strconv.FormatInt(int64(g), 16) +
		strconv.FormatInt(int64(b), 16)
}

// SupportsAlpha checks if color format supports alpha channel
func SupportsAlpha(color string) bool {
	return strings.HasPrefix(color, "rgba") ||
		strings.HasPrefix(color, "hsla") ||
		strings.Contains(color, "/")
}

// GetFontFamilyString returns concatenated font family options
func GetFontFamilyString(families []string) string {
	return strings.Join(families, ", ")
}

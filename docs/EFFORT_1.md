# Gotei Engine - Effort 1 Documentation

## Project Overview

Gotei is a Go-based HTML to PDF rendering engine that processes HTML documents and converts them to PDF format with comprehensive style support.

## Architecture

### Core Components

#### 1. **DOM Parser** (`internal/dom/`)
- **parser.go**: Parses HTML input into a DOM tree structure
  - Uses `golang.org/x/net/html` for standard HTML parsing
  - Converts HTML nodes to internal `Node` representation
  - Handles text nodes, element nodes, and attributes
  - Preserves HTML attributes for style resolution

- **node.go**: Defines the DOM node structure
  - `NodeType`: `ElementNode`, `TextNode` classification
  - `Node`: Stores tag, content, attributes, and child nodes
  - Attributes stored as `map[string]string` for easy lookup

#### 2. **Style Engine** (`internal/style/`)
- **resolver.go**: Resolves styles from HTML attributes and CSS
  - Parses HTML `align` attribute for text alignment
  - Parses inline `style` attribute with CSS properties
  - Applies CSS class styles from predefined class map
  - Supports text-align, font-size, color, font-weight, font-style
  - Validates alignment values: left, center, right, justify
  - Supports size parsing with CSS units: px, pt, em, rem, %

- **style.go**: Defines style structure and defaults
  - `Style` struct: FontSize, Margin, Padding, Align, Color, Bold, Italic, BgColor, Width, Height
  - `DefaultStyle`: Base styles applied to all elements
  - `TagDefaults`: Element-specific style defaults (h1-h6, p, div, span, strong, em, b, i)
  - `ClassStyles`: CSS class mappings for styling

#### 3. **Layout Engine** (`internal/layout/`)
- **tree.go**: Builds layout tree from DOM with style inheritance
  - `BuildLayoutTree(n *dom.Node, parentStyle style.Style) *Box`
  - **KEY FIX**: Implements style inheritance for alignment
    - Elements with explicit `align` attribute preserve it
    - Child text nodes inherit parent's `align` if they don't have explicit alignment
    - Text nodes under `<h1 align="center">` inherit center alignment
  - Recursively builds box tree with parent-child relationships

- **box.go**: Defines layout box structure
  - `Box` struct: X, Y, Width, Height positioning
  - Stores resolved `Style` and DOM `Node` reference
  - Maintains Parent and Children pointers for tree structure

- **layout.go**: Calculates box dimensions and positions
  - `Layout(box *Box, x, y, maxWidth float64) float64`
  - Applies margins and padding
  - Handles text node wrapping with `wrapText()`
  - Recursively lays out child boxes
  - Returns total height for vertical positioning

- **text.go**: Text wrapping and measurement
  - `wrapText()`: Breaks text into lines based on available width
  - Character-based width estimation

#### 4. **PDF Renderer** (`internal/render/`)
- **pdf.go**: Converts layout tree to PDF output
  - `RenderPDF(root *layout.Box) ([]byte, error)`
  - `renderBox()`: Recursively renders boxes to PDF
  - **Text Node Rendering**:
    - Reads `resolvedStyle.Align` for text alignment
    - Implements alignment positioning:
      - `"center"`: Centers text within box width
      - `"right"`: Right-aligns text with padding
      - `"left"`: Default left alignment
      - `"justify"`: Currently maps to left (can be extended)
  - Font styling: Bold (B), Italic (I), combinations
  - Color support: Named colors + hex color parsing
  - Debug box borders for layout visualization
  - Uses `codeberg.org/go-pdf/fpdf` library for PDF generation

#### 5. **Engine** (`internal/engine/`)
- **engine.go**: Orchestrates the rendering pipeline
  - `Render(htmlStr string) ([]byte, error)`
  - Pipeline: Parse HTML → Build Layout Tree → Calculate Layout → Render to PDF
  - Passes `style.DefaultStyle` as root parent style for inheritance

#### 6. **Server** (`cmd/server/`)
- **main.go**: HTTP server for PDF rendering
  - Endpoint: `POST /render`
  - Request format: `{"html": "<html>..."}`
  - Response: PDF file with `Content-Type: application/pdf`
  - Runs on `localhost:8080`

## Features Implemented

### Style Support

✅ **Text Alignment**
- HTML `align` attribute: left, center, right, justify
- CSS `text-align` style property
- Alignment inheritance from parent elements to child text nodes

✅ **Typography**
- Bold: `<strong>`, `<b>`, `style="font-weight: bold"`, `font-weight: 700/800/900`
- Italic: `<em>`, `<i>`, `style="font-style: italic"`
- Font sizes: Default tag styles + custom `style="font-size: Xpx"`
- Font size units: px, pt, em, rem, %

✅ **Colors**
- Named colors: red, blue, green, yellow, cyan, magenta, gray, lightgray, darkgray, orange, purple, brown, pink, black, white
- Hex colors: `#RRGGBB` format
- `style="color: colorname"` or `style="color: #hex"`

✅ **Spacing**
- Margins: Default per element, customizable
- Padding: Default per element, customizable
- `style="margin: Xpx"` and `style="padding: Xpx"`

✅ **Box Model**
- Width and height support
- Margin and padding application
- Box positioning and sizing

### Cascade & Inheritance

- **Cascade Order**: Defaults → Tag Defaults → HTML Attributes → Inline Styles → CSS Classes
- **Key Fix**: Child nodes inherit parent `align` unless explicitly overridden
- Class-based styling through CSS class mapping

## Test Coverage

### Examples Program (`examples/main.go`)

10 comprehensive test scenarios:

1. **alignment-combinations**: Left, center, right with bold/italic
2. **typography-styles**: Headings, bold, italic combinations
3. **color-and-alignment**: Colored text with different alignments
4. **mixed-sizes-and-alignment**: Various font sizes with alignment
5. **complex-layout**: Multi-section document with mixed alignments
6. **nested-styles**: Style inheritance and nesting
7. **all-alignments**: Comprehensive alignment testing
8. **bold-italic-colors**: Typography and color combinations
9. **headings-with-styles**: All heading levels with styles
10. **comprehensive-styles**: Full feature combination test

All examples generate PDF files in `examples/output/` for visual verification.

## Known Limitations

- Text wrapping uses character-based estimation (not precise font measurement)
- Justify alignment currently maps to left (not implemented)
- No support for background colors in text rendering
- Limited CSS property support (focused on core styling)
- No support for padding/margin on specific sides (all sides treated equally)
- Single-page PDF output (no multi-page support)

## Usage

### HTTP Server

```bash
go run ./cmd/server/main.go
```

Then POST HTML:
```bash
curl -X POST localhost:8080/render \
  -H "Content-Type: application/json" \
  -d '{"html":"<h1 align=\"center\">Title</h1><p>Content</p>"}' \
  --output output.pdf
```

### Direct API

```go
import "github.com/Viswesh934/gotei/internal/engine"

pdf, err := engine.Render("<h1>Hello</h1><p>World</p>")
if err != nil {
    log.Fatal(err)
}

ioutil.WriteFile("output.pdf", pdf, 0644)
```

### Running Examples

```bash
cd examples
go run main.go
```

## File Structure

```
Gotei/
├── cmd/
│   └── server/
│       └── main.go          # HTTP server
├── internal/
│   ├── dom/
│   │   ├── node.go         # DOM structure
│   │   └── parser.go       # HTML parser
│   ├── engine/
│   │   └── engine.go       # Render pipeline
│   ├── layout/
│   │   ├── box.go          # Box structure
│   │   ├── layout.go       # Layout calculation
│   │   ├── text.go         # Text wrapping
│   │   └── tree.go         # Layout tree builder
│   ├── render/
│   │   └── pdf.go          # PDF rendering
│   └── style/
│       ├── resolver.go     # Style resolution
│       └── style.go        # Style definitions
├── examples/
│   ├── main.go            # Test scenarios
│   ├── README.md          # Examples documentation
│   └── output/            # Generated PDFs (10 files)
├── docs/
│   └── EFFORT_1.md        # This document
├── go.mod
├── go.sum
├── Dockerfile
├── LICENSE
└── README.md
```

## Dependencies

- `golang.org/x/net/html`: Standard HTML parsing
- `codeberg.org/go-pdf/fpdf`: PDF generation library

## Recent Fixes

### Alignment Inheritance Fix

**Problem**: Text inside aligned elements (e.g., `<p align="center">text</p>`) was not inheriting the parent's alignment.

**Root Cause**: Each DOM node had its style resolved independently without inheritance. Text nodes would get default "left" alignment even when their parent had explicit center/right alignment.

**Solution**: Modified `BuildLayoutTree` to pass parent style as parameter:
```go
func BuildLayoutTree(n *dom.Node, parentStyle style.Style) *Box {
    currentStyle := style.Resolve(n)
    
    // Only inherit parent align if node doesn't have explicit align
    if n.Type == dom.ElementNode {
        if _, hasAlignAttr := n.Attr["align"]; !hasAlignAttr {
            currentStyle.Align = parentStyle.Align
        }
    } else if n.Type == dom.TextNode {
        // Text nodes always inherit from parent
        currentStyle.Align = parentStyle.Align
    }
    
    box := &Box{...}
    
    // Pass current style as parent style to children
    for _, child := range n.Children {
        childBox := BuildLayoutTree(child, box.Style)
        ...
    }
    return box
}
```

This ensures:
- Elements with explicit `align` preserve their alignment
- Elements without explicit `align` inherit from parent
- Text nodes always inherit parent's alignment

## Next Steps (Future Efforts)

- Multi-page PDF support
- Implement justify alignment
- Background color rendering
- More CSS properties (text-decoration, line-height, etc.)
- Better text measurement for accurate wrapping
- CSS selector support (beyond classes)
- Image support
- Table support
- Form elements rendering

# Gotei CSS Engine - EFFORT_1 Documentation

## Overview

This document summarizes the first engineering effort on the Gotei HTML-to-PDF rendering engine, focusing on establishing a W3C-compliant CSS architecture with comprehensive style support.

**Timeline**: Single session
**Status**: Partial implementation - Core foundation working, advanced features need debugging

---

## What Got Done ✅

### Phase 1: Foundation (CSS Properties & Box Model)
**Status**: ✅ COMPLETE

- **50+ CSS properties** defined in `internal/style/style.go`
- **BoxDimensions struct** for CSS box model (margin, padding, border)
- **W3C InheritableFields map** defining which properties cascade
- **Tag defaults** for semantic elements (h1-h6, p, div, span, strong, em, code, pre, etc.)
- **Class styles** for common patterns (.title, .center, .text-danger, etc.)

**Working**:
- Default styles applied to elements
- HTML alignment attribute (`align="center"`) works
- Inline styles (`style="color: red"`) work
- Tag-specific defaults apply (h1 is bold and large, etc.)

### Phase 2: CSS Cascade & Selectors
**Status**: ✅ FOUNDATION COMPLETE (but not wired in)

- **Selector engine** (`internal/css/selectors.go`)
  - Parses CSS selectors: element, class, ID, attribute selectors
  - Calculates specificity (a,b,c) per W3C spec
  - Matches selectors against DOM nodes

- **Cascade engine** (`internal/css/cascade.go`)
  - `ComputeStyle()` - Calculates final style with cascade rules
  - `applyProperty()` - Parses 50+ CSS properties
  - Specificity sorting (lowest to highest wins)
  - Inheritance application

- **CSS parser** (`internal/css/parser.go`) - NEW
  - Extracts CSS from `<style>` tags
  - Parses CSS rules and declarations
  - Builds StyleSheet object

- **Cascade tree applier** (`internal/css/tree.go`) - NEW
  - Applies cascade to entire DOM tree
  - Returns style map for all nodes

**Status Issue**: CSS cascade is implemented but needs debugging - `<style>` tag styles not fully appearing in rendered output

### Phase 3: Flexbox Layout Engine
**Status**: ✅ IMPLEMENTED (not fully verified with CSS cascade)

- **Flexbox algorithm** (`internal/layout/flexbox.go`)
  - Supports `display: flex` property routing
  - `flex-direction`: row/column layout
  - `justify-content`: center, flex-end, space-between, space-around
  - `align-items`: flex-start, center, flex-end alignment
  - `gap`: spacing between items
  - `flex-grow`/`flex-shrink`: item sizing

- **Font metrics** (`internal/layout/metrics.go`)
  - Estimates font metrics for Arial, Georgia, Courier
  - Calculates proper line heights
  - Measures text dimensions

**Status Issue**: Flexbox is implemented but CSS `display: flex` from `<style>` tags not routing correctly

### Phase 4: Visual Effects (Borders, Backgrounds, Shadows)
**Status**: ✅ IMPLEMENTED (not fully verified)

- **Border rendering** (`internal/render/borders.go`)
  - Solid/dashed/dotted border styles
  - `DrawBorder()` method with customizable width/color

- **Background rendering**
  - `DrawBackground()` fills boxes with CSS colors
  - `background-color` property parsing

- **Box shadow**
  - `DrawBoxShadow()` with offset, blur, opacity
  - `box-shadow` CSS property support

- **Integration into render pipeline**
  - Rendering order: background → shadow → text → borders → children

**Status Issue**: Background colors and borders not appearing in PDFs from `<style>` tag CSS

### Phase 5: CSS Cascade Integration (Just Added)
**Status**: ✅ WIRED (needs testing & debugging)

- **Engine pipeline updated** (`internal/engine/engine.go`)
  ```
  Parse HTML → Extract CSS → Apply Cascade → Build Layout Tree → Layout → Render
  ```

- **BuildLayoutTreeWithCascade()** - Uses pre-computed styles from cascade

- **CSS property parsing expanded** - 50+ properties now parse including:
  - All flexbox properties
  - Individual margin/padding sides
  - Border properties
  - Shadow properties
  - Display property

---

## Current Architecture

```
HTML Input
    ↓
┌─────────────────────────────────┐
│  Phase 1: Parse HTML            │  ✅ WORKING
│  dom.ParseHTML() → DOM Tree     │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│  Phase 2: Extract & Cascade CSS │  ⚠️  NEEDS DEBUG
│  css.ParseStyleSheet()          │  css.ApplyCascadeToTree()
│  → StyleSheet + Style Map       │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│  Phase 3: Build Layout Tree     │  ⚠️  NEEDS DEBUG
│  layout.BuildLayoutTreeWithCascade()
│  Uses cascade-computed styles   │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│  Phase 4: Calculate Layout      │  ✅ WORKING (block + flex)
│  layout.Layout()                │
│  → Dimensions & positions       │
└─────────────────────────────────┘
    ↓
┌─────────────────────────────────┐
│  Phase 5: Render to PDF         │  ✅ WORKING (core)
│  render.RenderPDF()             │  ⚠️  (colors/borders partial)
│  → Binary PDF bytes             │
└─────────────────────────────────┘
```

---

## What's Working ✅

### Examples That Generate Successfully

**Phase 1-2: Alignment Examples**
```bash
cd /workspaces/Gotei/examples/alignment && go run main.go
```
✅ All 10 examples work:
- alignment-combinations
- typography-styles (bold/italic)
- color-and-alignment (inline styles only)
- mixed-sizes-and-alignment
- complex-layout
- nested-styles
- all-alignments
- bold-italic-colors (inline styles only)
- headings-with-styles
- comprehensive-styles

**Why these work**: Use `align="center"` HTML attribute + inline `style="..."` attributes. These bypass the CSS cascade - they go directly through `style.Resolve()` which handles inline styles.

### Core Infrastructure
- ✅ HTML parsing (golang.org/x/net/html)
- ✅ DOM tree building with proper node hierarchy
- ✅ Inline style parsing and application
- ✅ Font/size/weight/style properties working
- ✅ Block layout algorithm (stacking children vertically)
- ✅ PDF generation with fpdf

---

## What's NOT Working ⚠️

### Features Implemented But Not Rendering

1. **CSS `<style>` colors and backgrounds**
   - Phase 4 examples generate but colors not visible
   - Background colors not filling boxes
   - Cause: CSS cascade not flowing into PDF rendering

2. **Flexbox from `<style>` tags**
   - Phase 3 examples generate but layout appears as block
   - `display: flex` not routing to flexbox algorithm
   - Cause: CSS cascade not setting display property correctly

3. **Borders from CSS**
   - Borders implemented but not appearing
   - BorderRendering struct created but not getting style data

### Root Causes to Debug

**Issue #1**: Cascade not populating style map
- `css.ParseStyleSheet()` may not extracting CSS correctly
- May not finding `<style>` tags in DOM tree

**Issue #2**: Style map not being used
- `BuildLayoutTreeWithCascade()` may not applying computed styles
- Fallback to `style.Resolve()` happening instead

**Issue #3**: Text color inheritance broken
- `inheritAllProperties()` may have reflection issues
- Properties may not updating correctly

---

## Test Commands

```bash
# Build everything
cd /workspaces/Gotei && go build ./...

# Test alignment (works)
cd /workspaces/Gotei/examples/alignment && go run main.go

# Test flexbox (partial - need debug)
cd /workspaces/Gotei/examples/flex && go run flex-test.go

# Test borders/backgrounds/colors (partial - need debug)
cd /workspaces/Gotei/examples/bor-bg-shad && go run phase4-test.go

# Start server and test with curl
go run ./cmd/server/main.go
# In another terminal:
curl -X POST localhost:8080/render \
  -H "Content-Type: application/json" \
  -d '{"html":"<h1 style=\"color:red\">Test</h1>"}' \
  --output test.pdf
```

---

## Key Files Structure

### Core Engine
- `internal/engine/engine.go` - Main orchestration (6 phases)
- `internal/dom/parser.go` - HTML parsing
- `internal/style/style.go` - 50+ CSS properties definitions

### CSS System
- `internal/css/selectors.go` - Selector matching & specificity
- `internal/css/cascade.go` - Cascade engine & property parsing
- `internal/css/parser.go` - CSS extraction from `<style>` tags
- `internal/css/tree.go` - Apply cascade to DOM tree

### Layout System
- `internal/layout/tree.go` - Build layout tree
- `internal/layout/layout.go` - Main layout orchestration
- `internal/layout/flexbox.go` - Flexbox algorithm
- `internal/layout/metrics.go` - Font metrics & text measurement
- `internal/layout/inherit.go` - Property inheritance

### Rendering
- `internal/render/pdf.go` - PDF rendering orchestration
- `internal/render/borders.go` - Borders, shadows, backgrounds
- `internal/render/pages.go` - Multi-page support (prepared)

### Examples
- `examples/alignment/main.go` - ✅ Working alignment tests
- `examples/flex/flex-test.go` - ⚠️ Flexbox tests (need debugging)
- `examples/bor-bg-shad/phase4-test.go` - ⚠️ Visual effects tests (need debugging)

---

## Debugging Checklist

When debugging why CSS styles don't show in PDFs:

### Step 1: Parse CSS
```go
root, _ := dom.ParseHTML(html)
sheet := css.ParseStyleSheet(root)
fmt.Printf("Parsed %d CSS rules\n", len(sheet.Rules))
for i, rule := range sheet.Rules {
    fmt.Printf("Rule %d: selector=%s, props=%v\n", i, rule.Selector, len(rule.Properties))
}
```

### Step 2: Apply Cascade
```go
styleMap := css.ApplyCascadeToTree(root, sheet)
// Check if styles are computed
if style, ok := styleMap[someNode]; ok {
    fmt.Printf("Computed style: Color=%s, Display=%s, FontSize=%.1f\n",
        style.Color, style.Display, style.FontSize)
}
```

### Step 3: Build Layout Tree
```go
layoutTree := layout.BuildLayoutTreeWithCascade(root, styleMap, style.Style{})
// Inspect if styles flowed through
func inspectLayout(box *layout.Box, depth int) {
    indent := strings.Repeat("  ", depth)
    fmt.Printf("%s<%s> Color:%s, Display:%s\n",
        indent, box.Node.Tag, box.Style.Color, box.Style.Display)
    for _, child := range box.Children {
        inspectLayout(child, depth+1)
    }
}
inspectLayout(layoutTree, 0)
```

### Step 4: Render PDF
```go
pdfBytes, _ := render.RenderPDF(layoutTree)
fmt.Printf("PDF generated: %d bytes\n", len(pdfBytes))
// Check if non-zero
```

---

## Known Limitations

1. **No media queries** - CSS media queries not implemented
2. **No pseudo-classes** - :hover, :nth-child, etc. not supported
3. **No tables** - `<table>`, `<tr>`, `<td>` not implemented
4. **No forms** - `<input>`, `<button>`, `<textarea>` not implemented
5. **Limited units** - Only px/pt/em/rem/% parsed, conversion may be approximate
6. **No grid layout** - CSS Grid infrastructure exists but not implemented
7. **Limited font support** - Only Arial, Georgia, Courier mapped
8. **No animations/transitions** - Keyframes not supported (PDF limitation)

---

## What's Next (Phase 5+ Pending)

### Immediate (To Fix Current Issues)
1. Debug CSS cascade not populating styles
2. Debug style map not flowing into layout tree
3. Verify text color inheritance in text nodes
4. Test flexbox routing with display:flex from CSS

### Medium Priority
1. Multi-page PDF support (infrastructure exists)
2. Margin collapsing (helpers exist)
3. Overflow handling

### Future Phases
- Media queries for responsive layouts
- Pseudo-classes (:hover, :nth-child)
- Table rendering engine
- Form element rendering
- Viewport units (vw, vh)

---

## Testing Strategy

All examples should compile with `go build ./...` and run without errors.

**Green Tests** (fully working):
```bash
cd /workspaces/Gotei/examples/alignment && go run main.go
# Output: ✓ All 10 PDFs generated
```

**Yellow Tests** (implemented but not rendering correctly):
```bash
cd /workspaces/Gotei/examples/flex && go run flex-test.go
cd /workspaces/Gotei/examples/bor-bg-shad && go run phase4-test.go
# Output: ✓ PDFs generated but colors/layout may not display
```

---

## Developer Notes

### Key Insights
1. **HTML attributes vs CSS**: Inline `style="..."` and `align="..."` work because they bypass cascade
2. **Reflection-based inheritance**: Using reflection to dynamically handle 50+ properties is flexible but needs careful debugging
3. **Pipeline architecture**: Each phase depends on previous phase output; errors propagate downstream
4. **Text node semantics**: Text nodes must inherit ALL properties; element nodes inherit only defaults

### Technical Debt
- CSS parser is simple, may not handle all valid CSS syntax
- Error handling could be more verbose (many silent failures)
- No validation that nodes exist before accessing in cascade
- BoxDimensions struct used but legacy float64 fields still exist in some places

### Code Quality
- Well-organized package structure
- Clear separation of concerns (parse → cascade → layout → render)
- Comprehensive style definitions with W3C defaults
- Proper use of struct composition (Border, Shadow, Background, BoxDimensions)

---

## Conclusion

**Status Summary**:
- ✅ **Foundation solid**: 50+ CSS properties, proper architecture, all infrastructure in place
- ✅ **Core rendering works**: PDFs generate correctly for supported features
- ⚠️ **Advanced features need debugging**: CSS cascade, flexbox, colors from `<style>` not fully functional
- 🟢 **Easy to debug**: Clear architecture makes it straightforward to trace issues

**Effort Grade**: B+ (Strong foundation, execution incomplete)

**Next Session**: Debug CSS cascade integration and verify color/flexbox rendering.

---

## Reference: CSS Properties Supported

### Text & Font
- font-size, font-family, font-weight, font-style
- line-height, letter-spacing, word-spacing
- text-align, text-decoration, text-transform

### Color & Background
- color, background-color, opacity

### Box Model
- margin (shorthand + sides), padding (shorthand + sides)
- width, height, max-width, min-width
- border (width/color/style), border-radius

### Display & Layout
- display (block, flex, inline, inline-block)
- flex-direction, flex-wrap, justify-content, align-items
- gap, flex-grow, flex-shrink, flex-basis

### Visual Effects
- box-shadow, text-shadow
- overflow, visibility, position

### Advanced
- z-index, float, clear
- grid-* properties (prepared but not implemented)

---

**Document Created**: April 8, 2026
**Status**: Snapshot after Phase 5 cascade integration
**Next Review**: After debugging session

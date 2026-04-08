# Gotei Style Examples

This folder contains test programs demonstrating Gotei's style rendering capabilities.

## Running the Examples

```bash
cd examples
go run main.go
```

This generates PDFs in the `output/` folder testing:

- **alignment-combinations** - left, center, right alignment with bold and italic
- **typography-styles** - bold, italic, mixed combinations with headings
- **color-and-alignment** - colored text with different alignments
- **mixed-sizes-and-alignment** - various font sizes with alignment
- **complex-layout** - multi-section document with different alignments
- **nested-styles** - style inheritance and nesting
- **all-alignments** - comprehensive alignment test with both attributes and CSS
- **bold-italic-colors** - typography and color combinations
- **headings-with-styles** - all heading levels with various styles
- **comprehensive-styles** - comprehensive test combining all features

## Generated PDFs

Each test case generates a separate PDF file for visual inspection.

## Features Tested

### Text Alignment
- `align="left"` attribute
- `align="center"` attribute  
- `align="right"` attribute
- `style="text-align: center"` CSS property

### Typography
- `<strong>` and `<b>` tags for bold
- `<em>` and `<i>` tags for italic
- Combined bold and italic styles

### Colors
- HTML color names (red, blue, green, purple, etc.)
- `style="color: colorname"` syntax

### Font Sizes
- Default sizes for various heading levels
- Custom sizes with `style="font-size: Xpx"`

### Combinations
- Aligned text with bold
- Colored text with alignment
- Mixed typography with alignment and colors
- Nested styles and inheritance

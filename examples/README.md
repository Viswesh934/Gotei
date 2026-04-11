# Gotei Style Examples

This folder contains test programs demonstrating Gotei's style rendering capabilities.

## Running the Examples

```bash
go run ./cmd/examples-runner

# Or run individual suites:
go run ./examples/alignment/main.go
go run ./examples/flex/flex-test.go
go run ./examples/tables/table-test.go
go run ./examples/phase2/phase2-test.go
go run ./examples/fullscale/form2-test.go
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

Table-specific runs generate PDFs in `examples/tables/output/` including:

- basic bordered table
- colspan table layout
- mixed table styling with heading/body cell variations

Phase 2 runs generate PDFs in `examples/phase2/output/` including:

- hyperlink click target verification
- local image rendering verification
- pagination stress verification (multi-page flow)

Full-scale runs generate production-style form outputs in `examples/fullscale/output/` including:

- multi-section Form 2 style document
- dynamic nominee/family table row generation
- embedded logo/signature/seal images
- long-form layout that exercises page flow and table rendering

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

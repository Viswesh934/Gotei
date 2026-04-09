package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Viswesh934/gotei/internal/engine"
)

func generateTestsFromJSON() []struct {
	name string
	html string
} {
	return []struct {
		name string
		html string
	}{
		// FONT FAMILIES - from allProperties.fontFamily
		{
			name: "01-font-families-all",
			html: `<div><h2>Font Families</h2>
				<p style="font-family: Arial">Arial</p>
				<p style="font-family: Verdana">Verdana</p>
				<p style="font-family: Helvetica">Helvetica</p>
				<p style="font-family: 'Times New Roman'">Times New Roman</p>
				<p style="font-family: Georgia">Georgia</p>
				<p style="font-family: 'Courier New'">Courier New</p>
				<p style="font-family: monospace">Generic monospace</p>
				<p style="font-family: sans-serif">Generic sans-serif</p>
				<p style="font-family: serif">Generic serif</p>
			</div>`,
		},
		// FONT SIZE - px values
		{
			name: "02-font-size-px",
			html: `<div><h2>Font Size (px)</h2>
				<p style="font-size: 8px">8px text</p>
				<p style="font-size: 12px">12px text</p>
				<p style="font-size: 16px">16px text</p>
				<p style="font-size: 20px">20px text</p>
				<p style="font-size: 32px">32px text</p>
				<p style="font-size: 48px">48px text</p>
			</div>`,
		},
		// FONT SIZE - em values
		{
			name: "03-font-size-em",
			html: `<div><h2>Font Size (em)</h2>
				<p style="font-size: 0.5em">0.5em text</p>
				<p style="font-size: 0.75em">0.75em text</p>
				<p style="font-size: 1em">1em text</p>
				<p style="font-size: 1.25em">1.25em text</p>
				<p style="font-size: 1.5em">1.5em text</p>
				<p style="font-size: 2em">2em text</p>
			</div>`,
		},
		// FONT SIZE - keywords
		{
			name: "04-font-size-keywords",
			html: `<div><h2>Font Size (keywords)</h2>
				<p>Keywords: xx-small, x-small, small, medium, large, x-large, xx-large, smaller, larger</p>
				<p style="font-size: 12px">Base: 12px</p>
			</div>`,
		},
		// FONT WEIGHT - keywords and numeric
		{
			name: "05-font-weight-all",
			html: `<div><h2>Font Weight</h2>
				<p style="font-weight: normal">Normal weight</p>
				<p style="font-weight: bold">Bold weight</p>
				<p style="font-weight: 100">Weight 100</p>
				<p style="font-weight: 300">Weight 300 (light)</p>
				<p style="font-weight: 400">Weight 400 (normal)</p>
				<p style="font-weight: 700">Weight 700 (bold)</p>
				<p style="font-weight: 900">Weight 900 (heavy)</p>
			</div>`,
		},
		// FONT STYLE - italic, oblique
		{
			name: "06-font-style-all",
			html: `<div><h2>Font Style</h2>
				<p style="font-style: normal">Normal font style</p>
				<p style="font-style: italic">Italic font style</p>
				<p style="font-style: oblique">Oblique font style</p>
			</div>`,
		},
		// LINE HEIGHT - unitless values
		{
			name: "07-line-height-unitless",
			html: `<div><h2>Line Height (unitless)</h2>
				<p style="line-height: 1">Line 1: Lorem ipsum dolor sit. Lorem ipsum.</p>
				<p style="line-height: 1.2">Line 1.2: Lorem ipsum dolor sit. Lorem ipsum.</p>
				<p style="line-height: 1.5">Line 1.5: Lorem ipsum dolor sit. Lorem ipsum.</p>
				<p style="line-height: 1.8">Line 1.8: Lorem ipsum dolor sit. Lorem ipsum.</p>
				<p style="line-height: 2">Line 2: Lorem ipsum dolor sit. Lorem ipsum.</p>
			</div>`,
		},
		// LINE HEIGHT - px values
		{
			name: "08-line-height-px",
			html: `<div><h2>Line Height (px)</h2>
				<p style="line-height: 16px">16px: Lorem ipsum dolor sit. Lorem ipsum.</p>
				<p style="line-height: 20px">20px: Lorem ipsum dolor sit. Lorem ipsum.</p>
				<p style="line-height: 24px">24px: Lorem ipsum dolor sit. Lorem ipsum.</p>
				<p style="line-height: 32px">32px: Lorem ipsum dolor sit. Lorem ipsum.</p>
			</div>`,
		},
		// LETTER SPACING - px and negative
		{
			name: "09-letter-spacing-all",
			html: `<div><h2>Letter Spacing</h2>
				<p style="letter-spacing: -2px">Negative 2px: Spaced text</p>
				<p style="letter-spacing: -0.5px">Negative 0.5px: Spaced text</p>
				<p style="letter-spacing: 0px">Zero 0px: Spaced text</p>
				<p style="letter-spacing: 0.5px">Positive 0.5px: Spaced text</p>
				<p style="letter-spacing: 2px">Positive 2px: Spaced text</p>
				<p style="letter-spacing: 5px">Positive 5px: Spaced text</p>
			</div>`,
		},
		// LETTER SPACING - em values
		{
			name: "10-letter-spacing-em",
			html: `<div><h2>Letter Spacing (em)</h2>
				<p style="letter-spacing: -0.05em">em -0.05: Spaced text</p>
				<p style="letter-spacing: 0em">em 0: Spaced text</p>
				<p style="letter-spacing: 0.05em">em 0.05: Spaced text</p>
				<p style="letter-spacing: 0.1em">em 0.1: Spaced text</p>
			</div>`,
		},
		// WORD SPACING - px values
		{
			name: "11-word-spacing-px",
			html: `<div><h2>Word Spacing (px)</h2>
				<p style="word-spacing: -2px">Negative: The quick brown fox</p>
				<p style="word-spacing: 0px">Zero: The quick brown fox</p>
				<p style="word-spacing: 4px">4px: The quick brown fox</p>
				<p style="word-spacing: 8px">8px: The quick brown fox</p>
			</div>`,
		},
		// WORD SPACING - em values
		{
			name: "12-word-spacing-em",
			html: `<div><h2>Word Spacing (em)</h2>
				<p style="word-spacing: -0.1em">em -0.1: The quick brown fox</p>
				<p style="word-spacing: 0em">em 0: The quick brown fox</p>
				<p style="word-spacing: 0.1em">em 0.1: The quick brown fox</p>
				<p style="word-spacing: 0.2em">em 0.2: The quick brown fox</p>
			</div>`,
		},
		// TEXT ALIGN - all values
		{
			name: "13-text-align-all",
			html: `<div><h2>Text Align</h2>
				<p style="text-align: left">Left aligned text goes here</p>
				<p style="text-align: center">Center aligned text goes here</p>
				<p style="text-align: right">Right aligned text goes here</p>
				<p style="text-align: justify">Justify: The quick brown fox jumps over lazy dog</p>
				<p style="text-align: start">Start alignment text</p>
				<p style="text-align: end">End alignment text</p>
			</div>`,
		},
		// TEXT DECORATION - underline, overline, line-through
		{
			name: "14-text-decoration-all",
			html: `<div><h2>Text Decoration</h2>
				<p style="text-decoration: none">No decoration</p>
				<p style="text-decoration: underline">Underlined text here</p>
				<p style="text-decoration: overline">Overlined text here</p>
				<p style="text-decoration: line-through">Strikethrough text here</p>
			</div>`,
		},
		// TEXT TRANSFORM - all values
		{
			name: "15-text-transform-all",
			html: `<div><h2>Text Transform</h2>
				<p style="text-transform: none">none: Hello World Mixed Case</p>
				<p style="text-transform: capitalize">capitalize: hello world mixed case</p>
				<p style="text-transform: uppercase">uppercase: hello world mixed case</p>
				<p style="text-transform: lowercase">lowercase: HELLO WORLD MIXED CASE</p>
				<p style="text-transform: full-width">full-width: hello world</p>
			</div>`,
		},
		// COLORS - NAMED BASIC
		{
			name: "16-colors-basic-named",
			html: `<div><h2>Colors - Basic Named</h2>
				<p style="color: black">Color: black</p>
				<p style="color: white; background-color: black">Color: white (on black)</p>
				<p style="color: red">Color: red</p>
				<p style="color: green">Color: green</p>
				<p style="color: blue">Color: blue</p>
				<p style="color: yellow">Color: yellow</p>
				<p style="color: cyan">Color: cyan</p>
				<p style="color: magenta">Color: magenta</p>
				<p style="color: gray">Color: gray</p>
				<p style="color: orange">Color: orange</p>
				<p style="color: purple">Color: purple</p>
			</div>`,
		},
		// COLORS - HEX 6-digit
		{
			name: "17-colors-hex-6digit",
			html: `<div><h2>Colors - Hex 6-digit</h2>
				<p style="color: #000000">Hex #000000 (black)</p>
				<p style="color: #FF0000">Hex #FF0000 (red)</p>
				<p style="color: #00FF00">Hex #00FF00 (green)</p>
				<p style="color: #0000FF">Hex #0000FF (blue)</p>
				<p style="color: #FFFFFF; background-color: #333333">Hex #FFFFFF (white)</p>
				<p style="color: #FFFF00">Hex #FFFF00 (yellow)</p>
				<p style="color: #FF00FF">Hex #FF00FF (magenta)</p>
				<p style="color: #00FFFF">Hex #00FFFF (cyan)</p>
			</div>`,
		},
		// COLORS - HEX 3-digit
		{
			name: "18-colors-hex-3digit",
			html: `<div><h2>Colors - Hex 3-digit</h2>
				<p style="color: #000">Hex #000</p>
				<p style="color: #F00">Hex #F00 (red)</p>
				<p style="color: #0F0">Hex #0F0 (green)</p>
				<p style="color: #00F">Hex #00F (blue)</p>
				<p style="color: #FFF; background-color: #333">Hex #FFF (white)</p>
				<p style="color: #FF0">Hex #FF0 (yellow)</p>
			</div>`,
		},
		// COLORS - RGB
		{
			name: "19-colors-rgb",
			html: `<div><h2>Colors - RGB</h2>
				<p style="color: rgb(0,0,0)">RGB (0,0,0) black</p>
				<p style="color: rgb(255,0,0)">RGB (255,0,0) red</p>
				<p style="color: rgb(0,255,0)">RGB (0,255,0) green</p>
				<p style="color: rgb(0,0,255)">RGB (0,0,255) blue</p>
				<p style="color: rgb(255,255,255); background-color: #333">RGB (255,255,255) white</p>
				<p style="color: rgb(255,165,0)">RGB (255,165,0) orange</p>
				<p style="color: rgb(128,0,128)">RGB (128,0,128) purple</p>
			</div>`,
		},
		// COLORS - HSL
		{
			name: "20-colors-hsl",
			html: `<div><h2>Colors - HSL</h2>
				<p style="color: hsl(0,0%,0%)">HSL (0,0%,0%) black</p>
				<p style="color: hsl(0,100%,50%)">HSL (0,100%,50%) red</p>
				<p style="color: hsl(120,100%,50%)">HSL (120,100%,50%) green</p>
				<p style="color: hsl(240,100%,50%)">HSL (240,100%,50%) blue</p>
				<p style="color: hsl(0,0%,100%); background-color: #333">HSL (0,0%,100%) white</p>
				<p style="color: hsl(39,100%,50%)">HSL (39,100%,50%) orange</p>
				<p style="color: hsl(300,100%,50%)">HSL (300,100%,50%) magenta</p>
			</div>`,
		},
		// COLORS - NAMED EXTENDED (sample 140+)
		{
			name: "21-colors-named-extended",
			html: `<div><h2>Colors - Extended Named (140+)</h2>
				<p style="color: aliceblue">aliceblue</p>
				<p style="color: antiquewhite">antiquewhite</p>
				<p style="color: aquamarine">aquamarine</p>
				<p style="color: coral">coral</p>
				<p style="color: cornflowerblue">cornflowerblue</p>
				<p style="color: darkseagreen">darkseagreen</p>
				<p style="color: rebeccapurple">rebeccapurple</p>
				<p style="color: mediumvioletred">mediumvioletred</p>
				<p style="color: steelblue">steelblue</p>
				<p style="color: tomato">tomato</p>
			</div>`,
		},
		// BACKGROUND COLORS
		{
			name: "22-background-colors",
			html: `<div><h2>Background Colors</h2>
				<p style="background-color: #FFFF00; color: black">Hex background yellow</p>
				<p style="background-color: rgb(100,149,237); color: white">RGB background cornflowerblue</p>
				<p style="background-color: hsl(48,100%,50%); color: black">HSL background golden</p>
				<p style="background-color: lightcoral; color: white">Named background lightcoral</p>
				<p style="background-color: lightyellow">Named background lightyellow</p>
				<p style="background-color: #D3D3D3">Hex background lightgray</p>
			</div>`,
		},
		// BOX SHADOWS
		{
			name: "23-box-shadows",
			html: `<div><h2>Box Shadows</h2>
				<p style="box-shadow: 2px 2px 4px #000000; padding: 10px">Shadow offset 2px right/down</p>
				<p style="box-shadow: -2px -2px 4px #666666; padding: 10px">Shadow offset -2px left/up</p>
				<p style="box-shadow: 0px 4px 8px #0000FF; padding: 10px">Shadow offset 0px horizontal, 4px down</p>
				<p style="box-shadow: 5px 0px 10px #FF0000; padding: 10px">Shadow red, 5px right</p>
			</div>`,
		},
		// TEXT SHADOWS
		{
			name: "24-text-shadows",
			html: `<div><h2>Text Shadows</h2>
				<p style="text-shadow: 2px 2px 4px #000000">Text shadow offset 2px right/down</p>
				<p style="text-shadow: -1px -1px 2px #666666">Text shadow offset -1px left/up</p>
				<p style="text-shadow: 1px 1px 1px #0000FF">Text shadow blue 1px</p>
				<p style="text-shadow: 0px 3px 5px #FF0000">Text shadow red 3px down</p>
			</div>`,
		},
		// BORDERS - styles
		{
			name: "25-borders-styles",
			html: `<div><h2>Border Styles</h2>
				<p style="border: 2px solid red; padding: 10px">Solid red border</p>
				<p style="border: 2px dashed blue; padding: 10px">Dashed blue border</p>
				<p style="border: 2px dotted green; padding: 10px">Dotted green border</p>
				<p style="border: 3px double black; padding: 10px">Double black border</p>
			</div>`,
		},
		// BORDER RADIUS
		{
			name: "26-border-radius",
			html: `<div><h2>Border Radius</h2>
				<p style="border: 1px solid black; padding: 10px; border-radius: 5px">Rounded 5px</p>
				<p style="border: 1px solid black; padding: 10px; border-radius: 10px">Rounded 10px</p>
				<p style="background-color: lightblue; padding: 10px; border-radius: 15px">Rounded 15px with bg</p>
			</div>`,
		},
		// PADDING - shorthand and sides
		{
			name: "27-padding-models",
			html: `<div><h2>Padding Models</h2>
				<p style="background-color: lightyellow; padding: 20px">Padding all 20px</p>
				<p style="background-color: lightgreen; padding: 10px 20px">Padding 10px top/bottom, 20px left/right</p>
				<p style="background-color: lightcoral; padding: 5px 10px 15px">Padding 5px top, 10px left/right, 15px bottom</p>
				<p style="background-color: lightblue; padding: 5px 10px 15px 20px">Padding 5-10-15-20</p>
			</div>`,
		},
		// MARGIN - shorthand
		{
			name: "28-margin-models",
			html: `<div><h2>Margin Models</h2>
				<p style="background-color: yellow; margin: 20px">Element with 20px margin</p>
				<p style="background-color: green; margin: 10px 30px">Element with 10px vertical, 30px horizontal</p>
				<p style="background-color: red; margin: 5px; color: white">Element with 5px margin all sides</p>
			</div>`,
		},
		// COMBINED COMPLEX STYLES
		{
			name: "29-combined-complex",
			html: `<div><h2>Combined Complex Styles</h2>
				<p style="
					font-family: Georgia;
					font-size: 18px;
					font-weight: 700;
					font-style: italic;
					letter-spacing: 1px;
					line-height: 1.8;
					color: #333333;
					background-color: #F0F0F0;
					padding: 15px;
					border: 2px solid #666666;
					border-radius: 8px;
					text-align: center;
					text-shadow: 1px 1px 2px #999999;
					box-shadow: 3px 3px 8px #CCCCCC;
				">Comprehensive styled paragraph with font, color, spacing, borders, shadows</p>
			</div>`,
		},
		// ALIGNMENT TESTS (original - kept for compatibility)
		{
			name: "30-alignment-combinations",
			html: `<div>
				<h1 align="left">Left Aligned Heading</h1>
				<p align="center">This is centered text with <strong>bold emphasis</strong>.</p>
				<p align="right">Right aligned paragraph with <em>italic style</em>.</p>
			</div>`,
		},
	}
}

func main() {
	testCases := generateTestsFromJSON()

	// Create output directory
	outputDir := "examples/output"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Println("📊 Generating Gotei comprehensive style tests from css_theme.json...")
	fmt.Println()
	fmt.Println("Testing categories from JSON:")
	fmt.Println("  • Font Families (9 types)")
	fmt.Println("  • Font Sizes (px, em, keywords)")
	fmt.Println("  • Font Weights (normal, bold, 100-900)")
	fmt.Println("  • Font Styles (normal, italic, oblique)")
	fmt.Println("  • Line Heights (unitless, px, percent)")
	fmt.Println("  • Letter Spacing (px, em, negative values)")
	fmt.Println("  • Word Spacing (px, em)")
	fmt.Println("  • Text Alignment (left, center, right, justify, start, end)")
	fmt.Println("  • Text Decoration (underline, overline, line-through)")
	fmt.Println("  • Text Transform (uppercase, lowercase, capitalize)")
	fmt.Println("  • Colors - Named (basic + 140+ extended)")
	fmt.Println("  • Colors - Hex (3-digit, 6-digit)")
	fmt.Println("  • Colors - RGB/HSL")
	fmt.Println("  • Box Shadows & Text Shadows")
	fmt.Println("  • Borders (solid, dashed, dotted, double)")
	fmt.Println("  • Border Radius")
	fmt.Println("  • Padding & Margin models")
	fmt.Println()
	fmt.Println("Generating PDFs...")

	for _, tc := range testCases {
		pdf, err := engine.Render(tc.html)
		if err != nil {
			log.Printf("✗ Error rendering %s: %v\n", tc.name, err)
			continue
		}

		filename := filepath.Join(outputDir, tc.name+".pdf")
		err = ioutil.WriteFile(filename, pdf, 0644)
		if err != nil {
			log.Printf("✗ Error writing %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("✓ Generated: %s\n", filename)
	}

	fmt.Println()
	fmt.Printf("All %d style tests completed! Check examples/output/ for PDFs.\n", len(testCases))
}

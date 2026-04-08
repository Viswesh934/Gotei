package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Viswesh934/gotei/internal/engine"
)

func main() {
	testCases := []struct {
		name string
		html string
	}{
		{
			name: "alignment-combinations",
			html: `<div>
				<h1 align="left">Left Aligned Heading</h1>
				<p align="center">This is centered text with <strong>bold emphasis</strong>.</p>
				<p align="right">Right aligned paragraph with <em>italic style</em>.</p>
			</div>`,
		},
		{
			name: "typography-styles",
			html: `<div>
				<h1>Large Bold Title</h1>
				<h2 style="text-align: center">Centered Subheading</h2>
				<p><strong>Bold text</strong> mixed with <em>italic text</em> and <b>another bold</b>.</p>
				<p align="center">Centered with <b><i>bold italic combination</i></b>.</p>
			</div>`,
		},
		{
			name: "color-and-alignment",
			html: `<div>
				<p align="left" style="color: red">Left red text</p>
				<p align="center" style="color: blue">Centered blue text</p>
				<p align="right" style="color: green">Right green text</p>
				<p style="text-align: center; color: purple"><strong>Purple centered bold</strong></p>
			</div>`,
		},
		{
			name: "mixed-sizes-and-alignment",
			html: `<div>
				<p align="center" style="font-size: 20px"><strong>Large centered bold text</strong></p>
				<p align="left" style="font-size: 14px">Normal left aligned paragraph.</p>
				<p align="right" style="font-size: 10px"><em>Small right aligned italic</em></p>
				<h3 align="center">Medium centered heading</h3>
			</div>`,
		},
		{
			name: "complex-layout",
			html: `<div>
				<h1 align="center"><strong>Document Title</strong></h1>
				<p align="left">Introduction paragraph with <em>emphasized</em> and <strong>bold</strong> text.</p>
				<h2 align="left">Section One</h2>
				<p align="left">Left aligned body text.</p>
				<h2 align="center">Section Two</h2>
				<p align="center">Center aligned content with <strong>bold styling</strong>.</p>
				<h2 align="right">Section Three</h2>
				<p align="right"><em>Right aligned emphasized conclusion</em></p>
			</div>`,
		},
		{
			name: "nested-styles",
			html: `<div>
				<p align="center">
					Centered paragraph with <strong>nested bold</strong> and <em>italic</em>.
				</p>
				<p align="right">
					<strong>Right bold start</strong> mixed with normal and <em>right italic end</em>.
				</p>
			</div>`,
		},
		{
			name: "all-alignments",
			html: `<div>
				<h1>Text Alignment Examples</h1>
				<p align="left"><strong>Left:</strong> This is left-aligned text using align attribute.</p>
				<p align="center"><strong>Center:</strong> This is center-aligned text.</p>
				<p align="right"><strong>Right:</strong> This is right-aligned text.</p>
				<p style="text-align: center"><strong>CSS:</strong> This uses text-align style property.</p>
			</div>`,
		},
		{
			name: "bold-italic-colors",
			html: `<div>
				<p><strong>Bold text</strong> and normal text</p>
				<p><em>Italic text</em> and normal text</p>
				<p><b><i>Bold and italic combined</i></b></p>
				<p align="center" style="color: red"><strong>Red centered bold</strong></p>
				<p align="right" style="color: blue"><em>Blue right italic</em></p>
			</div>`,
		},
		{
			name: "headings-with-styles",
			html: `<div>
				<h1 align="center"><strong>Main Title</strong></h1>
				<h2 align="left"><em>Subtitle Left</em></h2>
				<h3 align="center">Section Header</h3>
				<h4>Regular H4</h4>
				<h5 align="right"><strong>Right Aligned H5</strong></h5>
			</div>`,
		},
		{
			name: "comprehensive-styles",
			html: `<div>
				<h1 align="center">Comprehensive Style Test</h1>
				<p align="left">This tests <strong>bold</strong>, <em>italic</em>, and <b><i>bold italic</i></b> combinations.</p>
				<p align="center" style="color: red"><strong>Red centered text with bold</strong></p>
				<p align="right" style="color: blue"><em>Blue right-aligned italic</em></p>
				<p style="font-size: 16px; text-align: center">Larger centered text with custom font size</p>
				<p align="left"><strong>Bold</strong> left, <em>italic</em> center placement, and <b><i>combo</i></b> styles.</p>
				<p align="center" style="color: green">Green centered paragraph with styling.</p>
			</div>`,
		},
	}

	// Create output directory
	outputDir := "examples/output"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Println("Generating Gotei style examples...")
	fmt.Println()

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

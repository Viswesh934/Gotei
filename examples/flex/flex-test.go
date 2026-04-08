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
			name: "flexbox-row-basic",
			html: `<div style="display: flex; gap: 10px;">
				<div style="flex-grow: 1; background-color: lightblue;">Flex 1</div>
				<div style="flex-grow: 2; background-color: lightgreen;">Flex 2</div>
				<div style="flex-grow: 1; background-color: lightcoral;">Flex 3</div>
			</div>`,
		},
		{
			name: "flexbox-column",
			html: `<div style="display: flex; flex-direction: column;">
				<h2>Header</h2>
				<p>Main content area with flex column layout</p>
				<p>Footer area</p>
			</div>`,
		},
		{
			name: "flexbox-centering",
			html: `<div style="display: flex; justify-content: center; align-items: center;">
				<h1>Centered Content</h1>
			</div>`,
		},
		{
			name: "flexbox-space-between",
			html: `<div style="display: flex; justify-content: space-between;">
				<p>Item 1</p>
				<p>Item 2</p>
				<p>Item 3</p>
			</div>`,
		},
		{
			name: "line-height-test",
			html: `<div>
				<p style="line-height: 1.5;">Line height at 1.5x</p>
				<p style="line-height: 2;">Line height at 2.0x</p>
				<p style="line-height: 1;">Line height at 1.0x (tight)</p>
			</div>`,
		},
	}

	// Create output directory
	outputDir := "examples/flex-output"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Println("Generating Phase 3 flexbox & typography examples...")
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
	fmt.Printf("All %d flexbox tests completed! Check examples/flex-output/ for PDFs.\n", len(testCases))
}

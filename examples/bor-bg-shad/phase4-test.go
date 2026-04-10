package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Viswesh934/gotei/internal/engine"
)

func generatePhase4Examples() error {
	outputDir := "phase4-output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	tests := []struct {
		name string
		html string
	}{
		{
			name: "borders-solid",
			html: `
<html>
<head>
<style>
body { font-size: 14px; margin: 20px; }
.box { 
  width: 200px; 
  height: 100px; 
  border: 2px solid #333333;
  padding: 10px;
  margin: 20px;
}
.thick { border: 5px solid #ff0000; }
.thin { border: 1px solid #0000ff; }
</style>
</head>
<body>
<h1>Border Styles</h1>
<div class="box">Solid Black Border</div>
<div class="box thick">Thick Red Border</div>
<div class="box thin">Thin Blue Border</div>
</body>
</html>
			`,
		},
		{
			name: "backgrounds",
			html: `
<html>
<head>
<style>
body { font-size: 14px; margin: 20px; }
.bg-box {
  width: 200px;
  height: 100px;
  padding: 10px;
  margin: 20px;
  color: white;
  font-weight: bold;
}
.red { background-color: #ff6b6b; }
.blue { background-color: #4ecdc4; }
.green { background-color: #95e77d; }
</style>
</head>
<body>
<h1>Background Colors</h1>
<div class="bg-box red">Red Background</div>
<div class="bg-box blue">Teal Background</div>
<div class="bg-box green">Green Background</div>
</body>
</html>
			`,
		},
		{
			name: "borders-with-background",
			html: `
<html>
<head>
<style>
body { font-size: 14px; margin: 20px; }
.card {
  width: 250px;
  padding: 15px;
  margin: 20px;
  border: 3px solid #333333;
  background-color: #f0f0f0;
}
.card-header { font-weight: bold; color: #333333; }
.card-body { color: #666666; }
</style>
</head>
<body>
<h1>Cards with Borders and Background</h1>
<div class="card">
  <div class="card-header">Card Title</div>
  <div class="card-body">Card content with border and background</div>
</div>
<div class="card">
  <div class="card-header">Another Card</div>
  <div class="card-body">More content here</div>
</div>
</body>
</html>
			`,
		},
		{
			name: "mixed-styles",
			html: `
<html>
<head>
<style>
body { font-size: 14px; margin: 20px; }
.section {
  margin: 20px 0;
  padding: 15px;
  border: 2px solid #cccccc;
  background-color: #f9f9f9;
}
.highlight {
  padding: 10px;
  background-color: #fff3cd;
  border-left: 4px solid #ffc107;
  margin: 10px 0;
}
.info {
  padding: 10px;
  background-color: #d1ecf1;
  border: 1px solid #0c5460;
  color: #0c5460;
}
</style>
</head>
<body>
<h1>Mixed Styles</h1>
<div class="section">
  <p>Regular section with border and background</p>
</div>
<div class="highlight">
  Important note highlighted
</div>
<div class="info">
  Information box with blue theme
</div>
</body>
</html>
			`,
		},
		{
			name: "padding-borders",
			html: `
<html>
<head>
<style>
body { font-size: 14px; margin: 20px; }
.padded {
  width: 200px;
  padding: 20px;
  margin: 20px;
  border: 2px solid #333333;
  background-color: #eeeeee;
}
.small-padding {
  width: 200px;
  padding: 5px;
  margin: 20px;
  border: 2px solid #ff0000;
  background-color: #ffe0e0;
}
.large-padding {
  width: 200px;
  padding: 30px;
  margin: 20px;
  border: 2px solid #0000ff;
  background-color: #e0e0ff;
}
</style>
</head>
<body>
<h1>Padding and Borders</h1>
<div class="padded">Standard padding</div>
<div class="small-padding">Small padding (5px)</div>
<div class="large-padding">Large padding (30px)</div>
</body>
</html>
			`,
		},
	}

	for _, test := range tests {
		fmt.Printf("Generating phase 4 test: %s\n", test.name)

		// Use the shared engine path so CSS in <style> tags is resolved consistently.
		pdfBytes, err := engine.Render(test.html)
		if err != nil {
			fmt.Printf("  Error rendering PDF: %v\n", err)
			continue
		}

		// Write PDF
		outputPath := filepath.Join(outputDir, test.name+".pdf")
		if err := os.WriteFile(outputPath, pdfBytes, 0644); err != nil {
			fmt.Printf("  Error writing PDF: %v\n", err)
			continue
		}

		fmt.Printf("  ✓ Generated: %s\n", outputPath)
	}

	return nil
}

func main() {
	fmt.Println("=== Phase 4 Examples: Borders, Backgrounds, Shadows ===")
	if err := generatePhase4Examples(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n✓ Phase 4 examples completed")
}

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/Viswesh934/gotei/internal/engine"
)

func main() {
	outputDir := "examples/phase2/output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	imagePath := filepath.Join(outputDir, "phase2-sample.png")
	if err := generateSampleImage(imagePath); err != nil {
		panic(err)
	}

	tests := []struct {
		name string
		html string
	}{
		{
			name: "phase2-hyperlinks",
			html: `<div>
				<h2>Hyperlink Support</h2>
				<p>Visit <a href="https://github.com/Viswesh934/Gotei">Gotei Repository</a> for source code.</p>
				<p>Documentation link: <a href="https://www.w3.org/TR/CSS22/">CSS 2.2 Spec</a></p>
			</div>`,
		},
		{
			name: "phase2-images",
			html: fmt.Sprintf(`<div>
				<h2>Image Support</h2>
				<p>Image loaded from local path:</p>
				<img src="%s" style="width: 280px; height: 120px; border: 1px solid #111827;" />
			</div>`, imagePath),
		},
		{
			name: "phase2-pagination",
			html: buildPaginationHTML(),
		},
		{
			name: "phase2-table-pagination",
			html: buildTablePaginationHTML(),
		},
	}

	for _, tc := range tests {
		pdfBytes, err := engine.Render(tc.html)
		if err != nil {
			fmt.Printf("x render failed for %s: %v\n", tc.name, err)
			continue
		}

		outPath := filepath.Join(outputDir, tc.name+".pdf")
		if err := os.WriteFile(outPath, pdfBytes, 0644); err != nil {
			fmt.Printf("x write failed for %s: %v\n", outPath, err)
			continue
		}

		fmt.Printf("ok %s\n", outPath)
	}
}

func buildPaginationHTML() string {
	var b strings.Builder
	b.WriteString(`<div><h2>Pagination Stress Test</h2>`)
	for i := 1; i <= 120; i++ {
		b.WriteString(fmt.Sprintf(`<p style="margin: 4px 0;">Line %d: This content is used to force multi-page PDF output with repeated flow checks.</p>`, i))
	}
	b.WriteString(`</div>`)
	return b.String()
}

func buildTablePaginationHTML() string {
	var b strings.Builder
	b.WriteString(`<div><h2>Table Pagination Stress Test</h2>`)
	b.WriteString(`<table border="1" cellpadding="6" cellspacing="0" style="width: 560px;">`)
	b.WriteString(`<thead><tr><th>ID</th><th>Owner</th><th>Summary</th></tr></thead><tbody>`)
	for i := 1; i <= 140; i++ {
		b.WriteString(fmt.Sprintf(`<tr><td>%d</td><td>User %d</td><td>Row %d content for pagination validation with consistent width and row height.</td></tr>`, i, i, i))
	}
	b.WriteString(`</tbody></table></div>`)
	return b.String()
}

func generateSampleImage(path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 560, 240))
	for y := 0; y < 240; y++ {
		for x := 0; x < 560; x++ {
			r := uint8(20 + (x * 180 / 560))
			g := uint8(60 + (y * 120 / 240))
			b := uint8(140)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	for y := 90; y < 150; y++ {
		for x := 120; x < 440; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

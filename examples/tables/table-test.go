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
			name: "table-basic-border",
			html: `<table border="1" cellpadding="6" cellspacing="0" style="width: 520px;">
				<tr>
					<th>Item</th>
					<th>Qty</th>
					<th>Price</th>
				</tr>
				<tr>
					<td>Notebook</td>
					<td>2</td>
					<td>$12.50</td>
				</tr>
				<tr>
					<td>Pen Set</td>
					<td>1</td>
					<td>$8.00</td>
				</tr>
			</table>`,
		},
		{
			name: "table-colspan",
			html: `<table border="1" cellpadding="5" cellspacing="3" style="width: 540px;">
				<tr>
					<th colspan="3" style="background-color: #dbeafe;">Monthly Summary</th>
				</tr>
				<tr>
					<th>Department</th>
					<th>Tickets</th>
					<th>Resolution</th>
				</tr>
				<tr>
					<td>Support</td>
					<td>182</td>
					<td>94%</td>
				</tr>
				<tr>
					<td>Billing</td>
					<td>66</td>
					<td>91%</td>
				</tr>
			</table>`,
		},
		{
			name: "table-mixed-styling",
			html: `<div>
				<h2 style="margin-bottom: 8px;">Table Styling Test</h2>
				<table border="1" cellpadding="8" cellspacing="4" style="width: 560px; border-color: #1f2937;">
					<tr>
						<th style="background-color: #111827; color: white;">Name</th>
						<th style="background-color: #111827; color: white;">Role</th>
						<th style="background-color: #111827; color: white;">Status</th>
					</tr>
					<tr>
						<td style="background-color: #f8fafc;">Ari</td>
						<td style="background-color: #f8fafc;">Frontend</td>
						<td style="color: #059669; font-weight: bold;">Active</td>
					</tr>
					<tr>
						<td style="background-color: #f8fafc;">Mira</td>
						<td style="background-color: #f8fafc;">Backend</td>
						<td style="color: #d97706; font-weight: bold;">On Leave</td>
					</tr>
				</table>
			</div>`,
		},
	}

	outputDir := "examples/tables/output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	fmt.Println("Generating table support examples...")
	for _, tc := range testCases {
		pdfBytes, err := engine.Render(tc.html)
		if err != nil {
			log.Printf("x render failed for %s: %v", tc.name, err)
			continue
		}

		outPath := filepath.Join(outputDir, tc.name+".pdf")
		if err := ioutil.WriteFile(outPath, pdfBytes, 0644); err != nil {
			log.Printf("x write failed for %s: %v", outPath, err)
			continue
		}

		fmt.Printf("ok %s\n", outPath)
	}

	fmt.Printf("Done. Generated %d table test PDFs in %s\n", len(testCases), outputDir)
}

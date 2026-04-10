package render

import (
	"codeberg.org/go-pdf/fpdf"
	"github.com/Viswesh934/gotei/internal/layout"
	"github.com/Viswesh934/gotei/internal/style"
)

// BorderRendering handles all border drawing
type BorderRendering struct {
	pdf *fpdf.Fpdf
}

// DrawBorder renders a border around a box
func (br *BorderRendering) DrawBorder(box *layout.Box, style style.Style) {
	if style.Border.Width == 0 {
		return
	}

	// Set border properties
	br.pdf.SetLineWidth(style.Border.Width)
	setColorFromString(br.pdf, style.Border.Color)

	switch style.Border.Style {
	case "solid":
		if style.Border.Radius > 0 {
			br.pdf.RoundedRect(box.X, box.Y, box.Width, box.Height, style.Border.Radius, "1234", "")
		} else {
			br.pdf.Rect(box.X, box.Y, box.Width, box.Height, "")
		}
	case "dashed":
		br.drawDashedRect(box, 5, 5)
	case "dotted":
		br.drawDottedRect(box, 2)
	}

	// Reset to default
	br.pdf.SetLineWidth(0.5)
	br.pdf.SetDrawColor(0, 0, 0)
}

// DrawBoxShadow renders a shadow effect around a box
func (br *BorderRendering) DrawBoxShadow(box *layout.Box, shadow style.Shadow) {
	if shadow.Blur == 0 && shadow.OffsetX == 0 && shadow.OffsetY == 0 {
		return
	}

	// Save original color
	setColorFromString(br.pdf, shadow.Color)

	// Draw shadow slightly offset
	shadowX := box.X + shadow.OffsetX
	shadowY := box.Y + shadow.OffsetY

	// Reduce opacity for shadow effect
	br.pdf.SetAlpha(0.3, "")
	br.pdf.Rect(shadowX, shadowY, box.Width, box.Height, "F")
	br.pdf.SetAlpha(1, "")

	// Reset color
	br.pdf.SetDrawColor(0, 0, 0)
}

// DrawBackground fills the background of a box
func (br *BorderRendering) DrawBackground(box *layout.Box, bgColor string) {
	if bgColor == "" || bgColor == "white" || bgColor == "transparent" {
		return
	}

	setColorFromString(br.pdf, bgColor)
	br.pdf.Rect(box.X, box.Y, box.Width, box.Height, "F")
	br.pdf.SetDrawColor(0, 0, 0) // Reset
}

// Helper methods

func (br *BorderRendering) drawDashedRect(box *layout.Box, dashLen, gapLen float64) {
	// Simplified dashed rectangle - draw dashes along each side
	br.pdf.SetDashPattern([]float64{dashLen, gapLen}, 0)
	br.pdf.Rect(box.X, box.Y, box.Width, box.Height, "")
	br.pdf.SetDashPattern([]float64{}, 0) // Reset to solid
}

func (br *BorderRendering) drawDottedRect(box *layout.Box, dotSize float64) {
	// Simplified dotted rectangle - draw small circles/dots
	br.pdf.SetFillColor(0, 0, 0)

	// Top and bottom lines
	for x := box.X; x < box.X+box.Width; x += dotSize * 2 {
		br.pdf.Circle(x, box.Y, dotSize/2, "F")
		br.pdf.Circle(x, box.Y+box.Height, dotSize/2, "F")
	}

	// Left and right lines
	for y := box.Y; y < box.Y+box.Height; y += dotSize * 2 {
		br.pdf.Circle(box.X, y, dotSize/2, "F")
		br.pdf.Circle(box.X+box.Width, y, dotSize/2, "F")
	}
}

// BorderRadius applies rounded corners (simplified)
func (br *BorderRendering) DrawBorderRadius(box *layout.Box, radius float64) {
	if radius == 0 {
		return
	}

	br.pdf.RoundedRect(box.X, box.Y, box.Width, box.Height, radius, "1234", "")
}

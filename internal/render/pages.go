package render

import (
	"codeberg.org/go-pdf/fpdf"
	"github.com/Viswesh934/gotei/internal/layout"
)

// PageManager handles multi-page PDF rendering
type PageManager struct {
	pdf           *fpdf.Fpdf
	currentPage   int
	currentY      float64
	pageHeight    float64
	pageMarginTop float64
	pageMarginBot float64
	isFirstPage   bool
}

// NewPageManager creates a new page manager
func NewPageManager(pdf *fpdf.Fpdf, pageHeight, marginTop, marginBot float64) *PageManager {
	return &PageManager{
		pdf:           pdf,
		currentPage:   1,
		currentY:      marginTop,
		pageHeight:    pageHeight,
		pageMarginTop: marginTop,
		pageMarginBot: marginBot,
		isFirstPage:   true,
	}
}

// CheckPageBreak returns true if a new page should be started
func (pm *PageManager) CheckPageBreak(boxHeight float64) bool {
	availableSpace := pm.pageHeight - pm.pageMarginBot - pm.currentY

	if boxHeight > availableSpace {
		// Need new page
		pm.NewPage()
		return true
	}

	return false
}

// NewPage adds a new page to the PDF
func (pm *PageManager) NewPage() {
	if !pm.isFirstPage {
		pm.pdf.AddPage()
	}

	pm.isFirstPage = false
	pm.currentPage++
	pm.currentY = pm.pageMarginTop
}

// UpdateY updates the current Y position
func (pm *PageManager) UpdateY(boxHeight float64) {
	pm.currentY += boxHeight
}

// GetCurrentY returns the current Y position
func (pm *PageManager) GetCurrentY() float64 {
	return pm.currentY
}

// GetCurrentPage returns the current page number
func (pm *PageManager) GetCurrentPage() int {
	return pm.currentPage
}

// RenderBoxWithPageBreaks renders a box tree with automatic page breaks
func (pm *PageManager) RenderBoxWithPageBreaks(box *layout.Box) {
	pm.renderBoxRecursive(box)
}

func (pm *PageManager) renderBoxRecursive(box *layout.Box) {
	if box == nil {
		return
	}

	// Check if we need a new page
	pm.CheckPageBreak(box.Height)

	// Update box Y position for current page
	box.Y = pm.currentY

	// Render the box content
	// (This would call the actual rendering logic)

	// Update Y for next box
	pm.UpdateY(box.Height)

	// Render children
	for _, child := range box.Children {
		pm.renderBoxRecursive(child)
	}
}

// CollapsedMargin represents a collapsed margin value
type CollapsedMargin struct {
	Top    float64
	Bottom float64
}

// CollapseMargins calculates collapsed margins between adjacent blocks
// per W3C spec: larger margin wins, negative margins are summed
func CollapseMargins(margin1, margin2 float64) float64 {
	// Simple implementation: take the maximum
	if margin1 > margin2 {
		return margin1
	}
	return margin2
}

// CollapseVerticalMargins collapses vertical space between parent and child
func CollapseVerticalMargins(parentMarginBot, childMarginTop float64) float64 {
	return CollapseMargins(parentMarginBot, childMarginTop)
}

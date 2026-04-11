package layout

import (
	"math"
	"strconv"
	"strings"

	"github.com/Viswesh934/gotei/internal/dom"
)

// LayoutTable handles table-like layout and falls back to block flow for non-table nodes.
func LayoutTable(box *Box, x, y, maxWidth float64) float64 {
	if box == nil {
		return 0
	}

	display := strings.ToLower(box.Style.Display)
	tag := ""
	if box.Node != nil {
		tag = strings.ToLower(box.Node.Tag)
	}

	isTable := display == "table" || tag == "table"
	if !isTable {
		return layoutBlock(box, x, y)
	}

	box.X = x
	box.Y = y
	box.Width = maxWidth
	if box.Style.WidthPercent > 0 {
		box.Width = maxWidth * box.Style.WidthPercent / 100.0
	}
	if box.Style.Width > 0 {
		box.Width = box.Style.Width
	}
	if box.Width < 0 {
		box.Width = 0
	}

	contentX := box.X + box.Style.Padding.Left
	contentY := box.Y + box.Style.Padding.Top
	contentWidth := box.Width - box.Style.Padding.Left - box.Style.Padding.Right
	if contentWidth < 0 {
		contentWidth = 0
	}

	pruneGeneratedHeaderRows(box)

	entries, headerTemplates := tableRowEntries(box)
	if len(entries) == 0 {
		box.Height = box.Style.Padding.Top + box.Style.Padding.Bottom
		return box.Height + box.Style.Margin.Top + box.Style.Margin.Bottom
	}

	rows := make([]*Box, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows, entry.row)
	}

	colCount := countTableColumns(rows)
	if colCount == 0 {
		colCount = 1
	}

	colWidths := make([]float64, colCount)
	colGap := box.Style.Gap
	effectiveContentWidth := contentWidth - colGap*float64(colCount-1)
	if effectiveContentWidth < 0 {
		effectiveContentWidth = 0
	}
	baseWidth := 0.0
	if colCount > 0 {
		baseWidth = effectiveContentWidth / float64(colCount)
	}
	for i := range colWidths {
		colWidths[i] = baseWidth
	}

	// Expand columns if any single-span cell has a stronger width requirement.
	for _, row := range rows {
		cells := tableCells(row)
		colIndex := 0
		for _, cell := range cells {
			span := cellColSpan(cell)
			if span < 1 {
				span = 1
			}
			if span == 1 && colIndex < len(colWidths) {
				minWidth := tableCellMinWidth(cell)
				if minWidth > colWidths[colIndex] {
					colWidths[colIndex] = minWidth
				}
			}
			colIndex += span
		}
	}

	fitTableToWidth(colWidths, effectiveContentWidth)

	currentY := contentY
	rowGap := box.Style.Gap
	for _, entry := range entries {
		row := entry.row
		rowHeight := layoutTableRow(row, box, contentX, contentWidth, colWidths, colGap, currentY)

		if shouldMoveRowToNextPage(currentY, rowHeight) {
			currentY = nextPageStartY(currentY)

			if !entry.isHeader && len(headerTemplates) > 0 {
				for _, template := range headerTemplates {
					repeatHeader := cloneBoxTree(template)
					repeatHeader.Parent = box
					markGeneratedHeaderRow(repeatHeader)
					box.Children = append(box.Children, repeatHeader)

					headerHeight := layoutTableRow(repeatHeader, box, contentX, contentWidth, colWidths, colGap, currentY)
					currentY += headerHeight + rowGap
				}
			}

			rowHeight = layoutTableRow(row, box, contentX, contentWidth, colWidths, colGap, currentY)
		}

		currentY += rowHeight + rowGap
	}

	if len(rows) > 0 {
		currentY -= rowGap
	}

	box.Height = (currentY - contentY) + box.Style.Padding.Top + box.Style.Padding.Bottom
	if box.Height < 0 {
		box.Height = 0
	}

	return box.Height + box.Style.Margin.Top + box.Style.Margin.Bottom
}

type tableRowEntry struct {
	row      *Box
	isHeader bool
}

func tableRowEntries(table *Box) ([]tableRowEntry, []*Box) {
	if table == nil {
		return nil, nil
	}

	entries := []tableRowEntry{}
	headers := []*Box{}

	for _, child := range table.Children {
		if child == nil || child.Node == nil {
			continue
		}

		tag := strings.ToLower(child.Node.Tag)
		switch tag {
		case "tr":
			entries = append(entries, tableRowEntry{row: child, isHeader: false})
		case "thead", "tbody", "tfoot":
			isHeaderSection := tag == "thead"
			for _, grandChild := range child.Children {
				if grandChild == nil || grandChild.Node == nil || strings.ToLower(grandChild.Node.Tag) != "tr" {
					continue
				}
				entries = append(entries, tableRowEntry{row: grandChild, isHeader: isHeaderSection})
				if isHeaderSection {
					headers = append(headers, grandChild)
				}
			}
		}
	}

	return entries, headers
}

func layoutTableRow(row, table *Box, contentX, contentWidth float64, colWidths []float64, colGap, rowY float64) float64 {
	if row == nil {
		return 0
	}

	row.X = contentX
	row.Y = rowY
	row.Width = contentWidth

	cells := tableCells(row)
	colIndex := 0
	cellX := contentX
	maxRowHeight := 0.0

	for _, cell := range cells {
		if table != nil && table.Style.Border.Width > 0 && cell.Style.Border.Width == 0 {
			cell.Style.Border = table.Style.Border
		}

		span := cellColSpan(cell)
		if span < 1 {
			span = 1
		}

		cellWidth := 0.0
		for i := 0; i < span && colIndex+i < len(colWidths); i++ {
			cellWidth += colWidths[colIndex+i]
		}
		if span > 1 {
			cellWidth += colGap * float64(span-1)
		}

		cellOuterHeight := Layout(cell, cellX, rowY, cellWidth)
		if cellOuterHeight > maxRowHeight {
			maxRowHeight = cellOuterHeight
		}

		cellX += cellWidth + colGap
		colIndex += span
	}

	if maxRowHeight < 0 {
		maxRowHeight = 0
	}

	row.Height = maxRowHeight
	for _, cell := range cells {
		cell.Height = maxRowHeight
	}

	return maxRowHeight
}

func shouldMoveRowToNextPage(rowTop, rowHeight float64) bool {
	if rowHeight <= 0 {
		return false
	}

	pageBottom := nextPageStartY(rowTop)
	pageTop := pageBottom - DefaultHeight
	if almostEqual(rowTop, pageTop) {
		return false
	}

	return rowTop+rowHeight > pageBottom
}

func nextPageStartY(y float64) float64 {
	if y < 0 {
		y = 0
	}
	return (math.Floor(y/DefaultHeight) + 1) * DefaultHeight
}

func almostEqual(a, b float64) bool {
	delta := a - b
	if delta < 0 {
		delta = -delta
	}
	return delta < 0.001
}

func pruneGeneratedHeaderRows(table *Box) {
	if table == nil {
		return
	}

	filtered := make([]*Box, 0, len(table.Children))
	for _, child := range table.Children {
		if child == nil || child.Node == nil || child.Node.Attr == nil {
			filtered = append(filtered, child)
			continue
		}

		if child.Node.Attr["data-gotei-repeat-header"] == "1" {
			continue
		}

		filtered = append(filtered, child)
	}

	table.Children = filtered
}

func markGeneratedHeaderRow(row *Box) {
	if row == nil || row.Node == nil {
		return
	}
	if row.Node.Attr == nil {
		row.Node.Attr = map[string]string{}
	}
	row.Node.Attr["data-gotei-repeat-header"] = "1"
}

func cloneBoxTree(src *Box) *Box {
	if src == nil {
		return nil
	}

	clone := &Box{
		X:      src.X,
		Y:      src.Y,
		Width:  src.Width,
		Height: src.Height,
		Style:  src.Style,
		Node:   cloneDOMNode(src.Node),
	}

	for _, child := range src.Children {
		childClone := cloneBoxTree(child)
		if childClone != nil {
			childClone.Parent = clone
			clone.Children = append(clone.Children, childClone)
		}
	}

	return clone
}

func cloneDOMNode(src *dom.Node) *dom.Node {
	if src == nil {
		return nil
	}

	clone := &dom.Node{
		Type:    src.Type,
		Tag:     src.Tag,
		Content: src.Content,
	}

	if src.Attr != nil {
		clone.Attr = make(map[string]string, len(src.Attr))
		for k, v := range src.Attr {
			clone.Attr[k] = v
		}
	}

	return clone
}

func tableRows(table *Box) []*Box {
	if table == nil {
		return nil
	}

	var rows []*Box
	for _, child := range table.Children {
		if child == nil || child.Node == nil {
			continue
		}

		tag := strings.ToLower(child.Node.Tag)
		switch tag {
		case "tr":
			rows = append(rows, child)
		case "thead", "tbody", "tfoot":
			for _, grandChild := range child.Children {
				if grandChild != nil && grandChild.Node != nil && strings.ToLower(grandChild.Node.Tag) == "tr" {
					rows = append(rows, grandChild)
				}
			}
		}
	}

	return rows
}

func tableCells(row *Box) []*Box {
	if row == nil {
		return nil
	}

	cells := make([]*Box, 0, len(row.Children))
	for _, child := range row.Children {
		if child == nil || child.Node == nil {
			continue
		}
		tag := strings.ToLower(child.Node.Tag)
		if tag == "td" || tag == "th" {
			cells = append(cells, child)
		}
	}
	return cells
}

func countTableColumns(rows []*Box) int {
	maxCols := 0
	for _, row := range rows {
		cols := 0
		for _, cell := range tableCells(row) {
			span := cellColSpan(cell)
			if span < 1 {
				span = 1
			}
			cols += span
		}
		if cols > maxCols {
			maxCols = cols
		}
	}
	return maxCols
}

func cellColSpan(cell *Box) int {
	if cell == nil || cell.Node == nil || cell.Node.Attr == nil {
		return 1
	}
	raw := strings.TrimSpace(cell.Node.Attr["colspan"])
	if raw == "" {
		return 1
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 1
	}
	return n
}

func tableCellMinWidth(cell *Box) float64 {
	if cell == nil {
		return 0
	}

	if cell.Style.Width > 0 {
		return cell.Style.Width
	}

	textLen := 0
	collectTextLen(cell, &textLen)
	fontSize := cell.Style.FontSize
	if fontSize <= 0 {
		fontSize = 12
	}

	approxTextWidth := float64(textLen) * fontSize * 0.55
	return approxTextWidth + cell.Style.Padding.Left + cell.Style.Padding.Right
}

func collectTextLen(box *Box, total *int) {
	if box == nil || total == nil {
		return
	}

	if box.Node != nil && box.Node.Type == dom.TextNode {
		*total += len([]rune(box.Node.Content))
	}

	for _, child := range box.Children {
		collectTextLen(child, total)
	}
}

func fitTableToWidth(colWidths []float64, maxWidth float64) {
	if len(colWidths) == 0 {
		return
	}

	total := 0.0
	for _, w := range colWidths {
		total += w
	}
	if total <= 0 {
		equal := 0.0
		if len(colWidths) > 0 {
			equal = maxWidth / float64(len(colWidths))
		}
		for i := range colWidths {
			colWidths[i] = equal
		}
		return
	}

	if maxWidth <= 0 {
		return
	}

	scale := maxWidth / total
	for i := range colWidths {
		colWidths[i] *= scale
	}
}

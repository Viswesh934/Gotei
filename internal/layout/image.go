package layout

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func layoutImage(box *Box, x, y, maxWidth float64) float64 {
	box.X = x
	box.Y = y

	w, h := inferredImageSize(box, maxWidth)
	if box.Style.WidthPercent > 0 {
		w = maxWidth * box.Style.WidthPercent / 100.0
	}
	if box.Style.Width > 0 {
		w = box.Style.Width
	}
	if box.Style.HeightPercent > 0 && h > 0 {
		h = h * box.Style.HeightPercent / 100.0
	}
	if box.Style.Height > 0 {
		h = box.Style.Height
	}

	if maxWidth > 0 && w > maxWidth {
		ratio := 1.0
		if w > 0 {
			ratio = h / w
		}
		w = maxWidth
		h = w * ratio
	}

	if w <= 0 {
		w = 180
	}
	if h <= 0 {
		h = 120
	}

	box.Width = w
	box.Height = h
	return box.Height + box.Style.Margin.Top + box.Style.Margin.Bottom
}

func inferredImageSize(box *Box, maxWidth float64) (float64, float64) {
	if box == nil || box.Node == nil || box.Node.Attr == nil {
		if maxWidth > 0 {
			return maxWidth, 120
		}
		return 180, 120
	}

	src := box.Node.Attr["src"]
	if src == "" {
		if maxWidth > 0 {
			return maxWidth, 120
		}
		return 180, 120
	}

	f, err := os.Open(src)
	if err != nil {
		if maxWidth > 0 {
			return maxWidth, 120
		}
		return 180, 120
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		if maxWidth > 0 {
			return maxWidth, 120
		}
		return 180, 120
	}

	return float64(cfg.Width), float64(cfg.Height)
}

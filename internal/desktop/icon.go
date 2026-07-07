package desktop

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"fyne.io/fyne/v2"
)

// generateIcon creates a 64x64 icon for Task Hub Desktop:
// an indigo circle with a white checkmark representing a completed task.
func generateIcon() fyne.Resource {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	indigo := color.RGBA{R: 102, G: 126, B: 234, A: 255}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	cx, cy, r := 32, 32, 28

	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			dx, dy := x-cx, y-cy
			// Draw filled circle
			if dx*dx+dy*dy <= r*r {
				img.Set(x, y, indigo)
			}
		}
	}

	// Draw a white checkmark inside the circle
	// using anti-aliased line drawing
	drawLine(img, 20, 34, 28, 42, white)
	drawLine(img, 28, 42, 44, 24, white)
	drawLine(img, 20, 34, 28, 42, white)
	drawLine(img, 28, 42, 44, 24, white)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		// Fallback: return nil resource, app will use default
		return nil
	}
	return fyne.NewStaticResource("taskhub-icon", buf.Bytes())
}

// drawLine draws a thick anti-aliased line between (x1,y1) and (x2,y2) in the given color.
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy
	thickness := 3

	for {
		// Draw a thick point
		for ty := -thickness; ty <= thickness; ty++ {
			for tx := -thickness; tx <= thickness; tx++ {
				px, py := x1+tx, y1+ty
				if px >= 0 && px < 64 && py >= 0 && py < 64 {
					img.Set(px, py, c)
				}
			}
		}

		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

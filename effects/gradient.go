package effects

import (
	"image"
	"math"
)

// Gradient effect
type Gradient struct {
	Top    int
	Mid    int
	Bottom int
}

// Apply effect
func (g *Gradient) Apply(img *image.NRGBA, offset *image.Point, adv *int) {
	top := math.Min(float64(g.Top)/100.0, 1.0)
	top = math.Max(top, 0.0)
	mid := math.Min(float64(g.Mid)/100.0, 1.0)
	mid = math.Max(mid, 0.0)
	bottom := math.Min(float64(g.Bottom)/100.0, 1.0)
	bottom = math.Max(bottom, 0.0)

	rect := img.Bounds()
	midY := rect.Dy() / 2

	step := (mid - top) / float64(midY)
	mult := top
	for y := rect.Min.Y; y < rect.Min.Y+midY; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			p := img.NRGBAAt(x, y)
			p.R = uint8(float64(p.R) * mult)
			p.G = uint8(float64(p.G) * mult)
			p.B = uint8(float64(p.B) * mult)
			img.SetNRGBA(x, y, p)
		}
		mult += step
	}

	step = (bottom - mid) / float64(rect.Dy()-midY)
	mult = mid
	for y := rect.Min.Y + midY; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			p := img.NRGBAAt(x, y)
			p.R = uint8(float64(p.R) * mult)
			p.G = uint8(float64(p.G) * mult)
			p.B = uint8(float64(p.B) * mult)
			img.SetNRGBA(x, y, p)
		}
		mult += step
	}
}

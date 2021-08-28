package effects

import (
	"image"
	"image/draw"
)

// Outline effect
type Outline struct {
	Width int
}

// Apply effect
func (e *Outline) Apply(img *image.NRGBA, offset *image.Point, adv *int) {

	srcRect := img.Bounds()
	outRect := image.Rectangle{
		srcRect.Min.Sub(image.Pt(e.Width, e.Width)),
		srcRect.Max.Add(image.Pt(e.Width, e.Width)),
	}
	outImg := image.NewNRGBA(outRect)
	shadowImg := createShadow(img, image.Pt(0, 0))
	for x := -e.Width; x <= e.Width; x++ {
		for y := -e.Width; y <= e.Width; y++ {
			rect := srcRect.Sub(image.Pt(x, y))
			draw.Draw(outImg, rect, shadowImg, image.Pt(0, 0), draw.Over)
		}
	}
	draw.Draw(outImg, srcRect, img, image.Pt(0, 0), draw.Over)

	// place top-left corner at (0, 0) after blur
	outImg.Rect = outImg.Rect.Sub(outImg.Rect.Min)

	*img = *outImg
	*offset = offset.Sub(image.Pt(e.Width, e.Width))
	*adv += e.Width
}

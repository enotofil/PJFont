package effects

import (
	"image"
	"image/draw"
)

// Shadow effect
type Shadow struct {
	Dist image.Point
	Blur int
}

// Apply effect
func (s *Shadow) Apply(img *image.NRGBA, offset *image.Point, adv *int) {
	if s.Blur == 0 && s.Dist.Eq(image.Point{}) {
		return
	}
	shadowImg := createShadow(img, s.Dist)
	for i := 0; i < s.Blur; i++ {
		shadowImg = addBlurStep(shadowImg)
	}
	rect := img.Bounds().Sub(s.Dist)
	draw.Draw(shadowImg, rect, img, image.Pt(0, 0), draw.Over)

	// place top-left corner at (0, 0) after blur
	shadowImg.Rect = shadowImg.Rect.Sub(shadowImg.Rect.Min)

	// replace original glyph image
	*img = *shadowImg

	*offset = offset.Sub(image.Pt(s.Blur, s.Blur))
}

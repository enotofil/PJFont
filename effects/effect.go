package effects

import (
	"image"
	"image/color"
	"image/draw"
)

// Effect generic
type Effect interface {
	Apply(img *image.NRGBA, offset *image.Point, adv *int)
}

func createShadow(img *image.NRGBA, Dist image.Point) *image.NRGBA {
	rect := img.Bounds()
	shadowRect := image.Rectangle{
		rect.Min.Sub(Dist),
		rect.Max,
	}
	blackImg := image.NewNRGBA(shadowRect)
	draw.Draw(blackImg, rect, img, image.Pt(0, 0), draw.Src)
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			a := img.NRGBAAt(x, y).A
			blackImg.SetNRGBA(x, y, color.NRGBA{0, 0, 0, a})
		}
	}
	return blackImg
}

func addBlurStep(src *image.NRGBA) *image.NRGBA {
	rect := src.Bounds()
	fxRect := image.Rectangle{
		rect.Min.Sub(image.Pt(1, 1)),
		rect.Max.Add(image.Pt(1, 1)),
	}
	fxImg := image.NewNRGBA(fxRect)
	for x := fxRect.Min.X; x < fxRect.Max.X; x++ {
		for y := fxRect.Min.Y; y < fxRect.Max.Y; y++ {
			aSum := 0
			for x1 := -1; x1 <= 1; x1++ {
				for y1 := -1; y1 <= 1; y1++ {
					aSum += int(src.NRGBAAt(x+x1, y+y1).A)
				}
			}
			fxImg.SetNRGBA(x, y, color.NRGBA{0, 0, 0, uint8(aSum / 9)})
		}
	}
	return fxImg
}

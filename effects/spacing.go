package effects

import "image"

// Spacing between characters increase
type Spacing struct {
	AdvPlus int
}

// Apply effect
func (sp *Spacing) Apply(img *image.NRGBA, offset *image.Point, adv *int) {
	*adv += sp.AdvPlus
}

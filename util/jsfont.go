package util

// engine independent, no Ebiten here
import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/enotofil/PJFont/effects"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type jsGlyph struct {
	Code                         rune
	X, Y, W, H, XOff, YOff, XAdv int
}

type jsFont struct {
	Name       string
	LineHeight int
	Chars      []jsGlyph
	img        *image.NRGBA
	fit        bool
}

func newJsFont(
	name string,
	face font.Face,
	runes []rune,
	effectList []effects.Effect,
	lineHeight int,
	imgSize image.Point,
) *jsFont {

	chars := []jsGlyph{}
	packImg := image.NewNRGBA(image.Rectangle{image.Pt(0, 0), imgSize})
	x := 1
	y := 1
	maxW := 0
	fit := true

	for _, r := range runes {
		rect, faceImg, p, adv26_6, _ := face.Glyph(fixed.P(0, 0), r)

		offset := rect.Min
		adv := adv26_6.Round()

		if rect.Dx() == 0 { // space bar
			rect.Max = rect.Max.Add(image.Pt(1, 1)) // ebiten panics if img width=0
		}

		// create NRGBA image from Grey, align image rect to (0, 0)
		zeroRect := rect.Sub(rect.Min)
		img := image.NewNRGBA(zeroRect)
		draw.Draw(img, zeroRect, faceImg, p, draw.Src)

		for _, e := range effectList {
			e.Apply(img, &offset, &adv)
		}

		if img.Rect.Dx() > maxW {
			maxW = img.Rect.Dx()
		}

		if (y + img.Rect.Dy()) > imgSize.Y {
			y = 1
			x += maxW
			maxW = 0
		}
		// font doesnt fit in image
		if x+img.Rect.Dx() > imgSize.X {
			fit = false
			break
		}

		packRect := img.Rect.Add(image.Pt(x, y))
		draw.Draw(packImg, packRect, img, img.Rect.Min, draw.Src)
		chars = append(chars, jsGlyph{
			Code: r,
			X:    x,
			Y:    y,
			W:    packRect.Dx(),
			H:    packRect.Dy(),
			XOff: offset.X,
			YOff: offset.Y,
			XAdv: adv,
		})
		y += img.Rect.Dy()
	}

	return &jsFont{
		Name:       name,
		LineHeight: lineHeight,
		Chars:      chars,
		img:        packImg,
		fit:        fit,
	}
}

func (jf *jsFont) getGlyph(r rune) jsGlyph {
	for _, gl := range jf.Chars {
		if r == gl.Code {
			return gl
		}
	}
	return jsGlyph{}
}

func (jf *jsFont) doesntFit() bool {
	return !jf.fit
}

func (jf *jsFont) save(outDir string) error {
	// create out dir if not exist
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, os.ModePerm)
	}

	// JSON save
	jsonBytes, _ := json.Marshal(jf)
	jsonstr := strings.ReplaceAll(string(jsonBytes), "{", "\n{")
	jsonBytes = []byte(jsonstr)
	fileName := filepath.Join(outDir, jf.Name+".json")
	err := ioutil.WriteFile(fileName, jsonBytes, os.ModePerm)
	if err != nil {
		fmt.Println("json write error: ", err)
		return err
	}

	// PNG save
	fileName = filepath.Join(outDir, jf.Name+".png")
	pngFile, err := os.Create(fileName)
	if err == nil {
		err = png.Encode(pngFile, jf.img)
	}
	fmt.Println(err)
	return err
}

package util

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"

	"golang.org/x/image/font"

	"github.com/enotofil/PJFont/effects"
)

// TODO multiple ttf folders, sort fonts in alphabet order
const (
	//WindowW : app screen width
	WindowW = 1280
	//WindowH : app screen height
	WindowH    = 720
	paramW     = 260 //384
	paramH     = 28
	previewW   = 512
	previewH   = 512
	uiFontSize = 14

	outDir       = "out"
	bgSquareSize = 32  // pixels
	sizeStep     = 256 // pixels
	uiFontName   = "Go Mono (built-in)"
)

var (
	sampleStrings = []string{}
	helpText      = "ARROWS: select and change  |  HOME/END: min/max value  |  " +
		"ENTER: apply  |  C: change background  |  S: save "
	statusText         = ""
	selectColor        = color.NRGBA{0x50, 0x60, 0x80, 0xff}
	yellowColor        = color.NRGBA{0xFF, 0xFF, 0x55, 0xFF}
	bgColor      uint8 = 0x7F
	previewScale       = 1.0
)

// FontUI change parameters and display preview
type FontUI struct {
	// fonts     []*truetype.Font
	params    []fontParam
	curParam  int
	uiFont    font.Face
	font      *jsFont
	packedImg *ebiten.Image
	bgImg     *ebiten.Image
	paramImg  *ebiten.Image
	config    *Config
	loader    *fontLoader
}

// NewFontUI creates UI
func NewFontUI() *FontUI {
	config := LoadConfig()
	loader := newFontLoader(config.TTFontsPath)
	uiFont := loader.getFace(0, uiFontSize)
	params := []fontParam{
		// name, value, min, max, step
		fontParam{"font №", 0, 0, len(loader.fontNames) - 1, 1},            // 0
		fontParam{"font size, px", 24, 8, 56, 1},                           // 1
		fontParam{"top gradient, %", 100, 50, 100, 5},                      // 2
		fontParam{"mid gradient, %", 100, 50, 100, 5},                      // 3
		fontParam{"bottom gradient, %", 70, 50, 100, 5},                    // 4
		fontParam{"outline width, px", 1, 0, 3, 1},                         // 5
		fontParam{"x shadow dist, px", 1, 0, 4, 1},                         // 6
		fontParam{"y shadow dist, px", 1, 0, 4, 1},                         // 7
		fontParam{"shadow blur", 2, 0, 4, 1},                               // 8
		fontParam{"add spacing, px", 1, 0, 10, 1},                          // 9
		fontParam{"add line height, px", 1, 0, 10, 1},                      // 10
		fontParam{"bitmap width, px", previewW, sizeStep, 1024, sizeStep},  // 11
		fontParam{"bitmap height, px", previewH, sizeStep, 1024, sizeStep}, // 12
	}

	paramImg, _ := ebiten.NewImage(paramW, paramH, ebiten.FilterDefault)
	paramImg.Fill(selectColor)

	return &FontUI{
		params:    params,
		curParam:  0,
		uiFont:    uiFont,
		font:      nil,
		packedImg: nil,
		paramImg:  paramImg,
		bgImg:     nil,
		config:    config,
		loader:    loader,
	}
}

// Update input and draw
func (ui *FontUI) Update(screen *ebiten.Image) bool {
	if ui.packedImg == nil {
		ui.submit()
	}

	// select parameters
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if ui.curParam == 0 {
			ui.curParam = len(ui.params)
		}
		ui.curParam--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ui.curParam = (ui.curParam + 1) % len(ui.params)
	}

	// change values
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		ui.params[ui.curParam].change(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		ui.params[ui.curParam].change(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		ui.params[ui.curParam].change(9999)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		ui.params[ui.curParam].change(-9999)
	}

	// submit, generate font
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ui.submit()
	}

	// change preview background
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		bgColor = uint8((int(bgColor) + 0x20) % 0x100)
	}

	// save font JSON & PNG
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if ui.font.save(ui.config.OutFontsPath) == nil {
			statusText = "FONT SAVED"
		} else {
			statusText = "SAVE FAILED"
		}
	}

	// Draw packed preview background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(previewScale*bgSquareSize, previewScale*bgSquareSize)
	op.GeoM.Translate(paramW, 0.0)
	screen.DrawImage(ui.bgImg, op)

	// Draw packed font preview
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(previewScale, previewScale)
	op.GeoM.Translate(paramW, 0.0)
	screen.DrawImage(ui.packedImg, op)

	// Draw parameters, highlight selected
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0.0, paramH*float64(ui.curParam))
	screen.DrawImage(ui.paramImg, op)
	for i, p := range ui.params {
		y := paramH * (i + 1)
		if i == ui.curParam {
			text.Draw(screen, p.toString(), ui.uiFont, 1, y-8, yellowColor)
		} else {
			text.Draw(screen, p.toString(), ui.uiFont, 0, y-9, color.White)
		}
	}

	// Draw /fonts.ttf content, highlight current
	fontsOnPage := previewH / paramH
	page := ui.params[0].value / fontsOnPage
	posOnPage := ui.params[0].value % fontsOnPage

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 1)
	op.GeoM.Translate(paramW+previewW, float64(paramH*posOnPage))
	screen.DrawImage(ui.paramImg, op)

	y := paramH
	fontN := page * fontsOnPage
	for y < previewH && fontN < len(ui.loader.fontNames) {
		f := ui.loader.fontNames[fontN]
		text.Draw(screen, fmt.Sprintf("%2d. %s", fontN, f), ui.uiFont, paramW+previewW+10, y-9, color.White)
		y += paramH
		fontN++
	}

	// Draw x1 sample string and background
	ebitenutil.DrawRect(
		screen,
		0.0,
		previewH+1,
		float64(paramW+previewW)-1,
		float64(WindowH-previewH-uiFontSize*2),
		color.RGBA{bgColor, bgColor, bgColor, 0xFF},
	)
	for i, s := range ui.config.SampleStrings {
		ui.drawString(screen, s, 8, previewH+ui.font.LineHeight*(i+1), 1)
	}

	// Draw x2 sample string and background
	ebitenutil.DrawRect(
		screen,
		float64(paramW+previewW),
		previewH+1,
		float64(WindowW-paramW-previewW),
		float64(WindowH-previewH-uiFontSize*2),
		color.RGBA{bgColor, bgColor, bgColor, 0xFF},
	)
	ui.drawString(screen, ui.font.Name, paramW+previewW+8, previewH+ui.font.LineHeight*2+8, 2)

	// Draw help string
	text.Draw(screen, helpText+ui.font.Name+".*", ui.uiFont, 16, WindowH-uiFontSize/2, color.White)
	// Draw saved indicator
	text.Draw(screen, statusText, ui.uiFont, WindowW-120, WindowH-uiFontSize/2, yellowColor)

	return false
}

// read parameters, generate font and prepare UI
func (ui *FontUI) submit() {
	fxList := []effects.Effect{
		&effects.Gradient{
			Top:    ui.params[2].value,
			Mid:    ui.params[3].value,
			Bottom: ui.params[4].value,
		},
		&effects.Outline{
			Width: ui.params[5].value,
		},
		&effects.Shadow{
			Dist: image.Pt(ui.params[6].value, ui.params[7].value),
			Blur: ui.params[8].value,
		},
		&effects.Spacing{
			AdvPlus: ui.params[9].value,
		},
	}

	face := ui.loader.getFace(ui.params[0].value, ui.params[1].value)
	if face == nil {
		statusText = "LOAD ERROR"
		return
	}

	lineHeight := face.Metrics().Height.Round() + ui.params[10].value

	// no path, no ext, name_X where X = font size
	name := filepath.Base(ui.loader.fontNames[ui.params[0].value])
	name = fmt.Sprintf(
		"%s_%d",
		strings.ReplaceAll(name[:len(name)-4], " ", ""),
		ui.params[1].value,
	)

	ui.font = newJsFont(
		name,
		face,
		ui.config.CollectRunes(),
		fxList,
		lineHeight,
		image.Pt(ui.params[11].value, ui.params[12].value),
	)

	ui.packedImg, _ = ebiten.NewImageFromImage(ui.font.img, ebiten.FilterDefault)
	imgW := ui.packedImg.Bounds().Dx()
	imgH := ui.packedImg.Bounds().Dy()
	previewScale = math.Min(
		float64(previewW)/float64(imgW),
		float64(previewH)/float64(imgH),
	)
	ui.bgImg, _ = ebiten.NewImage(imgW/bgSquareSize, imgH/bgSquareSize, ebiten.FilterDefault)
	if ui.font.doesntFit() {
		ui.bgImg.Fill(color.NRGBA{192, 64, 64, 255})
	} else {
		for x := 0; x < imgW/bgSquareSize; x++ {
			for y := 0; y < imgH/bgSquareSize; y++ {
				m := uint8((x+y)%2)*64 + 127
				ui.bgImg.Set(x, y, color.RGBA{m, m, m, 255})
			}
		}
	}

	statusText = "        *"
}

// draws string with current jsFont for preview. THIS IS SLOW!
// use map[rune]glyph and predefined src rectangles for speed
func (ui *FontUI) drawString(screen *ebiten.Image, s string, x, y, scale int) {
	for _, r := range []rune(s) {
		gl := ui.font.getGlyph(r)
		src := image.Rect(gl.X, gl.Y, gl.X+gl.W, gl.Y+gl.H)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(scale), float64(scale))
		op.GeoM.Translate(float64(x+gl.XOff*scale), float64(y+gl.YOff*scale))
		screen.DrawImage(ui.packedImg.SubImage(src).(*ebiten.Image), op)
		x += gl.XAdv * scale
	}
}

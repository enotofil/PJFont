package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten"

	"github.com/enotofil/PJFont/util"
)

type Game struct{}
var (
	ui *util.FontUI
)

func init() {
	ui = util.NewFontUI()
}


func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	screen.Fill(color.NRGBA{0x2F, 0x3F, 0x5F, 0xff})

	ui.Update(screen)

	return
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return util.WindowW, util.WindowH
}
func main() {
		
	ebiten.SetWindowSize(util.WindowW, util.WindowH)
	ebiten.SetWindowTitle("PJFont")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

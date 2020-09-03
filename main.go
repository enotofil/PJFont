package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten"

	"github.com/enotofil/PJFont/util"
)

// Game implements ebiten.Game interface.
type Game struct{}

var (
	ui *util.FontUI
)

func init() {
	ui = util.NewFontUI()
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0x2F, 0x3F, 0x5F, 0xff})

	ui.Update(screen)

	return
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	ui.Resize(outsideWidth, outsideHeight)
	return util.WindowWidth, util.WindowHeight
}

func main() {

	ebiten.SetWindowSize(util.WindowWidth, util.WindowHeight)
	ebiten.SetWindowTitle("PJFont")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

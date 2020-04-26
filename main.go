package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten"

	"github.com/enotofil/PJFont/util"
)

var (
	ui *util.FontUI
)

func init() {
	ui = util.NewFontUI()
}

func update(screen *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.Fill(color.NRGBA{0x2F, 0x3F, 0x5F, 0xff})

	ui.Update(screen)

	return nil
}

func main() {
	// ebiten rules
	if err := ebiten.Run(update, util.WindowW, util.WindowH, 1, "PJFont"); err != nil {
		log.Fatal(err)
	}
}

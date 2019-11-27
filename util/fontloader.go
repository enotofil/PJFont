package util

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/image/font"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomono"
)

type fontLoader struct {
	fontNames []string
}

func newFontLoader(ttfDir string) *fontLoader {
	fontNames := []string{}
	fontNames = append(fontNames, "Go Mono.ttf")
	fontNames = append(fontNames, getFontNames(ttfDir)...)

	fmt.Printf("%d fonts found in %s\n", len(fontNames)-1, ttfDir)
	return &fontLoader{
		fontNames: fontNames,
	}
}

func getFontNames(ttfDir string) []string {
	fontNames := []string{}
	ttfFiles, err := ioutil.ReadDir(ttfDir)
	if err != nil {
		fmt.Printf("Can't read %s dir: %s\n", ttfDir, err)
		return fontNames
	}
	for _, f := range ttfFiles {
		if f.IsDir() {
			fontNames = append(fontNames, getFontNames(filepath.Join(ttfDir, f.Name()))...)
		} else if strings.HasSuffix(strings.ToLower(f.Name()), ".ttf") {
			fontNames = append(fontNames, filepath.Join(ttfDir, f.Name()))
		}
	}
	return fontNames
}

func (fl *fontLoader) getFace(fontN, size int) font.Face {
	var ttfont *truetype.Font
	if fontN == 0 {
		ttfont, _ = truetype.Parse(gomono.TTF)
	} else {
		fontBytes, err := ioutil.ReadFile(fl.fontNames[fontN])
		if err != nil {
			fmt.Printf("Can't read file '%s':\n%s\n", fl.fontNames[fontN], err)
			return nil
		}
		ttfont, err = truetype.Parse(fontBytes)
		if err != nil {
			fmt.Printf("Can't parse ttf '%s':\n%s\n", fl.fontNames[fontN], err)
			return nil
		}
	}
	face := truetype.NewFace(ttfont, &truetype.Options{
		Size:    float64(size),
		Hinting: font.HintingFull,
	})
	return face
}

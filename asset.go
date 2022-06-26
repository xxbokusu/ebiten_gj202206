package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var FS embed.FS

var (
	images = map[string]*ebiten.Image{}
)

func loadImage(destName string, sourceName string) error {
	byteData, err := FS.ReadFile(sourceName)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(byteData))
	if err != nil {
		return err
	}

	images[destName] = ebiten.NewImageFromImage(img)
	return nil
}

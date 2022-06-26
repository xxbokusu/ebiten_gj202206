package main

import (
	"bytes"
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

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

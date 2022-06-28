package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type TitleScene struct {
	big_font    font.Face
	normal_font font.Face
}

func (s *TitleScene) Update(manager SceneTransitionManager) error {
	if len(inpututil.AppendPressedKeys(nil)) >= 2 {
		manager.SceneTransition(&PlayScene{})
	}
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {

	title_string := "MAGNET GO!"
	r := text.BoundString(s.big_font, title_string)
	x := screenX/2 - r.Dx()/2
	y := screenY/2 + r.Dy()/2
	text.Draw(screen, title_string, s.big_font, x, y, color.White)

	title_string = "PLEASE TAKE ANOTHER PLAYER AND PRESS ANY 2 KEY"
	r = text.BoundString(s.normal_font, title_string)
	x = screenX/2 - r.Dx()/2
	y = screenY/2 + 100 + r.Dy()/2
	text.Draw(screen, title_string, s.normal_font, x, y, color.White)
}

func (s *TitleScene) init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	s.normal_font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.big_font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    72,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

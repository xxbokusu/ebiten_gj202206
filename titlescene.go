package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type TitleScene struct {
}

func (s *TitleScene) Update(manager SceneTransitionManager) error {
	if len(inpututil.AppendPressedKeys(nil)) >= 2 {
		manager.SceneTransition(&PlayScene{})
	}
	return nil
}

func (s *TitleScene) Draw(screen *ebiten.Image) {

	title_string := "MAGNET GO!"
	r := text.BoundString(big_font, title_string)
	x := screenX/2 - r.Dx()/2
	y := screenY/2 + r.Dy()/2
	text.Draw(screen, title_string, big_font, x, y, color.White)

	title_string = "PLEASE TAKE ANOTHER PLAYER AND PRESS ANY 2 KEY"
	r = text.BoundString(normal_font, title_string)
	x = screenX/2 - r.Dx()/2
	y = screenY/2 + 100 + r.Dy()/2
	text.Draw(screen, title_string, normal_font, x, y, color.White)
}

func (s *TitleScene) init() {
	err := loadFont()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	debug   = false
	screenX = 800
	screenY = 600

	boardX = 9
	boardY = 9
)

var ()

type SceneTransitionManager interface {
	SceneTransition(scene Scene)
}

type Scene interface {
	Update(manager SceneTransitionManager) error
	Draw(screen *ebiten.Image)
	init()
}

type Game struct {
	now_scene  Scene
	next_scene Scene
}

func (g *Game) Update() error {
	if g.next_scene != nil {
		g.now_scene = g.next_scene
		g.next_scene = nil
	}
	if err := g.now_scene.Update(g); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.now_scene.Draw(screen)
}

func (g *Game) SceneTransition(scene Scene) {
	g.next_scene = scene
	g.next_scene.init()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenX, screenY
}

func (g *Game) init() {
	imageSourceMap := map[string]string{
		"white_n":       "assets/images/go_white_n.png",
		"black_n":       "assets/images/go_black_n.png",
		"frame_white_n": "assets/images/go_frame_white_n.png",
		"frame_black_n": "assets/images/go_frame_black_n.png",
		"white_s":       "assets/images/go_white_s.png",
		"black_s":       "assets/images/go_black_s.png",
		"frame_white_s": "assets/images/go_frame_white_s.png",
		"frame_black_s": "assets/images/go_frame_black_s.png",
	}
	for key, value := range imageSourceMap {
		if err := loadImage(key, value); err != nil {
			log.Fatal(err)
		}

	}

	audioSourceMap := map[string]string{
		"set_stone":   "assets/se/set_stone.mp3",
		"force_stone": "assets/se/force_stone.mp3",
	}
	for key, value := range audioSourceMap {
		if err := loadAudio(key, value); err != nil {
			log.Fatal(err)
		}

	}
	playAudio("set_stone")
}

// NewGame method
func NewGame() *Game {
	g := &Game{}
	g.SceneTransition(&TitleScene{})
	g.init()
	return g
}

func main() {
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("Magnet Go!")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

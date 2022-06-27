package main

import (
	"log"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	debug   = false
	screenX = 800
	screenY = 600

	boardX = 9
	boardY = 9
)

var ()

type board_panel struct {
	state int
}

type Game struct {
	board [boardX][boardY]board_panel

	panelSpan      int
	outboardSpaceX int
	outboardSpaceY int

	isBlackTurn bool

	canPlayAudio bool
}

func (g *Game) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}

func (g *Game) Update() error {
	if !g.canPlayAudio {
		select {
		case <-playAudioCompleteCh:
			g.canPlayAudio = true
		default:
		}
		return nil
	}

	if g.isKeyJustPressed() {
		cursorX, cursorY := ebiten.CursorPosition()
		selectedPosX := (cursorX - g.outboardSpaceX) / g.panelSpan
		selectedPosY := (cursorY - g.outboardSpaceY) / g.panelSpan
		g.UpdateBoardState(selectedPosX, selectedPosY)
	}
	return nil
}
func (g *Game) UpdateBoardState(posX, posY int) error {
	if posX >= boardX || posY >= boardY {
		return nil
	}

	if g.board[posX][posY].state != 0 {
		return nil
	}

	if g.isBlackTurn {
		g.board[posX][posY].state = 1
	} else {
		g.board[posX][posY].state = 2
	}
	g.isBlackTurn = !g.isBlackTurn

	g.canPlayAudio = false
	playAudio("set_stone")

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for i := 0; i < boardX; i = i + 1 {
		ebitenutil.DrawLine(screen,
			float64(g.outboardSpaceX+i*g.panelSpan),
			0,
			float64(g.outboardSpaceX+i*g.panelSpan),
			screenY,
			color.White)
		ebitenutil.DrawLine(screen,
			0,
			float64(g.outboardSpaceY+i*g.panelSpan),
			screenX,
			float64(g.outboardSpaceY+i*g.panelSpan),
			color.White)
	}

	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			if g.board[i][j].state == 2 {
				g.DrawStone(screen, "white", i, j)
			} else if g.board[i][j].state == 1 {
				g.DrawStone(screen, "black", i, j)
			}
		}
	}

	if !g.isKeyJustPressed() {
		cursorX, cursorY := ebiten.CursorPosition()
		selectedPosX := (cursorX - g.outboardSpaceX) / g.panelSpan
		selectedPosY := (cursorY - g.outboardSpaceY) / g.panelSpan
		g.DrawStone(screen, "frame", selectedPosX, selectedPosY)
	}
}

func (g *Game) DrawStone(screen *ebiten.Image, name string, posX, posY int) error {
	if posX >= boardX || posY >= boardY {
		return nil
	}

	img_opt := &ebiten.DrawImageOptions{}
	img := images[name]
	img_width, img_height := img.Size()
	img_opt.GeoM.Scale(
		float64(g.panelSpan)/float64(img_width),
		float64(g.panelSpan)/float64(img_height))
	img_opt.GeoM.Translate(
		float64(g.outboardSpaceX+posX*g.panelSpan-g.panelSpan/2),
		float64(g.outboardSpaceY+posY*g.panelSpan-g.panelSpan/2))
	screen.DrawImage(img, img_opt)
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenX, screenY
}

func (g *Game) init() {
	g.outboardSpaceX = screenX / (boardX + 2) / 2
	g.outboardSpaceY = screenY / (boardY + 2) / 2
	g.panelSpan = g.outboardSpaceY * 2
	// board initialize
	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			g.board[i][j].state = 0
		}
	}

	imageSourceMap := map[string]string{
		"white": "assets/images/go_white.png",
		"black": "assets/images/go_black.png",
		"frame": "assets/images/go_frame_yellow.png",
	}
	for key, value := range imageSourceMap {
		if err := loadImage(key, value); err != nil {
			log.Fatal(err)
		}

	}

	audioSourceMap := map[string]string{
		"set_stone": "assets/se/set_stone.mp3",
	}
	for key, value := range audioSourceMap {
		if err := loadAudio(key, value); err != nil {
			log.Fatal(err)
		}

	}
	playAudio("set_stone")

	g.canPlayAudio = true
	g.isBlackTurn = true
}

// NewGame method
func NewGame() *Game {
	g := &Game{}
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

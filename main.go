package main

import (
	"image/color"
	"log"

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

type board_panel struct {
	state int
}

type Game struct {
	board [boardX][boardY]board_panel
}

func (g *Game) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.board[1][1].state = 1
		return true
	}
	return false
}

func (g *Game) Update() error {
	if g.isKeyJustPressed() {
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			if g.board[i][j].state == 0 {
				ebitenutil.DebugPrintAt(screen, "+", i*10, j*10)
			} else {
				ebitenutil.DebugPrintAt(screen, "O", i*10, j*10)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *Game) init() {
	// board initialize
	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			g.board[i][j].state = 0
		}
	}
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

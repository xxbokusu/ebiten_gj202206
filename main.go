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

	panelSpan      int
	outboardSpaceX int
	outboardSpaceY int
}

func (g *Game) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}

func (g *Game) Update() error {
	if g.isKeyJustPressed() {
		cursorX, cursorY := ebiten.CursorPosition()
		selectedPosX := (cursorX - g.outboardSpaceX) / g.panelSpan
		selectedPosY := (cursorY - g.outboardSpaceY) / g.panelSpan
		g.board[selectedPosX][selectedPosY].state = 1

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			if g.board[i][j].state != 0 {
				ebitenutil.DebugPrintAt(screen, "O", g.outboardSpaceX+i*g.panelSpan, g.outboardSpaceY+j*g.panelSpan)
			}
		}
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

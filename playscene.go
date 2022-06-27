package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type board_panel struct {
	state int
}

type PlayScene struct {
	board [boardX][boardY]board_panel

	panelSpan      int
	outboardSpaceX int
	outboardSpaceY int

	isBlackTurn bool

	canPlayAudio bool
}

func (s *PlayScene) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}

func (s *PlayScene) Update(manager SceneTransitionManager) error {
	if !s.canPlayAudio {
		select {
		case <-playAudioCompleteCh:
			s.canPlayAudio = true
		default:
		}
		return nil
	}

	if s.isKeyJustPressed() {
		cursorX, cursorY := ebiten.CursorPosition()
		selectedPosX := (cursorX - s.outboardSpaceX) / s.panelSpan
		selectedPosY := (cursorY - s.outboardSpaceY) / s.panelSpan
		s.UpdateBoardState(selectedPosX, selectedPosY)
	}
	return nil
}

func (s *PlayScene) UpdateBoardState(posX, posY int) error {
	if posX >= boardX || posY >= boardY {
		return nil
	}

	if s.board[posX][posY].state != 0 {
		return nil
	}

	if s.isBlackTurn {
		s.board[posX][posY].state = 1
	} else {
		s.board[posX][posY].state = 2
	}
	s.isBlackTurn = !s.isBlackTurn

	s.canPlayAudio = false
	playAudio("set_stone")

	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for i := 0; i < boardX; i = i + 1 {
		ebitenutil.DrawLine(screen,
			float64(s.outboardSpaceX+i*s.panelSpan),
			0,
			float64(s.outboardSpaceX+i*s.panelSpan),
			screenY,
			color.White)
		ebitenutil.DrawLine(screen,
			0,
			float64(s.outboardSpaceY+i*s.panelSpan),
			screenX,
			float64(s.outboardSpaceY+i*s.panelSpan),
			color.White)
	}

	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			if s.board[i][j].state == 2 {
				s.DrawStone(screen, "white", i, j)
			} else if s.board[i][j].state == 1 {
				s.DrawStone(screen, "black", i, j)
			}
		}
	}

	if !s.isKeyJustPressed() {
		cursorX, cursorY := ebiten.CursorPosition()
		selectedPosX := (cursorX - s.outboardSpaceX) / s.panelSpan
		selectedPosY := (cursorY - s.outboardSpaceY) / s.panelSpan
		if s.isBlackTurn {
			s.DrawStone(screen, "frame_black", selectedPosX, selectedPosY)
		} else {
			s.DrawStone(screen, "frame_white", selectedPosX, selectedPosY)
		}
	}
}

func (s *PlayScene) DrawStone(screen *ebiten.Image, name string, posX, posY int) error {
	if posX >= boardX || posY >= boardY {
		return nil
	}

	img_opt := &ebiten.DrawImageOptions{}
	img := images[name]
	img_width, img_height := img.Size()
	img_opt.GeoM.Scale(
		float64(s.panelSpan)/float64(img_width),
		float64(s.panelSpan)/float64(img_height))
	img_opt.GeoM.Translate(
		float64(s.outboardSpaceX+posX*s.panelSpan-s.panelSpan/2),
		float64(s.outboardSpaceY+posY*s.panelSpan-s.panelSpan/2))
	screen.DrawImage(img, img_opt)
	return nil
}

func (s *PlayScene) init() {
	s.outboardSpaceX = screenX / (boardX + 2) / 2
	s.outboardSpaceY = screenY / (boardY + 2) / 2
	s.panelSpan = s.outboardSpaceY * 2
	// board initialize
	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			s.board[i][j].state = 0
		}
	}

	s.canPlayAudio = true
	s.isBlackTurn = true
}

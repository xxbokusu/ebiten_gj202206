package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GoStone struct {
	isBlack bool
	isNorth bool
}

func (gs *GoStone) getStoneImgString() string {
	var str string
	if gs.isBlack {
		str += "black"
	} else {
		str += "white"
	}

	if gs.isNorth {
		str += "_n"
	} else {
		str += "_s"
	}
	return str
}

type BoardPanel struct {
	stone *GoStone
}

type PlayScene struct {
	board [boardX][boardY]BoardPanel

	panelSpan      int
	outboardSpaceX int
	outboardSpaceY int

	isBlackTurn bool
	isNorth     bool

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

	if s.board[posX][posY].stone != nil {
		return nil
	}

	stone := GoStone{
		isBlack: s.isBlackTurn,
		isNorth: s.isNorth,
	}
	s.board[posX][posY].stone = &stone
	s.isBlackTurn = !s.isBlackTurn
	if s.isBlackTurn {
		s.isNorth = !s.isNorth
	}

	s.canPlayAudio = false
	playAudio("set_stone")

	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for i := 0; i < boardX; i = i + 1 {
		ebitenutil.DrawLine(screen,
			float64(s.outboardSpaceX+i*s.panelSpan),
			float64(s.outboardSpaceY),
			float64(s.outboardSpaceX+i*s.panelSpan),
			screenY-float64(s.outboardSpaceY),
			color.White)
		ebitenutil.DrawLine(screen,
			float64(s.outboardSpaceX),
			float64(s.outboardSpaceY+i*s.panelSpan),
			screenX-float64(s.outboardSpaceX),
			float64(s.outboardSpaceY+i*s.panelSpan),
			color.White)
	}

	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			if s.board[i][j].stone != nil {
				s.DrawStone(screen, s.board[i][j].stone.getStoneImgString(), i, j)
			}
		}
	}

	if !s.isKeyJustPressed() {
		cursorX, cursorY := ebiten.CursorPosition()
		selectedPosX := (cursorX - s.outboardSpaceX) / s.panelSpan
		selectedPosY := (cursorY - s.outboardSpaceY) / s.panelSpan

		setting_stone := &GoStone{
			isBlack: s.isBlackTurn,
			isNorth: s.isNorth,
		}
		s.DrawStone(screen, "frame_"+setting_stone.getStoneImgString(), selectedPosX, selectedPosY)
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
			s.board[i][j].stone = nil
		}
	}

	s.canPlayAudio = true
	s.isBlackTurn = true
	s.isNorth = true
}

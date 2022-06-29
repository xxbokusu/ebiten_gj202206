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

	accelX float64
	accelY float64
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
		s.SetGoStone(selectedPosX, selectedPosY)
	}
	return nil
}

func (s *PlayScene) SetGoStone(posX, posY int) error {
	if posX >= boardX || posY >= boardY {
		return nil
	}

	if s.board[posX][posY].stone != nil {
		return nil
	}

	stone := GoStone{
		isBlack: s.isBlackTurn,
		isNorth: s.isNorth,
		accelX:  0,
		accelY:  0,
	}
	s.board[posX][posY].stone = &stone
	s.isBlackTurn = !s.isBlackTurn
	if s.isBlackTurn {
		s.isNorth = !s.isNorth
	}

	s.canPlayAudio = false
	playAudio("set_stone")

	s.makeMagneticForce(posX, posY)

	return nil
}

func (s *PlayScene) makeMagneticForce(srcX, srcY int) {
	src_stone := s.board[srcX][srcY].stone
	if srcX-2 >= 0 { // 移動先スペースが必要
		if dest_stone := s.board[srcX-1][srcY].stone; dest_stone != nil {
			if s.board[srcX-2][srcY].stone == nil {
				if dest_stone.isNorth == src_stone.isNorth {
					dest_stone.accelX += -1 // 反発力
					s.board[srcX-2][srcY].stone = dest_stone
					s.board[srcX-1][srcY].stone = nil
					s.makeMagneticForce(srcX-2, srcY)
				}
			}
		} else if dest_stone := s.board[srcX-2][srcY].stone; dest_stone != nil {
			if dest_stone.isNorth != src_stone.isNorth {
				dest_stone.accelX += 1 // 引力
				s.board[srcX-1][srcY].stone = dest_stone
				s.board[srcX-2][srcY].stone = nil
				s.makeMagneticForce(srcX-1, srcY)
			}
		}
	}
	if srcY-2 >= 0 {
		if dest_stone := s.board[srcX][srcY-1].stone; dest_stone != nil {
			if s.board[srcX][srcY-2].stone == nil {
				if dest_stone.isNorth == src_stone.isNorth {
					dest_stone.accelY += -1 // 反発力
					s.board[srcX][srcY-2].stone = dest_stone
					s.board[srcX][srcY-1].stone = nil
					s.makeMagneticForce(srcX, srcY-2)
				}
			}
		} else if dest_stone := s.board[srcX][srcY-2].stone; dest_stone != nil {
			if dest_stone.isNorth != src_stone.isNorth {
				dest_stone.accelY += 1 // 引力
				s.board[srcX][srcY-1].stone = dest_stone
				s.board[srcX][srcY-2].stone = nil
				s.makeMagneticForce(srcX, srcY-1)
			}
		}
	}
	if srcX+2 < boardX {
		if dest_stone := s.board[srcX+1][srcY].stone; dest_stone != nil {
			if s.board[srcX+2][srcY].stone == nil {
				if dest_stone.isNorth == src_stone.isNorth {
					dest_stone.accelX += 1 // 反発力
					s.board[srcX+2][srcY].stone = dest_stone
					s.board[srcX+1][srcY].stone = nil
				}
			}
		} else if dest_stone := s.board[srcX+2][srcY].stone; dest_stone != nil {
			if dest_stone.isNorth != src_stone.isNorth {
				dest_stone.accelX += -1 // 引力
				s.board[srcX+1][srcY].stone = dest_stone
				s.board[srcX+2][srcY].stone = nil
			}
		}
	}
	if srcY+2 < boardY {
		if dest_stone := s.board[srcX][srcY+1].stone; dest_stone != nil {
			if s.board[srcX][srcY+2].stone == nil {
				if dest_stone.isNorth == src_stone.isNorth {
					dest_stone.accelY += 1 // 反発力
					s.board[srcX][srcY+2].stone = dest_stone
					s.board[srcX][srcY+1].stone = nil
				}
			}
		} else if dest_stone := s.board[srcX][srcY+2].stone; dest_stone != nil {
			if dest_stone.isNorth != src_stone.isNorth {
				dest_stone.accelY += -1 // 引力
				s.board[srcX][srcY+1].stone = dest_stone
				s.board[srcX][srcY+2].stone = nil
			}
		}
	}
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for i := 0; i < boardX; i = i + 1 {
		ebitenutil.DrawLine(screen,
			float64(s.outboardSpaceX+i*s.panelSpan),
			float64(s.outboardSpaceY),
			float64(s.outboardSpaceX+i*s.panelSpan),
			screenY-float64(s.outboardSpaceY)*2,
			color.White)
		ebitenutil.DrawLine(screen,
			float64(s.outboardSpaceX),
			float64(s.outboardSpaceY+i*s.panelSpan),
			screenX-float64(s.outboardSpaceX)-float64(s.panelSpan),
			float64(s.outboardSpaceY+i*s.panelSpan),
			color.White)
	}

	for i := 0; i < boardX; i = i + 1 {
		for j := 0; j < boardY; j = j + 1 {
			if s.board[i][j].stone != nil {
				s.board[i][j].stone.accelX *= 0.8
				s.board[i][j].stone.accelY *= 0.8
				s.DrawStone(screen,
					s.board[i][j].stone.getStoneImgString(),
					float64(i)-s.board[i][j].stone.accelX,
					float64(j)-s.board[i][j].stone.accelY,
				)
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
		s.DrawStone(screen, "frame_"+setting_stone.getStoneImgString(), float64(selectedPosX), float64(selectedPosY))
	}
}

func (s *PlayScene) DrawStone(screen *ebiten.Image, name string, posX, posY float64) error {
	if posX >= boardX || posY >= boardY {
		return nil
	}

	img_opt := &ebiten.DrawImageOptions{}
	img := images[name]
	img_width, img_height := img.Size()
	img_opt.GeoM.Scale(
		float64(s.panelSpan)/float64(img_width),
		float64(s.panelSpan)/float64(img_height),
	)
	img_opt.GeoM.Translate(
		float64(s.outboardSpaceX-s.panelSpan/2)+posX*float64(s.panelSpan),
		float64(s.outboardSpaceY-s.panelSpan/2)+posY*float64(s.panelSpan),
	)
	screen.DrawImage(img, img_opt)
	return nil
}

func (s *PlayScene) init() {
	s.outboardSpaceY = screenY / (boardY + 2)
	s.panelSpan = s.outboardSpaceY
	s.outboardSpaceX = (screenX - s.panelSpan*boardX) / 2
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

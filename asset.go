package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const (
	sampleRate = 48000
)

//go:embed assets
var FS embed.FS

var (
	images = map[string]*ebiten.Image{}

	audioContext        = audio.NewContext(sampleRate)
	audioBytes          = map[string][]byte{}
	playAudioCompleteCh chan struct{}
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

func loadAudio(destName, sourceName string) error {
	byteData, err := FS.ReadFile(sourceName)
	if err != nil {
		log.Fatal(err)
	}
	audioBytes[destName] = byteData
	return nil
}

func playAudio(name string) {
	go func() {
		playAudioCompleteCh = make(chan struct{})
		stream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(audioBytes[name]))
		if err != nil {
			log.Fatal(err)
		}

		audioPlayer, err := audioContext.NewPlayer(stream)
		if err != nil {
			log.Fatal(err)
		}
		audioPlayer.Play()
		close(playAudioCompleteCh)
	}()
}

package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth            = 640
	screenHeight           = 480
	maxLiveCellsPercentage = 0.1
)

var (
	catKiss *ebiten.Image
)

type World struct {
	width        int
	height       int
	area         []bool
	maxLiveCells int
}

func NewWorld(width int, height int) *World {
	w := &World{
		width:        width,
		height:       height,
		area:         make([]bool, width*height),
		maxLiveCells: int(math.Round(maxLiveCellsPercentage * float64(width) * float64(height))),
	}
	return w
}

func (w *World) Update() {

}

func (w *World) Draw() {

}

type Game struct {
	count int
}

func (g *Game) Update() error {
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	s := catKiss.Bounds().Size()
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-float64(s.X)/2, -float64(s.Y)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)

	screen.DrawImage(catKiss, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("kitties kissing!")

	catBytes, err := os.ReadFile("images/cat_love_30.jpg")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(catBytes))
	if err != nil {
		log.Fatal(err)
	}

	catKiss = ebiten.NewImageFromImage(img)

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

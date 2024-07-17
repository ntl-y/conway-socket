package main

import (
	_ "image/jpeg"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var color = []byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Float64())}

const (
	screenFrameWidth       = 640
	screenFrameHeight      = 480
	screenWidth            = 128
	screenHeight           = 96
	maxLiveCellsPercentage = 0.1
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

func (w *World) Draw(pix []byte) {
	for i := range w.area {
		randomBool := rand.Intn(2) == 1
		if randomBool {
			pix[4*i] = color[0]
			pix[4*i+1] = color[1]
			pix[4*i+2] = color[2]
			pix[4*i+3] = color[3]
		} else {
			pix[4*i] = 0
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0
		}
	}
}

type Game struct {
	world  *World
	pixels []byte
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
		g.world.Draw(g.pixels)

	}
	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenFrameWidth, screenFrameHeight)
	ebiten.SetWindowTitle("pixel")

	game := &Game{world: NewWorld(screenWidth, screenHeight)}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

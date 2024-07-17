package main

import (
	_ "image/jpeg"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	cellColor  = []byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Float64())}
	background = byte(0)
)

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

func indexOfPixel(x, y, screenWidth int) int {
	return (y*screenWidth + x) * 4
}

func (w *World) Update(pix []byte) {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			index := indexOfPixel(x, y, w.width)
			neighbourAmount := countNeighbours(pix, x, y, w.width, w.height)
			isLiving := status(pix, x, y, w.width)
			switch {
			case isLiving && (neighbourAmount == 2 || neighbourAmount == 3):
				w.paint(pix, index)
			case !isLiving && neighbourAmount == 3:
				w.paint(pix, index)
			default:
				pix[index] = background
				pix[index+1] = background
				pix[index+2] = background
				pix[index+3] = background

			}
		}
	}
}

func status(pix []byte, x int, y int, width int) bool {
	return pix[indexOfPixel(x, y, width)] == cellColor[0]
}

func countNeighbours(pix []byte, x int, y int, width int, height int) int {
	neighbours := 0

	directions := [][2]int{
		{0, 1},  // north
		{0, -1}, //south
		{-1, 0}, // west
		{1, 0},  //east
	}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if nx >= 0 && nx < width && ny >= 0 && ny < height {
			index := indexOfPixel(nx, ny, width)
			if pix[index] == cellColor[0] {
				neighbours++
			}
		}

	}
	return neighbours
}

func (w *World) Draw(pix []byte) {
	for i := range w.area {
		pix[4*i] = background
		pix[4*i+1] = background
		pix[4*i+2] = background
		pix[4*i+3] = background
	}
}

func (w *World) paint(pix []byte, pixelIndex int) {
	pix[pixelIndex] = cellColor[0]
	pix[pixelIndex+1] = cellColor[1]
	pix[pixelIndex+2] = cellColor[2]
	pix[pixelIndex+3] = cellColor[3]
}

type Game struct {
	world  *World
	pixels []byte
}

func (g *Game) paint() {
	mx, my := ebiten.CursorPosition()

	currPixel := indexOfPixel(mx, my, g.world.width)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.pixels[currPixel] == 0 {
		g.world.paint(g.pixels, currPixel)
	}
}

func (g *Game) Update() error {
	g.paint()
	g.world.Update(g.pixels)
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

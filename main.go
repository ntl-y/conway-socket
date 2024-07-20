package main

import (
	"log"
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	cellColor  = []byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Float64())}
	background = byte(0)
	pause      = true
)

const (
	screenFrameWidth       = 640
	screenFrameHeight      = 480
	screenWidth            = 64
	screenHeight           = 48
	maxLiveCellsPercentage = 0.5
)

func indexOfPixel(x, y, screenWidth int) int {
	return (y*screenWidth + x) * 4
}

func pixelToCoords(index, screenWidth int) (int, int) {
	pixelIndex := index / 4

	y := pixelIndex / screenWidth
	x := pixelIndex % screenWidth

	return x, y
}

func indexInArea(x, y, width int) int {
	return y*width + x
}

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

	/////
	// for i := 0; i < w.maxLiveCells; i++ {
	// 	x := rand.Intn(width)
	// 	y := rand.Intn(height)
	// 	w.area[indexInArea(x, y, width)] = true
	// }
	/////
	return w
}

func (w *World) Draw(pix []byte) {
	for i, v := range w.area {
		if v {
			w.paint(pix, i*4)
		} else {
			w.erase(pix, i*4)
		}
	}
}

func (w *World) paint(pix []byte, pixelIndex int) {
	if len(pix) > 0 {
		x, y := pixelToCoords(pixelIndex, w.width)
		w.area[indexInArea(x, y, w.width)] = true

		pix[pixelIndex] = cellColor[0]
		pix[pixelIndex+1] = cellColor[1]
		pix[pixelIndex+2] = cellColor[2]
		pix[pixelIndex+3] = cellColor[3]
	}
}

func (w *World) erase(pix []byte, pixelIndex int) {

	if len(pix) > 0 {
		x, y := pixelToCoords(pixelIndex, w.width)
		w.area[indexInArea(x, y, w.width)] = false

		pix[pixelIndex] = background
		pix[pixelIndex+1] = background
		pix[pixelIndex+2] = background
		pix[pixelIndex+3] = background
	}
}

func countNeighbours(area []bool, x int, y int, width int, height int) int {
	neighbours := 0

	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1},
	}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if nx >= 0 && nx < width && ny >= 0 && ny < height {
			index := indexInArea(nx, ny, width)
			if area[index] {
				neighbours++
			}
		}
	}
	return neighbours
}

func (w *World) Update(pix []byte) {
	newArea := make([]bool, w.width*w.height)

	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			index := indexInArea(x, y, w.width)
			neighbourAmount := countNeighbours(w.area, x, y, w.width, w.height)
			if w.area[index] {
				if neighbourAmount == 2 || neighbourAmount == 3 {
					newArea[index] = true
				} else {
					newArea[index] = false
				}
			} else {
				if neighbourAmount == 3 {
					newArea[index] = true
				} else {
					newArea[index] = false
				}
			}
		}
	}
	w.area = newArea
	w.Draw(pix)
}

type Game struct {
	world  *World
	pixels []byte
	mu     sync.Mutex
}

func (g *Game) paint() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		pause = !pause
	}
	mx, my := ebiten.CursorPosition()

	if mx >= 0 && mx < g.world.width && my >= 0 && my < g.world.height {
		currPixel := indexOfPixel(mx, my, g.world.width)
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.pixels[currPixel] == 0 {
			g.mu.Lock()
			g.world.paint(g.pixels, currPixel)
			g.mu.Unlock()
		}
	}
}

func (g *Game) Update() error {

	if !pause {
		g.mu.Lock()
		g.world.Update(g.pixels)
		g.mu.Unlock()
	}

	go g.paint()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}

	g.world.Draw(g.pixels)
	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetTPS(20)
	ebiten.SetWindowSize(screenFrameWidth, screenFrameHeight)
	ebiten.SetWindowTitle("Game of Life")

	game := &Game{world: NewWorld(screenWidth, screenHeight)}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

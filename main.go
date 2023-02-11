package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	WIDTH  = 640
	HEIGHT = 480
)

type Question struct {
	x, y float64
	soal string
	pg   []Answer
}

type Answer struct {
	jawab         string
	rectX, rectY  float64
	x, y          float64
	width, height float64
	correct       bool
}

type Game struct {
	mouse struct {
		x, y int
	}
	currentQuestion int
	choice          bool
	listSoal        []Question
}

func (g *Game) Update() error {
	g.mouse.x, g.mouse.y = ebiten.CursorPosition()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	switch currentScene {
	case MENU:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			currentScene = PLAY
			trackSoal = 0
			g.currentQuestion = 0
			correctCount = 0
		}
	case PLAY:
		for j := 0; j < len(g.listSoal[g.currentQuestion].pg); j++ {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && CheckPointCollision(g.mouse.x, g.mouse.y, g.listSoal[g.currentQuestion].pg[j]) {
				correct = g.listSoal[g.currentQuestion].pg[j].correct
				g.choice = true
				break
			}
		}

		if g.choice {
			if time.Since(lastUpdate) > time.Second && countDown > 0 {
				countDown -= 1
				lastUpdate = time.Now()
			}
		}

		if countDown == 0 {
			trackSoal++
			countDown = 2
			g.currentQuestion = (1 + g.currentQuestion) % len(g.listSoal)
			if correct {
				correctCount++
			}
			correct = false
			g.choice = false
		}

		if trackSoal == len(g.listSoal) {
			currentScene = GAMEOVER
		}

	case GAMEOVER:
		if ebiten.IsKeyPressed(ebiten.KeyR) && currentScene == GAMEOVER {
			currentScene = PLAY
			trackSoal = 0
			g.currentQuestion = 0
			correctCount = 0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// currentScene = GAMEOVER

	switch currentScene {
	case MENU:
		text.Draw(screen, "Quiz Go", f, WIDTH/2-60, HEIGHT/4, color.White)
		text.Draw(screen, "Tekan \"Space\" untuk play!", subF, WIDTH/2-120, HEIGHT/2, color.White)
		text.Draw(screen, "created by aji mustofa @pepega90", subF, 10, 460, color.RGBA{0, 255, 0, 255})
	case PLAY:
		text.Draw(screen, g.listSoal[g.currentQuestion].soal, f, int(g.listSoal[g.currentQuestion].x), int(g.listSoal[g.currentQuestion].y), color.White)

		for _, v := range g.listSoal[g.currentQuestion].pg {
			if g.choice {
				if v.correct {
					ebitenutil.DrawRect(screen, v.rectX, v.rectY, v.width, v.height, color.RGBA{0, 255, 0, 255})
				} else {
					ebitenutil.DrawRect(screen, v.rectX, v.rectY, v.width, v.height, color.RGBA{255, 0, 0, 255})
				}
			} else {
				ebitenutil.DrawRect(screen, v.rectX, v.rectY, v.width, v.height, color.White)
			}
			text.Draw(screen, v.jawab, subF, int(v.x), int(v.y), color.RGBA{0, 0, 0, 255})
		}
	case GAMEOVER:
		text.Draw(screen, "Game Over", gameOverF, WIDTH/2-100, HEIGHT/4, color.RGBA{255, 0, 0, 255})
		text.Draw(screen, fmt.Sprintf("Kamu benar %d/%d soal!", correctCount, len(g.listSoal)), subF, WIDTH/2-95, HEIGHT/2-30, color.White)
		text.Draw(screen, "Tekan \"R\" untuk restart", f, WIDTH/2-160, HEIGHT/2+80, color.White)
	}

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Mouse X: %v\nMouse Y: %v", g.mouse.x, g.mouse.y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

// global variabel
const (
	MENU = iota
	PLAY
	GAMEOVER
)

var (
	gameOverF    font.Face
	f            font.Face
	subF         font.Face
	lastUpdate   time.Time
	countDown    = 2
	currentScene = MENU
	trackSoal    int
	correctCount int
	correct      bool
)

func CheckPointCollision(mx int, my int, rect Answer) bool {
	if mx >= int(rect.rectX) && // right of the left edge AND
		mx <= int(rect.rectX+float64(rect.width)) && // left of the right edge AND
		my >= int(rect.rectY) && // below the top AND
		my <= int(rect.rectY+float64(rect.height)) { // above the bottom
		return true
	}
	return false
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Quiz Game")

	// load font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf) // <= wasm font
	if err != nil {
		log.Fatalf("opentype load font: %v", err)
		return
	}

	gameOverF, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	f, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	subF, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	g := new(Game)

	q1 := Question{
		soal: "Berapa hasil dari 5 x 5 ?",
		x:    160,
		y:    70,
		pg: []Answer{
			{
				x:       123,
				y:       257,
				rectX:   58,
				rectY:   223,
				width:   150,
				height:  50,
				jawab:   "2",
				correct: false,
			},
			{
				x:       491,
				y:       257,
				rectX:   430,
				rectY:   227,
				width:   150,
				height:  50,
				jawab:   "25",
				correct: true,
			},
			{
				x:       125,
				y:       400,
				rectX:   60,
				rectY:   370,
				width:   150,
				height:  50,
				jawab:   "10",
				correct: false,
			},
			{
				x:       499,
				y:       408,
				rectX:   433,
				rectY:   370,
				width:   150,
				height:  50,
				jawab:   "6",
				correct: false,
			},
		},
	}

	q2 := Question{
		soal: "Berapa hasil dari 112 + 8 ?",
		x:    140,
		y:    70,
		pg: []Answer{
			{
				x:       123,
				y:       257,
				rectX:   58,
				rectY:   223,
				width:   150,
				height:  50,
				jawab:   "120",
				correct: true,
			},
			{
				x:       491,
				y:       257,
				rectX:   430,
				rectY:   227,
				width:   150,
				height:  50,
				jawab:   "124",
				correct: false,
			},
			{
				x:       125,
				y:       400,
				rectX:   60,
				rectY:   370,
				width:   150,
				height:  50,
				jawab:   "131",
				correct: false,
			},
			{
				x:       499,
				y:       408,
				rectX:   433,
				rectY:   370,
				width:   150,
				height:  50,
				jawab:   "128",
				correct: false,
			},
		},
	}

	q3 := Question{
		soal: "Berapa hasil dari (5 x 4) + 12 - 16 ?",
		x:    65,
		y:    70,
		pg: []Answer{
			{
				x:       123,
				y:       257,
				rectX:   58,
				rectY:   223,
				width:   150,
				height:  50,
				jawab:   "58",
				correct: false,
			},
			{
				x:       491,
				y:       257,
				rectX:   430,
				rectY:   227,
				width:   150,
				height:  50,
				jawab:   "111",
				correct: false,
			},
			{
				x:       125,
				y:       400,
				rectX:   60,
				rectY:   370,
				width:   150,
				height:  50,
				jawab:   "123",
				correct: false,
			},
			{
				x:       499,
				y:       408,
				rectX:   433,
				rectY:   370,
				width:   150,
				height:  50,
				jawab:   "16",
				correct: true,
			},
		},
	}

	g.currentQuestion = 0
	lastUpdate = time.Now()
	g.listSoal = append(g.listSoal, q1, q2, q3)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

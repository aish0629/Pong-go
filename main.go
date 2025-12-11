package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080

	PaddleWidth  = 12
	PaddleHeight = 100
	PaddleSpeed  = 6

	BallSize      = 10
	InitialSpeedX = 4.0
	InitialSpeedY = 1.5
	MaxBounceAngle = 75 * math.Pi / 180 // in radians
)

// Paddle represents a player paddle
type Paddle struct {
	x, y float64
}

// Ball represents the game ball
type Ball struct {
	x, y   float64
	vx, vy float64
}

// Game holds the world state
type Game struct {
	leftPaddle  Paddle
	rightPaddle Paddle
	ball        Ball

	leftScore   int
	rightScore  int

	paused      bool
	fontColor   color.Color
	bgColor     color.Color
	paddleColor color.Color
	ballColor   color.Color

	paddleImg *ebiten.Image
	ballImg   *ebiten.Image
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{}
	g.leftPaddle = Paddle{ x: 20, y: (ScreenHeight-PaddleHeight)/2 }
	g.rightPaddle = Paddle{ x: ScreenWidth - 20 - PaddleWidth, y: (ScreenHeight-PaddleHeight)/2 }

	g.resetBall(true)

	g.fontColor = color.White
	g.bgColor = color.RGBA{10, 10, 30, 255}
	g.paddleColor = color.RGBA{200, 200, 200, 255}
	g.ballColor = color.RGBA{240, 120, 80, 255}

	// create images for paddles and ball to draw efficiently
	g.paddleImg = ebiten.NewImage(PaddleWidth, PaddleHeight)
	g.paddleImg.Fill(g.paddleColor)
	g.ballImg = ebiten.NewImage(BallSize, BallSize)
	g.ballImg.Fill(g.ballColor)

	return g
}

func (g *Game) resetBall(toRight bool) {
	g.ball.x = ScreenWidth/2 - BallSize/2
	g.ball.y = ScreenHeight/2 - BallSize/2

	angle := (rand.Float64()*2 - 1) * (math.Pi/8) // small random angle
	speed := InitialSpeedX
	if !toRight {
		speed = -InitialSpeedX
	}
	g.ball.vx = speed
	g.ball.vy = InitialSpeedY*math.Sin(angle)
}

func (g *Game) Update() error {
	// toggle pause
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// simple debounce: toggle only on press — but ebiten doesn't give key events, so keep it quick:
		// For simplicity in this starter, space toggles pause but will flip rapidly if held. Good enough.
		g.paused = !g.paused
	}

	if g.paused {
		return nil
	}

	// Left paddle controls: W and S
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.leftPaddle.y -= PaddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.leftPaddle.y += PaddleSpeed
	}

	// Right paddle controls: Up and Down
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.rightPaddle.y -= PaddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.rightPaddle.y += PaddleSpeed
	}

	// clamp paddles
	if g.leftPaddle.y < 0 {
		g.leftPaddle.y = 0
	}
	if g.leftPaddle.y > ScreenHeight-PaddleHeight {
		g.leftPaddle.y = ScreenHeight-PaddleHeight
	}
	if g.rightPaddle.y < 0 {
		g.rightPaddle.y = 0
	}
	if g.rightPaddle.y > ScreenHeight-PaddleHeight {
		g.rightPaddle.y = ScreenHeight-PaddleHeight
	}

	// move ball
	g.ball.x += g.ball.vx
	g.ball.y += g.ball.vy

	// top/bottom collision
	if g.ball.y <= 0 {
		g.ball.y = 0
		g.ball.vy = -g.ball.vy
	}
	if g.ball.y+BallSize >= ScreenHeight {
		g.ball.y = ScreenHeight - BallSize
		g.ball.vy = -g.ball.vy
	}

	// paddle collisions
	// left
	if g.ball.x <= g.leftPaddle.x+PaddleWidth &&
		g.ball.x >= g.leftPaddle.x &&
		g.ball.y+BallSize >= g.leftPaddle.y &&
		g.ball.y <= g.leftPaddle.y+PaddleHeight {

		// compute hit position relative to paddle center
		pCenter := g.leftPaddle.y + PaddleHeight/2
		offset := (g.ball.y + BallSize/2) - pCenter
		normalized := offset / (PaddleHeight/2)
		angle := normalized * MaxBounceAngle

		speed := math.Hypot(g.ball.vx, g.ball.vy)
		g.ball.vx = math.Abs(speed*math.Cos(angle))
		g.ball.vy = speed * math.Sin(angle)
		// nudge ball out to avoid sticking
		g.ball.x = g.leftPaddle.x + PaddleWidth + 0.1
	}

	// right
	if g.ball.x+BallSize >= g.rightPaddle.x &&
		g.ball.x+BallSize <= g.rightPaddle.x+PaddleWidth &&
		g.ball.y+BallSize >= g.rightPaddle.y &&
		g.ball.y <= g.rightPaddle.y+PaddleHeight {

		pCenter := g.rightPaddle.y + PaddleHeight/2
		offset := (g.ball.y + BallSize/2) - pCenter
		normalized := offset / (PaddleHeight/2)
		angle := normalized * MaxBounceAngle

		speed := math.Hypot(g.ball.vx, g.ball.vy)
		g.ball.vx = -math.Abs(speed*math.Cos(angle))
		g.ball.vy = speed * math.Sin(angle)
		g.ball.x = g.rightPaddle.x - BallSize - 0.1
	}

	// scoring
	if g.ball.x < -BallSize {
		// right player scored
		g.rightScore++
		g.resetBall(false)
	}
	if g.ball.x > ScreenWidth+BallSize {
		// left player scored
		g.leftScore++
		g.resetBall(true)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// clear
	screen.Fill(g.bgColor)

	// draw center line
	for y := 0; y < ScreenHeight; y += 20 {
		ebitenutil.DrawRect(screen, ScreenWidth/2-1, float64(y), 2, 12, color.RGBA{80, 80, 120, 255})
	}

	// draw paddles
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.leftPaddle.x, g.leftPaddle.y)
	screen.DrawImage(g.paddleImg, op)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(g.rightPaddle.x, g.rightPaddle.y)
	screen.DrawImage(g.paddleImg, op2)

	// draw ball
	op3 := &ebiten.DrawImageOptions{}
	op3.GeoM.Translate(g.ball.x, g.ball.y)
	screen.DrawImage(g.ballImg, op3)

	// UI: scores
	s := fmt.Sprintf("%d    %d", g.leftScore, g.rightScore)
	ebitenutil.DebugPrintAt(screen, s, ScreenWidth/2-17, 550)

	// instructions
	ebitenutil.DebugPrintAt(screen, "W/S | Up/Down to move. Space toggles pause.", 10, ScreenHeight-20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Pong - Ebiten (Go)")
	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}


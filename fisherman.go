package projectHailMary

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
)

const (
	lavaZoneRadius   = 200
	astrophageRadius = 10
	catchRadius      = 20
)

type Fisherman struct {
	X, Y       float64
	Caught     int
	Melted     bool
	Win        bool
	Finished   bool
	canvas     js.Value
	ctx        js.Value
	width      int
	height     int
	done       chan struct{}
	astrophage struct {
		X, Y float64
	}
}

func projectHailMary() {
	done := make(chan struct{}, 0)

	width := 800
	height := 600

	doc := js.Global().Get("document")
	canvas := doc.Call("createElement", "canvas")
	canvas.Set("width", width)
	canvas.Set("height", height)
	doc.Get("body").Call("appendChild", canvas)

	ctx := canvas.Call("getContext", "2d")

	fisherman := Fisherman{
		X:      float64(width / 2),
		Y:      float64(height / 2),
		canvas: canvas,
		ctx:    ctx,
		width:  width,
		height: height,
		done:   done,
	}

	fmt.Println("Welcome to the Fishing Astrophage game!")
	fmt.Println("You are a fisherman in a spaceship trying to catch astrophages in a hot lava zone.")
	fmt.Println("Use the arrow keys to move the spaceship and try to catch all the astrophages before getting too close to the lava.")

	// Generate random astrophage position
	astrophageX := getRandomNumber(astrophageRadius, width-astrophageRadius)
	astrophageY := getRandomNumber(astrophageRadius, height-astrophageRadius)

	fisherman.astrophage.X = astrophageX
	fisherman.astrophage.Y = astrophageY

	// Start game loop
	js.Global().Call("requestAnimationFrame", js.FuncOf(fisherman.update))

	<-done
}

func getRandomNumber(min, max int) float64 {
	return float64(min + rand.Intn(max-min))
}

func (f *Fisherman) update(this js.Value, args []js.Value) interface{} {
	f.render()

	if f.Finished {
		f.done <- struct{}{}
		return nil
	}

	doc := js.Global().Get("document")

	doc.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		keyCode := args[0].Get("keyCode").Int()

		switch keyCode {
		case 37: // Left arrow key
			f.X -= 5
		case 39: // Right arrow key
			f.X += 5
		case 38: // Up arrow key
			f.Y -= 5
		case 40: // Down arrow key
			f.Y += 5
		}

		if f.checkCollision() {
			f.Melted = true
			f.Finished = true
		}

		if f.Caught >= 3 {
			f.Win = true
			f.Finished = true
		}

		return nil
	}))

	js.Global().Call("requestAnimationFrame", js.FuncOf(f.update))
	return nil
}

func (f *Fisherman) checkCollision() bool {
	distance := math.Sqrt(math.Pow(f.X-f.astrophage.X, 2) + math.Pow(f.Y-f.astrophage.Y, 2))
	return distance <= catchRadius
}

func (f *Fisherman) render() {
	ctx := f.ctx
	ctx.Call("clearRect", 0, 0, f.width, f.height)

	// Render lava zone
	ctx.Set("fillStyle", "red")
	ctx.Call("beginPath")
	ctx.Call("arc", f.width/2, f.height/2, lavaZoneRadius, 0, 2*math.Pi)
	ctx.Call("fill")

	// Render astrophage
	ctx.Set("fillStyle", "blue")
	ctx.Call("beginPath")
	ctx.Call("arc", f.astrophage.X, f.astrophage.Y, astrophageRadius, 0, 2*math.Pi)
	ctx.Call("fill")

	// Render fisherman spaceship
	ctx.Set("fillStyle", "green")
	ctx.Call("beginPath")
	ctx.Call("arc", f.X, f.Y, 5, 0, 2*math.Pi)
	ctx.Call("fill")

	// Render game state
	ctx.Set("fillStyle", "black")
	ctx.Call("font", "14px sans-serif")
	ctx.Call("fillText", fmt.Sprintf("Astrophages caught: %d", f.Caught), 10, 20)

	if f.Melted {
		ctx.Call("fillText", "Oh no! The spaceship melted in the lava.", 10, 40)
	}

	if f.Win {
		ctx.Call("fillText", "Congratulations! You caught enough astrophages!", 10, 40)
	}
}

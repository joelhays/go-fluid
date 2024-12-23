package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"

	r "github.com/gen2brain/raylib-go/raylib"
	f "github.com/joelhays/go-fluid/fluid"
	g "github.com/joelhays/go-fluid/gui"
)

const windowWidth, windowHeight = 1024, 1024
const guiWidth = 200
const cellSize = 16
const halfCellSize = cellSize / 2
const gridWidth = windowWidth / cellSize
const gridHeight = windowHeight / cellSize
const gridSize = gridWidth * gridHeight
const cameraZoomIncrement float32 = 0.125

// render state
var camera = r.Camera2D{Zoom: 1.0}
var gui = g.NewGui(windowWidth, windowHeight, guiWidth)

// simulation state
var macFluid *f.MACFluid

func init() {
	macFluid = f.NewMACFluid(gridWidth)
}

func main() {
	profileFlag := flag.Bool("profile", false, "enable profiling")
	flag.Parse()

	fmt.Println("Profiling Enabled:", *profileFlag)

	if *profileFlag == true {
		f, err := os.Create("go-fluid.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}

	r.InitWindow(windowWidth+guiWidth, windowHeight, "Golang Fluid Simulation")
	defer r.CloseWindow()

	r.SetTargetFPS(120)

	for !r.WindowShouldClose() {
		handleMouseDrag()

		if gui.ResetSimulation {
			macFluid.Reset()
		}

		macFluid.DiffusionRate = gui.DiffusionRate
		macFluid.Viscosity = gui.Viscosity
		macFluid.FadeRate = gui.FadeRate

		macFluid.Simulate(gui.StepSize)

		r.BeginDrawing()
		{
			r.ClearBackground(r.Black)
			r.BeginMode2D(camera)
			{
				drawDensityField()
				drawVelocityField()
				drawGrid()
			}
			r.EndMode2D()

			r.DrawFPS(0, 0)

			gui.Run()
		}
		r.EndDrawing()
	}

}

func handleMouseDrag() {
	mouseDelta := r.GetMouseDelta()

	mousePosition := r.GetMousePosition()
	// Get the world point that is under the mouse
	mouseWorldPos := r.GetScreenToWorld2D(mousePosition, camera)

	// convert mouse position to grid cell
	x := int(mouseWorldPos.X) / cellSize
	y := int(mouseWorldPos.Y) / cellSize

	if int(x) < 1 || x > macFluid.Size || y < 1 || y > macFluid.Size {
		return
	}

	if r.IsMouseButtonDown(r.MouseButtonRight) || r.IsKeyDown(r.KeyLeftSuper) {
		force := gui.Force
		val := r.Vector2ClampValue(mouseDelta, -halfCellSize, halfCellSize)
		macFluid.AddVelocity(x, y, val.X*force, val.Y*force)
	}
	if r.IsMouseButtonDown(r.MouseButtonLeft) {
		macFluid.AddDensity(x, y, 1)
		brushRadius := gui.BrushRadius
		for x1 := x - int(brushRadius); x1 <= x+int(brushRadius); x1++ {
			for y1 := y - int(brushRadius); y1 <= y+int(brushRadius); y1++ {
				if int(x) == int(x1) && int(y) == int(y1) {
					continue
				}
				if int(x1) < 1 || x1 > macFluid.Size || y1 < 1 || y1 > macFluid.Size {
					continue
				}

				brushCenter := r.Vector2{X: float32(x), Y: float32(y)}
				brushPoint := r.Vector2{X: float32(x1), Y: float32(y1)}
				brushRadius := float32(gui.BrushRadius) + 0.5

				if r.CheckCollisionPointCircle(brushPoint, brushCenter, brushRadius) {
					macFluid.AddDensity(int(x1), int(y1), 1)
				}
			}
		}
	}
}

func drawGrid() {
	if gui.ShowGrid == false {
		return
	}

	color := r.DarkGray
	for x := range macFluid.Size + 1 {
		for y := range macFluid.Size + 1 {
			startX := int32(x * cellSize)
			startY := int32(y * cellSize)
			r.DrawLine(startX, startY, startX+cellSize, startY, color)
			r.DrawLine(startX+cellSize, startY, startX+cellSize, startY+cellSize, color)
			r.DrawLine(startX+cellSize, startY+cellSize, startX, startY+cellSize, color)
			r.DrawLine(startX, startY+cellSize, startX, startY, color)
		}
	}
}

func drawVelocityField() {
	if gui.ShowVelocityField == false {
		return
	}

	for x := 1; x <= macFluid.Size; x++ {
		for y := 1; y <= macFluid.Size; y++ {
			xv := macFluid.XVelocities[x][y] * cellSize
			yv := macFluid.YVelocities[x][y] * cellSize

			pos := r.NewVector2(float32(x*cellSize+halfCellSize), float32(y*cellSize+halfCellSize))
			dir := r.NewVector2(xv, yv)
			end := r.Vector2Add(pos, dir)

			r.DrawLineV(pos, end, r.Red)
		}
	}
}

func drawDensityField() {
	for x := range macFluid.Size + 1 {
		for y := range macFluid.Size + 1 {
			density := macFluid.DensityField[x][y]
			density1 := macFluid.DensityField[x+1][y]
			density2 := macFluid.DensityField[x][y+1]
			density3 := macFluid.DensityField[x+1][y+1]

			var vel, vel1, vel2, vel3 float32 = 0, 0, 0, 0
			if gui.ColorizeTurbulence {
				vel = (macFluid.XVelocities[x][y] + macFluid.YVelocities[x][y]) * 2
				vel1 = (macFluid.XVelocities[x+1][y] + macFluid.YVelocities[x+1][y]) * 2
				vel2 = (macFluid.XVelocities[x][y+1] + macFluid.YVelocities[x][y+1]) * 2
				vel3 = (macFluid.XVelocities[x+1][y+1] + macFluid.YVelocities[x+1][y+1]) * 2
			}

			pos := r.NewVector2(float32(x*cellSize), float32(y*cellSize))
			size := r.NewVector2(cellSize, cellSize)

			// hue 0-360, value 0-1
			baseColor := r.ColorToHSV(gui.FluidColor)
			topLeftColor := r.ColorFromHSV(float32(math.Mod(float64(baseColor.X+vel), 360)), baseColor.Y, baseColor.Z*density)
			bottomLeftColor := r.ColorFromHSV(float32(math.Mod(float64(baseColor.X+vel2), 360)), baseColor.Y, baseColor.Z*density2)
			topRightColor := r.ColorFromHSV(float32(math.Mod(float64(baseColor.X+vel3), 360)), baseColor.Y, baseColor.Z*density3)
			bottomRightColor := r.ColorFromHSV(float32(math.Mod(float64(baseColor.X+vel1), 360)), baseColor.Y, baseColor.Z*density1)

			r.DrawRectangleGradientEx(r.Rectangle{
				X:      pos.X,
				Y:      pos.Y,
				Width:  size.X,
				Height: size.Y,
			}, topLeftColor, bottomLeftColor, topRightColor, bottomRightColor)
		}
	}
}

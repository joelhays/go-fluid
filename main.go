package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	r "github.com/gen2brain/raylib-go/raylib"
	f "github.com/joelhays/go-fluid/fluid"
)

const windowWidth, windowHeight = 1024, 1024
const cellSize = 16
const halfCellSize = cellSize / 2
const gridWidth = windowWidth / cellSize
const gridHeight = windowHeight / cellSize
const gridSize = gridWidth * gridHeight
const cameraZoomIncrement float32 = 0.125

// render state
var camera = r.Camera2D{Zoom: 1.0}

// simulation state
// var fluid *f.Fluid
var macFluid *f.MACFluid

func init() {
	// fluid = f.NewFluid(gridWidth, gridHeight)
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

	r.InitWindow(windowWidth+150, windowHeight, "Golang Fluid Simulation")
	defer r.CloseWindow()

	r.SetTargetFPS(120)

	for !r.WindowShouldClose() {
		// handlePanAndZoom()
		handleMouseDrag()

		// fluid.Simulate(0.1) //r.GetFrameTime())

		macFluid.DiffusionRate = diffusionRate
		macFluid.Viscosity = viscosity
		macFluid.Simulate(0.1) //r.GetFrameTime())

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

			handleGui()
		}
		r.EndDrawing()

	}

}

func handlePanAndZoom() {
	mousePosition := r.GetMousePosition()
	// Get the world point that is under the mouse
	mouseWorldPos := r.GetScreenToWorld2D(mousePosition, camera)

	if r.IsMouseButtonDown(r.MouseButtonLeft) {
		camera.Offset.X += r.GetMouseDelta().X
		camera.Offset.Y += r.GetMouseDelta().Y
	}

	wheel := r.GetMouseWheelMove()
	if wheel != 0 {
		// Set the offset to where the mouse is
		camera.Offset = mousePosition

		// Set the target to match, so that the camera maps the world space point
		// under the cursor to the screen space point under the cursor at any zoom
		camera.Target = mouseWorldPos

		// Zoom increment
		camera.Zoom += (wheel * cameraZoomIncrement)
		camera.Zoom = r.Clamp(camera.Zoom, 1.0, 20.0)
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
		val := r.Vector2ClampValue(mouseDelta, -halfCellSize, halfCellSize)
		macFluid.AddVelocity(x, y, val.X*force, val.Y*force)
		// for x1 := x - int(brushRadius); x1 <= x+int(brushRadius); x1++ {
		// 	for y1 := y - int(brushRadius); y1 <= y+int(brushRadius); y1++ {
		// 		if int(x) == int(x1) && int(y) == int(y1) {
		// 			continue
		// 		}
		// 		if int(x1) < 1 || x1 > macFluid.Size || y1 < 1 || y1 > macFluid.Size {
		// 			continue
		// 		}
		// 		macFluid.AddVelocity(x1, y1, val.X, val.Y)
		// 	}
		// }
	}
	if r.IsMouseButtonDown(r.MouseButtonLeft) {
		macFluid.AddDensity(x, y, 1)
		for x1 := x - int(brushRadius); x1 <= x+int(brushRadius); x1++ {
			for y1 := y - int(brushRadius); y1 <= y+int(brushRadius); y1++ {
				if int(x) == int(x1) && int(y) == int(y1) {
					continue
				}
				if int(x1) < 1 || x1 > macFluid.Size || y1 < 1 || y1 > macFluid.Size {
					continue
				}
				macFluid.AddDensity(int(x1), int(y1), 1)
			}
		}
	}
}

func drawGrid() {
	if showGrid == false {
		return
	}

	color := r.DarkGray
	// for x := range gridWidth {
	// 	for y := range gridHeight {
	for x := range macFluid.Size + 1 {
		for y := range macFluid.Size + 1 {
			startX := int32(x * cellSize)
			startY := int32(y * cellSize)
			r.DrawLine(startX, startY, startX+cellSize, startY+0, color)
			r.DrawLine(startX+cellSize, startY, startX+cellSize, startY+cellSize, color)
			r.DrawLine(startX+cellSize, startY+cellSize, startX, startY+cellSize, color)
			r.DrawLine(startX, startY+cellSize, startX, startY, color)
		}
	}
}

func drawVelocityField() {
	if showVelocityField == false {
		return
	}

	// for x := range fluid.Width {
	// 	for y := range fluid.Height {
	// 		xv, _ := fluid.XVelocities.Get(x, y)
	// 		yv, _ := fluid.YVelocities.Get(x, y)
	//
	// 		pos := r.NewVector2(float32(x*cellSize+halfCellSize), float32(y*cellSize+halfCellSize))
	// 		dir := r.NewVector2(xv.Value, yv.Value)
	// 		end := r.Vector2Add(pos, dir)
	//
	// 		r.DrawLineV(pos, end, r.Red)
	// 		r.DrawCircleV(end, 1, r.Red)
	// 	}
	// }

	// endSize := r.Vector2{X: 2, Y: 2}
	for x := 1; x <= macFluid.Size; x++ {
		for y := 1; y <= macFluid.Size; y++ {
			// for x := range macFluid.Size + 2 {
			// 	for y := range macFluid.Size + 2 {
			xv := macFluid.XVelocities[x][y] * halfCellSize
			yv := macFluid.YVelocities[x][y] * halfCellSize

			pos := r.NewVector2(float32(x*cellSize+halfCellSize), float32(y*cellSize+halfCellSize))
			dir := r.NewVector2(xv, yv)
			end := r.Vector2Add(pos, dir)

			r.DrawLineV(pos, end, r.Red)
			// if cellSize >= 16 {
			// 	r.DrawRectangleV(end, endSize, r.Red)
			// }
			// r.DrawCircleV(end, 1, r.Red)
		}
	}
}

func drawDensityField() {
	// for x := range fluid.Width {
	// 	for y := range fluid.Height {
	// 		density, _ := fluid.DensityField.Get(x, y)
	// 		density1, _ := fluid.DensityField.Get(x+1, y)
	// 		density2, _ := fluid.DensityField.Get(x, y+1)
	// 		density3, _ := fluid.DensityField.Get(x+1, y+1)
	//
	// 		pos := r.NewVector2(float32(x*cellSize), float32(y*cellSize))
	// 		size := r.NewVector2(cellSize, cellSize)
	//
	// 		baseColor := r.ColorToHSV(fluidColor)
	// 		topLeftColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density.Value)
	// 		bottomLeftColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density2.Value)
	// 		topRightColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density3.Value)
	// 		bottomRightColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density1.Value)
	//
	// 		r.DrawRectangleGradientEx(r.Rectangle{
	// 			X:      pos.X,
	// 			Y:      pos.Y,
	// 			Width:  size.X,
	// 			Height: size.Y,
	// 		}, topLeftColor, bottomLeftColor, topRightColor, bottomRightColor)
	// 	}
	// }
	// for x := 1; x <= macFluid.Size; x++ {
	// 	for y := 1; y <= macFluid.Size; y++ {
	for x := range macFluid.Size + 1 {
		for y := range macFluid.Size + 1 {
			density := macFluid.DensityField[x][y]
			density1 := macFluid.DensityField[x+1][y]
			density2 := macFluid.DensityField[x][y+1]
			density3 := macFluid.DensityField[x+1][y+1]

			pos := r.NewVector2(float32(x*cellSize), float32(y*cellSize))
			size := r.NewVector2(cellSize, cellSize)

			baseColor := r.ColorToHSV(fluidColor)
			topLeftColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density)
			bottomLeftColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density2)
			topRightColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density3)
			bottomRightColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density1)

			r.DrawRectangleGradientEx(r.Rectangle{
				X:      pos.X,
				Y:      pos.Y,
				Width:  size.X,
				Height: size.Y,
			}, topLeftColor, bottomLeftColor, topRightColor, bottomRightColor)
		}
	}
}

package main

import (
	_ "github.com/gen2brain/raylib-go/raygui"
	r "github.com/gen2brain/raylib-go/raylib"
)

const windowWidth, windowHeight = 1024, 768
const cellSize = 32
const halfCellSize = cellSize / 2
const gridWidth = windowWidth / cellSize
const gridHeight = windowHeight / cellSize
const gridSize = gridWidth * gridHeight
const cameraZoomIncrement float32 = 0.125

// render state
var camera = r.Camera2D{Zoom: 1.0}

// simulation state
var fluid *Fluid

func init() {
	fluid = NewFluid()
}

func main() {

	r.InitWindow(windowWidth, windowHeight, "Golang Fluid Simulation")
	defer r.CloseWindow()

	r.SetTargetFPS(120)

	for !r.WindowShouldClose() {
		// handlePanAndZoom()
		handleMouseDrag()

		fluid.Simulate(r.GetFrameTime())

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

	if r.IsMouseButtonDown(r.MouseButtonRight) {
		val := r.Vector2ClampValue(mouseDelta, -halfCellSize, halfCellSize)
		fluid.AddVelocity(x, y, val.X, val.Y)
		for x1 := x - int(brushRadius); x1 <= x+int(brushRadius); x1++ {
			for y1 := y - int(brushRadius); y1 <= y+int(brushRadius); y1++ {
				if int(x) == int(x1) && int(y) == int(y1) {
					continue
				}
				fluid.AddVelocity(x1, y1, val.X, val.Y)
			}
		}
	}
	if r.IsMouseButtonDown(r.MouseButtonLeft) {
		fluid.AddDensity(x, y, .1)
		for x1 := x - int(brushRadius); x1 <= x+int(brushRadius); x1++ {
			for y1 := y - int(brushRadius); y1 <= y+int(brushRadius); y1++ {
				if int(x) == int(x1) && int(y) == int(y1) {
					continue
				}
				fluid.AddDensity(int(x1), int(y1), .1)
			}
		}
	}
}

func drawGrid() {
	if showGrid == false {
		return
	}

	color := r.DarkGray
	for x := range gridWidth {
		for y := range gridHeight {
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

	for x := range fluid.xVelocities.Width {
		for y := range fluid.xVelocities.Height {
			xv, _ := fluid.xVelocities.Get(x, y)
			yv, _ := fluid.yVelocities.Get(x, y)

			pos := r.NewVector2(float32(x*cellSize+halfCellSize), float32(y*cellSize+halfCellSize))
			dir := r.NewVector2(xv.Value, yv.Value)
			end := r.Vector2Add(pos, dir)

			r.DrawLineV(pos, end, r.Red)
			r.DrawCircleV(end, 1, r.Red)
		}
	}
}

func drawDensityField() {
	for x := range fluid.densityField.Width {
		for y := range fluid.densityField.Height {
			density, _ := fluid.densityField.Get(x, y)
			density1, _ := fluid.densityField.Get(x+1, y)
			density2, _ := fluid.densityField.Get(x, y+1)
			density3, _ := fluid.densityField.Get(x+1, y+1)

			pos := r.NewVector2(float32(x*cellSize), float32(y*cellSize))
			size := r.NewVector2(cellSize, cellSize)

			baseColor := r.ColorToHSV(fluidColor)
			topLeftColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density.Value)
			bottomLeftColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density2.Value)
			topRightColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density3.Value)
			bottomRightColor := r.ColorFromHSV(baseColor.X, baseColor.Y, baseColor.Z*density1.Value)

			r.DrawRectangleGradientEx(r.Rectangle{
				X:      pos.X,
				Y:      pos.Y,
				Width:  size.X,
				Height: size.Y,
			}, topLeftColor, bottomLeftColor, topRightColor, bottomRightColor)
		}
	}
}

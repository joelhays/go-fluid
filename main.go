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
		handlePanAndZoom()
		handleMouseDrag()

		fluid.Simulate()

		r.BeginDrawing()
		{
			r.ClearBackground(r.Black)
			r.BeginMode2D(camera)
			{
				drawGrid()
				drawVelocityField()
			}
			r.EndMode2D()

			r.DrawFPS(0, 0)
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
	if !r.IsMouseButtonDown(r.MouseButtonRight) {
		return
	}

	// add the mouse vector to the velocity vectors

	mouseDelta := r.GetMouseDelta()

	mousePosition := r.GetMousePosition()
	// Get the world point that is under the mouse
	mouseWorldPos := r.GetScreenToWorld2D(mousePosition, camera)

	// convert mouse position to grid cell
	i := int(mouseWorldPos.Y)/cellSize*gridWidth + int(mouseWorldPos.X)/cellSize

	// user may click outside of game area, so disregard
	if i > 0 && i < gridSize-1 {
		adjustedVector := r.Vector2Add(fluid.velocityField[i], mouseDelta)
		fluid.velocityField[i] = r.Vector2ClampValue(adjustedVector, -halfCellSize, halfCellSize)
	}
}

func drawGrid() {
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
	for x := range gridWidth {
		for y := range gridHeight {
			vector := fluid.velocityField[y*gridWidth+x]

			pos := r.NewVector2(float32(x*cellSize+halfCellSize), float32(y*cellSize+halfCellSize))
			dir := r.NewVector2(vector.X, vector.Y)
			end := r.Vector2Add(pos, dir)

			r.DrawLineV(pos, end, r.Red)
			r.DrawCircleV(end, 1, r.Red)
		}
	}
}

package main

import (
	rg "github.com/gen2brain/raylib-go/raygui"
	r "github.com/gen2brain/raylib-go/raylib"
)

var showVelocityField bool = false
var showGrid bool = false
var brushRadius int32 = 1
var diffusionRate float32 = 0.5
var fluidColor r.Color = r.Blue
var prevRect *r.Rectangle = nil
var panelRect *r.Rectangle = nil

func handleGui() {
	prevRect = nil

	var panelMargin float32 = 0
	var panelWidth float32 = 150.0
	var panelHeight float32 = windowHeight - panelMargin*2
	var panelX float32 = windowWidth - panelWidth - panelMargin
	panelRect = &r.Rectangle{X: panelX, Y: panelMargin, Width: panelWidth, Height: panelHeight}
	rg.Panel(*panelRect, "Fluid Simulation")

	rg.Line(getControlRect(), "Show Velocity Field")
	showVelocityField = rg.CheckBox(getControlRect(), "", showVelocityField)

	rg.Line(getControlRect(), "Show Grid")
	showGrid = rg.CheckBox(getControlRect(), "", showGrid)

	rg.Line(getControlRect(), "Brush Radius")
	brushRadius = rg.Spinner(getControlRect(), "", &brushRadius, 0, 5, false)

	rg.Line(getControlRect(), "Fluid Color")
	colorRect := getControlRect()
	colorRect.Height = 60
	colorRect.Width = prevRect.Width - 25
	prevRect.Height += 30
	prevRect.Y += 30

	fluidColor = rg.ColorPicker(colorRect, "", fluidColor)

	// rg.Line(getControlRect(), "Brush Radius")
	// brushRadius = rg.Spinner(getControlRect(), "", &brushRadius, 0, 5, false)
}

func getControlRect() r.Rectangle {
	if prevRect != nil {
		newRect := r.Rectangle{
			X:      prevRect.X,
			Y:      prevRect.Y + 30,
			Width:  prevRect.Width,
			Height: 30,
		}
		prevRect = &newRect
		return newRect
	}

	var padding float32 = 5
	prevRect = &r.Rectangle{
		X:      panelRect.X + padding,
		Y:      panelRect.Y + 30,
		Width:  panelRect.Width - padding*2,
		Height: 30,
	}
	return *prevRect
}

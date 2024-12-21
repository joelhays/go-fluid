package main

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	r "github.com/gen2brain/raylib-go/raylib"
)

var resetSimulation bool = false
var showVelocityField bool = false
var showGrid bool = false
var brushRadius int32 = 2
var diffusionRate float32 = 0.0000125
var viscosity float32 = 0.0001
var force float32 = 8
var stepSize float32 = 0.1
var fadeRate float32 = 0.0025
var fluidColor r.Color = r.Blue
var prevRect *r.Rectangle = nil
var panelRect *r.Rectangle = nil

func handleGui() {
	prevRect = nil

	var panelMargin float32 = 0
	var panelWidth float32 = 200.0
	var panelHeight float32 = windowHeight - panelMargin*2
	var panelX float32 = windowWidth - panelMargin
	panelRect = &r.Rectangle{X: panelX, Y: panelMargin, Width: panelWidth, Height: panelHeight}
	rg.Panel(*panelRect, "Fluid Simulation")

	rg.Line(getControlRect(), "Controls")
	resetSimulation = rg.Button(getControlRect(), "Reset Simulation")

	rg.Line(getControlRect(), "Show Velocity Field")
	showVelocityField = rg.CheckBox(getControlRect(), "", showVelocityField)

	rg.Line(getControlRect(), "Show Grid")
	showGrid = rg.CheckBox(getControlRect(), "", showGrid)

	rg.Line(getControlRect(), "Brush Radius")
	brushRadius = rg.Spinner(getControlRect(), "", &brushRadius, 0, 5, false)

	rg.Line(getControlRect(), fmt.Sprintf("Diffusion Rate - %1.7f", diffusionRate))
	diffusionRate = rg.Slider(getControlRect(), "", "", diffusionRate, 0.0, 0.0001)

	rg.Line(getControlRect(), fmt.Sprintf("Viscosity - %1.3f", viscosity))
	viscosity = rg.Slider(getControlRect(), "", "", viscosity, 0.0, 0.005)

	rg.Line(getControlRect(), fmt.Sprintf("Force - %1.2f", force))
	force = rg.Slider(getControlRect(), "", "", force, 1, 30)

	rg.Line(getControlRect(), fmt.Sprintf("Step Size - %1.2f", stepSize))
	stepSize = rg.Slider(getControlRect(), "", "", stepSize, 0.1, 1)

	rg.Line(getControlRect(), fmt.Sprintf("Fade Rate - %1.4f", fadeRate))
	fadeRate = rg.Slider(getControlRect(), "", "", fadeRate, 0.0, 0.02)

	rg.Line(getControlRect(), "Fluid Color")
	colorRect := getControlRect()
	colorRect.Height = 60
	colorRect.Width = prevRect.Width - 25
	prevRect.Height += 30
	prevRect.Y += 30

	fluidColor = rg.ColorPicker(colorRect, "", fluidColor)
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

	var padding float32 = 15
	prevRect = &r.Rectangle{
		X:      panelRect.X + padding,
		Y:      panelRect.Y + 30,
		Width:  panelRect.Width - padding*2,
		Height: 30,
	}
	return *prevRect
}

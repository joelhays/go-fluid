package gui

import (
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	r "github.com/gen2brain/raylib-go/raylib"
)

type Gui struct {
	// control state
	ResetSimulation   bool
	ShowVelocityField bool
	ShowGrid          bool
	BrushRadius       int32
	DiffusionRate     float32
	Viscosity         float32
	Force             float32
	StepSize          float32
	FadeRate          float32
	FluidColor        r.Color

	// internal state
	windowHeight float32
	windowWidth  float32
	panelWidth   float32
	prevRect     *r.Rectangle
	panelRect    *r.Rectangle
}

func NewGui(windowWidth, windowHeight, panelWidth int) *Gui {
	return &Gui{
		ResetSimulation:   false,
		ShowVelocityField: false,
		ShowGrid:          false,
		BrushRadius:       2,
		DiffusionRate:     0.0000125,
		Viscosity:         0.0001,
		Force:             10,
		StepSize:          0.1,
		FadeRate:          0.0025,
		FluidColor:        r.Blue,

		windowHeight: float32(windowHeight),
		windowWidth:  float32(windowWidth),
		panelWidth:   float32(panelWidth),
		prevRect:     nil,
		panelRect:    nil,
	}
}

func (g *Gui) Run() {
	g.prevRect = nil

	g.panelRect = &r.Rectangle{X: g.windowWidth, Y: 0, Width: g.panelWidth, Height: g.windowHeight}
	rg.Panel(*g.panelRect, "Fluid Simulation")

	rg.Line(g.getControlRect(), "Controls")
	g.ResetSimulation = rg.Button(g.getControlRect(), "Reset Simulation")

	rg.Line(g.getControlRect(), "Show Velocity Field")
	g.ShowVelocityField = rg.CheckBox(g.getControlRect(), "", g.ShowVelocityField)

	rg.Line(g.getControlRect(), "Show Grid")
	g.ShowGrid = rg.CheckBox(g.getControlRect(), "", g.ShowGrid)

	rg.Line(g.getControlRect(), "Brush Radius")
	g.BrushRadius = rg.Spinner(g.getControlRect(), "", &g.BrushRadius, 0, 5, false)

	rg.Line(g.getControlRect(), fmt.Sprintf("Diffusion Rate - %1.7f", g.DiffusionRate))
	g.DiffusionRate = rg.Slider(g.getControlRect(), "", "", g.DiffusionRate, 0.0, 0.0001)

	rg.Line(g.getControlRect(), fmt.Sprintf("Viscosity - %1.3f", g.Viscosity))
	g.Viscosity = rg.Slider(g.getControlRect(), "", "", g.Viscosity, 0.0, 0.005)

	rg.Line(g.getControlRect(), fmt.Sprintf("Force - %1.2f", g.Force))
	g.Force = rg.Slider(g.getControlRect(), "", "", g.Force, 1, 50)

	rg.Line(g.getControlRect(), fmt.Sprintf("Step Size - %1.2f", g.StepSize))
	g.StepSize = rg.Slider(g.getControlRect(), "", "", g.StepSize, 0.1, 1)

	rg.Line(g.getControlRect(), fmt.Sprintf("Fade Rate - %1.4f", g.FadeRate))
	g.FadeRate = rg.Slider(g.getControlRect(), "", "", g.FadeRate, 0.0, 0.02)

	rg.Line(g.getControlRect(), "Fluid Color")
	colorRect := g.getControlRect()
	colorRect.Height = 60
	colorRect.Width = g.prevRect.Width - 25
	g.prevRect.Height += 30
	g.prevRect.Y += 30

	g.FluidColor = rg.ColorPicker(colorRect, "", g.FluidColor)
}

func (g *Gui) getControlRect() r.Rectangle {
	if g.prevRect != nil {
		newRect := r.Rectangle{
			X:      g.prevRect.X,
			Y:      g.prevRect.Y + 30,
			Width:  g.prevRect.Width,
			Height: 30,
		}
		g.prevRect = &newRect
		return newRect
	}

	var padding float32 = 15
	g.prevRect = &r.Rectangle{
		X:      g.panelRect.X + padding,
		Y:      g.panelRect.Y + 30,
		Width:  g.panelRect.Width - padding*2,
		Height: 30,
	}
	return *g.prevRect
}

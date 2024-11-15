package main

import (
	r "github.com/gen2brain/raylib-go/raylib"
	"math/rand"
)

type Fluid struct {
	velocityField []r.Vector2
}

func NewFluid() *Fluid {
	fluid := &Fluid{
		velocityField: make([]r.Vector2, gridSize),
	}

	for i := range fluid.velocityField {
		x := float32(rand.Intn(cellSize) - halfCellSize)
		y := float32(rand.Intn(cellSize) - halfCellSize)
		fluid.velocityField[i] = r.NewVector2(x, y)
	}

	return fluid
}

func (f *Fluid) Simulate() {
	_ = r.GetFrameTime()
	f.advect()
	f.diffuse()
	f.project()
}

func (f *Fluid) advect() {
}
func (f *Fluid) diffuse() {
}
func (f *Fluid) project() {
}

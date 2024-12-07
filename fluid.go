package main

import (
	"math/rand"

	r "github.com/gen2brain/raylib-go/raylib"
)

type Fluid struct {
	velocityField    *Grid2[r.Vector2]
	densityField     *Grid2[float32]
	densityFieldPrev *Grid2[float32]
}

func NewFluid() *Fluid {
	fluid := &Fluid{
		velocityField:    NewGrid2[r.Vector2](gridWidth, gridHeight),
		densityField:     NewGrid2[float32](gridWidth, gridHeight),
		densityFieldPrev: NewGrid2[float32](gridWidth, gridHeight),
	}

	for x := range gridWidth {
		for y := range gridHeight {
			vx := float32(rand.Intn(cellSize) - halfCellSize)
			vy := float32(rand.Intn(cellSize) - halfCellSize)
			vx = float32(halfCellSize / 2)
			vy = float32(halfCellSize / 2)

			fluid.velocityField.Set(x, y, r.NewVector2(vx, vy))

			d := rand.Float32() * .75
			// d := float32(0)
			fluid.densityField.Set(x, y, d)
			fluid.densityFieldPrev.Set(x, y, d)
		}
	}

	return fluid
}

func (f *Fluid) Simulate() {
	f.simulateDensity()

	// f.advect()
	// f.copyState()
	// f.diffuse()
	// f.copyState()
	f.project()
}

func (f *Fluid) simulateDensity() {
	f.advect()
	f.copyState()
	f.diffuse()
	f.copyState()
}

func (f *Fluid) copyState() {
	// copy new state to previous
	for x := range f.densityField.Width {
		for y := range f.densityField.Height {
			val, _ := f.densityField.Get(x, y)
			f.densityFieldPrev.Set(x, y, val.Value)
		}
	}
}

func (f *Fluid) advect() {
	// move the density field through the velocity field

	for i := range gridWidth {
		for j := range gridHeight {
			v, _ := f.velocityField.Get(i, j)

			dt := r.GetFrameTime()

			// calculate previous x and y positions of the current grid particle
			// by moving backwards along the velocity field
			xBacktrace := float32(i) - v.Value.X*dt
			yBacktrace := float32(j) - v.Value.Y*dt

			prevX := r.Clamp(xBacktrace, 0.5, gridWidth+0.5)  // clamp x at the edges
			prevY := r.Clamp(yBacktrace, 0.5, gridHeight+0.5) // clamp y at the edges

			val := f.bilinearInterpolate(prevX, prevY)
			f.densityField.Set(i, j, val)
		}
	}
}

func (f *Fluid) bilinearInterpolate(x, y float32) float32 {
	// truncate x and y and get the indexes for the 4 adjacent cells at this position
	x0 := int(x)
	y0 := int(y)
	x1 := x0 + 1
	y1 := y0 + 1

	// calculate the floating point distance between the cell center and interpolation position
	// resulting in a value in the range 0.0 - 1.0, which represents the x and y contributions
	dx := x - float32(x0)
	dy := y - float32(y0)

	// get the density values at the 4 adjacent cells that will be interpolated
	v00, _ := f.densityFieldPrev.Get(x0, y0)
	v01, _ := f.densityFieldPrev.Get(x0, y1)
	v10, _ := f.densityFieldPrev.Get(x1, y0)
	v11, _ := f.densityFieldPrev.Get(x1, y1)

	// calculate the new density using the unit square method of bilinear interpolation
	// -- on a unit square, the four points are interpolated as:
	// 	 f(x,y) is appromixately f(0,0)(1-x)(1-y)+f(0,1)(1-x)y+f(1,0)x(1-y)+f(1,1)xy
	return v00.Value*(1-dx)*(1-dy) +
		v01.Value*(1-dx)*dy +
		v10.Value*dx*(1-dy) +
		v11.Value*dx*dy
}

func (f *Fluid) diffuse() {
	// diffuse the density field
	// high density cells diffuse to low density cells
	var relaxationSteps int = 20
	// var diffusionRate float32 = 0.0001
	var diffusionRate float32 = 0.05

	// diffusion delta
	diffusionFactor := r.GetFrameTime() * diffusionRate // * float32(f.densityFieldPrev.Width*f.densityFieldPrev.Height)

	// Gauss-Seidel Relaxation
	for _ = range relaxationSteps {
		for x := range f.densityFieldPrev.Width {
			for y := range f.densityFieldPrev.Height {
				self, _ := f.densityFieldPrev.Get(x, y)
				neighborDensities := f.densityField.GetNeighbors(x, y)
				numNeighbors := float32(len(neighborDensities))

				var sumOfNeighborDensities float32 = 0.0
				for _, neighbor := range neighborDensities {
					sumOfNeighborDensities += neighbor.Value
				}
				diffusedDensity := (self.Value + sumOfNeighborDensities*diffusionFactor) / (1 + numNeighbors*diffusionFactor)
				f.densityField.Set(x, y, diffusedDensity)
			}
		}
	}
}

func (f *Fluid) project() {
}

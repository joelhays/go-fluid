package main

import (
	"math/rand"

	r "github.com/gen2brain/raylib-go/raylib"
)

type Fluid2 struct {
	densityField     *Grid2[float32]
	densityFieldPrev *Grid2[float32]
	xVelocities      *Grid2[float32]
	xVelocitiesPrev  *Grid2[float32]
	yVelocities      *Grid2[float32]
	yVelocitiesPrev  *Grid2[float32]
}

func NewFluid2() *Fluid2 {
	fluid := &Fluid2{
		densityField:     NewGrid2[float32](gridWidth, gridHeight),
		densityFieldPrev: NewGrid2[float32](gridWidth, gridHeight),
		xVelocities:      NewGrid2[float32](gridWidth, gridHeight),
		xVelocitiesPrev:  NewGrid2[float32](gridWidth, gridHeight),
		yVelocities:      NewGrid2[float32](gridWidth, gridHeight),
		yVelocitiesPrev:  NewGrid2[float32](gridWidth, gridHeight),
	}

	for x := range gridWidth {
		for y := range gridHeight {
			vx := float32(rand.Intn(cellSize) - halfCellSize)
			vy := float32(rand.Intn(cellSize) - halfCellSize)

			fluid.xVelocities.Set(x, y, vx)
			fluid.xVelocitiesPrev.Set(x, y, vx)
			fluid.yVelocities.Set(x, y, vy)
			fluid.yVelocitiesPrev.Set(x, y, vy)

			d := float32(0)
			fluid.densityField.Set(x, y, d)
			fluid.densityFieldPrev.Set(x, y, d)
		}
	}

	return fluid
}

func (f *Fluid2) Simulate(dt float32) {
	f.simulateVelocity(dt)
	f.simulateDensity(dt)
}

func (f *Fluid2) simulateVelocity(dt float32) {
	var viscosity float32 = 0.025
	f.diffuse(dt, f.xVelocities, f.xVelocitiesPrev, viscosity)
	f.diffuse(dt, f.yVelocities, f.yVelocitiesPrev, viscosity)
	f.project(f.xVelocities, f.yVelocities, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.swapState()

	f.advect(dt, f.xVelocities, f.xVelocitiesPrev, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.advect(dt, f.yVelocities, f.yVelocitiesPrev, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.project(f.xVelocities, f.yVelocities, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.swapState()
}

func (f *Fluid2) simulateDensity(dt float32) {
	var diffusionRate float32 = 0.5
	f.advect(dt, f.densityField, f.densityFieldPrev, f.xVelocities, f.yVelocities)
	f.swapState()
	f.diffuse(dt, f.densityField, f.densityFieldPrev, diffusionRate)
	f.swapState()
}

func (f *Fluid2) swapState() {
	f.densityField, f.densityFieldPrev = f.densityFieldPrev, f.densityField
	f.xVelocities, f.xVelocitiesPrev = f.xVelocitiesPrev, f.xVelocities
	f.yVelocities, f.yVelocitiesPrev = f.yVelocitiesPrev, f.yVelocities
}

func (f *Fluid2) advect(dt float32, grid *Grid2[float32], gridPrev *Grid2[float32], xVelocities, yVelocities *Grid2[float32]) {
	for i := range grid.Width {
		for j := range grid.Height {

			xv, _ := xVelocities.Get(i, j)
			yv, _ := yVelocities.Get(i, j)

			// calculate previous x and y positions of the current grid particle
			// by moving backwards along the velocity field.
			// by calculating new density based on previous particle position,
			// the simulation becomes bounded.
			px := float32(i) - xv.Value*dt
			py := float32(j) - yv.Value*dt

			px = r.Clamp(px, 0.5, gridWidth+0.5)  // clamp x at the edges
			py = r.Clamp(py, 0.5, gridHeight+0.5) // clamp y at the edges

			val := f.bilinearInterpolate(px, py, gridPrev)

			grid.Set(i, j, val)
		}
	}
}

func (f *Fluid2) diffuse(dt float32, grid *Grid2[float32], gridPrev *Grid2[float32], diffusionRate float32) {
	// diffuse the density field
	// high density cells diffuse to low density cells
	var relaxationSteps int = 20

	// diffusion delta
	diffusionFactor := dt * diffusionRate

	// Gauss-Seidel Relaxation
	for _ = range relaxationSteps {
		for x := range grid.Width {
			for y := range grid.Height {
				self, _ := gridPrev.Get(x, y)
				neighborValues := grid.GetNeighbors(x, y)
				numNeighbors := float32(len(neighborValues))

				var sumOfNeighborValues float32 = 0.0
				for _, neighbor := range neighborValues {
					sumOfNeighborValues += neighbor.Value
				}
				diffusedValue := (self.Value + sumOfNeighborValues*diffusionFactor) / (1 + numNeighbors*diffusionFactor)
				grid.Set(x, y, diffusedValue)
			}
		}
	}
}

func (f *Fluid2) project(grid *Grid2[float32], gridPrev *Grid2[float32], xVelocities, yVelocities *Grid2[float32]) {
}

func (f *Fluid2) bilinearInterpolate(x, y float32, grid *Grid2[float32]) float32 {
	// truncate x and y and get the indexes for the 4 adjacent cells at this position
	x0 := int(x)
	y0 := int(y)
	x1 := x0 + 1
	y1 := y0 + 1

	// calculate the floating point distance between the cell center and interpolation position
	// resulting in a value in the range 0.0 - 1.0, which represents the x and y contributions
	dx := x - float32(x0)
	dy := y - float32(y0)

	// get the values at the 4 adjacent cells that will be interpolated
	v00, _ := grid.Get(x0, y0)
	v01, _ := grid.Get(x0, y1)
	v10, _ := grid.Get(x1, y0)
	v11, _ := grid.Get(x1, y1)

	// calculate the new density using the unit square method of bilinear interpolation
	// -- on a unit square, the four points are interpolated as:
	// 	 f(x,y) is appromixately f(0,0)(1-x)(1-y)+f(0,1)(1-x)y+f(1,0)x(1-y)+f(1,1)xy
	return v00.Value*(1-dx)*(1-dy) +
		v01.Value*(1-dx)*dy +
		v10.Value*dx*(1-dy) +
		v11.Value*dx*dy
}

package Fluid

import (
	"math/rand"
)

type Fluid struct {
	Width            int
	Height           int
	DensityField     *Grid2[float32]
	densityFieldPrev *Grid2[float32]
	XVelocities      *Grid2[float32]
	xVelocitiesPrev  *Grid2[float32]
	YVelocities      *Grid2[float32]
	yVelocitiesPrev  *Grid2[float32]
}

func NewFluid(width, height int) *Fluid {
	fluid := &Fluid{
		Width:            width,
		Height:           height,
		DensityField:     NewGrid2[float32](width, height),
		densityFieldPrev: NewGrid2[float32](width, height),
		XVelocities:      NewGrid2[float32](width, height),
		xVelocitiesPrev:  NewGrid2[float32](width, height),
		YVelocities:      NewGrid2[float32](width, height),
		yVelocitiesPrev:  NewGrid2[float32](width, height),
	}

	for x := range width {
		for y := range height {
			// vx := float32(rand.Intn(cellSize) - halfCellSize)
			// vy := float32(rand.Intn(cellSize) - halfCellSize)
			vx := float32(rand.Float32()-0.5) * 64
			vy := float32(rand.Float32()-0.5) * 64

			fluid.XVelocities.Set(x, y, vx)
			fluid.xVelocitiesPrev.Set(x, y, vx)
			fluid.YVelocities.Set(x, y, vy)
			fluid.yVelocitiesPrev.Set(x, y, vy)

			d := float32(0)
			fluid.DensityField.Set(x, y, d)
			fluid.densityFieldPrev.Set(x, y, d)
		}
	}

	return fluid
}

func (f *Fluid) Simulate(dt float32) {
	f.simulateVelocity(dt)
	f.simulateDensity(dt)
}

func (f *Fluid) AddDensity(x, y int, val float32) {
	if v, ok := f.densityFieldPrev.Get(x, y); ok {
		newVal := v.Value + val
		if newVal > 1 {
			newVal = 1
		}
		f.densityFieldPrev.Set(x, y, newVal)
	}
}

func (f *Fluid) AddVelocity(x, y int, xval, yval float32) {
	if xv, ok := f.xVelocitiesPrev.Get(x, y); ok {
		newxval := xv.Value + xval
		if yv, ok := f.yVelocitiesPrev.Get(x, y); ok {
			newyval := yv.Value + yval
			f.xVelocitiesPrev.Set(x, y, newxval)
			f.yVelocitiesPrev.Set(x, y, newyval)
		}
	}
}

func (f *Fluid) simulateVelocity(dt float32) {
	var viscosity float32 = 0.025
	diffuse(dt, f.XVelocities, f.xVelocitiesPrev, viscosity)
	diffuse(dt, f.YVelocities, f.yVelocitiesPrev, viscosity)
	f.project(f.XVelocities, f.YVelocities, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.swapState()

	advect(dt, f.XVelocities, f.xVelocitiesPrev, f.xVelocitiesPrev, f.yVelocitiesPrev)
	advect(dt, f.YVelocities, f.yVelocitiesPrev, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.project(f.XVelocities, f.YVelocities, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.swapState()
}

func (f *Fluid) simulateDensity(dt float32) {
	var diffusionRate float32 = 0.5
	advect(dt, f.DensityField, f.densityFieldPrev, f.XVelocities, f.YVelocities)
	f.swapState()
	diffuse(dt, f.DensityField, f.densityFieldPrev, diffusionRate)
	f.swapState()
}

func (f *Fluid) swapState() {
	f.DensityField, f.densityFieldPrev = f.densityFieldPrev, f.DensityField
	f.XVelocities, f.xVelocitiesPrev = f.xVelocitiesPrev, f.XVelocities
	f.YVelocities, f.yVelocitiesPrev = f.yVelocitiesPrev, f.YVelocities
}

func advect(dt float32, grid *Grid2[float32], gridPrev *Grid2[float32], xVelocities, yVelocities *Grid2[float32]) {
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

			if px < 0.5 {
				px = 0.5
			} else if px > float32(grid.Width)+0.5 {
				px = float32(grid.Width) + 0.5
			}

			if py < 0.5 {
				py = 0.5
			} else if py > float32(grid.Height)+0.5 {
				py = float32(grid.Height) + 0.5
			}

			val := bilinearInterpolate(px, py, gridPrev)

			grid.Set(i, j, val)
		}
	}
	setBoundaries(grid)
}

func diffuse(dt float32, grid *Grid2[float32], gridPrev *Grid2[float32], diffusionRate float32) {
	// diffuse the density field
	// high density cells diffuse to low density cells
	var relaxationSteps int = 20

	// diffusion delta
	diffusionFactor := dt * diffusionRate

	// Gauss-Seidel Relaxation
	for range relaxationSteps {
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
		setBoundaries(grid)
	}
}

func (f *Fluid) project(xVelocities, yVelocities, xVelocitiesPrev, yVelocitiesPrev *Grid2[float32]) {
	// var N float32 = float32(f.Height) // float32(XVelocities.Width * XVelocities.Height)
	// var h float32 = 1.0 / N

	// calculate grid spacing
	// dx := 1024 / float32(f.Width)
	// dy := 768 / float32(f.Height)
	dx := float32(32) //float32(f.Width) / 32
	dy := float32(32) // float32(f.Height) / 32

	for i := range xVelocities.Width {
		for j := range xVelocities.Height {
			a, _ := xVelocities.Get(i+1, j)
			b, _ := xVelocities.Get(i-1, j)
			c, _ := yVelocities.Get(i, j+1)
			d, _ := yVelocities.Get(i, j-1)
			// e := -0.5 * h * (a.Value - b.Value + c.Value - d.Value)
			// e := -0.5 * ((a.Value-b.Value)/float32(f.Width) + (c.Value-d.Value)/float32(f.Height))
			// e := -0.5 * (a.Value - b.Value + c.Value - d.Value) / float32(f.Height)

			du_dx := a.Value + b.Value/(2.0*dx)
			dv_dy := c.Value + d.Value/(2.0*dy)
			divergence := du_dx + dv_dy

			yVelocitiesPrev.Set(i, j, divergence)
			// yVelocitiesPrev.Set(i, j, e)
			xVelocitiesPrev.Set(i, j, 0.0)
		}
	}
	setBoundaries(yVelocitiesPrev)
	setBoundaries(xVelocitiesPrev)

	var relaxationSteps int = 20
	for range relaxationSteps {
		for i := range xVelocities.Width {
			for j := range xVelocities.Height {
				a, _ := yVelocitiesPrev.Get(i, j)
				b, _ := xVelocitiesPrev.Get(i-1, j)
				c, _ := xVelocitiesPrev.Get(i+1, j)
				d, _ := xVelocitiesPrev.Get(i, j-1)
				e, _ := xVelocitiesPrev.Get(i, j+1)
				// todo: just get all neighbor velocities?
				f := (a.Value + b.Value + c.Value + d.Value + e.Value) / 4.0
				xVelocitiesPrev.Set(i, j, f)
			}
		}
		setBoundaries(xVelocities)
	}

	// for range relaxationSteps {
	// 	for x := range xVelocitiesPrev.Width {
	// 		for y := range xVelocitiesPrev.Height {
	// 			yv, _ := yVelocitiesPrev.Get(x, y)
	//
	// 			neighborValues := xVelocitiesPrev.GetNeighbors(x, y)
	// 			numNeighbors := float32(len(neighborValues))
	//
	// 			var sumOfNeighborValues float32 = 0.0
	// 			for _, neighbor := range neighborValues {
	// 				sumOfNeighborValues += neighbor.Value
	// 			}
	// 			newVal := yv.Value + (sumOfNeighborValues / numNeighbors)
	// 			xVelocitiesPrev.Set(x, y, newVal)
	// 		}
	// 	}
	// }

	for i := range xVelocities.Width {
		for j := range xVelocities.Height {
			a, _ := xVelocities.Get(i+1, j)
			b, _ := xVelocities.Get(i-1, j)
			c, _ := xVelocities.Get(i, j+1)
			d, _ := xVelocities.Get(i, j-1)

			// e := 0.5 * float32(f.Width) * (a.Value - b.Value)
			// f := 0.5 * float32(f.Height) * (c.Value - d.Value)
			e := (a.Value - b.Value) / (2.0 * dx) // / float32(f.Width)
			f := (c.Value - d.Value) / (2.0 * dy) // / float32(f.Height)

			xv, _ := xVelocities.Get(i, j)
			yv, _ := yVelocities.Get(i, j)

			xVelocities.Set(i, j, xv.Value-e)
			yVelocities.Set(i, j, yv.Value-f)
		}
	}
	setBoundaries(yVelocities)
	setBoundaries(xVelocities)
}

func setBoundaries(grid *Grid2[float32]) {
	h := grid.Height - 1
	for i := range grid.Width {
		vtop, _ := grid.Get(i, 1)
		vbot, _ := grid.Get(i, h-1)
		grid.Set(i, 0, -vtop.Value)
		grid.Set(i, h, -vbot.Value)
	}

	w := grid.Width - 1
	for j := range grid.Height {
		vleft, _ := grid.Get(1, j)
		vright, _ := grid.Get(w-1, j)
		grid.Set(0, j, -vleft.Value)
		grid.Set(w, j, -vright.Value)
	}

	v10, _ := grid.Get(1, 0)
	v01, _ := grid.Get(0, 1)
	vtopleft := -(v10.Value + v01.Value) / 2.0
	grid.Set(0, 0, -vtopleft)

	vw0, _ := grid.Get(w-1, 0)
	vwj, _ := grid.Get(w, 1)
	vtopright := -(vw0.Value + vwj.Value) / 2.0
	grid.Set(w, 0, -vtopright)

	v0h, _ := grid.Get(0, h-1)
	v1h, _ := grid.Get(1, h)
	vbottomleft := -(v0h.Value + v1h.Value) / 2.0
	grid.Set(0, h, -vbottomleft)

	vw1h, _ := grid.Get(w-1, h)
	vwh1, _ := grid.Get(w, h-1)
	vbottomright := -(vw1h.Value + vwh1.Value) / 2.0
	grid.Set(w, h, -vbottomright)
}

func bilinearInterpolate(x, y float32, grid *Grid2[float32]) float32 {
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

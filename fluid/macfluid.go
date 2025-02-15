package Fluid

type boundaryAction int

const (
	copyBoundary boundaryAction = iota
	reflectBoundary
)

type swapStateAction int

const (
	swapVelocities swapStateAction = iota
	swapDensities
)

type MACFluid struct {
	Size             int
	DensityField     [][]float32
	densityFieldPrev [][]float32
	XVelocities      [][]float32
	xVelocitiesPrev  [][]float32
	YVelocities      [][]float32
	yVelocitiesPrev  [][]float32
	DiffusionRate    float32
	Viscosity        float32
	FadeRate         float32
}

func NewMACFluid(size int) *MACFluid {
	// allocate additional space for boundary conditions
	paddedSize := size + 2

	makeArray2d := func() [][]float32 {
		arr := make([][]float32, paddedSize)
		for i := range paddedSize {
			arr[i] = make([]float32, paddedSize)
		}
		return arr
	}

	fluid := &MACFluid{
		Size:             size,
		DensityField:     makeArray2d(),
		densityFieldPrev: makeArray2d(),
		XVelocities:      makeArray2d(),
		xVelocitiesPrev:  makeArray2d(),
		YVelocities:      makeArray2d(),
		yVelocitiesPrev:  makeArray2d(),
		DiffusionRate:    0,
		Viscosity:        0,
		FadeRate:         0,
	}

	return fluid
}

func (f *MACFluid) Reset() {
	for x := 0; x < f.Size+2; x++ {
		for y := 0; y < f.Size+2; y++ {
			f.DensityField[x][y] = 0
			f.densityFieldPrev[x][y] = 0
			f.XVelocities[x][y] = 0
			f.xVelocitiesPrev[x][y] = 0
			f.YVelocities[x][y] = 0
			f.yVelocitiesPrev[x][y] = 0
		}
	}
}

func (f *MACFluid) Simulate(dt float32) {
	f.simulateVelocity(dt)
	f.simulateDensity(dt)
}

func (f *MACFluid) AddDensity(x, y int, val float32) {
	f.DensityField[x][y] = val
}

func (f *MACFluid) AddVelocity(x, y int, xval, yval float32) {
	f.XVelocities[x][y] = xval
	f.YVelocities[x][y] = yval
}

func (f *MACFluid) swapState(s swapStateAction) {
	switch s {
	case swapVelocities:
		f.XVelocities, f.xVelocitiesPrev = f.xVelocitiesPrev, f.XVelocities
		f.YVelocities, f.yVelocitiesPrev = f.yVelocitiesPrev, f.YVelocities
	case swapDensities:
		f.DensityField, f.densityFieldPrev = f.densityFieldPrev, f.DensityField
	}
}

func (f *MACFluid) simulateVelocity(dt float32) {
	var viscosity float32 = f.Viscosity

	f.swapState(swapVelocities)
	f.diffuse(reflectBoundary, dt, f.XVelocities, f.xVelocitiesPrev, viscosity)
	f.diffuse(reflectBoundary, dt, f.YVelocities, f.yVelocitiesPrev, viscosity)
	f.project(f.XVelocities, f.YVelocities, f.xVelocitiesPrev, f.yVelocitiesPrev)

	f.swapState(swapVelocities)
	f.advect(reflectBoundary, dt, f.XVelocities, f.xVelocitiesPrev, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.advect(reflectBoundary, dt, f.YVelocities, f.yVelocitiesPrev, f.xVelocitiesPrev, f.yVelocitiesPrev)
	f.project(f.XVelocities, f.YVelocities, f.xVelocitiesPrev, f.yVelocitiesPrev)
}

func (f *MACFluid) simulateDensity(dt float32) {
	f.fade(dt, f.DensityField, f.FadeRate)

	f.swapState(swapDensities)
	f.diffuse(copyBoundary, dt, f.DensityField, f.densityFieldPrev, f.DiffusionRate)

	f.swapState(swapDensities)
	f.advect(copyBoundary, dt, f.DensityField, f.densityFieldPrev, f.XVelocities, f.YVelocities)
}

func (f *MACFluid) fade(dt float32, grid [][]float32, fadeRate float32) {
	for x := 0; x < f.Size+2; x++ {
		for y := 0; y < f.Size+2; y++ {
			grid[x][y] -= dt * fadeRate
			if grid[x][y] < 0 {
				grid[x][y] = 0
			}
		}
	}

}

func (f *MACFluid) advect(b boundaryAction, dt float32, grid [][]float32, gridPrev [][]float32, xVelocities, yVelocities [][]float32) {
	for x := 1; x <= f.Size; x++ {
		for y := 1; y <= f.Size; y++ {

			xv := xVelocities[x][y]
			yv := yVelocities[x][y]

			// calculate previous x and y positions of the current grid particle
			// by moving backwards along the velocity field.
			// by calculating new density based on previous particle position,
			// the simulation becomes bounded.
			px := float32(x) - xv*dt
			py := float32(y) - yv*dt

			if px < 0.5 {
				px = 0.5
			} else if px > float32(f.Size)+0.5 {
				px = float32(f.Size) + 0.5
			}

			if py < 0.5 {
				py = 0.5
			} else if py > float32(f.Size)+0.5 {
				py = float32(f.Size) + 0.5
			}

			val := f.bilinearInterpolate(px, py, gridPrev)

			grid[x][y] = val
		}
	}
	f.setBoundaries(b, grid)
}

func (f *MACFluid) diffuse(b boundaryAction, dt float32, grid [][]float32, gridPrev [][]float32, diffusionRate float32) {
	// diffuse the density field
	// high density cells diffuse to low density cells
	var relaxationSteps int = 20

	// diffusion delta
	diffusionFactor := dt * diffusionRate * float32(f.Size) * float32(f.Size)

	// Gauss-Seidel Relaxation
	for range relaxationSteps {
		for x := 1; x <= f.Size; x++ {
			for y := 1; y <= f.Size; y++ {
				self := gridPrev[x][y]

				right := grid[x+1][y]
				left := grid[x-1][y]
				bottom := grid[x][y+1]
				top := grid[x][y-1]

				sumOfNeighborValues := right + left + bottom + top
				var numNeighbors float32 = 4.0

				diffusedValue := (self + sumOfNeighborValues*diffusionFactor) / (1 + numNeighbors*diffusionFactor)
				grid[x][y] = diffusedValue
			}
		}
		f.setBoundaries(b, grid)
	}
}

func (f *MACFluid) project(xVelocities, yVelocities, xVelocitiesPrev, yVelocitiesPrev [][]float32) {
	for x := 1; x <= f.Size; x++ {
		for y := 1; y <= f.Size; y++ {
			a := xVelocities[x+1][y]
			b := xVelocities[x-1][y]
			c := yVelocities[x][y+1]
			d := yVelocities[x][y-1]

			divergence := -0.5 * (a - b + c - d) / float32(f.Size)

			yVelocitiesPrev[x][y] = divergence
			xVelocitiesPrev[x][y] = 0.0
		}
	}
	f.setBoundaries(copyBoundary, yVelocitiesPrev)
	f.setBoundaries(copyBoundary, xVelocitiesPrev)

	var relaxationSteps int = 20
	for range relaxationSteps {
		for x := 1; x <= f.Size; x++ {
			for y := 1; y <= f.Size; y++ {
				a := yVelocitiesPrev[x][y]
				b := xVelocitiesPrev[x-1][y]
				c := xVelocitiesPrev[x+1][y]
				d := xVelocitiesPrev[x][y-1]
				e := xVelocitiesPrev[x][y+1]
				f := (a + b + c + d + e) / 4.0
				xVelocitiesPrev[x][y] = f
			}
		}
		f.setBoundaries(copyBoundary, xVelocities)
	}

	for x := 1; x <= f.Size; x++ {
		for y := 1; y <= f.Size; y++ {
			a := xVelocitiesPrev[x+1][y]
			b := xVelocitiesPrev[x-1][y]
			c := xVelocitiesPrev[x][y+1]
			d := xVelocitiesPrev[x][y-1]

			xVelocities[x][y] -= 0.5 * float32(f.Size) * (a - b)
			yVelocities[x][y] -= 0.5 * float32(f.Size) * (c - d)
		}
	}
	f.setBoundaries(reflectBoundary, yVelocities)
	f.setBoundaries(reflectBoundary, xVelocities)
}

func (f *MACFluid) setBoundaries(b boundaryAction, grid [][]float32) {
	// copy or reflect cell values to boundary cells as needed
	// to ensure the simulation is properly contained

	for i := 1; i <= f.Size; i++ {
		if b == reflectBoundary {
			grid[0][i] = -grid[1][i]
			grid[f.Size+1][i] = -grid[f.Size][i]
		} else {
			grid[0][i] = grid[1][i]
			grid[f.Size+1][i] = grid[f.Size][i]
		}

		if b == reflectBoundary {
			grid[i][0] = -grid[i][1]
			grid[i][f.Size+1] = -grid[i][f.Size]
		} else {
			grid[i][0] = grid[i][1]
			grid[i][f.Size+1] = grid[i][f.Size]
		}

		grid[0][0] = 0.5 * (grid[1][0] + grid[0][1])
		grid[0][f.Size+1] = 0.5 * (grid[1][f.Size+1] + grid[0][f.Size])
		grid[f.Size+1][0] = 0.5 * (grid[f.Size][0] + grid[f.Size+1][1])
		grid[f.Size+1][f.Size+1] = 0.5 * (grid[f.Size][f.Size+1] + grid[f.Size+1][f.Size])
	}
}

func (f *MACFluid) bilinearInterpolate(x, y float32, grid [][]float32) float32 {
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
	v00 := grid[x0][y0]
	v01 := grid[x0][y1]
	v10 := grid[x1][y0]
	v11 := grid[x1][y1]

	// calculate the new density using the unit square method of bilinear interpolation
	// -- on a unit square, the four points are interpolated as:
	// 	 f(x,y) is appromixately f(0,0)(1-x)(1-y)+f(0,1)(1-x)y+f(1,0)x(1-y)+f(1,1)xy
	return v00*(1-dx)*(1-dy) +
		v01*(1-dx)*dy +
		v10*dx*(1-dy) +
		v11*dx*dy
}

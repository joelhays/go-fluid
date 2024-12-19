package Fluid

type GridCell[T any] struct {
	X     int
	Y     int
	Value T
}

type Grid2[T any] struct {
	Width  int
	Height int
	cells  []GridCell[T]
}

func NewGrid2[T any](width, height int) *Grid2[T] {
	return &Grid2[T]{
		Width:  width,
		Height: height,
		cells:  make([]GridCell[T], width*height),
	}
}

func (g *Grid2[T]) Get(x, y int) (GridCell[T], bool) {
	idx := y*g.Width + x
	if idx < 0 || idx > g.Width*g.Height-1 {
		var result GridCell[T]
		result.X = -1
		result.Y = -1
		return result, false
	}
	return g.cells[idx], true
}

func (g *Grid2[T]) Set(x, y int, val T) {
	idx := y*g.Width + x
	if idx < 0 || idx > g.Width*g.Height-1 {
		return
	}
	g.cells[idx] = GridCell[T]{X: x, Y: y, Value: val}
}

func (g *Grid2[T]) GetNeighbors(x, y int) []GridCell[T] {
	neighbors := make([]GridCell[T], 0)
	if topLeft, ok := g.Get(x-1, y-1); ok {
		neighbors = append(neighbors, topLeft)
	}
	if top, ok := g.Get(x, y-1); ok {
		neighbors = append(neighbors, top)
	}
	if topRight, ok := g.Get(x+1, y-1); ok {
		neighbors = append(neighbors, topRight)
	}
	if left, ok := g.Get(x-1, y); ok {
		neighbors = append(neighbors, left)
	}
	if right, ok := g.Get(x+1, y); ok {
		neighbors = append(neighbors, right)
	}
	if bottomLeft, ok := g.Get(x-1, y+1); ok {
		neighbors = append(neighbors, bottomLeft)
	}
	if bottom, ok := g.Get(x, y+1); ok {
		neighbors = append(neighbors, bottom)
	}
	if bottomRight, ok := g.Get(x+1, y+1); ok {
		neighbors = append(neighbors, bottomRight)
	}
	return neighbors
}

func (g *Grid2[T]) SetNeighbors(x, y int, val T) {
	g.Set(x-1, y-1, val)
	g.Set(x, y-1, val)
	g.Set(x+1, y-1, val)
	g.Set(x-1, y, val)
	g.Set(x+1, y, val)
	g.Set(x-1, y+1, val)
	g.Set(x, y+1, val)
	g.Set(x+1, y+1, val)
}

package ray

import "math"

type Camera struct {
	Position     Quat
	Direction    Quat
	CanvasWidth  float64
	CanvasHeight float64
	TMin         float64
	TMax         float64
	Distance     float64
	Fov          float64
}

func (c *Camera) Rotate(q Quat) {
	c.Direction = Mul(c.Direction, q).Norm()
}

func (c Camera) TranslateCoords(x float64, y float64) Quat {

	viewportWidth := math.Tan(c.Fov * math.Pi / 180)
	viewportHeight := viewportWidth * (c.CanvasWidth / c.CanvasHeight)

	x -= c.CanvasWidth / 2
	y -= c.CanvasHeight / 2

	vx := x * (viewportWidth / c.CanvasWidth)
	vy := y * (viewportHeight / c.CanvasHeight)

	return Quat{0, vx, vy, c.Distance}
}

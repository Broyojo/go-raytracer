package ray

import (
	"image/color"
	"math"
)

type Ray struct {
	Origin    Quat
	Direction Quat
}

func (r Ray) CalculatePoint(t float64) Quat {
	return Add(r.Origin, Scale(r.Direction, t))
}

type Mat struct {
	Color        color.RGBA64
	Reflectivity float64
	Specular     float64
}

type Light struct {
	LightType int
	Intensity float64
	Position  Quat
	Direction Quat
}

// Sphere
type Sphere struct {
	Center   Quat
	Radius   float64
	Material Mat
}

func (s Sphere) GetT(ray Ray) (float64, float64) {
	oc := Sub(ray.Origin, s.Center)
	a := Dot(ray.Direction, ray.Direction)
	b := Dot(oc, ray.Direction) * 2
	c := Dot(oc, oc) - s.Radius*s.Radius
	d := b*b - (4 * a * c)
	if d > 0 {
		s := math.Pow(d, 0.5)
		h := 2 * a
		t0 := (-b + s) / h
		t1 := (-b - s) / h
		return t0, t1
	}
	return -1.0, -1.0
}

func (s Sphere) GetNormal(pos Quat) Quat {
	return Sub(pos, s.Center).Norm()
}

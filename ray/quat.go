package ray

import (
	"math"
)

type Quat struct {
	W float64
	X float64
	Y float64
	Z float64
}

func Add(qs ...Quat) Quat {
	for _, v := range qs[1:] {
		qs[0].W += v.W
		qs[0].X += v.X
		qs[0].Y += v.Y
		qs[0].Z += v.Z
	}
	return qs[0]
}

func Sub(qs ...Quat) Quat {
	for _, v := range qs[1:] {
		qs[0].W -= v.W
		qs[0].X -= v.X
		qs[0].Y -= v.Y
		qs[0].Z -= v.Z
	}
	return qs[0]
}

func Mul(qs ...Quat) Quat {
	for _, v := range qs {
		w := qs[0].W*v.W - qs[0].X*v.X - qs[0].Y*v.Y - qs[0].Z*v.Z
		x := qs[0].W*v.X + qs[0].X*v.W + qs[0].Y*v.Z - qs[0].Z*v.Y
		y := qs[0].W*v.Y + qs[0].Y*v.W + qs[0].Z*v.X - qs[0].X*v.Z
		z := qs[0].W*v.Z + qs[0].Z*v.W + qs[0].X*v.Y - qs[0].Y*v.X
		qs[0] = Quat{w, x, y, z}
	}
	return qs[0]
}

func Scale(q Quat, n float64) Quat {
	return Quat{q.W * n, q.X * n, q.Y * n, q.Z * n}
}

func Dot(q Quat, q0 Quat) float64 {
	return q.W*q0.W + q.X*q0.X + q.Y*q0.Y + q.Z*q0.Z
}

func (q Quat) Len() float64 {
	return math.Pow(Dot(q, q), 0.5)
}

func (q Quat) Norm() Quat {
	l := q.Len()
	return Quat{q.W / l, q.X / l, q.Y / l, q.Z / l}
}

func (q Quat) Conj() Quat {
	return Quat{q.W, -q.X, -q.Y, -q.Z}
}

func Euler(pitch float64, roll float64, yaw float64) Quat {
	w := math.Cos(roll/2)*math.Cos(pitch/2)*math.Cos(yaw/2) + math.Sin(roll/2)*math.Sin(pitch/2)*math.Sin(yaw/2)
	x := math.Sin(roll/2)*math.Cos(pitch/2)*math.Cos(yaw/2) - math.Cos(roll/2)*math.Sin(pitch/2)*math.Sin(yaw/2)
	y := math.Cos(roll/2)*math.Sin(pitch/2)*math.Cos(yaw/2) + math.Sin(roll/2)*math.Cos(pitch/2)*math.Sin(yaw/2)
	z := math.Cos(roll/2)*math.Cos(pitch/2)*math.Sin(yaw/2) - math.Sin(roll/2)*math.Sin(pitch/2)*math.Cos(yaw/2)
	return Quat{w, x, y, z}
}

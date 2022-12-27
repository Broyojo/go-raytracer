package main

import (
	"image/color"

	"github.com/broyojo/raytracer-go/ray"
)

var camera = ray.Camera{
	Position:     ray.Quat{W: 0, X: 0, Y: 0, Z: 0},
	Direction:    ray.Euler(0, 0, 0),
	CanvasWidth:  1024,
	CanvasHeight: 1024,
	TMin:         0.1,
	TMax:         100000,
	Distance:     1,
	Fov:          45,
}

var scene = ray.Scene{
	Camera: camera,
	Objects: []ray.Sphere{
		{
			Center: ray.Quat{W: 0, X: 0, Y: -1, Z: 3},
			Radius: 1,
			Material: ray.Mat{
				Color:        color.RGBA64{R: 65535, G: 0, B: 0, A: 65535},
				Reflectivity: 0.5,
				Specular:     625,
			},
		},
		{
			Center: ray.Quat{W: 0, X: -2, Y: 1, Z: 3},
			Radius: 1,
			Material: ray.Mat{
				Color:        color.RGBA64{R: 0, G: 0, B: 65535, A: 65535},
				Reflectivity: 0.7,
				Specular:     200,
			},
		},
		{
			Center: ray.Quat{W: 0, X: 2, Y: 1, Z: 3},
			Radius: 1,
			Material: ray.Mat{
				Color:        color.RGBA64{R: 0, G: 65535, B: 0, A: 65535},
				Reflectivity: 0.2,
				Specular:     845,
			},
		},
		{
			Center: ray.Quat{W: 0, X: 0, Y: -5001, Z: 0},
			Radius: 5000,
			Material: ray.Mat{
				Color:        color.RGBA64{R: 65535, G: 65535, B: 0, A: 65535},
				Reflectivity: 0.6,
				Specular:     100,
			},
		},
	},
	Lights: []ray.Light{
		{
			LightType: 0,
			Intensity: 0.3,
		},
		{
			LightType: 1,
			Intensity: 0.2,
			Position:  ray.Quat{W: 0, X: 2, Y: 1, Z: 0},
		},
		{
			LightType: 2,
			Intensity: 0.5,
			Direction: ray.Quat{W: 0, X: 1, Y: 4, Z: 4},
		},
	},
	BackgroundColor: color.RGBA64{R: 0, G: 0, B: 0, A: 65535},
}

func main() {
	scene.Render("balls/balls", 100, 10)
}

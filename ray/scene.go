package ray

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Scene struct {
	Camera          Camera
	Objects         []Sphere
	Lights          []Light
	BackgroundColor color.RGBA64
}

type closest struct {
	closestT      float64
	closestObject Sphere
	yes           bool
}

func (s Scene) RayTrace(ray Ray, depth int) color.RGBA64 {
	closest := s.CalculateIntersection(ray)

	if !closest.yes {
		return s.BackgroundColor
	}

	intersect := ray.CalculatePoint(closest.closestT)

	normal := closest.closestObject.GetNormal(intersect)

	direction := ray.Direction.Conj()

	localColor := scaleColor(closest.closestObject.Material.Color, s.Shade(intersect, normal, direction, closest.closestObject.Material.Specular))

	reflectivity := closest.closestObject.Material.Reflectivity

	if depth <= 0 || reflectivity <= 0 {
		return localColor
	}

	reflectionDirection := Reflect(direction, normal)
	reflectionColor := s.RayTrace(Ray{intersect, reflectionDirection}, depth-1)

	return addColor(scaleColor(localColor, 1-reflectivity), scaleColor(reflectionColor, reflectivity))
}

func scaleColor(mc color.RGBA64, s float64) color.RGBA64 {
	return color.RGBA64{constrain(float64(mc.R) * s), constrain(float64(mc.G) * s), constrain(float64(mc.B) * s), 65535}
}

func addColor(c0 color.RGBA64, c1 color.RGBA64) color.RGBA64 {
	return color.RGBA64{uint16(c0.R + c1.R), uint16(c0.G + c1.G), uint16(c0.B + c1.B), 65535}
}

func constrain(c float64) uint16 {
	if c > 65535.0 {
		c = 65535.0
	}
	return uint16(c)
}

func (s Scene) CalculateIntersection(ray Ray) closest {
	closestT := s.Camera.TMax + 1
	var closestObject Sphere
	yes := false

	for _, object := range s.Objects {
		t0, t1 := object.GetT(ray)

		if s.Camera.TMin <= t0 && t0 <= s.Camera.TMax && t0 < closestT {
			closestT = t0
			closestObject = object
			yes = true
		}
		if s.Camera.TMin <= t1 && t1 <= s.Camera.TMax && t1 < closestT {
			closestT = t1
			closestObject = object
			yes = true
		}
	}
	return closest{closestT, closestObject, yes}
}

func (s Scene) Shade(point Quat, normal Quat, vector Quat, specular float64) float64 {
	intensity := 0.0
	var direction Quat
	for _, light := range s.Lights {
		if light.LightType == 0 {
			intensity += light.Intensity
		} else {
			if light.LightType == 1 {
				direction = Sub(light.Position, point)
			} else {
				direction = light.Direction
			}

			object := s.CalculateIntersection(Ray{point, direction})

			if object.yes {
				continue
			}

			nd := Dot(normal, direction)

			if nd > 0.0 {
				intensity += light.Intensity * nd / (normal.Len() * direction.Len())
			}

			if specular != -1.0 {
				r := Reflect(direction, normal)
				rv := Dot(r, vector)
				if rv > 0.0 {
					intensity += light.Intensity * math.Pow(rv/(r.Len()*vector.Len()), specular)
				}
			}
		}
	}
	return intensity
}

func Reflect(direction Quat, normal Quat) Quat {
	return Sub(Scale(Scale(normal, 2), Dot(normal, direction)), direction)
}

func (s Scene) Render(file_name string, n_frames int, depth int) {
	total := time.Now()
	count := 0
	n := runtime.NumCPU()
	log.Printf("running with %d cores", n)
	for count < n_frames {
		t := time.Now()
		img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(s.Camera.CanvasWidth), int(s.Camera.CanvasHeight)}})
		ch := make(chan func())
		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for f := range ch {
					f()
				}
			}()
		}
		for x := 0.0; x < s.Camera.CanvasWidth; x++ {
			for y := 0.0; y < s.Camera.CanvasHeight; y++ {
				ch <- s.work(x, y, img, depth)
			}
		}
		close(ch)
		wg.Wait()
		log.Printf("done rendering frame %v in %v", count, time.Since(t))
		s.Camera.Position.Z = calculatePath(s.Camera.Position.X)
		s.Camera.Position.X += 0.1
		s.Camera.Rotate(Euler(0.00001, 0, 0))
		f, _ := os.Create(file_name + "-" + strconv.Itoa(count) + ".png")
		png.Encode(f, img)
		count++
	}

	/*
		img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(s.Camera.CanvasWidth), int(s.Camera.CanvasHeight)}})
		count := 0
		for count < n_frames {
			t := time.Now()
			for x := 0.0; x < s.Camera.CanvasWidth; x++ {
				for y := 0.0; y < s.Camera.CanvasHeight; y++ {
					origin := Mul(s.Camera.Direction, s.Camera.TranslateCoords(x, y))
					color := s.RayTrace(Ray{origin, Sub(origin, s.Camera.Position).Norm()}, depth)
					img.Set(int(x), int(y), color)
				}
			}
			f, _ := os.Create(file_name + "-" + strconv.Itoa(count) + ".png")
			png.Encode(f, img)
			update(time.Since(t))

			s.Camera.Rotate(Euler(0.01, 0, 0))

			count++
		}
	*/
	log.Printf("Took %v to render %v frames", time.Since(total), n_frames)
}

func (s Scene) work(x float64, y float64, img *image.RGBA, depth int) func() {
	return func() {
		origin := Mul(s.Camera.Direction, s.Camera.TranslateCoords(x, y))
		color := s.RayTrace(Ray{origin, Sub(origin, s.Camera.Position).Norm()}, depth)
		img.Set(int(x), int(y), color)
	}
}

func calculatePath(x float64) float64 {
	z := math.Pow(9-10*x*x, 0.5) / math.Pow(10, 0.5)
	if x >= 0 {
		return z
	}
	return -z
}

package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 640

	// G is the gravitational constant
	G = float64(6.67408e-11)
	c = float64(299792458)
)

var (
	ps = newSunMercurySystem()

	img *image.RGBA
)

// Particle is a struct containing all information about a physics particle
type Particle struct {
	x  float64
	y  float64
	vx float64
	vy float64
	m  float64
}

// ParticleSystem is a wrapper for a system of particles
type ParticleSystem struct {
	particles []*Particle
	num       int
}

func newParticleSystem() ParticleSystem {
	var s []*Particle
	for i := 0; i < 500; i++ {
		s = append(s, &Particle{(rand.Float64() - 0.5) * 8e11, (rand.Float64() - 0.5) * 8e11, 0, 0, math.Pow(10, 26+rand.Float64()*2)})
	}
	return ParticleSystem{s, 500}
}

func newSunEarthSystem() ParticleSystem {
	var s []*Particle
	s = append(s, &Particle{0, 0, 0, 0, 2e30})
	s = append(s, &Particle{1.5e11, 0, 0, 1e4, 6e24})
	return ParticleSystem{s, 2}
}
func newSunMercurySystem() ParticleSystem {
	var s []*Particle
	s = append(s, &Particle{0, 0, 0, 0, 2e30})
	s = append(s, &Particle{6.98e10, 0, 0, 3.9e4, 3e23})
	return ParticleSystem{s, 2}
}

func newBinarySystem() ParticleSystem {
	var s []*Particle
	s = append(s, &Particle{-1.5e11, 0, 0, -1e4, 2e30})
	s = append(s, &Particle{1.5e11, 0, 0, 1e4, 2e30})
	s = append(s, &Particle{0, 0, -5e3, 1e4, 2e24})
	return ParticleSystem{s, 2}
}

// Update does all calculations to update a single particle
func (p *Particle) Update(pid int, s *ParticleSystem, deltaT float64) Particle {
	sumX, sumY := float64(0), float64(0)
	for i, x := range s.particles {
		if i != pid {
			bx := (p.x - x.x) * x.m / (p.m + x.m)
			by := (x.y - p.y) * x.m / (p.m + x.m)
			factor := (G * x.m * p.m) / p.DistSquare(x)
			var angle float64
			//if (x.x - p.x) < 0 {
			angle = math.Atan2((x.y - p.y), (p.x - x.x))
			/* } else {
				angle = math.Pi + math.Tanh((x.y-p.y)/(x.x-p.x))
			} */
			if pid == 1 {

			}
			//log.Println(x.y, p.y, x.x, p.x)
			sumX -= math.Cos(angle) * factor
			sumY += math.Sin(angle) * factor
		}
	}

	// log.Printf("sumX: %v, sumY: %v\n", sumX, sumY)
	q := Particle{p.x, p.y, p.vx, p.vy, p.m}
	//log.Printf("q: %+v", q)
	q.vx += sumX / (q.m) * deltaT
	q.vy += sumY / (q.m) * deltaT
	q.x += q.vx * deltaT
	q.y += q.vy * deltaT
	return q
}

// DistSquare returns the distance between two particles squared
func (p *Particle) DistSquare(x *Particle) float64 {
	return math.Pow(p.x-x.x, 2) + math.Pow(p.y-x.y, 2)
}
func (p *Particle) Dist(x *Particle) float64 {
	return math.Sqrt(math.Pow(p.x-x.x, 2) + math.Pow(p.y-x.y, 2))
}

func update(screen *ebiten.Image) error {
	// img = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	newPs := ps
	if ebiten.IsRunningSlowly() {
		return nil
	}
	for i, x := range ps.particles {
		//log.Printf("p[%v] \t %v, %v", i, x.x, x.y)
		newX := x.Update(i, &ps, 3e3)
		newPs.particles[i] = &newX

		if x.x >= -4e11 && x.x <= 4e11 && x.y >= -4e11 && x.y <= 4e11 {
			rX := (float64(x.x/8e11) + .5) * screenWidth
			rY := (float64(x.y/8e11) + .5) * screenHeight
			pos := 4*int(rY)*screenWidth + 4*int(rX)
			img.Pix[pos] = 0xff
			img.Pix[pos+1] = 0xff
			img.Pix[pos+2] = 0xff
			img.Pix[pos+3] = 0xff
		}
	}
	screen.ReplacePixels(img.Pix)
	d := math.Sqrt((ps.particles[1].y-ps.particles[0].y)*(ps.particles[1].y-ps.particles[0].y) + (ps.particles[0].x-ps.particles[1].x)*(ps.particles[0].x-ps.particles[1].x))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f\nr: %e", ebiten.CurrentFPS(), d))
	ps = newPs
	return nil
}

func main() {
	img = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "N-Body Simulation"); err != nil {
		log.Fatal(err)
	}
}

package noise

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	permutationSize = 256
	minStretch      = 0.1
	maxStretch      = 100
)

var gradients = [8][2]float64{
	{1, 0},
	{math.Sqrt2 / 2, math.Sqrt2 / 2},
	{0, 1},
	{-math.Sqrt2 / 2, math.Sqrt2 / 2},
	{-1, 0},
	{-math.Sqrt2 / 2, -math.Sqrt2 / 2},
	{0, -1},
	{math.Sqrt2 / 2, -math.Sqrt2 / 2},
}

// Sampler generates deterministic 2D gradient Perlin noise.
// Its permutation table and directional transform are immutable after creation,
// so a Sampler can be shared safely between goroutines.
type Sampler struct {
	permutation [permutationSize * 2]uint8
	cosAngle    float64
	sinAngle    float64
	invStretch  float64
}

// New creates a noise sampler for seed and directional sampling settings.
// A stretch greater than 1 lengthens noise features along the given angle.
func New(seed int64, directionAngleDegrees, directionStretch float64) (*Sampler, error) {
	if math.IsNaN(directionAngleDegrees) || math.IsInf(directionAngleDegrees, 0) {
		return nil, fmt.Errorf("direction angle must be a finite number")
	}
	if math.IsNaN(directionStretch) || math.IsInf(directionStretch, 0) ||
		directionStretch < minStretch || directionStretch > maxStretch {
		return nil, fmt.Errorf("direction stretch must be between %.1f and %d", minStretch, maxStretch)
	}

	permutation := [permutationSize]uint8{}
	for i := range permutation {
		permutation[i] = uint8(i)
	}
	random := rand.New(rand.NewSource(seed))
	for i := len(permutation) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		permutation[i], permutation[j] = permutation[j], permutation[i]
	}

	sampler := &Sampler{invStretch: 1 / directionStretch}
	for i, value := range permutation {
		sampler.permutation[i] = value
		sampler.permutation[i+permutationSize] = value
	}
	angle := math.Mod(directionAngleDegrees, 360) * (math.Pi / 180)
	sampler.cosAngle = math.Cos(angle)
	sampler.sinAngle = math.Sin(angle)
	return sampler, nil
}

// Sample returns the noise value at the given noise-space coordinates.
// Values are in the range [-1, 1]. The directional transform is applied before
// sampling, leaving the underlying Perlin noise algorithm unchanged.
func (s *Sampler) Sample(x, y float64) float64 {
	along := (x*s.cosAngle + y*s.sinAngle) * s.invStretch
	across := -x*s.sinAngle + y*s.cosAngle
	x, y = along, across

	xFloor, yFloor := math.Floor(x), math.Floor(y)
	x0, y0 := int(xFloor)&255, int(yFloor)&255
	dx, dy := x-xFloor, y-yFloor
	x1, y1 := (x0+1)&255, (y0+1)&255

	aa := s.hash(x0, y0)
	ab := s.hash(x0, y1)
	ba := s.hash(x1, y0)
	bb := s.hash(x1, y1)

	lower := lerp(dotGradient(aa, dx, dy), dotGradient(ba, dx-1, dy), fade(dx))
	upper := lerp(dotGradient(ab, dx, dy-1), dotGradient(bb, dx-1, dy-1), fade(dx))
	return lerp(lower, upper, fade(dy))
}

func (s *Sampler) hash(x, y int) uint8 {
	return s.permutation[(int(s.permutation[x])+y)&255]
}

func dotGradient(hash uint8, x, y float64) float64 {
	gradient := gradients[hash&7]
	return gradient[0]*x + gradient[1]*y
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

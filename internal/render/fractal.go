package render

import "github.com/austintaylor/perlin-noise-go/internal/noise"

// Fractal noise adds progressively smaller features to the broad shapes.
func fractalField(sampler *noise.Sampler, x, y float64) float64 {
	value, amplitude := 0.0, 1.0
	for octave := 0; octave < 7; octave++ {
		value += amplitude * sampler.Sample(x, y)
		x, y = x*2.03+17.1, y*2.03+9.2
		amplitude *= 0.5
	}
	return value
}

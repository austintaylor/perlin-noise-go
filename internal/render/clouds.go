package render

import (
	"image/color"
	"math"

	"github.com/austintaylor/perlin-noise-go/internal/noise"
)

// Fractal noise adds progressively smaller features to the broad cloud shapes.
func cloudField(sampler *noise.Sampler, x, y float64) float64 {
	value, amplitude := 0.0, 1.0
	for octave := 0; octave < 7; octave++ {
		value += amplitude * sampler.Sample(x, y)
		x, y = x*2.03+17.1, y*2.03+9.2
		amplitude *= 0.5
	}
	return value
}

func cloudColor(sampler *noise.Sampler, x, y, depth float64) color.NRGBA {
	field := cloudField(sampler, x, y)
	density := max(0, field+0.12)
	opacity := 1 - math.Exp(-density*7)
	// A nearby sample toward the light estimates how much cloud blocks it.
	occlusion := max(0, cloudField(sampler, x-0.06, y-0.09)-field)
	light := 0.68 + 0.32*math.Exp(-occlusion*18)
	sky := [3]float64{38 + depth*52, 112 + depth*45, 193 + depth*29}
	cloud := [3]float64{255 * light, 253 * light, 248 * light}
	mix := func(channel int) uint8 {
		return uint8(math.Round(sky[channel]*(1-opacity) + cloud[channel]*opacity))
	}
	return color.NRGBA{R: mix(0), G: mix(1), B: mix(2), A: 255}
}

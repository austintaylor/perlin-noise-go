package palette

import (
	"fmt"
	"image/color"
	"math"
	"sort"
)

// Stop pairs a position in the normalized palette range [0, 1] with a color.
type Stop struct {
	Position float64
	Color    color.NRGBA
}

// Palette maps normalized values to colors by interpolating between stops.
// Its stops are copied on construction and remain immutable.
type Palette struct {
	stops []Stop
}

// New creates a palette from stops at strictly increasing positions.
// Stops must include both endpoints: 0 and 1.
func New(stops []Stop) (Palette, error) {
	if len(stops) < 2 {
		return Palette{}, fmt.Errorf("palette requires at least two stops")
	}
	copied := append([]Stop(nil), stops...)
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Position < copied[j].Position
	})
	if copied[0].Position != 0 || copied[len(copied)-1].Position != 1 {
		return Palette{}, fmt.Errorf("palette stops must include positions 0 and 1")
	}
	for i, stop := range copied {
		if math.IsNaN(stop.Position) || math.IsInf(stop.Position, 0) ||
			stop.Position < 0 || stop.Position > 1 {
			return Palette{}, fmt.Errorf("palette stop position must be finite and between 0 and 1")
		}
		if i > 0 && copied[i-1].Position == stop.Position {
			return Palette{}, fmt.Errorf("palette stop positions must be unique")
		}
	}
	return Palette{stops: copied}, nil
}

// Preset returns a named built-in palette.
func Preset(name string) (Palette, error) {
	var stops []Stop
	switch name {
	case "clouds":
		stops = []Stop{
			{Position: 0, Color: color.NRGBA{R: 12, G: 20, B: 38, A: 255}},
			{Position: 0.3, Color: color.NRGBA{R: 48, G: 72, B: 105, A: 255}},
			{Position: 0.62, Color: color.NRGBA{R: 132, G: 163, B: 190, A: 255}},
			{Position: 0.84, Color: color.NRGBA{R: 218, G: 225, B: 226, A: 255}},
			{Position: 1, Color: color.NRGBA{R: 255, G: 250, B: 232, A: 255}},
		}
	case "fire":
		stops = []Stop{
			{Position: 0, Color: color.NRGBA{R: 24, G: 8, B: 18, A: 255}},
			{Position: 0.28, Color: color.NRGBA{R: 112, G: 16, B: 12, A: 255}},
			{Position: 0.55, Color: color.NRGBA{R: 220, G: 48, B: 8, A: 255}},
			{Position: 0.78, Color: color.NRGBA{R: 255, G: 152, B: 12, A: 255}},
			{Position: 1, Color: color.NRGBA{R: 255, G: 248, B: 176, A: 255}},
		}
	default:
		return Palette{}, fmt.Errorf("unknown palette %q (choose clouds or fire)", name)
	}
	return New(stops)
}

// ColorAt returns the interpolated color for normalized position [0, 1].
// Values outside the range are clamped to the nearest endpoint.
func (p Palette) ColorAt(position float64) color.NRGBA {
	if position <= 0 {
		return p.stops[0].Color
	}
	if position >= 1 {
		return p.stops[len(p.stops)-1].Color
	}
	index := sort.Search(len(p.stops), func(i int) bool {
		return p.stops[i].Position >= position
	})
	left, right := p.stops[index-1], p.stops[index]
	fraction := (position - left.Position) / (right.Position - left.Position)
	return color.NRGBA{
		R: interpolate(left.Color.R, right.Color.R, fraction),
		G: interpolate(left.Color.G, right.Color.G, fraction),
		B: interpolate(left.Color.B, right.Color.B, fraction),
		A: interpolate(left.Color.A, right.Color.A, fraction),
	}
}

// ColorForNoise maps a Perlin noise value from [-1, 1] onto this palette.
// Values beyond that range are clamped to the endpoint colors.
func (p Palette) ColorForNoise(value float64) color.NRGBA {
	return p.ColorAt((value + 1) / 2)
}

func interpolate(a, b uint8, fraction float64) uint8 {
	return uint8(math.Round(float64(a) + (float64(b)-float64(a))*fraction))
}

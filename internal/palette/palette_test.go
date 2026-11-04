package palette

import (
	"image/color"
	"math"
	"testing"
)

func TestColorAtInterpolatesStopsAndEndpoints(t *testing.T) {
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	p, err := New([]Stop{
		{Position: 0, Color: black},
		{Position: 1, Color: white},
	})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		position float64
		want     color.NRGBA
	}{
		{position: -1, want: black},
		{position: 0, want: black},
		{position: 0.5, want: color.NRGBA{R: 128, G: 128, B: 128, A: 128}},
		{position: 1, want: white},
		{position: 2, want: white},
	}
	for _, tt := range tests {
		if got := p.ColorAt(tt.position); got != tt.want {
			t.Errorf("ColorAt(%v) = %v, want %v", tt.position, got, tt.want)
		}
	}
}

func TestColorAtUsesIntermediateStop(t *testing.T) {
	red := color.NRGBA{R: 255, A: 255}
	green := color.NRGBA{G: 255, A: 255}
	blue := color.NRGBA{B: 255, A: 255}
	p, err := New([]Stop{
		{Position: 1, Color: blue},
		{Position: 0, Color: red},
		{Position: 0.5, Color: green},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := p.ColorAt(0.5); got != green {
		t.Errorf("ColorAt(0.5) = %v, want intermediate stop %v", got, green)
	}
	if got, want := p.ColorAt(0.25), (color.NRGBA{R: 128, G: 128, A: 255}); got != want {
		t.Errorf("ColorAt(0.25) = %v, want %v", got, want)
	}
}

func TestColorForNoiseMapsSamplerRange(t *testing.T) {
	start := color.NRGBA{R: 10, A: 255}
	end := color.NRGBA{B: 200, A: 255}
	p, err := New([]Stop{
		{Position: 0, Color: start},
		{Position: 1, Color: end},
	})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		value float64
		want  color.NRGBA
	}{
		{value: -2, want: start},
		{value: -1, want: start},
		{value: 0, want: color.NRGBA{R: 5, B: 100, A: 255}},
		{value: 1, want: end},
		{value: 2, want: end},
	}
	for _, tt := range tests {
		if got := p.ColorForNoise(tt.value); got != tt.want {
			t.Errorf("ColorForNoise(%v) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestPresets(t *testing.T) {
	tests := []struct {
		name string
		low  color.NRGBA
		high color.NRGBA
	}{
		{
			name: "clouds",
			low:  color.NRGBA{R: 12, G: 20, B: 38, A: 255},
			high: color.NRGBA{R: 255, G: 250, B: 232, A: 255},
		},
		{
			name: "fire",
			low:  color.NRGBA{R: 24, G: 8, B: 18, A: 255},
			high: color.NRGBA{R: 255, G: 248, B: 176, A: 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := Preset(tt.name)
			if err != nil {
				t.Fatal(err)
			}
			if got := p.ColorAt(0); got != tt.low {
				t.Errorf("low endpoint = %v, want %v", got, tt.low)
			}
			if got := p.ColorAt(1); got != tt.high {
				t.Errorf("high endpoint = %v, want %v", got, tt.high)
			}
			if middle := p.ColorAt(0.5); middle == tt.low || middle == tt.high {
				t.Errorf("middle color should interpolate between endpoints, got %v", middle)
			}
		})
	}
}

func TestPresetRejectsUnknownName(t *testing.T) {
	if _, err := Preset("lava"); err == nil {
		t.Fatal("expected unknown palette to be rejected")
	}
}

func TestNewRejectsInvalidStops(t *testing.T) {
	valid := color.NRGBA{A: 255}
	tests := []struct {
		name  string
		stops []Stop
	}{
		{name: "too few", stops: []Stop{{Position: 0, Color: valid}}},
		{name: "missing zero endpoint", stops: []Stop{{Position: 0.1, Color: valid}, {Position: 1, Color: valid}}},
		{name: "missing one endpoint", stops: []Stop{{Position: 0, Color: valid}, {Position: 0.9, Color: valid}}},
		{name: "duplicate position", stops: []Stop{{Position: 0, Color: valid}, {Position: 0.5, Color: valid}, {Position: 0.5, Color: valid}, {Position: 1, Color: valid}}},
		{name: "NaN position", stops: []Stop{{Position: 0, Color: valid}, {Position: math.NaN(), Color: valid}, {Position: 1, Color: valid}}},
		{name: "infinite position", stops: []Stop{{Position: 0, Color: valid}, {Position: math.Inf(1), Color: valid}}},
		{name: "out of range", stops: []Stop{{Position: -0.1, Color: valid}, {Position: 1, Color: valid}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := New(tt.stops); err == nil {
				t.Fatal("expected invalid stops to be rejected")
			}
		})
	}
}

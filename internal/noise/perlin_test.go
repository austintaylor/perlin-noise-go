package noise

import (
	"math"
	"testing"
)

func TestSampleIsRepeatableForSameSeed(t *testing.T) {
	first, err := New(42, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := New(42, 0, 1)
	if err != nil {
		t.Fatal(err)
	}

	points := [][2]float64{{0.25, 0.75}, {-1.3, 5.8}, {12.125, -7.5}}
	for _, point := range points {
		if got, want := first.Sample(point[0], point[1]), second.Sample(point[0], point[1]); got != want {
			t.Errorf("Sample(%v) = %v, want %v", point, got, want)
		}
	}
}

func TestDifferentSeedsChangeSamples(t *testing.T) {
	first, err := New(1, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	second, err := New(2, 0, 1)
	if err != nil {
		t.Fatal(err)
	}

	different := false
	for y := 0.13; y < 5 && !different; y += 0.31 {
		for x := 0.17; x < 5; x += 0.29 {
			if first.Sample(x, y) != second.Sample(x, y) {
				different = true
				break
			}
		}
	}
	if !different {
		t.Fatal("different seeds produced identical samples at all test points")
	}
}

func TestSampleIsBoundedAndZeroAtLatticePoints(t *testing.T) {
	sampler, err := New(1234, 17, 3.5)
	if err != nil {
		t.Fatal(err)
	}
	for y := -10.0; y <= 10; y += 0.17 {
		for x := -10.0; x <= 10; x += 0.19 {
			if got := sampler.Sample(x, y); got < -1 || got > 1 {
				t.Fatalf("Sample(%v, %v) = %v, outside [-1, 1]", x, y, got)
			}
		}
	}

	untransformed, err := New(1234, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	for i := -10; i <= 10; i++ {
		for j := -10; j <= 10; j++ {
			if got := untransformed.Sample(float64(i), float64(j)); got != 0 {
				t.Errorf("Sample(%d, %d) = %v, want 0", i, j, got)
			}
		}
	}
}

func TestDirectionalStretchTransformsSamplingCoordinates(t *testing.T) {
	base, err := New(8, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	horizontal, err := New(8, 0, 3)
	if err != nil {
		t.Fatal(err)
	}
	vertical, err := New(8, 90, 3)
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1.7, -2.3
	if got, want := horizontal.Sample(x, y), base.Sample(x/3, y); math.Abs(got-want) > 1e-14 {
		t.Errorf("horizontal stretch sample = %v, want %v", got, want)
	}
	if got, want := vertical.Sample(x, y), base.Sample(y/3, -x); math.Abs(got-want) > 1e-14 {
		t.Errorf("vertical stretch sample = %v, want %v", got, want)
	}
}

func TestNewRejectsInvalidDirectionSettings(t *testing.T) {
	tests := []struct {
		name  string
		angle float64
		scale float64
	}{
		{name: "NaN angle", angle: math.NaN(), scale: 1},
		{name: "infinite angle", angle: math.Inf(1), scale: 1},
		{name: "NaN stretch", angle: 0, scale: math.NaN()},
		{name: "infinite stretch", angle: 0, scale: math.Inf(1)},
		{name: "stretch below minimum", angle: 0, scale: 0.09},
		{name: "stretch above maximum", angle: 0, scale: 100.1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := New(1, tt.angle, tt.scale); err == nil {
				t.Fatal("expected invalid direction settings to be rejected")
			}
		})
	}
}

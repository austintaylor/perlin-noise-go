package render

import (
	"bytes"
	"math"
	"testing"

	"github.com/austintaylor/perlin-noise-go/internal/noise"
	"github.com/austintaylor/perlin-noise-go/internal/palette"
)

func TestCloudImageIsRepeatableAndHasSkyAndHighlights(t *testing.T) {
	sampler, err := noise.New(42, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	colors, err := palette.Preset("clouds")
	if err != nil {
		t.Fatal(err)
	}
	first, err := Image(160, 100, 4, 3, sampler, colors, Options{Clouds: true})
	if err != nil {
		t.Fatal(err)
	}
	second, err := Image(160, 100, 4, 3, sampler, colors, Options{Clouds: true})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first.Pix, second.Pix) {
		t.Fatal("cloud rendering must be deterministic")
	}
	sky, highlights := 0, 0
	for y := 0; y < 100; y++ {
		for x := 0; x < 160; x++ {
			c := first.NRGBAAt(x, y)
			if c.A != 255 {
				t.Fatal("cloud image must be opaque")
			}
			if int(c.B)-int(c.R) > 80 {
				sky++
			}
			if c.R > 230 && c.G > 230 && c.B > 230 {
				highlights++
			}
		}
	}
	if sky < 1000 || highlights < 1000 {
		t.Fatalf("expected visible sky and bright clouds, got %d sky and %d highlight pixels", sky, highlights)
	}
	otherSampler, err := noise.New(43, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	other, err := Image(160, 100, 4, 3, otherSampler, colors, Options{Clouds: true})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(first.Pix, other.Pix) {
		t.Fatal("changing the seed must change cloud shapes")
	}
	if _, err := Image(10, 10, 4, 3, sampler, colors, Options{Clouds: true, Flame: true}); err == nil {
		t.Fatal("expected conflicting shaping modes to fail")
	}
}

func TestCloudFieldIncludesFineDetail(t *testing.T) {
	sampler, err := noise.New(42, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	x, y := 0.37, 1.21
	value := cloudField(sampler, x, y)
	if math.IsNaN(value) || math.IsInf(value, 0) {
		t.Fatal("cloud field must be finite")
	}
	if math.Abs(value-sampler.Sample(x, y)) < 0.001 {
		t.Fatal("cloud field must include detail beyond the first octave")
	}
}

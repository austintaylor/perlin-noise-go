package render

import (
	"bytes"
	"math"
	"testing"

	"github.com/austintaylor/perlin-noise-go/internal/noise"
	"github.com/austintaylor/perlin-noise-go/internal/palette"
)

func TestCloudImageIsRepeatableAndMapsLayeredNoise(t *testing.T) {
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
	darkest, brightest := uint8(255), uint8(0)
	for y := 0; y < 100; y++ {
		for x := 0; x < 160; x++ {
			c := first.NRGBAAt(x, y)
			if c.A != 255 {
				t.Fatal("cloud image must be opaque")
			}
			want := colors.ColorForNoise(1.7 * cloudField(sampler, float64(x)*4/160, float64(y)*3/100))
			if c != want {
				t.Fatalf("pixel (%d, %d) = %v, want direct noise mapping %v", x, y, c, want)
			}
			darkest = min(darkest, c.R)
			brightest = max(brightest, c.R)
		}
	}
	if int(brightest)-int(darkest) < 150 {
		t.Fatalf("expected broad tonal variation, got red-channel range %d..%d", darkest, brightest)
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

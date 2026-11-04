package render

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/austintaylor/perlin-noise-go/internal/noise"
	"github.com/austintaylor/perlin-noise-go/internal/palette"
)

func TestImageDimensionsAndPNGEncoding(t *testing.T) {
	img := renderTestImage(t, 37, 23, 5, 3, 42)
	if got, want := img.Bounds(), image.Rect(0, 0, 37, 23); got != want {
		t.Fatalf("image bounds = %v, want %v", got, want)
	}

	var encoded bytes.Buffer
	if err := png.Encode(&encoded, img); err != nil {
		t.Fatalf("encode PNG: %v", err)
	}
	decoded, err := png.Decode(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		t.Fatalf("decode PNG: %v", err)
	}
	if got, want := decoded.Bounds(), img.Bounds(); got != want {
		t.Fatalf("decoded image bounds = %v, want %v", got, want)
	}
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			got := color.NRGBAModel.Convert(decoded.At(x, y)).(color.NRGBA)
			if want := img.NRGBAAt(x, y); got != want {
				t.Fatalf("decoded pixel (%d, %d) = %v, want %v", x, y, got, want)
			}
		}
	}
}

func TestImageIsDeterministic(t *testing.T) {
	first := renderTestImage(t, 31, 19, 7, 4, 17)
	second := renderTestImage(t, 31, 19, 7, 4, 17)
	for y := 0; y < first.Bounds().Dy(); y++ {
		for x := 0; x < first.Bounds().Dx(); x++ {
			if got, want := first.NRGBAAt(x, y), second.NRGBAAt(x, y); got != want {
				t.Fatalf("pixel (%d, %d) differs: %v != %v", x, y, got, want)
			}
		}
	}
}

func TestFlameImageHasHotBaseAndCoolTop(t *testing.T) {
	sampler, err := noise.New(42, 90, 4)
	if err != nil {
		t.Fatal(err)
	}
	colors, err := palette.Preset("fire")
	if err != nil {
		t.Fatal(err)
	}
	for _, height := range []int{1, 2, 101} {
		img, err := Image(97, height, 14, 10, sampler, colors, Options{Flame: true})
		if err != nil {
			t.Fatal(err)
		}
		for x := 0; x < 97; x++ {
			if got, want := img.NRGBAAt(x, height-1), colors.ColorAt(1); got != want {
				t.Fatalf("height %d: base pixel %d = %v, want %v", height, x, got, want)
			}
			if height > 1 {
				if got, want := img.NRGBAAt(x, 0), colors.ColorAt(0); got != want {
					t.Fatalf("height %d: top pixel %d = %v, want %v", height, x, got, want)
				}
			}
		}
	}
}

func TestFlameValueFalloffAndNoise(t *testing.T) {
	for y := 0; y < 101; y++ {
		want := 2*float64(y)/100 - 1
		if got := flameValue(0, y, 101); got != want {
			t.Fatalf("zero-noise value at row %d = %v, want %v", y, got, want)
		}
		for _, noiseValue := range []float64{-1, -0.5, 0, 0.5, 1} {
			if got := flameValue(noiseValue, y, 101); got < -1 || got > 1 {
				t.Fatalf("value at row %d = %v, outside [-1, 1]", y, got)
			}
		}
	}
	if flameValue(-0.5, 50, 101) >= flameValue(0.5, 50, 101) {
		t.Fatal("noise must influence intensity in the flame interior")
	}
}

func TestImageMatchesUnshapedSampler(t *testing.T) {
	img := renderTestImage(t, 31, 19, 7, 4, 17)
	sampler, err := noise.New(17, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	colors, err := palette.Preset("fire")
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < 19; y++ {
		for x := 0; x < 31; x++ {
			want := colors.ColorForNoise(sampler.Sample(float64(x)*7/31, float64(y)*4/19))
			if got := img.NRGBAAt(x, y); got != want {
				t.Fatalf("unshaped pixel (%d, %d) = %v, want %v", x, y, got, want)
			}
		}
	}
}

func TestImageRejectsInvalidInputs(t *testing.T) {
	sampler, err := noise.New(1, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	colors, err := palette.Preset("clouds")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		width      int
		height     int
		gridWidth  int
		gridHeight int
		sampler    *noise.Sampler
	}{
		{name: "zero width", width: 0, height: 10, gridWidth: 1, gridHeight: 1, sampler: sampler},
		{name: "negative height", width: 10, height: -1, gridWidth: 1, gridHeight: 1, sampler: sampler},
		{name: "too many pixels", width: maxPixels, height: 2, gridWidth: 1, gridHeight: 1, sampler: sampler},
		{name: "zero grid", width: 10, height: 10, gridWidth: 0, gridHeight: 1, sampler: sampler},
		{name: "nil sampler", width: 10, height: 10, gridWidth: 1, gridHeight: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Image(tt.width, tt.height, tt.gridWidth, tt.gridHeight, tt.sampler, colors, Options{}); err == nil {
				t.Fatal("expected invalid input to be rejected")
			}
		})
	}
	if _, err := Image(10, 10, 1, 1, sampler, palette.Palette{}, Options{}); err == nil {
		t.Fatal("expected an uninitialized palette to be rejected")
	}
}

func renderTestImage(t *testing.T, width, height, gridWidth, gridHeight int, seed int64) *image.NRGBA {
	t.Helper()
	sampler, err := noise.New(seed, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	colors, err := palette.Preset("fire")
	if err != nil {
		t.Fatal(err)
	}
	img, err := Image(width, height, gridWidth, gridHeight, sampler, colors, Options{})
	if err != nil {
		t.Fatal(err)
	}
	return img
}

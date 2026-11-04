package main

import (
	"bytes"
	"errors"
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/austintaylor/perlin-noise-go/internal/palette"
)

func TestRunWritesShapedFlamePNG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "flame.png")
	if err := run([]string{"-output", path, "-width", "23", "-height", "31",
		"-palette", "fire", "-flame", "-direction-angle", "90", "-direction-stretch", "4"},
		&bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	colors, err := palette.Preset("fire")
	if err != nil {
		t.Fatal(err)
	}
	for x := 0; x < 23; x++ {
		if got := color.NRGBAModel.Convert(img.At(x, 0)); got != colors.ColorAt(0) {
			t.Fatalf("top pixel %d = %v, want coolest palette color", x, got)
		}
		if got := color.NRGBAModel.Convert(img.At(x, 30)); got != colors.ColorAt(1) {
			t.Fatalf("bottom pixel %d = %v, want hottest palette color", x, got)
		}
	}
}

func TestRunWritesPNG(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "clouds.png")
	var stdout, stderr bytes.Buffer
	err := run([]string{
		"-output", outputPath,
		"-width", "19", "-height", "13",
		"-grid-width", "4", "-grid-height", "3",
		"-seed", "42",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr: %s)", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), outputPath) {
		t.Errorf("success message %q does not mention output path", stdout.String())
	}
	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("open output: %v", err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		t.Fatalf("decode output PNG: %v", err)
	}
	if got, want := img.Bounds(), image.Rect(0, 0, 19, 13); got != want {
		t.Fatalf("output bounds = %v, want %v", got, want)
	}
}

func TestRunWritesCloudPNG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "clouds.png")
	if err := run([]string{"-output", path, "-clouds", "-width", "80", "-height", "50",
		"-grid-width", "4", "-grid-height", "3"}, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds() != image.Rect(0, 0, 80, 50) {
		t.Fatalf("unexpected cloud image bounds: %v", img.Bounds())
	}
}

func TestRunRejectsUnsupportedDemo(t *testing.T) {
	err := run([]string{"-demo"}, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "demo mode is not implemented yet") {
		t.Fatalf("run error = %v, want clear unsupported demo error", err)
	}
}

func TestRunPropagatesHelp(t *testing.T) {
	err := run([]string{"-help"}, &bytes.Buffer{}, &bytes.Buffer{})
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("run error = %v, want flag.ErrHelp", err)
	}
}

package main

import (
	"bytes"
	"errors"
	"flag"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

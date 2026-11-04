package config

import (
	"io"
	"testing"
)

func TestParseDefaults(t *testing.T) {
	cfg, err := Parse(nil, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Output != "perlin.png" || cfg.Width != 512 || cfg.Height != 512 ||
		cfg.GridWidth != 8 || cfg.GridHeight != 8 || cfg.Seed != 1 ||
		cfg.Palette != "clouds" || cfg.Flame || cfg.Clouds || cfg.DirectionAngle != 0 ||
		cfg.DirectionStretch != 1 || cfg.Demo || cfg.DemoScale != 32 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestParseCloudMode(t *testing.T) {
	cfg, err := Parse([]string{"-clouds"}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Clouds {
		t.Fatal("cloud mode was not enabled")
	}
}

func TestParseOptions(t *testing.T) {
	args := []string{
		"-output", "flames.png",
		"-width", "640", "-height", "480",
		"-grid-width", "20", "-grid-height", "12",
		"-seed", "42", "-palette", "fire", "-flame",
		"-direction-angle", "90", "-direction-stretch", "2.5",
		"-demo", "-demo-scale", "16",
	}
	cfg, err := Parse(args, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Output != "flames.png" || cfg.Width != 640 || cfg.Height != 480 ||
		cfg.GridWidth != 20 || cfg.GridHeight != 12 || cfg.Seed != 42 ||
		cfg.Palette != "fire" || !cfg.Flame || cfg.DirectionAngle != 90 ||
		cfg.DirectionStretch != 2.5 || !cfg.Demo || cfg.DemoScale != 16 {
		t.Fatalf("options not applied: %+v", cfg)
	}
}

func TestParseRejectsInvalidOptions(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "unknown palette", args: []string{"-palette", "lava"}},
		{name: "conflicting modes", args: []string{"-clouds", "-flame"}},
		{name: "clouds with fire palette", args: []string{"-clouds", "-palette", "fire"}},
		{name: "zero width", args: []string{"-width", "0"}},
		{name: "excessive pixel count", args: []string{"-width", "8192", "-height", "8192"}},
		{name: "invalid stretch", args: []string{"-direction-stretch", "0"}},
		{name: "non-finite angle", args: []string{"-direction-angle", "NaN"}},
		{name: "demo exceeds pixel limit", args: []string{"-demo", "-grid-width", "1024", "-grid-height", "1024", "-demo-scale", "32"}},
		{name: "positional argument", args: []string{"extra"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.args, io.Discard); err == nil {
				t.Fatal("expected invalid options to be rejected")
			}
		})
	}
}

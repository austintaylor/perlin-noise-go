package config

import (
	"flag"
	"fmt"
	"io"
	"math"
)

const (
	maxDimension = 8192
	maxPixels    = 64_000_000
)

// Config contains the user-selectable options for generating an image.
type Config struct {
	Output           string
	Width            int
	Height           int
	GridWidth        int
	GridHeight       int
	Seed             int64
	Palette          string
	Flame            bool
	Clouds           bool
	DirectionAngle   float64
	DirectionStretch float64
	Demo             bool
	DemoScale        int
}

// Parse parses CLI arguments and validates the resulting configuration.
func Parse(args []string, output io.Writer) (Config, error) {
	cfg := Config{}
	flags := flag.NewFlagSet("perlin-noise-go", flag.ContinueOnError)
	flags.SetOutput(output)
	flags.StringVar(&cfg.Output, "output", "perlin.png", "output PNG file")
	flags.IntVar(&cfg.Width, "width", 512, "output image width in pixels")
	flags.IntVar(&cfg.Height, "height", 512, "output image height in pixels")
	flags.IntVar(&cfg.GridWidth, "grid-width", 8, "noise grid width in cells")
	flags.IntVar(&cfg.GridHeight, "grid-height", 8, "noise grid height in cells")
	flags.Int64Var(&cfg.Seed, "seed", 1, "random seed")
	flags.StringVar(&cfg.Palette, "palette", "clouds", "color palette: clouds or fire")
	flags.BoolVar(&cfg.Flame, "flame", false, "add a hot base and bottom-to-top cooling (combine with -palette fire)")
	flags.BoolVar(&cfg.Clouds, "clouds", false, "render layered clouds with sky, density and directional lighting")
	flags.Float64Var(&cfg.DirectionAngle, "direction-angle", 0, "noise direction angle in degrees")
	flags.Float64Var(&cfg.DirectionStretch, "direction-stretch", 1, "stretch along the noise direction")
	flags.BoolVar(&cfg.Demo, "demo", false, "render an enlarged grid demonstration")
	flags.IntVar(&cfg.DemoScale, "demo-scale", 32, "pixels per cell in demo mode")
	if err := flags.Parse(args); err != nil {
		return Config{}, err
	}
	if flags.NArg() != 0 {
		return Config{}, fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate checks configuration values and resource limits.
func (cfg Config) Validate() error {
	if cfg.Flame && cfg.Clouds {
		return fmt.Errorf("flame and clouds modes cannot be combined")
	}
	if cfg.Clouds && cfg.Palette != "clouds" {
		return fmt.Errorf("clouds mode requires -palette clouds")
	}
	if cfg.Output == "" {
		return fmt.Errorf("output path must not be empty")
	}
	if err := validateDimensions("image", cfg.Width, cfg.Height); err != nil {
		return err
	}
	if cfg.GridWidth < 1 || cfg.GridWidth > maxDimension {
		return fmt.Errorf("grid width must be between 1 and %d", maxDimension)
	}
	if cfg.GridHeight < 1 || cfg.GridHeight > maxDimension {
		return fmt.Errorf("grid height must be between 1 and %d", maxDimension)
	}
	switch cfg.Palette {
	case "clouds", "fire":
	default:
		return fmt.Errorf("unknown palette %q (choose clouds or fire)", cfg.Palette)
	}
	if math.IsNaN(cfg.DirectionAngle) || math.IsInf(cfg.DirectionAngle, 0) {
		return fmt.Errorf("direction angle must be a finite number")
	}
	if math.IsNaN(cfg.DirectionStretch) || math.IsInf(cfg.DirectionStretch, 0) ||
		cfg.DirectionStretch < 0.1 || cfg.DirectionStretch > 100 {
		return fmt.Errorf("direction stretch must be between 0.1 and 100")
	}
	if cfg.DemoScale < 1 || cfg.DemoScale > maxDimension {
		return fmt.Errorf("demo scale must be between 1 and %d", maxDimension)
	}
	if cfg.Demo {
		if err := validateDimensions("demo image", cfg.GridWidth*cfg.DemoScale, cfg.GridHeight*cfg.DemoScale); err != nil {
			return err
		}
	}
	return nil
}

func validateDimensions(name string, width, height int) error {
	if width < 1 || width > maxDimension {
		return fmt.Errorf("%s width must be between 1 and %d", name, maxDimension)
	}
	if height < 1 || height > maxDimension {
		return fmt.Errorf("%s height must be between 1 and %d", name, maxDimension)
	}
	if width > maxPixels/height {
		return fmt.Errorf("%s dimensions exceed the %d-pixel limit", name, maxPixels)
	}
	return nil
}

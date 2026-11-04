package main

import (
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/austintaylor/perlin-noise-go/internal/config"
	"github.com/austintaylor/perlin-noise-go/internal/noise"
	"github.com/austintaylor/perlin-noise-go/internal/palette"
	"github.com/austintaylor/perlin-noise-go/internal/render"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); errors.Is(err, flag.ErrHelp) {
		return
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "perlin-noise-go: %v\n", err)
		os.Exit(2)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	cfg, err := config.Parse(args, stderr)
	if errors.Is(err, flag.ErrHelp) {
		return err
	}
	if err != nil {
		return err
	}
	if cfg.Demo {
		return fmt.Errorf("demo mode is not implemented yet")
	}
	sampler, err := noise.New(cfg.Seed, cfg.DirectionAngle, cfg.DirectionStretch)
	if err != nil {
		return fmt.Errorf("create noise sampler: %w", err)
	}
	colors, err := palette.Preset(cfg.Palette)
	if err != nil {
		return fmt.Errorf("select palette: %w", err)
	}
	img, err := render.Image(cfg.Width, cfg.Height, cfg.GridWidth, cfg.GridHeight, sampler, colors, render.Options{Flame: cfg.Flame})
	if err != nil {
		return fmt.Errorf("render image: %w", err)
	}
	file, err := os.Create(cfg.Output)
	if err != nil {
		return fmt.Errorf("create output %q: %w", cfg.Output, err)
	}
	if err := png.Encode(file, img); err != nil {
		encodeErr := fmt.Errorf("encode PNG %q: %w", cfg.Output, err)
		closeErr := file.Close()
		if closeErr != nil {
			return errors.Join(encodeErr, fmt.Errorf("close output %q: %w", cfg.Output, closeErr))
		}
		return encodeErr
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close output %q: %w", cfg.Output, err)
	}

	_, err = fmt.Fprintf(stdout, "Wrote %s (%d x %d PNG)\n", cfg.Output, cfg.Width, cfg.Height)
	return err
}

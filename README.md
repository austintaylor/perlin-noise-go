# perlin-noise-go

A standard-library-only Go program for generating PNG images with 2D Perlin
noise. It is intended as a small, approachable demonstration of the algorithm.

## Status

The CLI configuration and validation are in place. Noise generation and PNG
rendering are not implemented yet.

## Run

```sh
go run . -help
```

Once rendering is implemented, generate a regular image with:

```sh
go run . -output clouds.png -width 800 -height 600 \
  -grid-width 12 -grid-height 9 -seed 42 -palette clouds
```

Select the fire palette and orient/stretch the noise direction with:

```sh
go run . -output fire.png -palette fire \
  -direction-angle 90 -direction-stretch 2.5
```

The `-demo` option will render the noise grid enlarged, with each cell scaled
to `-demo-scale` pixels. Demo image dimensions are derived from the grid size
and scale; `-width` and `-height` apply to regular images.

## Options

| Option | Default | Description |
| --- | --- | --- |
| `-output` | `perlin.png` | Output PNG path |
| `-width`, `-height` | `512` | Regular image dimensions in pixels |
| `-grid-width`, `-grid-height` | `8` | Noise grid dimensions in cells |
| `-seed` | `1` | Seed for deterministic noise generation |
| `-palette` | `clouds` | Palette preset: `clouds` or `fire` |
| `-direction-angle` | `0` | Noise direction angle in degrees |
| `-direction-stretch` | `1` | Stretch factor along the noise direction |
| `-demo` | `false` | Render an enlarged grid demonstration |
| `-demo-scale` | `32` | Pixels per cell in demo mode |

Image dimensions must be between 1 and 8192 pixels per side, and an image may
contain at most 64 million pixels. Grid dimensions must be between 1 and 8192
cells per side. Demo output is subject to the same image limits. Direction
stretch must be between 0.1 and 100. Invalid values and unexpected positional
arguments are reported as errors.

Font rendering and numeric cell labels are out of scope for v1.

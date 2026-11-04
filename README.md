# perlin-noise-go

A standard-library-only Go program for generating PNG images with 2D Perlin
noise. It is intended as a small, approachable demonstration of the algorithm.

## Status

The CLI configuration, seeded and directional Perlin noise sampler, and
`clouds`/`fire` palette mapping are in place. The CLI writes regular PNG images;
enlarged demo rendering is not implemented yet.

## Run

```sh
go run . -help
```

Generate a regular image with:

```sh
go run . -output clouds.png -width 800 -height 600 \
  -grid-width 12 -grid-height 9 -seed 42 -palette clouds
```

Select the fire palette and orient/stretch the noise direction with:

```sh
go run . -output fire.png -palette fire \
  -direction-angle 90 -direction-stretch 2.5
```

For a hot base that cools toward the top, add `-flame`:

```sh
go run . -output flames.png -width 800 -height 1000 -palette fire -flame \
  -grid-width 14 -grid-height 10 -seed 42 \
  -direction-angle 90 -direction-stretch 4
```

Flame shaping combines a bottom-to-top intensity gradient with noise, tapering
the noise at both ends. The bottom row uses the highest palette value and the
top row the lowest; a one-row image uses the highest. Intermediate values
fluctuate to create irregular tongues rather than a uniform fade. This is a
static procedural approximation, not a physical fire simulation. The palette
and direction remain independent controls; `-palette fire` alone is unchanged.

The `-demo` option is planned but not implemented yet. Once available, it will
render the noise grid enlarged, with each cell scaled to `-demo-scale` pixels.
Demo image dimensions will be derived from the grid size and scale; `-width`
and `-height` apply to regular images.

## Options

| Option | Default | Description |
| --- | --- | --- |
| `-output` | `perlin.png` | Output PNG path |
| `-width`, `-height` | `512` | Regular image dimensions in pixels |
| `-grid-width`, `-grid-height` | `8` | Noise grid dimensions in cells |
| `-seed` | `1` | Seed for deterministic noise generation |
| `-palette` | `clouds` | Palette preset: `clouds` or `fire` |
| `-flame` | `false` | Hot base with intensity decreasing toward the top |
| `-direction-angle` | `0` | Noise direction angle in degrees |
| `-direction-stretch` | `1` | Stretch factor along the noise direction |
| `-demo` | `false` | Render an enlarged grid demonstration |
| `-demo-scale` | `32` | Pixels per cell in demo mode |

Image dimensions must be between 1 and 8192 pixels per side, and an image may
contain at most 64 million pixels. Grid dimensions must be between 1 and 8192
cells per side. Demo output is subject to the same image limits. Direction
stretch must be between 0.1 and 100. Invalid values and unexpected positional
arguments are reported as errors.

Noise values in `[-1, 1]` map linearly to palette positions `[0, 1]`; values
outside that range are clamped to the nearest palette endpoint. Palette colors
are interpolated between stops in standard RGBA channels.

Font rendering and numeric cell labels are out of scope for v1.

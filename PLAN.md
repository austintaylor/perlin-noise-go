# Version 1 plan

## Goal

Build a deterministic, standard-library-only Go CLI that generates PNGs from
2D Perlin noise. The project is intended to demonstrate how the algorithm
works, not to provide a production texture-generation toolkit.

## V1 scope

- Configure output path and image width/height.
- Configure the noise grid, seed, palette, and directional sampling.
- Include `clouds` and `fire` palette presets.
- Generate regular images and an enlarged demonstration image with visible
  cell boundaries.
- Keep the implementation and CLI in Go's standard library.

## Out of scope

- Web UI.
- Font rendering or numeric values printed inside cells.
- External packages, animation, and advanced texture controls.

## Tasks

### 1. Establish the Go project

- [x] Add `go.mod` and a documented command-line entry point.
- [x] Define and validate a configuration for output path, pixel dimensions,
  grid dimensions, seed, palette, direction, and demo enlargement.
- [x] Document defaults, examples, and invalid input behavior in `README.md`.

### 2. Implement the noise sampler

- [x] Implement seeded 2D gradient Perlin noise with an immutable permutation
  table after initialization.
- [x] Add directional sampling controls so the sampled pattern can be stretched
  or oriented without changing the underlying noise algorithm.
- [x] Test repeatability for a fixed seed, variation across seeds, and expected
  sample bounds.

### 3. Add palette mapping

- [x] Define palette stops and interpolation independently from noise sampling.
- [x] Add `clouds` and `fire` presets, with clear handling of noise-to-palette
  value mapping.
- [x] Test interpolation endpoints and preset selection.

### 4. Render PNG output

- [ ] Map noise samples to colors in an `image.NRGBA` and encode PNGs using
  `image/png`.
- [ ] Add optional row-chunk parallel rendering. Workers share only immutable
  sampler/palette data and write to disjoint output rows; avoid per-pixel
  goroutines and shared mutable random state.
- [ ] Add regular-render tests for dimensions and PNG encoding.

### 5. Add demonstration mode

- [ ] Render a small noise grid and enlarge it with nearest-neighbor scaling.
- [ ] Draw cell boundaries so the grid is visible in the output.
- [ ] Keep font rendering and numeric cell labels out of v1.
- [ ] Test enlargement dimensions and grid-boundary placement.

### 6. Finish CLI integration

- [ ] Wire configuration, validation, rendering, and file output together.
- [ ] Return useful errors for invalid settings and file/encoding failures.
- [ ] Add end-to-end CLI examples to `README.md` and verify generated PNGs.

## Completion criteria

- A user can run the program with standard Go tooling and generate a PNG without
  third-party dependencies.
- Repeating a command with the same settings and seed produces the same image.
- Both palette presets and directional sampling are configurable.
- Demo output visibly enlarges the grid and draws cell boundaries, without text.

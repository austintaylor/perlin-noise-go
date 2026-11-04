package render

import (
	"fmt"
	"image"
	"runtime"
	"sync"

	"github.com/austintaylor/perlin-noise-go/internal/noise"
	"github.com/austintaylor/perlin-noise-go/internal/palette"
)

const maxPixels = 64_000_000

// Image renders a noise field into an NRGBA image. The grid dimensions specify
// the number of Perlin lattice cells sampled across the image.
func Image(width, height, gridWidth, gridHeight int, sampler *noise.Sampler, colors palette.Palette) (*image.NRGBA, error) {
	if width < 1 || height < 1 || width > maxPixels/height {
		return nil, fmt.Errorf("image dimensions must be positive and contain at most %d pixels", maxPixels)
	}
	if gridWidth < 1 || gridHeight < 1 {
		return nil, fmt.Errorf("grid dimensions must be positive")
	}
	if sampler == nil {
		return nil, fmt.Errorf("noise sampler must not be nil")
	}
	if !colors.Valid() {
		return nil, fmt.Errorf("palette must be initialized with valid stops")
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	workers := min(runtime.GOMAXPROCS(0), height)
	var group sync.WaitGroup
	group.Add(workers)
	for worker := 0; worker < workers; worker++ {
		yStart := worker * height / workers
		yEnd := (worker + 1) * height / workers
		go func(yStart, yEnd int) {
			defer group.Done()
			for y := yStart; y < yEnd; y++ {
				noiseY := float64(y) * float64(gridHeight) / float64(height)
				for x := 0; x < width; x++ {
					noiseX := float64(x) * float64(gridWidth) / float64(width)
					img.SetNRGBA(x, y, colors.ColorForNoise(sampler.Sample(noiseX, noiseY)))
				}
			}
		}(yStart, yEnd)
	}
	group.Wait()
	return img, nil
}

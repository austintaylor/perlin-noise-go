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

// Options selects optional shaping of the noise field before color mapping.
type Options struct {
	Flame  bool
	Clouds bool
}

// Image renders a noise field into an NRGBA image. The grid dimensions specify
// the number of Perlin lattice cells sampled across the image.
func Image(width, height, gridWidth, gridHeight int, sampler *noise.Sampler, colors palette.Palette, options Options) (*image.NRGBA, error) {
	if options.Flame && options.Clouds {
		return nil, fmt.Errorf("flame and clouds modes cannot be combined")
	}
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
					value := sampler.Sample(noiseX, noiseY)
					if options.Clouds || options.Flame {
						value = 1.7 * fractalField(sampler, noiseX, noiseY)
					}
					if options.Flame {
						value = flameValue(value, y, height)
					}
					img.SetNRGBA(x, y, colors.ColorForNoise(value))
				}
			}
		}(yStart, yEnd)
	}
	group.Wait()
	return img, nil
}

func flameValue(value float64, y, height int) float64 {
	// A one-row image represents the hot base.
	depth := 1.0
	if height > 1 {
		depth = float64(y) / float64(height-1)
	}
	// Taper noise at both ends to keep the base hot and the top cool.
	intensity := depth + 2*depth*(1-depth)*value
	return 2*max(0, min(1, intensity)) - 1
}

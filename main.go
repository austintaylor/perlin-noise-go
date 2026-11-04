package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/austintaylor/perlin-noise-go/internal/config"
)

func main() {
	cfg, err := config.Parse(os.Args[1:], os.Stderr)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "perlin-noise-go: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("Configuration valid for %q; image rendering is not implemented yet.\n", cfg.Output)
}

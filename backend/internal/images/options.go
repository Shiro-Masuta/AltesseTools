package images

import (
	"image/png"

	"golang.org/x/image/tiff"
)

type Options struct {
	Format       string
	Quality      int                  // JPEG, WebP, AVIF
	Lossless     bool                 // WebP, AVIF
	PNGLevel     png.CompressionLevel // PNG compression 0–9
	TIFFCompress tiff.CompressionType // TIFF compression (Deflate, LZW…)
}

func applyDefaults(opts *Options) *Options {
	if opts == nil {
		return &Options{
			Format:       "jpeg",
			Quality:      85,
			Lossless:     false,
			PNGLevel:     png.DefaultCompression,
			TIFFCompress: tiff.Deflate,
		}
	}

	if opts.Quality == 0 {
		opts.Quality = 85
	}
	if opts.PNGLevel == 0 {
		opts.PNGLevel = png.DefaultCompression
	}
	if opts.TIFFCompress == 0 {
		opts.TIFFCompress = tiff.Deflate
	}

	return opts
}

package images

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/Kagami/go-avif"
	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func ConvertFromReader(r io.Reader, opts *Options) ([]byte, error) {
	// 1. Décodage
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("erreur de décodage de l'image : %w", err)
	}

	// 2. Application des valeurs par défaut
	opts = applyDefaults(opts)
	format := strings.ToUpper(opts.Format)

	// 3. Buffer de sortie
	var buf bytes.Buffer

	// 4. Encodage selon le format + options
	switch format {
	case "PNG":
		encoder := png.Encoder{CompressionLevel: opts.PNGLevel}
		err = encoder.Encode(&buf, img)

	case "JPEG", "JPG":
		// Qualité de 1 à 100 (85 par défaut)
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: opts.Quality})

	case "WEBP":
		err = webp.Encode(&buf, img, &webp.Options{
			Quality:  float32(opts.Quality),
			Lossless: opts.Lossless,
		})

	case "AVIF":
		err = avif.Encode(&buf, img, &avif.Options{
			Quality: opts.Quality,
		})

	case "BMP":
		// Pas de compression disponible pour BMP
		err = bmp.Encode(&buf, img)

	case "TIFF":
		err = tiff.Encode(&buf, img, &tiff.Options{
			Compression: opts.TIFFCompress,
		})

	default:
		return nil, fmt.Errorf("format '%s' non supporté", format)
	}

	// 5. Gestion d'erreurs d'encodage
	if err != nil {
		return nil, fmt.Errorf("échec de l'encodage %s : %w", format, err)
	}

	return buf.Bytes(), nil
}

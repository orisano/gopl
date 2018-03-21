package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

func main() {
	format := flag.String("f", "jpeg", "output image format")
	flag.Parse()

	if err := convertFormat(*format, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "ex01: %v\n", err)
		os.Exit(1)
	}
}

func convertFormat(format string, in io.Reader, out io.Writer) error {
	encoders := map[string]func(image.Image, io.Writer) error{
		"jpeg": toJPEG,
		"png":  toPNG,
		"gif":  toGIF,
	}
	encode, ok := encoders[format]
	if !ok {
		return fmt.Errorf("unsupported format: %v", format)
	}

	image, _, err := image.Decode(in)
	if err != nil {
		return err
	}
	if err := encode(image, out); err != nil {
		return err
	}
	return nil
}

func toJPEG(img image.Image, out io.Writer) error {
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

func toPNG(img image.Image, out io.Writer) error {
	return png.Encode(out, img)
}

func toGIF(img image.Image, out io.Writer) error {
	return gif.Encode(out, img, nil)
}

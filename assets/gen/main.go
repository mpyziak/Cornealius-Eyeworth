//go:build ignore

// Logo.png -> Logo64.png, run by `go generate ./assets/...`.
//
// 64 because every consumer is then a power-of-two reduction: 16px title bar,
// 32px taskbar, exact at 200% DPI. Build time, not startup - this is ~20ms.
package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	xdraw "golang.org/x/image/draw"
)

const iconSize = 64

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	srcPath := filepath.Join(root, "Logo.png")
	dstPath := filepath.Join(root, "Logo64.png")

	raw, err := os.ReadFile(srcPath)
	if err != nil {
		log.Fatalf("read %s: %v", srcPath, err)
	}

	src, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		log.Fatalf("decode %s: %v", srcPath, err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), xdraw.Over, nil)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		log.Fatalf("encode: %v", err)
	}
	if err := os.WriteFile(dstPath, buf.Bytes(), 0644); err != nil {
		log.Fatalf("write %s: %v", dstPath, err)
	}

	b := src.Bounds()
	log.Printf("wrote %s: %dx%d -> %dx%d, %d bytes",
		filepath.Base(dstPath), b.Dx(), b.Dy(), iconSize, iconSize, buf.Len())
}

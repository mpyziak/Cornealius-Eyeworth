package assets

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"testing"

	xdraw "golang.org/x/image/draw"
)

const iconSize = 64

func decode(t *testing.T, name string, raw []byte) image.Image {
	t.Helper()
	img, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("%s does not decode as PNG: %v", name, err)
	}
	return img
}

func TestIconIsPreScaled(t *testing.T) {
	b := decode(t, "Logo64.png", iconBytes64).Bounds()

	if b.Dx() != iconSize || b.Dy() != iconSize {
		t.Fatalf("Icon is %dx%d, want %dx%d - run: go generate ./assets/...",
			b.Dx(), b.Dy(), iconSize, iconSize)
	}

	lb := decode(t, "Logo.png", logoBytes).Bounds()
	if b.Dx()*4 > lb.Dx() {
		t.Errorf("Icon (%dx%d) is not meaningfully smaller than Logo (%dx%d)",
			b.Dx(), b.Dy(), lb.Dx(), lb.Dy())
	}
}

func TestIconMatchesLogo(t *testing.T) {
	committed := decode(t, "Logo64.png", iconBytes64)
	source := decode(t, "Logo.png", logoBytes)

	fresh := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	xdraw.CatmullRom.Scale(fresh, fresh.Bounds(), source, source.Bounds(), xdraw.Over, nil)

	var sum, n float64
	for y := range iconSize {
		for x := range iconSize {
			cr, cg, cb, ca := committed.At(x, y).RGBA()
			fr, fg, fb, fa := fresh.At(x, y).RGBA()
			for _, d := range [4]float64{
				float64(cr) - float64(fr),
				float64(cg) - float64(fg),
				float64(cb) - float64(fb),
				float64(ca) - float64(fa),
			} {
				sum += math.Abs(d)
				n++
			}
		}
	}

	const tolerance = 0.01 * 65535
	if mean := sum / n; mean > tolerance {
		t.Fatalf("Logo64.png does not match a fresh downscale of Logo.png "+
			"(mean channel difference %.0f, tolerance %.0f).\n"+
			"Logo.png probably changed - run: go generate ./assets/...",
			mean, tolerance)
	}
}

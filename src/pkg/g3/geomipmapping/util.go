package geomipmapping

import (
	"os"
	"io"
	"image"
)

type ImageHeightMap struct {
	img image.Image
}

// Generates a new HeighMap object from an image stream
func NewHeightMapFromImageStream(reader io.Reader) (*ImageHeightMap, os.Error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return &ImageHeightMap{img}, nil
}

// Generates a new HeighMap object from an image file
func NewHeightMapFromImageFile(fileName string) (*ImageHeightMap, os.Error) {
	file, err := os.Open(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewHeightMapFromImageStream(file)
}

func (hm *ImageHeightMap) Size() (width, height int) {
	rect := hm.img.Bounds()
	return rect.Max.X - rect.Min.Y, rect.Max.Y - rect.Min.Y
}

func (hm *ImageHeightMap) Height(x, y float32) float32 {
	w, h := hm.Size()
	color1 := hm.img.At(int(x)%w, int(y)%h)
	color2 := hm.img.At(int(x+1)%w, int(y)%h)
	color3 := hm.img.At(int(x)%w, int(y+1)%h)
	color4 := hm.img.At(int(x+1)%w, int(y+1)%h)
	fx := x - float32(int(x))
	fy := y - float32(int(y))
	h1, _, _, _ := color1.RGBA()
	h2, _, _, _ := color2.RGBA()
	h3, _, _, _ := color3.RGBA()
	h4, _, _, _ := color4.RGBA()
	hx1 := float32(h1)*(1.0-fx) + fx*float32(h2)
	hx2 := float32(h3)*(1.0-fx) + fx*float32(h4)
	he := hx1*(1.0-fy) + fy*hx2
	return float32(he) / (255.0 * 255.0)
}

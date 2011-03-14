package g3

import (
	"os"
	"io"
	"io/ioutil"
	"image"
)

func ReadStringsFromFiles(fileNames ...string) ([]string, os.Error) {
	strings := make([]string, 0)
	for _, fileName := range fileNames {
		str, err := ReadStringFromFile(fileName)
		if err != nil {
			return strings, err
		}
		strings = append(strings, str)
	}
	return strings, nil
}

func ReadStringFromStream(reader io.Reader) (string, os.Error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ReadStringFromFile(fileName string) (string, os.Error) {
	file, err := os.Open(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return ReadStringFromStream(file)
}

// Images

func ReadImagesFromStreams(readers ...io.Reader) ([]image.Image, os.Error) {
	images := make([]image.Image, 0)
	for _, reader := range readers {
		img, _, err := image.Decode(reader)
		if err != nil {
			return images, err
		}
		images = append(images, img)
	}
	return images, nil
}

func ReadImagesFromFiles(fileNames ...string) ([]image.Image, os.Error) {
	images := make([]image.Image, 0, len(fileNames))
	for _, fileName := range fileNames {
		reader, err := os.Open(fileName, os.O_RDONLY, 0666)
		if err != nil {
			return images, err
		}
		defer reader.Close()

		img, _, err := image.Decode(reader)
		if err != nil {
			return images, err
		}
		images = append(images, img)
	}
	return images, nil
}


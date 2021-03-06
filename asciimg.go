package main

import (
	"flag"
	"fmt"
	"github.com/buger/goterm"
	"image"
	"image/color"
	_ "image/jpeg"
	"io"

	// Side-effect import.
	// Сайд-эффект — добавление декодера PNG в пакет image.
	_ "image/png"
	"os"
	// Внешняя зависимость.
	"github.com/olekukonko/ts"
	"golang.org/x/image/draw"
)

func getTermWidth() int {
	term, _ := ts.GetSize()
	if term.Col()-term.PosX() > 0 {
		return term.Col() - term.PosX()
	}
	return goterm.Width()
}

func getTermHeigh() int {
	if getTermWidth() > 0 {
		return -1
	}
	return goterm.Height()
}

var (
	o         = flag.String("o", "", "output into file")
	is_scaled = flag.Bool("noscale", false, "scale option")
	w         = flag.Int("w", getTermWidth(), "wigth")
	h         = flag.Int("h", getTermHeigh(), "heigh")
)

func scale(img image.Image, w int, h int) image.Image {
	dstImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dstImg, dstImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dstImg
}

func decodeImageFile(imgName string) (image.Image, error) {
	imgFile, err := os.Open(imgName)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imgFile)

	return img, err
}

func processPixel(c color.Color) rune {
	gc := color.GrayModel.Convert(c)
	r, _, _, _ := gc.RGBA()
	r = r >> 8

	symbols := []rune("@80GCLft1i;:,. ")
	return symbols[r/(256/uint32(len(symbols)-1))]
}

func convertToAscii(img image.Image) [][]rune {
	textImg := make([][]rune, img.Bounds().Dy())
	for i := range textImg {
		textImg[i] = make([]rune, img.Bounds().Dx())
	}

	for i := range textImg {
		for j := range textImg[i] {
			textImg[i][j] = processPixel(img.At(j, i))
		}
	}
	return textImg
}

func exportToFile(textImg [][]rune, writer io.Writer) error {
	for i := range textImg {
		for j := range textImg[i] {
			fmt.Fprint(writer, string(textImg[i][j]))
		}
		fmt.Fprint(writer, "\n")
	}
	return nil
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: asciimg <image.jpg>")
		os.Exit(0)
	}
	imgName := flag.Arg(0)

	img, err := decodeImageFile(imgName)
	kef := img.Bounds().Dy() / img.Bounds().Dx()

	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	if *is_scaled == false && len(*o) == 0 {
		if *h == -1 {
			*h = *w * kef * 45 / 100
		}
		newimg := scale(img, *w, *h)
		img = newimg
	}

	textImg := convertToAscii(img)

	if len(*o) > 0 {
		w, err := os.Create(*o)
		if err != nil {
			os.Exit(1)
		}
		exportToFile(textImg, w)
	} else {
		for i := range textImg {
			for j := range textImg[i] {
				fmt.Printf("%c", textImg[i][j])
			}
			fmt.Println()
		}
	}
}

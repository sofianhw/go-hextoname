package hextoname

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type RGBHSL struct {
	Red   int
	Blue  int
	Green int
	Hue   int
	Sat   int
	Light int
	Name  string
}

type listRGBHSL struct {
	Objects []RGBHSL
}

func Setup(path string) []RGBHSL {
	var hexRGBHSL listRGBHSL
	hexRGBHSL.Objects = make([]RGBHSL, 0)
	csvFile, _ := os.Open(path)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		rgb := ToRGB(line[0])
		hsl := ToHSL(line[0])
		rgbhsl := RGBHSL{rgb[0], rgb[1], rgb[2], hsl[0], hsl[1], hsl[2], line[1]}
		hexRGBHSL.Objects = append(hexRGBHSL.Objects, rgbhsl)
	}
	return hexRGBHSL.Objects
}

func ToRGB(hexcode string) []int {
	hexCodeString := []rune(hexcode)
	red, _ := strconv.ParseInt(string(hexCodeString[0:2]), 16, 32)
	blue, _ := strconv.ParseInt(string(hexCodeString[2:4]), 16, 32)
	green, _ := strconv.ParseInt(string(hexCodeString[4:6]), 16, 32)
	rgb := []int{int(red), int(blue), int(green)}
	return rgb
}

func ToHSL(hexcode string) []int {
	rgb := ToRGB(hexcode)
	r, g, b := float64(rgb[0])/255, float64(rgb[1])/255, float64(rgb[2])/255
	min := math.Min(r, math.Min(g, b))
	max := math.Max(r, math.Max(g, b))
	delta := max - min
	l := (min + max) / 2
	s := float64(0.0)
	if l > 0 && l < 1 {
		if l < 0.5 {
			s = delta / (2 * l)
		} else {
			s = delta / (2 - 2*l)
		}
	}

	h := float64(0.0)
	if delta > 0 {
		if max == r && max != g {
			h += (g - b) / delta
		}
		if max == g && max != b {
			h += (2.0 + (b-r)/delta)
		}
		if max == b && max != r {
			h += (4.0 + (r-g)/delta)
		}
	}
	hsl := []int{int(h * 60), int(s * 100), int(l * 100)}
	return hsl
}

func GetName(color string, rgbhsls []RGBHSL) (name string) {
	rgb := ToRGB(color)
	r := rgb[0]
	g := rgb[1]
	b := rgb[2]

	hsl := ToHSL(color)
	h := hsl[0]
	s := hsl[1]
	l := hsl[2]

	ndf1 := 0.0
	ndf2 := 0.0
	ndf := 0.0
	// cl := -1
	df := -1.0

	for _, rgbhsl := range rgbhsls {
		ndf1 = math.Pow(float64(r-rgbhsl.Red), 2.0) + math.Pow(float64(g-rgbhsl.Green), 2.0) + math.Pow(float64(b-rgbhsl.Blue), 2.0)
		ndf2 = math.Pow(float64(h-rgbhsl.Hue), 2.0) + math.Pow(float64(s-rgbhsl.Sat), 2.0) + math.Pow(float64(l-rgbhsl.Light), 2.0)
		ndf = ndf1 + ndf2*2
		if df < 0 || df > ndf {
			df = ndf
			// cl = idx
			name = rgbhsl.Name
		}
	}

	return name
}

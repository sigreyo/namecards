package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	// Read the image file
	imgFile, _ := os.Open("namecard.png")
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load the font
	fontBytes, err := ioutil.ReadFile("Allura-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Create a PDF
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Create a new face to calculate text width and height
	opts := truetype.Options{}
	opts.Size = 12
	// face := truetype.NewFace(f, &opts)

	// Card size
	cardWidth := 80.0
	cardHeight := 46.0

	// Margins
	marginX := (210 - 2*cardWidth) / 3
	marginY := (297 - 6*cardHeight) / 7

	// Images per page
	imagesPerPage := 12

	// Loop over the names and create name cards
	for i, name := range names {
		bounds := img.Bounds()
		rgba := image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

		// Add a new page if it's a new set of 12 images
		if i%imagesPerPage == 0 {
			pdf.AddPage()
		}

		// Calculate the current row and column
		row := float64((i % imagesPerPage) / 2)
		col := float64((i % imagesPerPage) % 2)

		// ... other parts of the code

		// Calculate text width and height
		c := freetype.NewContext()
		c.SetDPI(600)
		c.SetFont(f)
		c.SetFontSize(12)

		width, err := c.DrawString(name, freetype.Pt(0, 0))
		if err != nil {
			log.Fatal(err)
		}

		// Convert freetype's fixed point width to int
		widthInt := int(width.X >> 6)

		// Calculate the point to start drawing so the text is centered
		startX := (rgba.Bounds().Dx() - widthInt) / 2
		startY := int(float64(rgba.Bounds().Dy()) * 0.5) // position the text at 40% of the Y-axis

		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)
		c.SetSrc(image.Black)

		pt := freetype.Pt(startX, startY) // position of the text
		_, err = c.DrawString(name, pt)
		if err != nil {
			log.Fatal(err)
		}

		// ... other parts of the code

		imgName := "output_" + name + ".png"
		imgFile, _ := os.Create(imgName)
		defer imgFile.Close()
		err = png.Encode(imgFile, rgba)
		if err != nil {
			log.Fatal(err)
		}

		// Add the image to the PDF at the correct position
		posX := marginX + col*(cardWidth+marginX)
		posY := marginY + row*(cardHeight+marginY)
		pdf.Image(imgName, posX, posY, cardWidth, cardHeight, false, "", 0, "")

		os.Remove(imgName)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Save the PDF
	err = pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		log.Fatal(err)
	}
}

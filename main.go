package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type Config struct {
	quote      string
	author     string
	width      int
	height     int
	outputPath string
	fontPath   string
	scale      float64
}

func parseFlags() Config {
	quote := flag.String("quote", "[ Hello World ]", "The quote text")
	author := flag.String("author", "", "The author of the quote")
	width := flag.Int("width", 3840, "Width of the wallpaper")
	height := flag.Int("height", 2160, "Height of the wallpaper")
	outputPath := flag.String("output", "wallpaper.png", "Output file path")
	fontPath := flag.String("font", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", "Path to font file")
	scale := flag.Float64("scale", 1.0, "Scale factor for font size")

	flag.Parse()

	if *quote == "" {
		fmt.Println("No quote provided...")
	}

	return Config{
		quote:      *quote,
		author:     *author,
		width:      *width,
		height:     *height,
		outputPath: *outputPath,
		fontPath:   *fontPath,
		scale:      *scale,
	}
}

func createGradientBackground(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with black background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	return img
}

func addText(img *image.RGBA, config Config) error {
    c := freetype.NewContext()
    c.SetDPI(72)

    fontData, err := os.ReadFile(config.fontPath)
    if err == nil {
        font, err := truetype.Parse(fontData)
        if err != nil {
            return fmt.Errorf("failed to parse font: %v", err)
        } else {
            c.SetFont(font)
        }
    }

    c.SetClip(img.Bounds())
    c.SetDst(img)
    c.SetSrc(image.White)

    quoteSize := (float64(config.height) / 30) * config.scale
    authorSize := quoteSize * 0.6 * config.scale

    c.SetFontSize(quoteSize)
    lines := wrapText(config.quote, config.width/int(quoteSize)*2)

    y := float64(config.height)/2 - (float64(len(lines))*quoteSize)/2
    for _, line := range lines {
        textWidth := int(c.PointToFixed(quoteSize * float64(len(line)) * 0.6) >> 6)
        x := (config.width - textWidth) / 2
        pt := freetype.Pt(x, int(y))
        _, err = c.DrawString(line, pt)
        if err != nil {
            return fmt.Errorf("failed to draw quote: %v", err)
        }
        y += quoteSize * 1.5
    }

    if config.author != "" {
        c.SetFontSize(authorSize)
        authorText := fmt.Sprintf("- %s  ", config.author)
        textWidth := int(c.PointToFixed(authorSize * float64(len(authorText)) * 0.6) >> 6)
        x := (config.width - textWidth) / 2
        pt := freetype.Pt(x, int(y+authorSize))
        _, err = c.DrawString(authorText, pt)
        if err != nil {
            return fmt.Errorf("failed to draw author: %v", err)
        }
    }

    return nil
}

func wrapText(text string, maxChars int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= maxChars {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}

func main() {
	config := parseFlags()

	// Create background
	img := createGradientBackground(config.width, config.height)

	// Add text
	err := addText(img, config)
	if err != nil {
		log.Fatalf("Failed to add text: %v", err)
	}

	// Save the image
	f, err := os.Create(config.outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	fmt.Printf("Wallpaper generated to: %s\n", config.outputPath)
}
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"strings"
	"time"

	ex "github.com/markus-wa/demoinfocs-golang/v5/examples"
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
)

const DOT_SIZE = 15

var (
	lastMapImg  []byte
	demoPath    string
	mapName     string
	initialTime time.Time
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("CS2 Demo-Based Real-Time Radar Visualizer")

	fmt.Printf("Demo Path: ")
	demoPath, _ = reader.ReadString('\n')
	demoPath = strings.TrimSpace(demoPath)
	demoPath = strings.ReplaceAll(demoPath, "\\", "/") // normalize path
	demoPath = strings.ReplaceAll(demoPath, "\"", "")  // remove quotes if any

	fmt.Printf("Map Name (e.g. de_mirage): ")
	mapName, _ = reader.ReadString('\n')
	mapName = strings.TrimSpace(mapName)

	fmt.Println("\nUsing demo path:", demoPath)
	fmt.Println("Using map name:", mapName)

	fmt.Println("\nWarning: Using this in CS2 may be against the game's rules. Use at your own risk.")
	fmt.Println("Now go to CS2, switch to the specified map and start a demo recording with 'record <demoname>'. Once this is done press Enter here to start the visualizer.")
	fmt.Scanln()

	fileInfo, err := os.Stat(demoPath)
	if err != nil {
		initialTime = time.Now()
	} else {
		initialTime = fileInfo.ModTime()
	}
	fmt.Println(fileInfo)

	contents, _ := os.ReadFile("./index.html")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(contents))
	})

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(lastMapImg)
	})

	go http.ListenAndServe(":8080", nil)
	fmt.Println("HTTP server started on http://localhost:8080")

	for {
		fileInfo, err := os.Stat(demoPath)
		if err != nil {
			fmt.Println("Demo file not found. waiting...")
			initialTime = time.Now()
			continue
		}

		// Check if file has been modified since last check and is not empty
		if fileInfo.ModTime().After(initialTime) && fileInfo.Size() > 0 {
			fmt.Println("Demo file modified. setting initial time to 0")
			initialTime = fileInfo.ModTime()
			f, err := os.Open(demoPath)
			if err != nil {
				fmt.Println("Failed to open demo file:", err)
			} else {
				processDemo(f)
				f.Close()
			}
		}
		time.Sleep(100 * time.Millisecond) // Polling interval
	}
}

func processDemo(f *os.File) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic during parsing (likely demoinfocs bug):", r)
		}
	}()

	p := demoinfocs.NewParser(f)

	var (
		mapMetadata ex.Map      = ex.GetMapMetadata(mapName)
		mapRadarImg image.Image = ex.GetMapRadar(mapName)
	)

	// Parse File to last tick to get latest positions
	err := p.ParseToEnd()
	if err != nil {
		// fmt.Println("Failed to parse demo:", err)
		// No need for it to spam the console, it's gonna appear quite a bit when demo is being recorded
	}

	// Check if map radar image was loaded
	if mapRadarImg == nil {
		fmt.Println("Map radar image not loaded - skipping visualization")
		return
	}

	img := image.NewRGBA(mapRadarImg.Bounds())
	draw.Draw(img, mapRadarImg.Bounds(), mapRadarImg, image.Point{}, draw.Over)

	for _, player := range p.GameState().Participants().Playing() {
		fmt.Println("Processing player:", player.Name)
		if !player.IsAlive() {
			continue // Skip dead players
		}

		pos := player.Position()
		x, y := mapMetadata.TranslateScale(pos.X, pos.Y)

		// Determine color based on team
		var col color.RGBA
		switch player.Team {
		case 2: // T
			col = color.RGBA{255, 0, 0, 255}
		case 3: // CT
			col = color.RGBA{0, 0, 255, 255}
		}

		// Draw dot
		drawCircle(img, int(x), int(y), DOT_SIZE/2, col)
	}

	buffer := bytes.NewBuffer(nil)
	png.Encode(buffer, img)
	lastMapImg = buffer.Bytes()
}

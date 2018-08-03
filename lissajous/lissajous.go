// Lissajous gera animações GIF de figuras de Lissajous aleatórias
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // primeira cor da paleta
	blackIndex = 1 // próxima cor da paleta
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	handler := func(w http.ResponseWriter, r *http.Request) {
		cycle := r.URL.Query().Get("cycles")
		size := r.URL.Query().Get("size")
		lissajous(w, cycle, size)
	}

	http.HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func lissajous(out io.Writer, cycleValue, sizeValue string) {
	const (
		res     = 0.001
		nframes = 64
		delay   = 8
	)

	cycles, err := strconv.Atoi(cycleValue)
	if err != nil {
		out.Write([]byte(err.Error()))
		return
	}

	size, err := strconv.Atoi(sizeValue)
	if err != nil {
		out.Write([]byte(err.Error()))
		return
	}

	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim) // NOTA: ignorando erros de codificação
}

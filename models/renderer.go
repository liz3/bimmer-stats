package models

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"image/color"
	"math"
	"strconv"
)

type CacheStruct struct {
	Surface *sdl.Surface
	Dest *sdl.Rect
}

type Renderer struct {
	Instance *AppInstance
	mouseX int32
	RenderCache map[string]*CacheStruct
	CacheReady bool
}

func (r *Renderer) Update() {
	var surface = r.Instance.Surface
	surface.FillRect(nil, 0)
	r.renderGraphs()

	if r.Instance.Options.RenderLine {
		l := sdl.Rect{r.mouseX, 0, 1,surface.H }
		surface.FillRect(&l, 0xffff0000)

	}


	r.RenderText(r.Instance.Status, 10,surface.H - 30)

	r.Instance.Window.UpdateSurface()
}

func (r *Renderer) setMousePos(x int32) {
	r.mouseX = x
	if r.Instance.Options.RenderLine {
		r.CacheReady = false
	}
}

func (r *Renderer) RenderText(s string, x int32, y int32) {
	var surface = r.Instance.Surface
	r.RenderTextSurface(s,x,y,surface)
}

func (r *Renderer) renderGraphs() {
	if r.Instance.File == nil {
		return
	}
	if r.CacheReady {
		r.blitGraphs()
		return
	}
	fmt.Println("Rendering graphs")
	for _, val := range r.RenderCache {
		val.Surface.Free()
	}
	var file = r.Instance.File
	r.RenderCache = make(map[string]*CacheStruct)
	var w = r.Instance.Surface.W
	var h = r.Instance.Surface.H - 150
	for {
		if h % 4 == 0 {
			break
		}
		h--
	}
	var graphHeight = (h / 4) - 45
	var graphCenter = graphHeight / 2
	var indexStart = r.Instance.GraphOffset
	var trueScale = float64(r.Instance.ScaleFactor) * float64(w) / float64(len(file.Rows)-1)
	for i := indexStart; i < indexStart + 4; i++ {
		if i > int32(len(file.Names)) {
			break
		}

		var name = file.Names[i]
		var dictEntry = file.MaxDict[name]
		var absMax = math.Abs(dictEntry.Max)
		var absMin = math.Abs(dictEntry.Min)
		var x = math.Max(absMin, absMax)
		var factor float64
		if x == 0 {
			factor = 1
		} else {
			factor = float64(graphCenter) / x
		}
		surface, err := sdl.CreateRGBSurface(0, w, graphHeight + 32, 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)
		if err != nil {
			panic(err)
		}
		r.RenderTextSurface(name + " | Min: " + strconv.FormatFloat(dictEntry.Min, 'f', 6, 64) +" | Max: " + strconv.FormatFloat(dictEntry.Max, 'f', 6, 64) + " | Average: " + strconv.FormatFloat(dictEntry.Average, 'f', 6, 64), 0,0, surface)
		surface.FillRect(&sdl.Rect{
			X: 0,
			Y: graphCenter + 30,
			W: w,
			H: 1,
		}, 0xffff0000)
		for x := 0; int32(x) < w; x++ {

			var entryIndex = int32(math.Floor(float64(x)/trueScale)) + r.Instance.Offset
			if entryIndex >= int32(len(file.Rows)) || entryIndex < 0{
				continue
			}
			var entry = file.Rows[entryIndex]
			var value = (entry.Entries[i] * factor) * -1
			var trueY = (graphCenter + int32(value))
			if int32(x) == r.mouseX && r.Instance.Options.RenderLine{
				r.RenderTextSurface("Value: " + strconv.FormatFloat(entry.Entries[i], 'f', 6, 64), w, 0, surface)
			}
			surface.Set(x, int(trueY) + 30, color.White)
			surface.Set(x, int(trueY) + 29, color.White)
		}
		var actualIndex = i - indexStart
		fmt.Println(actualIndex)
			r.RenderCache[name] = &CacheStruct{Surface: surface, Dest: &sdl.Rect{W: w, H: graphHeight + 45, X: 0, Y: (actualIndex * (graphHeight + 32)) + (10 * (actualIndex + 1))}}


	}
	r.CacheReady = true
	r.blitGraphs()
}

func (r *Renderer) RenderTextSurface(s string, x int32, y int32, surface *sdl.Surface) {

	solid,_ := r.Instance.Font.RenderUTF8Blended(s, sdl.Color{255, 255, 255, 255})
	defer solid.Free()
	var tx = x
	if x == surface.W {
		tx -= solid.W
	}
	solid.Blit(nil, surface, &sdl.Rect{
		X: tx,
		Y: y,
		W: solid.W,
		H: solid.H,
	})
}

func (r *Renderer) blitGraphs() {
	var target = r.Instance.Surface
	for _, val := range r.RenderCache {
		val.Surface.Blit(nil, target, val.Dest)
	}
}


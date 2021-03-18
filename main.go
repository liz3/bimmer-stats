package main

import (
	"example.com/liz3/bimmer_stats/models"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)


func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	var font *ttf.Font

	window, err := sdl.CreateWindow("BimmerLink Reader", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		1280, 720, sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE | sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		panic(err)
	}
	window.SetMinimumSize(500, 500)
	defer window.Destroy()
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	if font, err = ttf.OpenFont("Roboto-Regular.ttf", 18); err != nil {
		panic(err)
	}
	var instance = &models.AppInstance{
		Window:  window,
		Surface: surface,
		Font: font,
	}

	defer font.Close()

	instance.Setup()
	instance.RunLoop()

}

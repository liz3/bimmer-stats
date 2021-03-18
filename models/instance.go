package models

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"io/ioutil"
	"strconv"
	"time"
)

type Options struct {
	OpenPath string
	LeftCtrl bool
	ToShow string
	Until time.Time
	RenderLine bool
}

type AppInstance struct {
	Ready       bool
	File        *BimmerFile
	Running     bool
	Window      *sdl.Window
	Surface     *sdl.Surface
	Renderer    *Renderer
	Font        *ttf.Font
	GraphOffset int32
	ScaleFactor int32
	Offset      int32
	EditMode    bool
	Status      string
	Mode        int
	Options     *Options
}

/*
Modes:
0 ready
1 open file
 */

func (instance *AppInstance) Setup() {
	if instance.Ready {
		return
	}
	instance.File = nil
	instance.Running = false
	instance.Ready = true
	instance.GraphOffset = 0
	instance.ScaleFactor = 1
	instance.Offset = 0
	instance.EditMode = false
	instance.Status = ""
	instance.Mode = 0
	instance.Renderer = &Renderer{Instance: instance, RenderCache: make(map[string]*CacheStruct), CacheReady: false}
	instance.Options = &Options{OpenPath: "", LeftCtrl: false, RenderLine: true}
}
func (instance *AppInstance) RunLoop() {
	instance.Running = true
	for instance.Running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				instance.HandleKeyPress(t)
				break
			case *sdl.MouseMotionEvent:
				instance.Renderer.setMousePos(t.X)
				break
			case *sdl.WindowEvent:
				s, _ := instance.Window.GetSurface()
				instance.Surface = s
				instance.Renderer.CacheReady = false
				break
			case *sdl.TextInputEvent:
				instance.HandleTextInput(t.GetText())
				break

			case *sdl.QuitEvent:
				instance.Running = false
				break
			}
		}
		instance.ComputeStatus()
		instance.Renderer.Update()
	}
}
func (instance *AppInstance) LoadFile(url string) {
	result, err := ioutil.ReadFile(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := string(result)
	instance.File = ParseBimmerFile(content)
	instance.Window.SetTitle("BimmerLink Reader: " + url)
	instance.File.FileName = url
}

func (instance *AppInstance) HandleKeyPress(event *sdl.KeyboardEvent) {
	if event.Keysym.Sym == sdl.K_LCTRL {
		instance.Options.LeftCtrl = !instance.Options.LeftCtrl
		return
	}
	if event.State != 1 {
		return
	}
	var code = event.Keysym.Scancode
	var m = instance.Mode
	var filePresent = instance.File != nil
	//fmt.Println(code)

	// o
	if code == 18 {
		if m == 0 {
			instance.Mode = 1
		}
	}
	//Backspace
	if code == 42 {
		if m == 1 {
			if instance.Options.LeftCtrl {
				instance.Options.OpenPath = ""
			} else if instance.Options.OpenPath != "" {
				instance.Options.OpenPath = instance.Options.OpenPath[0:len(instance.Options.OpenPath) -1]
			}
		}
	}
	// e
	if code == 8 {
		if instance.Options.LeftCtrl && filePresent {
			instance.ClearFile()
		} else if filePresent && m == 0 {
			if instance.ScaleFactor < 200 {
				instance.ScaleFactor++
				instance.MinorStatus("Zoom In " + strconv.Itoa(int(instance.ScaleFactor)), time.Second * 3)
				instance.Renderer.CacheReady = false
			}
		}
	}
	//q
	if code == 20 {
		if filePresent && m == 0 {
			if instance.ScaleFactor > 0 {
				instance.ScaleFactor--
				instance.MinorStatus("Zoom Out " + strconv.Itoa(int(instance.ScaleFactor)), time.Second * 3)
				instance.Renderer.CacheReady = false
			}
		}
	}
	//escape
	if code == 41 {
		if m == 1 {
			instance.Mode = 0
		}
	}
	//l
	if code == 15 {
		if  m == 0 {
			instance.Options.RenderLine = !instance.Options.RenderLine
			instance.MinorStatus("Toggle time line", time.Second * 5)
		}
	}

	//s
	if code == 22 {
		if m == 0 && filePresent{
			if instance.GraphOffset < int32( len(instance.File.Names) - 4) {
				instance.GraphOffset++
				instance.MinorStatus("Graphs Down " + strconv.Itoa(int(instance.GraphOffset)), time.Second * 10)
				instance.Renderer.CacheReady = false
			}
		}
	}
	//w
	if code == 26 {
		if m == 0 && filePresent{
			if instance.GraphOffset > 0 {
				instance.GraphOffset--
				instance.MinorStatus("Graphs Up " + strconv.Itoa(int(instance.GraphOffset)), time.Second * 10)
				instance.Renderer.CacheReady = false
			}
		}
	}
	//v
	if code == 25 {
			fmt.Println(event.Keysym.Sym)
		if instance.Options.LeftCtrl {
			if m == 1 {
				text, _ := sdl.GetClipboardText()
				instance.Options.OpenPath += text
			}
		}
	}
	//a
	if code == 4 {
		if m == 0 && filePresent {
			instance.Offset--
			instance.MinorStatus("Offset Left " + strconv.Itoa(int(instance.Offset)), time.Second * 10)
			instance.Renderer.CacheReady = false


		}
	}
	//l
	if code == 7 {
		if m == 0 && filePresent {
			instance.Offset++
			instance.MinorStatus("Offset Right " + strconv.Itoa(int(instance.Offset)), time.Second * 10)
			instance.Renderer.CacheReady = false


		}
	}
	//Enter
	if code == 40 {
		if m == 1 {
			instance.LoadFile(instance.Options.OpenPath)
			instance.Renderer.CacheReady = false
			instance.Mode = 0
		}
	}
}

func (instance *AppInstance) HandleTextInput(text string) {
	if instance.Mode == 1 {
		instance.Options.OpenPath += text
	}
}

func (instance *AppInstance) ComputeStatus() {
	var m = instance.Mode
	var status = ""
	now := time.Now()
	if m == 0 {
		if instance.Options.ToShow != "" && now.Before(instance.Options.Until) {
			status += "[" + instance.Options.ToShow + " | READY]"
		} else {
			status += "[READY] "

		}
		if instance.File != nil {
			status += instance.File.FileName
		} else {
			status += "No file"
		}
	}
	if m == 1 {
		status = "[Open File] " + instance.Options.OpenPath
	}
	instance.Status = status
}

func (instance *AppInstance) MinorStatus(s string, duration time.Duration) {
	instance.Options.ToShow = s
	instance.Options.Until = time.Now().Add(duration)
}

func (instance *AppInstance) ClearFile() {
	instance.File = nil
	instance.Options.ToShow = ""
	instance.Window.SetTitle("BimmerLink Reader")

}

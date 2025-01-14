package main

import (
	imgui "imgui/gui"
	"log"
)

type MyApp struct {
}

func (m *MyApp) Render(gui *imgui.GUI) {
	debug := false
	gui.StartRow()
	gui.StartCellWithHeader("Buttons")
	gui.Text("Hello world")
	if gui.Button("Press me") {
		log.Println("button pressed")
	}
	gui.Text("Here is more text")
	if gui.Button("Next button") {
		log.Println("next button pressed")
	}
	gui.EndCell()

	gui.StartCell()
	debug = gui.Button("Debug")
	gui.EndCell()

	gui.EndRow()
	if debug {
		gui.Debug()
	}
}

func main() {
	app := &MyApp{}
	imgui.Run(app)
}

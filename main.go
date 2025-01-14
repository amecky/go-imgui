package main

import (
	"fmt"
	imgui "imgui/gui"
	"log"

	"github.com/amecky/table/table"
)

var INTERVALLS = []string{"5m", "15m", "30m", "1h", "4h", "1d", "1w"}

type View interface {
	Render(gui *imgui.GUI)
}

type MyApp struct {
	views      []View
	activeView int
}

type InputView struct {
	input      string
	inputState bool
}

func (iv *InputView) Render(gui *imgui.GUI) {
	gui.StartRow()
	gui.StartCell()
	iv.input, iv.inputState = gui.Input("Input:", iv.input, iv.inputState, 30)
	gui.EndCell()
	gui.EndRow()
}

type TableView struct{}

func (tv *TableView) Render(gui *imgui.GUI) {
	gui.StartRow()
	gui.StartCell()
	tbl := table.New().Headers("One", "Two", "Three", "Block")
	for i := 0; i < 10; i++ {
		r := tbl.CreateRow()
		r.AddDefaultText(fmt.Sprintf("%d", i+1))
		st := 1
		if i < 5 {
			st = -1
		}
		r.AddFloat(float64(i)*5.0, st)
		r.AddInt(i+1, 0)
		r.AddBlock(i%2 == 0)
	}
	gui.Table(tbl)
	gui.EndCell()
	gui.EndRow()
}

type TickerView struct {
	selectedIntervall int
	steps             int
	radio             bool
	input             string
	inputState        bool
}

func (m *TickerView) Render(gui *imgui.GUI) {
	gui.StartRow()
	gui.StartCell()
	gui.StartGroup()
	gui.Text("Make some selections")
	ret := gui.Selection("Intervall:", INTERVALLS, m.selectedIntervall)
	if ret == -1 {
		m.selectedIntervall--
		if m.selectedIntervall < 0 {
			m.selectedIntervall = 0
		}
	}
	if ret == 1 {
		m.selectedIntervall++
		if m.selectedIntervall >= len(INTERVALLS) {
			m.selectedIntervall = len(INTERVALLS) - 1
		}
	}
	m.steps = gui.IntSlider("Num:", 0, 100, m.steps, 10)
	m.radio = gui.Checbox("Toggle Me", m.radio)
	if gui.Button("Reload") {
		log.Println("RELOAD")
	}
	gui.EndGroup()
	gui.EndCell()
	gui.EndRow()

	gui.StartRow()
	gui.StartCell()
	m.input, m.inputState = gui.Input("Input:", m.input, m.inputState, 30)
	gui.EndCell()
	gui.EndRow()
}

func (m *MyApp) Render(gui *imgui.GUI) {

	debug := false

	gui.BeginMenuBar()
	if gui.BeginMenu("File") {
		if gui.MenuItem("Open..") {
			log.Println("Open")
		}
		if gui.MenuItem("Save") {
			log.Println("Save")
		}
		if gui.MenuItem("Close") {
			log.Println("Close")
		}

	}
	gui.EndMenu()
	if gui.BeginMenu("Views") {

		if gui.MenuItem("InputView") {
			m.activeView = 0
		}
		if gui.MenuItem("TickerView") {
			m.activeView = 1
		}
		if gui.MenuItem("TableView") {
			m.activeView = 2
		}
	}
	gui.EndMenu()
	gui.EndMenuBar()

	m.views[m.activeView].Render(gui)

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
	/*
		gui.StartRow()
		gui.StartCell()
		m.selected, m.selectedDropDown = gui.DropDown("Intervalls:", INTERVALLS, m.selected, m.selectedDropDown)
		gui.Text("Selected Intervall: " + INTERVALLS[m.selected])
		gui.EndCell()
		gui.EndRow()
	*/

	if debug {
		gui.Debug()
	}
}

func main() {
	app := &MyApp{}
	app.views = append(app.views, &InputView{})
	app.views = append(app.views, &TickerView{
		input: "TESLA",
	})
	app.views = append(app.views, &TableView{})
	imgui.Run(app)
}

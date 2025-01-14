package imgui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type appModel struct {
	gui *GUI
	app App
}

type App interface {
	Render(gui *GUI)
}

// Init initializes the model
func (m appModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.gui.SendKey(msg.String())
		if !m.gui.saveInput {
			switch msg.String() {
			case "q": // Quit the app
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		log.Println("Resizing  == > X", msg.Width, "Y", msg.Height)
		m.gui = NewGUI(msg.Width, msg.Height)
	case tea.MouseMsg:
		mouseEvent := tea.MouseEvent(msg)
		if m.gui != nil && mouseEvent.Action == tea.MouseActionMotion {
			m.gui.SetMousePos(mouseEvent)
		}
		if m.gui != nil && mouseEvent.Action != tea.MouseActionMotion && mouseEvent.Action == tea.MouseActionPress {
			m.gui.SetMouseEvent(mouseEvent)
		}
		if mouseEvent.Action != tea.MouseActionMotion && mouseEvent.Action == tea.MouseActionPress {
			log.Printf("%+v at %d %d\n", mouseEvent, mouseEvent.X, mouseEvent.Y)
		}
	}
	return m, nil
}

// View renders the UI
func (m appModel) View() string {
	if m.gui != nil {
		m.gui.Begin()
		m.app.Render(m.gui)
		return m.gui.End()
	}
	return ""
}

func Run(app App) {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	// Initialize the program
	p := tea.NewProgram(appModel{app: app}, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting app: %v\n", err)
		os.Exit(1)
	}
}

package imgui

import (
	"fmt"
	"log"
	"strings"

	"github.com/amecky/table/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Menu struct {
	active   string
	curX     int
	curY     int
	menuPos  int
	itemPos  int
	menuSize int
}

type rect struct {
	x int
	y int
	w int
	h int
}

func (r rect) Inside(px, py int) bool {
	if px == r.x && py == r.y {
		return true
	}
	if px == r.x && py == r.y+r.h {
		return true
	}
	if px == r.x+r.w && py == r.y+r.h {
		return true
	}
	if px == r.x+r.w && py == r.y {
		return true
	}
	return px >= r.x && px <= r.x+r.w && py >= r.y && py <= r.y+r.h
}

type GUI struct {
	width      int
	height     int
	mouseEvent tea.MouseEvent
	processed  int
	buffer     *Buffer
	input      string
	saveInput  bool
	mouseX     int
	mouseY     int

	useMenu  bool
	menuPos  int
	itemPos  int
	menuSize int

	menu Menu

	started bool
}

func NewGUI(w, h int) *GUI {
	return &GUI{
		width:     w,
		height:    h,
		processed: -1,
		buffer:    NewBuffer(w, h),
	}
}

func (g *GUI) Begin() {
	g.buffer.Clear()
	g.started = true
	g.StartRow()
	g.buffer.StartCell()
	g.menuPos = 0
}

func (g *GUI) End() string {
	if g.buffer.grouping {
		g.EndGroup()
	}
	g.buffer.EndCell()
	g.buffer.EndRow()
	g.processed = -1
	return g.buffer.String()
}

func (g *GUI) SetMouseEvent(e tea.MouseEvent) {
	g.mouseEvent = e
	g.processed = 1
}

func (g *GUI) SetMousePos(e tea.MouseEvent) {
	g.mouseX = e.X
	g.mouseY = e.Y
}

func (g *GUI) SendKey(s string) {
	if g.saveInput {
		if s == "backspace" {
			if len(g.input) > 0 {
				g.input = g.input[0 : len(g.input)-1]
			}
		} else if s == "enter" {
			g.saveInput = false
		} else {
			g.input += s
		}
		log.Println(g.input)
	}
}

func (g *GUI) Debug() {
	g.buffer.Debug()
}

func (g *GUI) Text(text string) {
	g.buffer.PushID(text)
	g.buffer.Write(text, 0, false)
	g.buffer.PopID()
}

func (g *GUI) Button(text string) bool {
	g.buffer.PushID("BUTTON_" + text)
	g.buffer.Write(" "+text+" ", OK_BUTTON_STYLE, false)
	ret := false
	if g.processed == 1 && g.buffer.HasFocus(g.mouseEvent.X, g.mouseEvent.Y) {
		log.Println("Button " + text + " pressed")
		g.processed = -1
		ret = true
	}
	g.buffer.PopID()
	return ret
}

// Selection writes the label with an down and up arrow and the selected text
// Returns the index of the selected item
func (g *GUI) Selection(label string, lines []string, selected int) int {
	g.buffer.PushID("SELECTION_" + lines[0])
	curPos := g.buffer.CurrentPos()
	sel := selected
	curPos.x += internalLen(label) + 1
	if g.processed == 1 && curPos.Matches(g.mouseEvent.X, g.mouseEvent.Y) {
		g.processed = -1
		sel--
		if sel < 0 {
			sel = len(lines) - 1
		}
	}
	l := findMaxLen(lines) + 2
	curPos.x += l + 1
	if g.processed == 1 && curPos.Matches(g.mouseEvent.X, g.mouseEvent.Y) {
		g.processed = 1
		sel++
		if sel >= len(lines) {
			sel = 0
		}
	}

	g.buffer.Write(label+" ", 0, true)
	g.buffer.Write("⯇", ARROW_STYLE, true)
	txt := formatString(" "+lines[sel]+" ", l, table.AlignCenter)
	g.buffer.Write(txt, INPUT_STYLE, true)
	g.buffer.Write("⯈", ARROW_STYLE, true)
	g.buffer.PopID()
	return sel
}

func (g *GUI) IntSlider(label string, min, max, value, steps int) int {
	g.buffer.PushID("INT_SLIDE_" + fmt.Sprintf(" %d ", value))
	g.buffer.Write(label+" ", 0, true)
	g.buffer.Write("⯇", ARROW_STYLE, true)
	if g.processed == 1 && g.buffer.HasFocus(g.mouseEvent.X, g.mouseEvent.Y) {
		value -= steps
		if value < min {
			value = min
		}
		g.processed = -1
	}
	l := len(fmt.Sprintf("%d", max)) + 2
	g.buffer.Write(formatString(fmt.Sprintf(" %d ", value), l, table.AlignCenter), INPUT_STYLE, true)
	g.buffer.Write("⯈", ARROW_STYLE, true)
	if g.processed == 1 && g.buffer.HasFocus(g.mouseEvent.X, g.mouseEvent.Y) {
		value += steps
		if value > max {
			value = max
		}
		g.processed = -1
	}
	g.buffer.PopID()
	return value
}

func (g *GUI) Checkbox(label string, active bool) bool {
	g.buffer.PushID("CHECKBOX_" + label)
	if active {
		g.buffer.Write("■", ARROW_STYLE, true)
	} else {
		g.buffer.Write("▢", ARROW_STYLE, true)
	}
	if g.processed == 1 && g.buffer.HasFocus(g.mouseEvent.X, g.mouseEvent.Y) {
		if active {
			active = false
		} else {
			active = true
		}
		g.processed = -1
	}
	g.buffer.Write(" "+label, 0, true)
	g.buffer.PopID()
	return active
}

func (g *GUI) Radio(label string, entries []string, selected int) int {
	g.buffer.PushID("RADIO_" + label)
	g.buffer.Write(label+" ", 0, true)
	ret := selected
	curPos := g.buffer.CurrentPos()
	for i := 0; i < len(entries); i++ {
		if g.processed == 1 && curPos.Matches(g.mouseEvent.X, g.mouseEvent.Y) {
			ret = i
			g.processed = -1
		}
		if i == ret {
			g.buffer.Write("■", ARROW_STYLE, true)
		} else {
			g.buffer.Write("▢", ARROW_STYLE, true)
		}

		g.buffer.Write(" "+entries[i]+" ", 0, true)
		curPos.x += len(entries[i]) + 3
	}
	g.buffer.PopID()
	return ret
}

func (g *GUI) DropDown(label string, lines []string, selected int, active bool) (int, bool) {
	g.buffer.PushID("DROPDOWN_" + lines[0])
	g.buffer.Write(label+" ", 0, true)
	if active {
		g.buffer.Write("⯆", ARROW_STYLE, true)
	} else {
		g.buffer.Write("⯈", ARROW_STYLE, true)
	}
	if g.processed == 1 && g.buffer.HasFocus(g.mouseEvent.X, g.mouseEvent.Y) {
		if active {
			active = false
		} else {
			active = true
		}
		g.processed = -1
	}
	g.buffer.Write(" "+lines[selected], 0, false)
	ret := active
	sel := selected

	if active {
		for i, s := range lines {
			st := 0
			if g.buffer.IsInside(g.mouseX, g.mouseY, 10, 0) {
				st = 1
			}
			if g.processed == 1 && g.buffer.IsInside(g.mouseX, g.mouseY, 10, 0) {
				sel = i
				ret = false
				st = 1
				g.processed = -1
			}
			g.buffer.Write(" "+s, st, false)

		}
	}
	g.buffer.PopID()
	return sel, ret
}

func (g *GUI) Input(label, text string, active bool, size int) (string, bool) {
	g.buffer.PushID("INPUT_" + label)
	g.buffer.Write(label+" ", 0, true)
	ret := text
	if active && !g.saveInput {
		active = false
	}
	l := len(text)
	t := g.input
	if len(t) > size {
		t = t[0:size]
		g.input = t
	}
	if !active {
		t = text
	}
	d := size - l
	if d > 0 {
		t += strings.Repeat(" ", d)
	}
	if active {
		g.buffer.Write(t, INPUT_ACTIVE_STYLE, true)
		ret = g.input
	} else {
		g.buffer.Write(t, INPUT_STYLE, true)
	}
	if g.processed == 1 && g.buffer.HasFocus(g.mouseEvent.X, g.mouseEvent.Y) {
		if active {
			active = false
			g.saveInput = false
			g.input = text
		} else {
			g.input = text
			g.saveInput = true
			active = true
		}
		g.processed = -1
	}
	if active {
		ret = g.input
	}

	g.buffer.PopID()
	return ret, active
}

func (g *GUI) StartGroup() {
	g.buffer.grouping = true
}

func (g *GUI) EndGroup() {
	g.buffer.grouping = false
	g.buffer.curY++
	g.buffer.curX = 0
}

func (g *GUI) StartRow() {
	g.buffer.StartRow()
}

func (g *GUI) EndRow() {
	g.buffer.EndRow()
}

func (g *GUI) StartCell() {
	g.StartCellWithHeader("")
}

func (g *GUI) StartCellWithHeader(title string) {
	if g.started {
		g.started = false
		if len(g.buffer.cells) > 0 {
			c := &g.buffer.cells[0]
			c.title = title
		}
	} else {
		g.buffer.StartCellWithHeader(title)
	}
}

func (g *GUI) EndCell() {
	g.buffer.EndCell()
}

func (g *GUI) BeginMenuBar() {
	g.useMenu = true
	g.buffer.useMenu = true
}

func (g *GUI) EndMenuBar() {
}

func (g *GUI) BeginMenu(label string) bool {
	id := "MENU_" + label
	g.itemPos = 1
	g.buffer.WriteEx(g.menuPos, 0, " "+label+" ", 1)
	ret := false
	if g.menu.active == id {
		ret = true
	}
	if g.processed == 1 && g.buffer.HasFocus(g.mouseEvent.X, g.mouseEvent.Y) {
		g.processed = -1
		if ret {
			ret = false
			g.menu.active = ""
		} else {
			ret = true
			g.menu.active = id
		}
	}
	g.menuSize = internalLen(label) + 3
	return ret
}

func (g *GUI) EndMenu() {
	g.menuPos += g.menuSize

}

func (g *GUI) MenuItem(label string) bool {
	txt := " " + label + " "
	d := 20 - len(txt)
	if d > 0 {
		txt += strings.Repeat(" ", d)
	}
	g.buffer.WriteEx(g.menuPos, g.itemPos, txt, 1)
	ret := false
	if g.processed == 1 {
		r := rect{
			x: g.menuPos,
			y: g.itemPos,
			w: len(txt),
			h: 0,
		}
		if r.Inside(g.mouseX, g.mouseY) {
			g.processed = -1
			ret = true
			g.menu.active = ""
		}
	}
	g.itemPos++
	return ret
}

func (g *GUI) Table(rt *table.Table) {
	var sizes = make([]int, 0)
	for _, th := range rt.TableHeaders {
		sizes = append(sizes, internalLen(th.Text))
	}
	total := 0
	for _, r := range rt.Rows {
		for j, c := range r.Cells {
			if internalLen(c.Text) > sizes[j] {
				sizes[j] = internalLen(c.Text)
			}
		}
	}
	for _, s := range sizes {
		total += s + rt.PaddingSize*2
	}
	for i, h := range rt.TableHeaders {
		g.buffer.Write(rt.BorderStyle.H_LINE, 0, true)
		g.buffer.Write(strings.Repeat(" ", rt.PaddingSize), 0, true)
		g.buffer.Write(formatString(h.Text, sizes[i], table.AlignCenter), 0, true)
		g.buffer.Write(strings.Repeat(" ", rt.PaddingSize), 0, true)
	}
	g.buffer.Write(rt.BorderStyle.H_LINE, 0, false)

	g.buffer.Write(rt.BorderStyle.LEFT_DEL, 0, true)
	for i, s := range sizes {
		g.buffer.Write(strings.Repeat("-", s+rt.PaddingSize*2), 0, true)
		if i < len(sizes)-1 {
			g.buffer.Write(rt.BorderStyle.CROSS, 0, true)
		}
	}
	g.buffer.Write(rt.BorderStyle.RIGHT_DEL, 0, false)

	for _, r := range rt.Rows {
		for i, c := range r.Cells {
			st := c.Marker
			if st != 0 {
				if st == -1 {
					st = TABLE_RED
				} else if st == 1 {
					st = TABLE_LIGHT_GREEN
				} else {
					st = TABLE_RED - 2 + st
				}
			}
			g.buffer.Write(rt.BorderStyle.H_LINE, 0, true)
			g.buffer.Write(strings.Repeat(" ", rt.PaddingSize), 0, true)
			str := formatString(c.Text, sizes[i], c.Alignment)
			g.buffer.Write(str, st, true)
			g.buffer.Write(strings.Repeat(" ", rt.PaddingSize), 0, true)
		}
		g.buffer.Write(rt.BorderStyle.H_LINE, 0, false)
	}
}

func formatString(txt string, length int, align table.TextAlign) string {
	var ret string
	d := length - internalLen(txt)
	if d < 0 {
		d = 0
	}
	// left
	if align == table.AlignLeft {
		ret = txt
		if d > 0 {
			ret += strings.Repeat(" ", d)
		}
	}
	// right
	if align == table.AlignRight {
		if d > 0 {
			ret = strings.Repeat(" ", d)
		}
		ret += txt
	}
	// center
	if align == table.AlignCenter {
		d /= 2
		if d > 0 {
			ret = strings.Repeat(" ", d)
		}
		ret += txt
		d = length - internalLen(txt) - d
		if d > 0 {
			ret += strings.Repeat(" ", d)
		}
	}
	return ret
}

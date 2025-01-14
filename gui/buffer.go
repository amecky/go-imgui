package imgui

import (
	"log"
	"strings"
)

type Vec struct {
	x int
	y int
}

func (v Vec) Matches(x, y int) bool {
	return x == v.x && y == v.y
}

type Row struct {
	cells []int
}

type cell struct {
	title string
	rect
}

type DrawCommand struct {
	uid     string
	focus   bool
	text    string
	style   int
	x       int
	y       int
	size    int
	cellIdx int
}

type Buffer struct {
	width       int
	height      int
	size        int
	chars       []rune
	styles      []int
	currentUID  string
	uids        *Stack
	commands    []DrawCommand
	rows        []Row
	cells       []cell
	curX        int
	curY        int
	grouping    bool
	groupMargin int
	curCell     int
	useMenu     bool
	states      map[string]bool
}

func NewBuffer(w, h int) *Buffer {
	sz := w * h
	ret := &Buffer{
		width:       w,
		height:      h,
		size:        sz,
		chars:       make([]rune, sz),
		styles:      make([]int, sz),
		uids:        &Stack{},
		grouping:    false,
		groupMargin: 1,
		states:      make(map[string]bool),
	}
	return ret
}

func (b *Buffer) Clear() {
	for i := 0; i < b.size; i++ {
		b.chars[i] = ' '
		b.styles[i] = 0
	}
	b.currentUID = ""
	b.commands = b.commands[:0]
	b.cells = b.cells[:0]
	b.rows = b.rows[:0]
	b.curX = 0
	b.curY = 0
	b.curCell = 0
}

func (b *Buffer) PushID(s string) {
	b.uids.Push(s)
}

func (b *Buffer) PopID() {
	if b.grouping {
		b.curX += b.groupMargin
	} else {
		//b.curY++
		c := b.cells[len(b.cells)-1]
		b.curX = c.x
	}
	b.uids.Pop()
}

func (b *Buffer) SetFocus() {
	id := b.uids.Top()
	for _, c := range b.commands {
		c.focus = false
		if c.uid == id {
			c.focus = true
		}
	}
}

func (b *Buffer) CurrentPos() Vec {
	return Vec{
		x: b.curX,
		y: b.curY,
	}
}

func (b *Buffer) At(x, y int) (rune, int) {
	if x < b.width && y < b.height {
		idx := y*b.width + x
		if idx < b.size {
			return b.chars[idx], b.styles[idx]
		}
	}
	return ' ', 0
}

func (b *Buffer) Write(txt string, style int, inline bool) {
	id := b.uids.Top()
	b.commands = append(b.commands, DrawCommand{
		uid:     id,
		style:   style,
		text:    txt,
		focus:   false,
		x:       b.curX,
		y:       b.curY,
		size:    internalLen(txt),
		cellIdx: b.curCell,
	})
	if inline {
		b.curX += internalLen(txt)
	} else {
		if b.grouping {
			b.curX += internalLen(txt)
		} else {
			cc := b.cells[len(b.cells)-1]
			b.curX = cc.x
			b.curY++
		}
	}
}

func (b *Buffer) WriteEx(x, y int, txt string, style int) {
	id := b.uids.Top()
	b.commands = append(b.commands, DrawCommand{
		uid:     id,
		style:   style,
		text:    txt,
		focus:   true,
		x:       x,
		y:       y,
		size:    internalLen(txt),
		cellIdx: -1,
	})
}

func (b *Buffer) IsInside(x, y, w, h int) bool {
	r := rect{
		x: b.curX,
		y: b.curY,
		w: w,
		h: h,
	}
	return r.Inside(x, y)
}

func (b *Buffer) HasFocus(x, y int) bool {
	if len(b.commands) > 0 {
		cmd := b.commands[len(b.commands)-1]
		r := rect{
			x: cmd.x,
			y: cmd.y,
			w: cmd.size,
			h: 0,
		}
		if r.Inside(x, y) {
			return true
		}
	}
	return false
}

func (b *Buffer) Debug() {
	log.Println("---------------------")
	for _, c := range b.commands {
		log.Printf("%+v\n", c)
	}
	log.Println("---------------------")
	for _, c := range b.cells {
		log.Printf("%+v\n", c)
	}
	log.Println("---------------------")
}

func (b *Buffer) Set(x, y int, c rune, style int) {
	if x < b.width && y < b.height {
		idx := y*b.width + x
		if idx < b.size {
			b.chars[idx] = c
			b.styles[idx] = style
		}
	}
}

func (b *Buffer) String() string {
	// fill buffer
	for _, c := range b.cells {
		for i := c.x; i < c.x+c.w; i++ {
			b.Set(i, c.y-1, '─', BORDER)
		}
		if c.title != "" {
			b.WriteEx(c.x+1, c.y-1, c.title, HEADER_STYLE)
			//b.WriteEx(x, c.y-1, " "+c.title+" ", HEADER_STYLE)
		}
		for i := c.x; i < c.x+c.w-1; i++ {
			b.Set(i, c.y+c.h, '─', BORDER)
		}
		for i := c.y; i < c.y+c.h; i++ {
			b.Set(c.x-1, i, '│', BORDER)
		}
		for i := c.y; i < c.y+c.h; i++ {
			b.Set(c.x+c.w-1, i, '│', BORDER)
		}
		b.Set(c.x-1, c.y-1, '┌', BORDER)
		b.Set(c.x+c.w-1, c.y-1, '┐', BORDER)
		b.Set(c.x-1, c.y+c.h, '└', BORDER)
		b.Set(c.x+c.w-1, c.y+c.h, '┘', BORDER)

	}

	for _, c := range b.commands {
		if !c.focus {
			for i, ch := range c.text {
				b.Set(c.x+i, c.y, ch, c.style)
			}
		}
	}

	for _, c := range b.commands {
		if c.focus {
			for i, ch := range c.text {
				b.Set(c.x+i, c.y, ch, c.style)
			}
		}
	}
	// convert buffer to string
	sb := strings.Builder{}

	for y := 0; y < b.height-1; y++ {
		cy := y * b.width
		for x := 0; x < b.width; x++ {
			if b.styles[cy+x] == 0 {
				sb.WriteString(string(b.chars[cy+x]))
			} else {
				st := STYLES[b.styles[cy+x]-1]
				sb.WriteString(st.Convert(string(b.chars[cy+x])))
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (b *Buffer) StartRow() {
	b.rows = append(b.rows, Row{})

}
func (b *Buffer) StartCell() {
	b.StartCellWithHeader("")
}

func (b *Buffer) StartCellWithHeader(title string) {
	b.cells = append(b.cells, cell{
		title: title,
	})
	if len(b.cells) == 1 {
		b.curX = 1
		b.curY = 1
		if b.useMenu {
			b.curY++
		}
	} else {
		prev := b.cells[len(b.cells)-2]
		b.curX = prev.w + 2
		b.curY = prev.y
		if len(b.rows) > 1 {
			b.curX = 1
			b.curY = b.cells[b.curCell].y + b.cells[b.curCell].h + 2
		}
	}
	cur := &b.cells[len(b.cells)-1]
	cur.x = b.curX
	cur.y = b.curY
	b.curCell = len(b.cells) - 1
	cr := &b.rows[len(b.rows)-1]
	cr.cells = append(cr.cells, b.curCell)
}

func (b *Buffer) EndCell() {
	if len(b.cells) > 0 {
		cidx := b.curCell
		cur := &b.cells[cidx]
		for _, c := range b.commands {
			if c.cellIdx == cidx {
				if c.size+c.x > cur.w {
					cur.w = c.x + c.size + 1
				}
				h := c.y - cur.y
				if h >= cur.h {
					cur.h = h + 1
				}
			}
		}
		if cur.x+internalLen(cur.title) > cur.w {
			cur.w = cur.x + internalLen(cur.title)
		}

	}
}

func (b *Buffer) EndRow() {
	cr := &b.rows[len(b.rows)-1]
	my := 0
	for _, c := range cr.cells {
		cell := b.cells[c]
		if cell.h > my {
			my = cell.h
		}
	}
	for _, c := range cr.cells {
		cell := &b.cells[c]
		cell.h = my
	}
}

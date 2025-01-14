package imgui

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/amecky/table/table"
	tea "github.com/charmbracelet/bubbletea"
)

func TestFormatString(t *testing.T) {
	txt := formatString("12345", 7, table.AlignCenter)
	assert.Equal(t, " 12345 ", txt)

	txt = formatString("12345", 8, table.AlignCenter)
	assert.Equal(t, " 12345  ", txt)

	txt = formatString("12", 8, table.AlignCenter)
	assert.Equal(t, "   12   ", txt)

	txt = formatString("12", 8, table.AlignLeft)
	assert.Equal(t, "12      ", txt)

	txt = formatString("12", 8, table.AlignRight)
	assert.Equal(t, "      12", txt)
}

func TestSelection(t *testing.T) {
	gui := NewGUI(20, 4)
	gui.Begin()
	gui.StartRow()
	gui.StartCell()
	entries := []string{"12", "34", "567"}
	me := tea.MouseEvent{
		X: 6,
		Y: 1,
	}
	gui.SetMouseEvent(me)
	ret := gui.Selection("Test", entries, 2)
	gui.EndCell()
	gui.EndRow()
	gui.End()
	assert.Equal(t, 1, ret)
}

func TestText(t *testing.T) {
	gui := NewGUI(20, 4)
	gui.Begin()
	txt := "Hello"
	gui.Text(txt)
	gui.End()
	for i, c := range txt {
		r, s := gui.buffer.At(1+i, 1)
		assert.Equal(t, r, c)
		assert.Equal(t, 0, s)
	}
	firstRow := []rune{'┌', '─', '─', '─', '─', '─', '─', '┐'}
	for i, c := range firstRow {
		r, s := gui.buffer.At(i, 0)
		assert.Equal(t, r, c)
		assert.Equal(t, 0, s)
	}
	middleRow := []rune{'│', 'H', 'e', 'l', 'l', 'o', ' ', '│'}
	for i, c := range middleRow {
		r, s := gui.buffer.At(i, 1)
		assert.Equal(t, r, c)
		assert.Equal(t, 0, s)
	}
	lastRow := []rune{'└', '─', '─', '─', '─', '─', '─', '┘'}
	for i, c := range lastRow {
		r, s := gui.buffer.At(i, 2)
		assert.Equal(t, r, c)
		assert.Equal(t, 0, s)
	}
}

func TestStartCellWithHeader(t *testing.T) {
	gui := NewGUI(20, 4)
	gui.Begin()
	txt := "Hello"
	gui.StartRow()
	gui.StartCellWithHeader("Testing")
	gui.Text(txt)
	gui.EndCell()
	gui.EndRow()
	gui.End()
	fmt.Println(gui.buffer)
	for i, c := range txt {
		r, s := gui.buffer.At(1+i, 1)
		assert.Equal(t, r, c)
		assert.Equal(t, 0, s)
	}
	firstRow := []rune{'┌', 'T', 'E', 'S', 'T', 'I', 'N', '┐'}
	for i, c := range firstRow {
		r, s := gui.buffer.At(i, 0)
		assert.Equal(t, r, c)
		assert.Equal(t, 0, s)
	}
	middleRow := []rune{'│', 'H', 'e', 'l', 'l', 'o', ' ', '│'}
	for i, c := range middleRow {
		r, s := gui.buffer.At(i, 1)
		assert.Equal(t, r, c)
		assert.Equal(t, 0, s)
	}
	lastRow := []rune{'└', '─', '─', '─', '─', '─', '─', '┘'}
	for i, c := range lastRow {
		r, s := gui.buffer.At(i, 2)
		assert.Equal(t, r, c)
		assert.Equal(t, 0, s)
	}
}

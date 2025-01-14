package imgui

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
)

func TestCells(t *testing.T) {
	b := NewBuffer(20, 5)
	b.Clear()
	b.StartRow()
	b.StartCell()
	b.Write("Test1", 0, false)
	b.EndCell()
	b.StartCell()
	b.Write("Test2", 0, false)
	b.Write("Test3", 0, false)
	b.EndCell()
	b.EndRow()
	b.Debug()
	fmt.Println(b)
	assert.Equal(t, 2, len(b.cells))
}
